package main

import (
	"gitlab.com/ruggieri/gonema/pkg/database/elastic"
	"log"
)

func main(){
	_,err := elastic.New("http://localhost","9200")
	if err != nil{
		log.Fatal(err)
	}
}