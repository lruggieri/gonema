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

//GetAggregations calls Elasticsearch to get data aggregations (eg. most shared movies/series)
//It returns either nil or a ResponseLayout.Response object
func GetAggregations(iAggType, iresType string) (interface{}, error) {
	gonemapiHost := os.Getenv("GONEMAES_API_HOST")
	gonemapiPort := os.Getenv("GONEMAES_API_PORT")
	mountPoint := "/aggregations"

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
	reqQuery.Add("aggtype", iAggType)
	reqQuery.Add("restype", iresType)

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
