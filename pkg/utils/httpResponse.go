package utils

import "net/http"

func DealWithInternalError(w http.ResponseWriter, iErr error){
	Logger.Error(iErr)
	http.Error(w, "" +
		"We are very sorry, but something on our side has broken. " +
		"This issue has been reported and will be dealt as soon as possible by our engineering team.", http.StatusInternalServerError)
}

func Respond(w http.ResponseWriter, iStatusCode int, iMessage []byte){
	w.WriteHeader(iStatusCode)
	_, err := w.Write(iMessage)
	if err != nil {
		DealWithInternalError(w,err)
	}
}
