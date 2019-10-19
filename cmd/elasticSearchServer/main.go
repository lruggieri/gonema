package main

import (
	"fmt"
	"github.com/lruggieri/utils/netutil"
	"github.com/lruggieri/gonema/pkg/database/elastic"
	"github.com/lruggieri/gonema/pkg/utils"
	"log"
	"net/http"
	"os"
	"strconv"
)


var elasticConnection *elastic.Connection
func main(){
	esHost := os.Getenv("GONEMAES_SERVER_HOST")
	esPort := os.Getenv("GONEMAES_SERVER_PORT")
	if len(esHost) == 0{
		esHost = "http://localhost"
	}
	if len(esPort) == 0{
		esPort = "9200"
	}

	esCon,err := elastic.New(esHost,esPort)
	//esCon,err := elastic.New("http://localhost","9200")
	if err != nil{
		log.Fatal(err)
	}
	elasticConnection = esCon

	mux := http.NewServeMux()
	mux.Handle("/complete",netutil.HandleWithError(AutocompleteHandler))

	port := os.Getenv("GONEMAES_API_PORT")
	if port == "" {
		port = "8080"
	}
	utils.Logger.Info("running API on port "+port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))

}

func AutocompleteHandler(w http.ResponseWriter, r *http.Request) netutil.ResponseLayout{

	urlParameters := r.URL.Query()
	requestedField, requestedText := urlParameters.Get("field"),urlParameters.Get("text")
	requestedSize := urlParameters.Get("size")
	var requestedSizeInt int
	if len(requestedField) == 0{
		return netutil.ResponseLayout{StatusCode:http.StatusBadRequest,Error:"'field' parameter not specified"}
	}
	if len(requestedText) == 0{
		return netutil.ResponseLayout{StatusCode:http.StatusBadRequest,Error:"'text' parameter not specified"}
	}
	if requestedSizeConverted, err := strconv.Atoi(requestedSize) ; err != nil || requestedSizeConverted <= 0{
		return netutil.ResponseLayout{StatusCode:http.StatusBadRequest,Error:"'size' parameter invalid"}
	}else{
		requestedSizeInt = requestedSizeConverted
	}

	suggestions, err := elasticConnection.SuggestQuery("gonema",requestedField,requestedText,requestedSizeInt)
	if err != nil{
		return netutil.ResponseLayout{StatusCode:http.StatusInternalServerError,Error:err.Error(),IsInternalError:true}
	}
	if err != nil{
		return netutil.ResponseLayout{StatusCode:http.StatusInternalServerError,Error:err.Error(),IsInternalError:true}
	}
	return netutil.ResponseLayout{StatusCode:http.StatusOK,Response:suggestions}
}