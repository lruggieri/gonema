package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/lruggieri/utils"
	"gitlab.com/ruggieri/gonema/pkg/database/initialdata"
	"log"
	"net/http"
	"os"
	"path"
)

var currentDirectoryPath = utils.GetCallerPaths(1)[0]

//maps required index => path to relative mapping (path elements to be joined)
var requiredIndices = map[string][]string{
	"gonema": {currentDirectoryPath,"mappings", "gonema"},
}

type elasticDB struct{
	connection *elastic.Client
}

//checks that the input ES has the proper indices. If not, initialize them
func checkEsDB(iElasticClient *elastic.Client) (oError error){

	for requiredIndexName, requiredIndexPath := range requiredIndices{
		res, err := iElasticClient.Indices.Get([]string{requiredIndexName})
		if err != nil{
			return err
		}

		var resBody  map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			return errors.New("Error parsing the response body: "+ err.Error())
		}

		if res.StatusCode != http.StatusOK{
			if res.StatusCode == http.StatusNotFound{
				//index not found, we have to create it
				indexTemplateFile, err := os.Open(path.Join(requiredIndexPath...))
				if err != nil{
					return err
				}
				indexCreationRequest := iElasticClient.Indices.Create.WithBody(indexTemplateFile)

				res, err = iElasticClient.Indices.Create("gonema",indexCreationRequest)
				if err != nil{
					return err
				}
				if res.StatusCode != 200{
					if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
						return errors.New("Error parsing the response body: "+ err.Error())
					}
					resBytes, _ := json.Marshal(resBody)
					return errors.New("cannot create index template for index '"+requiredIndexName+"'. Response: "+string(resBytes))
				}

				//now insert initial data
				basicDataBuilder := initialdata.Builder{}
				basicMovies, err := basicDataBuilder.GetMovies()
				if err != nil{
					log.Fatal(err)
				}

				fmt.Println("got",len(basicMovies),"basic movies")

			}else{
				resBytes, _ := json.Marshal(resBody)
				return errors.New("cannot get indices info. Resp: "+string(resBytes))
			}
		}else{
			//TODO check that at least we have some movie document
		}
	}



	return nil
}

func New(iHost, iPort string) (oElasticDB *elasticDB, oErr error){
	cfg := elastic.Config{
		Addresses: []string{
			iHost +":"+ iPort,
		},
	}
	es, err := elastic.NewClient(cfg)
	if err != nil{
		return nil, err
	}

	err = checkEsDB(es)
	if err != nil{
		return nil, err
	}

	return &elasticDB{connection:es}, nil
}