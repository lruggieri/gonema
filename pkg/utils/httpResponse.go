package utils

import (
	"encoding/json"
	"net/http"
)


const ResponseMessageForInternalError = "We are very sorry, but something on our side has broken. " +
	"This issue has been reported and will be dealt as soon as possible by our engineering team."

func DealWithInternalError(w http.ResponseWriter, iErr error){
	Logger.Error(iErr)
	http.Error(w, "" +
		"We are very sorry, but something on our side has broken. " +
		"This issue has been reported and will be dealt as soon as possible by our engineering team.", http.StatusInternalServerError)
}

type ResponseLayout struct{
	Response interface{} `json:"response"`
	Error string `json:"error,omitempty"`
	IsInternalError bool `json:"-"`
}

func Respond(w http.ResponseWriter, iStatusCode int, iResponse ResponseLayout){

	if iStatusCode > 299{
		if iResponse.IsInternalError{
			iStatusCode = http.StatusInternalServerError
			Logger.Error(iResponse.Error)
			iResponse.Error = ResponseMessageForInternalError
		}
		iResponse.Response = nil
	}
	w.Header().Set("content-type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(iStatusCode)

	jsonMessage, _ := json.Marshal(iResponse)
	_, err := w.Write(jsonMessage)
	if err != nil {
		DealWithInternalError(w, err)
	}

}