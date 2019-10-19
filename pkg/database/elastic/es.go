package elastic

import (
	"context"
	"errors"
	"github.com/lruggieri/gonema/pkg/database/initialdata"
	"github.com/lruggieri/utils"
	"github.com/olivere/elastic/v7"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

const(
	bulkSize = 2000
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
			utils.Logger.Info("creating index "+requiredIndexName)
			//index not found, we have to create it
			indexTemplateFile, err := os.Open(path.Join(requiredIndexBlocks.templatePath...))
			if err != nil{
				return err
			}

			indexTemplateBytes, err := ioutil.ReadAll(indexTemplateFile)
			if err != nil{
				return err
			}

			utils.Logger.Info("using template "+string(indexTemplateBytes))

			indexCreationResult, err := iElasticClient.CreateIndex(requiredIndexName).BodyString(string(indexTemplateBytes)).Do(context.Background())
			if err != nil{
				return err
			}
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
				utils.Logger.Info("no document found in index "+requiredIndexName+". Inserting...")
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
	currentBulkElements := 0
	totalInsertion := 0

	type expandedMovie struct{
		initialdata.Movie
		//SuggestionsName suggestion `json:"suggest_name"`
	}

	//index each movie
	for _,movieToIndex := range basicMovies{
		if movieToIndex.Id <= 0{continue}

		expandedMovie := expandedMovie{
			Movie:movieToIndex,
			//SuggestionsName:suggestion{Input:strings.Split(movieToIndex.Name," ")},
			//SuggestionsName:suggestion{Input:getStringGrams(movieToIndex.Name," ",2)},
		}

		bulkRequest.Add(elastic.NewBulkIndexRequest().Index(iIndexName).Id(movieToIndex.ImdbID).Doc(expandedMovie))
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

			totalInsertion += currentBulkElements
			utils.Logger.Info("inserted:" + strconv.Itoa(totalInsertion) + " movies")

			bulkRequest.Reset()
			currentBulkElements = 0
		}
	}

	return nil
}

/*
Return every possible n-gram of the input string after splitting it with the input separator.
Limit returned NGrams to iMinGrams.
Input string put as the first element of the resulting slice (if it passes checks like iMinGrams).
*/
func getStringGrams(iString ,iSeparator string, iMinGrams int) (oPermutations []string){
	if iMinGrams <= 0{
		iMinGrams = 1
	}

	initialSet := strings.Split(iString,iSeparator) //starting set

	/*
	Idea: insert every single element first; the check every element inserted before and add the current element
		between every word of that previous element.
	Complexity: O(n)
	*/

	finalNGrams := make([][]string,0)
	for _, initialElement := range initialSet{
		actualGramsLen := len(finalNGrams)
		for nGramIdx := 0; nGramIdx < actualGramsLen ; nGramIdx++{
			nGram := finalNGrams[nGramIdx]
			for nGramPartIdx := 0; nGramPartIdx <= len(nGram) ; nGramPartIdx++{
				//inserting initialElement in position nGramPartIdx
				partBefore := nGram[nGramPartIdx:]
				partAfter := nGram[:nGramPartIdx]
				newGram := append(append(partBefore,initialElement),partAfter...)
				finalNGrams = append(finalNGrams,newGram)
			}
		}
		finalNGrams = append(finalNGrams,[]string{initialElement})
	}

	finalStrings := make([]string,0,len(finalNGrams))
	for _, finalNGram := range finalNGrams{
		if len(finalNGram) >= iMinGrams{
			tempString := strings.Join(finalNGram, iSeparator)
			if tempString == iString{continue} //the initial string has to be the first element
			finalStrings = append(finalStrings, tempString)
		}
	}

	// we want the initial string at the beginning (so it's easier to find it out when getting suggestions form ES)
	if len(strings.Split(iString, iSeparator)) >= iMinGrams{
		finalStrings = append([]string{iString},finalStrings...)
	}


	return finalStrings

}

func New(iHost, iPort string) (oElasticDB *Connection, oErr error){
	esUrl := iHost
	if len(iPort) > 0 {
		esUrl += ":" + iPort
	}
	es, err := elastic.NewClient(elastic.SetSniff(false),elastic.SetURL(esUrl))
	if err != nil{
		return nil, err
	}

	err = checkEsDB(es)
	if err != nil{
		return nil, err
	}

	return &Connection{connection: es}, nil
}