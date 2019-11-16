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

func GetTorrents(iKeyword string) (interface{}, error) {
	gonemapiHost := os.Getenv("GONEMAES_API_HOST")
	gonemapiPort := os.Getenv("GONEMAES_API_PORT")
	mountPoint := "/torrents"

	requestHostPort := gonemapiHost
	if len(gonemapiPort) > 0{
		requestHostPort += ":" + gonemapiPort
	}
	requestHostPort += mountPoint

	client := http.Client{
		Timeout:60*3*time.Second,
	}
	req, err := http.NewRequest("GET",requestHostPort,nil)
	if err != nil{
		return nil, err
	}

	reqQuery := req.URL.Query()
	reqQuery.Add("key", iKeyword)

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

	var decodedResp netutil.ResponseLayout
	err = json.Unmarshal(body,&decodedResp)
	if err != nil{
		return nil,err
	}

	return decodedResp.Response, nil
}
