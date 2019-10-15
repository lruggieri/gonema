package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Connection struct{
	connection *elastic.Client
}

func (c *Connection) Suggest(iField string, iText string) (oSuggestions []string,oErr error){

	type suggestQuery struct{
		Suggest struct{

		}`json:"suggest"`
	}

	size := 10
	request := esapi.SearchRequest{
		Index:[]string{"gonema"},
		SuggestField:iField,
		SuggestText:iText,
		SuggestSize:&size,
	}

	var resBody  map[string]interface{}
	res, _ := request.Do(context.Background(),c.connection)
	defer res.Body.Close()
	if res.IsError(){
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			return nil,errors.New("Error parsing the response body: "+ err.Error())
		}
		resBytes, _ := json.Marshal(resBody)
		return nil,errors.New("Cannot perform suggestion. Response: "+string(resBytes))
	}

	fmt.Println(resBody)
	return nil, nil
}
