package controller

import (
	"encoding/json"
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

	requestHostPort := gonemapiHost
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

	var decodedResp interface{}
	err = json.Unmarshal(body,&decodedResp)
	if err != nil{
		return nil,err
	}

	return decodedResp, nil
}