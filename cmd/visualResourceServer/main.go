package main

import (
	"fmt"
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
			utils.Respond(w,http.StatusOK,[]byte(cachedResult.(string)))
		}else{
			resource,err := visual_resource.GetResources("",imdbID)
			resourceJson := resource.Json()
			if err != nil{
				utils.DealWithInternalError(w,err)
			}else{
				utils.Respond(w,http.StatusOK,[]byte(resourceJson))
			}
			localCache.Insert(resourceImdbIDElementCacheKey, utils.CacheElementKey(imdbID), resourceJson)
		}
	}else if resourceTitle := requestParameters.Get(resourceNameElementCacheKey) ; len(resourceTitle) > 0{
		if cachedResult := localCache.Fetch(resourceNameElementCacheKey, utils.CacheElementKey(resourceTitle)) ; cachedResult != nil{
			utils.Respond(w,http.StatusOK,[]byte(cachedResult.(string)))
		}else{
			utils.Respond(w,http.StatusOK,[]byte("Sorry, "+resourceNameElementCacheKey+" not handled yet"))
			//TODO save in cache
		}
	}else{
		utils.Respond(w,http.StatusBadRequest,[]byte(
			"Please, specify '"+resourceImdbIDElementCacheKey+"' or '"+resourceNameElementCacheKey+"'"))
	}
}