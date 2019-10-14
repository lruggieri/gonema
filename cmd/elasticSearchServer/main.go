package main

import (
	"gitlab.com/ruggieri/gonema/pkg/database/elastic"
	"log"
	"os"
)

func main(){
	_,err := elastic.New(os.Getenv("GONEMAES_HOST"),os.Getenv("GONEMAES_PORT"))
	if err != nil{
		log.Fatal(err)
	}
}