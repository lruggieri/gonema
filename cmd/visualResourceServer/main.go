package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gitlab.com/ruggieri/gonema/pkg/utils"
	"gitlab.com/ruggieri/gonema/pkg/visual_resource"
	"log"
	"net/http"
	"os"
	"time"
)


const(
	resourceImdbIDElementCacheKey = "imdbID"
	resourceNameElementCacheKey = "resourceName"
)
var(
	localCache = utils.NewCache()
)

func main(){
	utils.DebugActive = true
	utils.Logger.Level = logrus.DebugLevel

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: "+err.Error()+"")
	}

	localCache.SetNewRootElementDuration(resourceImdbIDElementCacheKey, time.Hour)
	localCache.SetNewRootElementDuration(resourceNameElementCacheKey, 2 * time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("/",emptyRequest)
	mux.HandleFunc("/resourceInfo",resourceInfo)

	port := os.Getenv("VRS_PORT")
	if port == "" {
		port = "8080"
	}
	utils.Logger.Info("running API on port "+port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))

}


func emptyRequest(w http.ResponseWriter, r *http.Request){
	http.NotFound(w,r)
}

func resourceInfo(w http.ResponseWriter, r *http.Request){
	requestParameters:= r.URL.Query()
	if imdbID := requestParameters.Get(resourceImdbIDElementCacheKey) ; len(imdbID) > 0{
		if cachedResult := localCache.Fetch(resourceImdbIDElementCacheKey, utils.CacheElementKey(imdbID)) ; cachedResult != nil{
			respond(w,http.StatusOK,[]byte(cachedResult.(string)))
		}else{
			resource,err := visual_resource.GetResources("",imdbID)
			resourceJson := resource.Json()
			if err != nil{
				dealWithInternalError(w,err)
			}else{
				respond(w,http.StatusOK,[]byte(resourceJson))
			}
			localCache.Insert(resourceImdbIDElementCacheKey, utils.CacheElementKey(imdbID), resourceJson)
		}
	}else if resourceTitle := requestParameters.Get(resourceNameElementCacheKey) ; len(resourceTitle) > 0{
		if cachedResult := localCache.Fetch(resourceNameElementCacheKey, utils.CacheElementKey(resourceTitle)) ; cachedResult != nil{
			respond(w,http.StatusOK,[]byte(cachedResult.(string)))
		}else{
			respond(w,http.StatusOK,[]byte("Sorry, "+resourceNameElementCacheKey+" not handled yet"))
			//TODO save in cache
		}
	}else{
		respond(w,http.StatusBadRequest,[]byte(
			"Please, specify '"+resourceImdbIDElementCacheKey+"' or '"+resourceNameElementCacheKey+"'"))
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
