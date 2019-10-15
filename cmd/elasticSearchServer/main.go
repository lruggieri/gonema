package main

import (
	"fmt"
	"gitlab.com/ruggieri/gonema/pkg/database/elastic"
	"log"
	"net/http"
	"os"
)


var elasticConnection *elastic.Connection
func main(){
	//esCon,err := elastic.New(os.Getenv("GONEMAES_HOST"),os.Getenv("GONEMAES_PORT"))
	esCon,err := elastic.New("http://localhost","9200")
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
	fmt.Println("running API on port "+port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))

}


func AutocompleteHandler(w http.ResponseWriter, r *http.Request){

	//TODO parse request and dela with the suggestion


	suggestions, err := elasticConnection.Suggest("gonema","suggest_name","harry po",20)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(suggestions)
}