package main

import (
	"gitlab.com/ruggieri/gonema/pkg/database/elastic"
	"log"
	"net/http"
)


var elasticConnection *elastic.Connection
func main(){
	//esCon,err := elastic.New(os.Getenv("GONEMAES_HOST"),os.Getenv("GONEMAES_PORT"))
	esCon,err := elastic.New("http://localhost","9200")
	if err != nil{
		log.Fatal(err)
	}
	elasticConnection = esCon

	http.HandleFunc("/complete",AutocompleteHandler)

	_, err = elasticConnection.Suggest("suggest_name","harry po")
	if err != nil{
		log.Fatal(err)
	}

}


func AutocompleteHandler(w http.ResponseWriter, r *http.Request){

}