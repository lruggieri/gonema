package elastic

import (
	"context"
	"encoding/json"
	"errors"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/lruggieri/utils"
	"gitlab.com/ruggieri/gonema/pkg/database/initialdata"
	"net/http"
	"os"
	"path"
	"strings"
)

var currentDirectoryPath = utils.GetCallerPaths(1)[0]

//maps required index => path to relative mapping (path elements to be joined)
var requiredIndices = map[string][]string{
	"gonema": {currentDirectoryPath,"mappings", "gonema.json"},
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
				defer res.Body.Close()
				if res.IsError(){
					if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
						return errors.New("Error parsing the response body: "+ err.Error())
					}
					resBytes, _ := json.Marshal(resBody)
					return errors.New("cannot create index template for index '"+requiredIndexName+"'. Response: "+string(resBytes))
				}

				//now insert initial data
				err = insertInitialData(iElasticClient, "gonema")
				if err != nil{
					return err
				}

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

func insertInitialData(iElasticClient *elastic.Client, iIndexName string) (oErr error){
	basicDataBuilder := initialdata.Builder{}
	basicMovies, err := basicDataBuilder.GetMovies()
	if err != nil{
		return err
	}

	bulkSize := 5000
	currentBulkElements := 0
	currentBulkBody := strings.Builder{}
	type bulkMovie struct{
		Movie initialdata.Movie `json:"-"`
		Id string `json:"_id"`
	}

	//index each movie
	for _,movieToIndex := range basicMovies{
		if movieToIndex.Id <= 0{continue}

		marshBMovie,err := json.Marshal(movieToIndex)
		if err != nil{
			return err
		}
		//action_and_meta_data
		currentBulkBody.WriteString(`{"index":{"_index":"`+iIndexName+`", "_id":"`+movieToIndex.ImdbID+`"}}`)
		currentBulkBody.WriteString("\n")
		//optional_source (not so optional in my opinion, but...)
		currentBulkBody.Write(marshBMovie)
		currentBulkBody.WriteString("\n")
		currentBulkElements++

		if currentBulkElements > 0 && currentBulkElements % bulkSize == 0{
			request := esapi.BulkRequest{
				Index:iIndexName,
				Body:strings.NewReader(currentBulkBody.String()),
				Refresh:"true",
			}
			resp, err := request.Do(context.Background(),iElasticClient)
			if err != nil{
				return err
			}
			defer resp.Body.Close()

			var r map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
				return errors.New("Error parsing the response body: "+err.Error())
			}
			if resp.IsError(){
				marshResp, err := json.Marshal(r)
				if err != nil{
					return err
				}
				return errors.New("error inserting movie "+movieToIndex.ImdbID+": "+resp.Status()+".Full resp: "+string(marshResp))
			}else{
			}

			//reset next bulk
			currentBulkBody = strings.Builder{}
			currentBulkElements = 0
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