package elastic

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"gitlab.com/ruggieri/gonema/pkg/database/initialdata"
)

type Connection struct{
	connection *elastic.Client
}

func (c *Connection) Suggest(iIndex string,iField string, iText string, iSize int) (oSuggestions []string,oErr error){

	suggester := elastic.NewCompletionSuggester(iField).
		Field(iField).
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
	for _, ops := range suggestions.Suggest[iField] {
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