package elastic

import (
	"context"
	"errors"
	"github.com/lruggieri/utils"
	"github.com/olivere/elastic/v7"
	"gitlab.com/ruggieri/gonema/pkg/database/initialdata"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

var currentDirectoryPath = utils.GetCallerPaths(1)[0]

type indexBlocks struct{
	templatePath []string
	insertionFunction func(iElasticClient *elastic.Client, iIndexName string) (oErr error)
}

//maps required index => path to relative mapping (path elements to be joined)
var requiredIndices = map[string]indexBlocks{
	"gonema": {
		templatePath:[]string{currentDirectoryPath,"mappings", "gonema.json"},
		insertionFunction:insertInitialMovieData,
	},
}



//checks that the input ES has the proper indices. If not, initialize them
func checkEsDB(iElasticClient *elastic.Client) (oError error){

	for requiredIndexName, requiredIndexBlocks := range requiredIndices{
		indexExist, err := iElasticClient.IndexExists(requiredIndexName).Do(context.Background())
		if err != nil{
			return err
		}

		if !indexExist{
			//index not found, we have to create it
			indexTemplateFile, err := os.Open(path.Join(requiredIndexBlocks.templatePath...))
			if err != nil{
				return err
			}

			indexTemplateBytes, err := ioutil.ReadAll(indexTemplateFile)
			if err != nil{
				return err
			}

			indexCreationResult, err := iElasticClient.CreateIndex(requiredIndexName).BodyString(string(indexTemplateBytes)).Do(context.Background())

			if !indexCreationResult.Acknowledged{
				return errors.New(requiredIndexName+" index creation not acknowledged")
			}

			//now insert initial data
			err = requiredIndexBlocks.insertionFunction(iElasticClient, "gonema")
			if err != nil{
				return err
			}
		}else {
			//checking if there is some document in the index

			documentsNumber, err := iElasticClient.Count(requiredIndexName).Do(context.Background())
			if err != nil{
				return err
			}
			if documentsNumber == 0{
				//if there is no document, insert initial ones
				err = requiredIndexBlocks.insertionFunction(iElasticClient, "gonema")
				if err != nil {
					return err
				}
			}
		}
	}



	return nil
}

//call functions to build initial data an insert the in ES
func insertInitialMovieData(iElasticClient *elastic.Client, iIndexName string) (oErr error){
	basicDataBuilder := initialdata.Builder{}
	basicMovies, err := basicDataBuilder.GetMovies()
	if err != nil{
		return err
	}

	bulkRequest := iElasticClient.Bulk()
	bulkSize := 5000 //number of movies for each bulk request
	currentBulkElements := 0

	//index each movie
	for _,movieToIndex := range basicMovies{
		if movieToIndex.Id <= 0{continue}

		bulkRequest.Add(elastic.NewBulkIndexRequest().Index(iIndexName).Id(movieToIndex.ImdbID).Doc(movieToIndex))
		currentBulkElements++

		if currentBulkElements > 0 && currentBulkElements % bulkSize == 0{
			bulkResponse, err := bulkRequest.Do(context.Background())
			if err != nil{
				return err
			}

			indexed := bulkResponse.Indexed()
			if len(indexed) != currentBulkElements{
				return errors.New("tried to index "+strconv.Itoa(currentBulkElements)+" but " +
					"successfully indexed "+strconv.Itoa(len(indexed)))
			}

			bulkRequest.Reset()
			currentBulkElements = 0
		}
	}

	return nil
}

func New(iHost, iPort string) (oElasticDB *Connection, oErr error){
	esUrl := iHost
	if len(iPort) > 0 {
		esUrl += ":" + iPort
	}
	es, err := elastic.NewClient(elastic.SetURL(esUrl))
	if err != nil{
		return nil, err
	}

	err = checkEsDB(es)
	if err != nil{
		return nil, err
	}

	return &Connection{connection: es}, nil
}