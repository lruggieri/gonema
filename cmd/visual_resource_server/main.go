package main

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/ruggieri/gonema/pkg/utils"
	"gitlab.com/ruggieri/gonema/pkg/visual_resource"
	"log"
	"net/http"
)

func main(){
	utils.DebugActive = true
	utils.Logger.Level = logrus.DebugLevel

	mux := http.NewServeMux()
	mux.HandleFunc("/",emptyRequest)
	mux.HandleFunc("/resourceInfo",resourceInfo)

	log.Fatal(http.ListenAndServe(":8080", mux))
}


func emptyRequest(w http.ResponseWriter, r *http.Request){
	http.NotFound(w,r)
}

func resourceInfo(w http.ResponseWriter, r *http.Request){
	requestParameters:= r.URL.Query()
	if imdbID:= requestParameters.Get("imdbID") ; len(imdbID) > 0{
		resource,err := visual_resource.GetResource("",imdbID)
		if err != nil{
			dealWithInternalError(w,err)
		}else{
			respond(w,http.StatusOK,[]byte(resource.String()))
		}
	}else if resourceTitle := requestParameters.Get("resourceTitle") ; len(resourceTitle) > 0{
		respond(w,http.StatusOK,[]byte("Sorry, resourceTitle not handled yet"))
	}else{
		respond(w,http.StatusBadRequest,[]byte("Please, specify 'imdbID' or 'resourceTitle'"))
	}
}

func dealWithInternalError(w http.ResponseWriter, iErr error){
	utils.Logger.Error(iErr)
	http.Error(w, "" +
		"We are very sorry, but something on our side has broken. " +
		"This issue has been reported and will be dealt as soon as possible by our engineering team.", http.StatusInternalServerError)
}
func respond(w http.ResponseWriter, iStatusCode int, iMessage []byte){
	w.WriteHeader(iStatusCode)
	_, err := w.Write(iMessage)
	if err != nil {
		dealWithInternalError(w,err)
	}
}

