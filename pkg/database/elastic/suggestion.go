package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/lruggieri/gonema/pkg/database/initialdata"
	"github.com/lruggieri/gonema/pkg/utils"
	"github.com/olivere/elastic/v7"
	"strings"
)

//to enable inserting tokenized suggestions
type suggestion struct{
	Input []string `json:"input"`
}

//uses ES query on the user input
func (c *Connection) SuggestQuery(iIndex string, iField string, iText string, iSize int) (oSuggestions []string,oErr error){
	oSuggestions = make([]string,0)

	esSuggestField := ""
	switch iField {
	case "name":{esSuggestField = "suggest_name"}
	case "imdbId":{esSuggestField = "suggest_imdb_id"}
	default:
		return nil,errors.New("requested suggest field '"+iField+"' is invalid")
	}

	query := elastic.NewMatchQuery(esSuggestField,iText)

	res, err := c.connection.Search(iIndex).Query(query).Size(iSize).Do(context.Background())
	if err != nil{
		return nil, err
	}

	if res.Hits != nil{
		if res.Hits.Hits != nil{
			for _, hit := range res.Hits.Hits{
				var resMovie initialdata.Movie
				err := json.Unmarshal(hit.Source, &resMovie)
				if err != nil {
					return nil, err
				}
				oSuggestions = append(oSuggestions, resMovie.Name)
			}
		}
	}

	return oSuggestions,nil
}

//uses standard ES suggestion
func (c *Connection) SuggestStandard(iIndex string, iField string, iText string, iSize int) (oSuggestions []string,oErr error){

	esSuggestField := ""
	switch iField {
	case "name":{esSuggestField = "suggest_name"}
	case "imdbId":{esSuggestField = "suggest_imdb_id"}
	default:
		return nil,errors.New("requested suggest field '"+iField+"' is invalid")
	}

	oSuggestions = make([]string,0)

	suggester := elastic.NewCompletionSuggester(esSuggestField).
		Field(esSuggestField).
		Text(iText).
		Size(iSize)
	suggesterSource := elastic.NewSearchSource().Suggester(suggester)
	suggestions, err := c.connection.Search().
		Index(iIndex).
		SearchSource(suggesterSource).
		Do(context.Background())
	if err != nil{
		return nil,err
	}

	for _, ops := range suggestions.Suggest[esSuggestField] {
		for _, op := range ops.Options {
			if op.Source == nil {
				continue
			}
			var suggestion map[string]suggestion
			err := json.Unmarshal(op.Source, &suggestion)
			if err != nil {
				return nil, err
			}
			var resultingMovieName string
			if suggestField, ok := suggestion[esSuggestField] ; ok{
				resultingMovieName = strings.Join(suggestField.Input," ")
			}else{
				utils.Logger.Error("cannot find field "+esSuggestField+" in suggestion result")
			}

			oSuggestions = append(oSuggestions, resultingMovieName)
		}
	}

	return oSuggestions, nil
}

