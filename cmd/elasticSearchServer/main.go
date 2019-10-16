package main

import (
	"fmt"
	"gitlab.com/ruggieri/gonema/pkg/database/elastic"
	"gitlab.com/ruggieri/gonema/pkg/utils"
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
	mux.HandleFunc("/complete",AutocompleteHandler)

	port := os.Getenv("GONEMAES_API_PORT")
	if port == "" {
		port = "8080"
	}
	utils.Logger.Info("running API on port "+port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))

}

type ResponseLayout struct{
	Response interface{} `json:"response"`
	Error error `json:"error,omitempty"`
}


func AutocompleteHandler(w http.ResponseWriter, r *http.Request){

	urlParameters := r.URL.Query()
	requestedField, requestedText := urlParameters.Get("field"),urlParameters.Get("text")
	requestedSize := urlParameters.Get("size")
	var requestedSizeInt int
	if len(requestedField) == 0{
		utils.Respond(w,http.StatusBadRequest,utils.ResponseLayout{Error:"'field' parameter not specified"})
		return
	}
	if len(requestedText) == 0{
		utils.Respond(w,http.StatusBadRequest,utils.ResponseLayout{Error:"'text' parameter not specified"})
		return
	}
	if requestedSizeConverted, err := strconv.Atoi(requestedSize) ; err != nil || requestedSizeConverted <= 0{
		utils.Respond(w,http.StatusBadRequest,utils.ResponseLayout{Error:"'size' parameter invalid"})
		return
	}else{
		requestedSizeInt = requestedSizeConverted
	}



	suggestions, err := elasticConnection.Suggest("gonema",requestedField,requestedText,requestedSizeInt)
	if err != nil{
		utils.Respond(w,http.StatusInternalServerError,utils.ResponseLayout{Error:err.Error(),IsInternalError:true})
		return
	}
	if err != nil{
		utils.Respond(w,http.StatusInternalServerError,utils.ResponseLayout{Error:err.Error(),IsInternalError:true})
		return
	}
	utils.Respond(w,http.StatusOK,utils.ResponseLayout{Response:suggestions})
}