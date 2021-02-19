package controller

import (
	"encoding/json"
	"errors"
	"github.com/lruggieri/utils/netutil"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//GetResourceInfo calls Elasticsearch to get info about a specific resors, either
//by its ImdbID or its name
//It returns either nil or a ResponseLayout.Response object
func GetResourceInfo(resourceName, resourceImdbID string) (interface{}, error) {
	gonemapiHost := os.Getenv("GONEMAES_API_HOST")
	gonemapiPort := os.Getenv("GONEMAES_API_PORT")
	mountPoint := "/resourceInfo"

	requestHostPort := gonemapiHost
	if len(gonemapiPort) > 0 {
		requestHostPort += ":" + gonemapiPort
	}
	requestHostPort += mountPoint

	client := http.Client{
		Timeout: 60 * 3 * time.Second,
	}
	req, err := http.NewRequest("GET", requestHostPort, nil)
	if err != nil {
		return nil, err
	}

	reqQuery := req.URL.Query()
	reqQuery.Add("resourceName", resourceName)
	reqQuery.Add("imdbId", resourceImdbID)

	req.URL.RawQuery = reqQuery.Encode()
	//req.Header.Set(...,...)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var decodedResp netutil.ResponseLayout
	err = json.Unmarshal(body, &decodedResp)
	if err != nil {
		return nil, err
	}

	return decodedResp.Response, nil
}
