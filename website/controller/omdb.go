package controller

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetResourceInfoFromOmdb(resourceName, resourceImdbID string) (interface{}, error) {
	endPoint := "http://www.omdbapi.com"
	mountPoint := "/"
	apiKey := os.Getenv("OMDB_APIKEY")

	requestUrl := endPoint+mountPoint

	client := http.Client{
		Timeout:5*time.Second,
	}
	req, err := http.NewRequest("GET",requestUrl,nil)
	if err != nil{
		return nil, err
	}

	reqQuery := req.URL.Query()
	reqQuery.Add("apikey",apiKey)
	if len(resourceImdbID) > 0{
		reqQuery.Add("i",resourceImdbID)
	}else{
		reqQuery.Add("s",resourceName)
	}

	req.URL.RawQuery = reqQuery.Encode()
	//req.Header.Set(...,...)

	resp, err := client.Do(req)
	if err != nil{
		return nil,err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	if resp.StatusCode != http.StatusOK{
		return nil, errors.New(string(body))
	}
	if len(resourceImdbID) > 0{
		var omdbResponse OmdbResponseByID
		err = json.Unmarshal(body,&omdbResponse)
		if err != nil{
			return nil,err
		}
		if omdbResponse.IsValid(){
			return omdbResponse,nil
		}else{
			return nil, errors.New("invalid response from Omdb with StatusCode "+strconv.Itoa(resp.StatusCode)+ "" +
				"and error '"+omdbResponse.Error+"'")
		}
	}else{
		var omdbResponse OmdbResponseBySearch
		err = json.Unmarshal(body,&omdbResponse)
		if err != nil{
			return nil,err
		}
		if omdbResponse.IsValid(){
			resultResponse := make([]interface{},0)
			alreadySeenID := make(map[string]bool) //sometime we get duplicated IDs and I do not know why...
			for _, searchResult := range omdbResponse.Search{
				if _,ok := alreadySeenID[searchResult.ImdbID] ; !ok{
					alreadySeenID[searchResult.ImdbID] = true
					resultResponse = append(resultResponse, searchResult)
				}
			}
			return resultResponse,nil
		}else{
			return nil, errors.New("invalid response from Omdb  for query "+req.URL.RawQuery+" " +
				"with StatusCode "+strconv.Itoa(resp.StatusCode)+ "" +
				"and error '"+omdbResponse.Error+"'")
		}
	}
}

type OmbdResponse struct{
	Response string `json:"Response"`
	Error string `json:"Error,omitempty"`
}
func(or *OmbdResponse) IsValid() bool{
	if or.Response == "True"{
		return true
	}
	return false
}
type OmdbResponseByID struct{
	*OmbdResponse
	Title string `json:"Title"`
	Year string `json:"Year"`
	Rated string `json:"Rated"`
	Released string `json:"Released"`
	Runtime string `json:"Runtime"`
	Genre string `json:"Genre"`
	Director string `json:"Director"`
	Writer string `json:"Writer"`
	Actors string `json:"Actors"`
	Plot string `json:"Plot"`
	Language string `json:"Language"`
	Country string `json:"Country"`
	Awards string `json:"Awards"`
	Poster string `json:"Poster"`
	Ratings []struct{
		Source string `json:"Source"`
		Value string `json:"Value"`
	} `json:"Ratings"`
	Metascore string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes string `json:"imdbVotes"`
	ImdbID string `json:"imdbId"`
	Type string `json:"Type"`
	DVD string `json:"Dvd"`
	BoxOffice string `json:"boxOffice"`
	Production string `json:"Production"`
	Website string `json:"Website"`
}
type OmdbResponseBySearch struct{
	*OmbdResponse
	Search []struct{
		Title string `json:"Title"`
		Year string `json:"Year"`
		ImdbID string `json:"imdbID"`
		Poster string `json:"Poster"`
	} `json:"Search"`
	TotalResult int `json:"totalResults,string"`
}