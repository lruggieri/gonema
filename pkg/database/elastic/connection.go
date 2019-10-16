package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"gitlab.com/ruggieri/gonema/pkg/database/initialdata"
)

type Connection struct{
	connection *elastic.Client
}

func (c *Connection) Suggest(iIndex string, iField string, iText string, iSize int) (oSuggestions []string,oErr error){

	esSuggestField := ""
	switch iField {
	case "name":{esSuggestField = "suggest_name"}
	case "imdbId":{esSuggestField = "suggest_imdb_id"}
	default:
		return nil,errors.New("requested suggest field '"+iField+"' is invalid")
	}

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
	oSuggestions = make([]string,0)
	for _, ops := range suggestions.Suggest[esSuggestField] {
		for _, op := range ops.Options {
			if op.Source == nil {
				continue
			}
			var movie initialdata.Movie
			err := json.Unmarshal(op.Source, &movie)
			if err != nil {
				return nil, err
			}
			oSuggestions = append(oSuggestions, movie.Name)
		}
	}

	return oSuggestions, nil
}