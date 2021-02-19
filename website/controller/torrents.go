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

//GetTorrents calls Elasticsearch to get info about torrents for a specific search.
//This search can be performed either providing a keyword or providing a specific resource ImdbID.
//The input type ('movies', 'series', etc.) will be use by ES as filter on results.
func GetTorrents(iKeyword, iImdbID, iType string) (interface{}, error) {
	gonemapiHost := os.Getenv("GONEMAES_API_HOST")
	gonemapiPort := os.Getenv("GONEMAES_API_PORT")
	mountPoint := "/torrents"

	requestHostPort := gonemapiHost
	if len(gonemapiPort) > 0 {
		requestHostPort += ":" + gonemapiPort
	}
	requestHostPort += mountPoint

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", requestHostPort, nil)
	if err != nil {
		return nil, err
	}

	reqQuery := req.URL.Query()
	reqQuery.Add("key", iKeyword)
	reqQuery.Add("imdbID", iImdbID)
	reqQuery.Add("type", iType)

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
