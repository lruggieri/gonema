package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func GetResourceInfo(resourceName, resourceImdbID string) (interface{}, error) {
	gonemapiHost := os.Getenv("GONEMAPI_HOST")
	gonemapiPort := os.Getenv("GONEMAPI_PORT")

	requestHostPort := "http://"+gonemapiHost
	if len(gonemapiPort) > 0{
		requestHostPort += ":" + gonemapiPort
	}
	requestHostPort += "/resourceInfo"
	fmt.Print(requestHostPort)

	client := http.Client{
		Timeout:60*3*time.Second,
	}
	req, err := http.NewRequest("GET",requestHostPort,nil)
	if err != nil{
		return nil, err
	}

	reqQuery := req.URL.Query()
	reqQuery.Add("resourceName",resourceName)
	reqQuery.Add("imdbID",resourceImdbID)

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

	if resp.StatusCode != 200{
		return nil, errors.New(string(body))
	}

	/*resp := `[{"imdb_id":"tt6146586","images":{"small":"https://m.media-amazon.com/images/M/MV5BMDg2YzI0ODctYjliMy00NTU0LTkxODYtYTNkNjQwMzVmOTcxXkEyXkFqcGdeQXVyNjg2NjQwMDQ@._V1_UX182_CR0,0,182,268_AL_.jpg","big":"https://m.media-amazon.com/images/M/MV5BMDg2YzI0ODctYjliMy00NTU0LTkxODYtYTNkNjQwMzVmOTcxXkEyXkFqcGdeQXVyNjg2NjQwMDQ@._V1_SY1000_CR0,0,648,1000_AL_.jpg"},"imdb_score":0,"title":"John Wick 3 - Parabellum","year":2019,"release_date":"0001-01-01T00:00:00Z","categories":["Badass","Immortal","Kratos with nice hair and guns, lots of guns"],"plot":"","stars":null,"writers":null,"directors":null,"available_torrents":[{"magnet_link":"the Magnet URI has no parameters","quality":"SuperDuper","resolution":"42K","sound":"H264","codec":"","name":"Keanu4Evah","seeders":42,"leechers":42},{"magnet_link":"the Magnet URI has no parameters","quality":"SuperDuper2","resolution":"42K","sound":"H264","codec":"","name":"Keanu4Evah2","seeders":42,"leechers":42},{"magnet_link":"the Magnet URI has no parameters","quality":"SuperDuper3","resolution":"42K","sound":"H264","codec":"","name":"Keanu4Evah3","seeders":42,"leechers":42},{"magnet_link":"the Magnet URI has no parameters","quality":"SuperDuper4","resolution":"42K","sound":"H264","codec":"","name":"Keanu4Evah4","seeders":42,"leechers":42},{"magnet_link":"the Magnet URI has no parameters","quality":"SuperDuper5","resolution":"42K","sound":"H264","codec":"","name":"Keanu4Evah5","seeders":42,"leechers":42}]}]`
	var decodedResp interface{}
	err := json.Unmarshal([]byte(resp),&decodedResp)
	if err != nil{
		return nil,err
	}*/

	return string(body), nil
}