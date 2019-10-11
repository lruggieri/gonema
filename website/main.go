package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/ruggieri/gonema/pkg/utils"
	"gitlab.com/ruggieri/gonema/website/controller"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const(
	environmentDir = "configs/env"
)

var(
	templatesDir string
	staticAssetDir string
)

// neuteredFileSystem is used to prevent directory listing of static assets
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	// Check if path exists
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	// If path exists, check if is a file or a directory.
	// If is a directory, stop here with an error saying that file
	// does not exist. So user will get a 404 error code for a file/directory
	// that does not exist, and for directories that exist.
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	// If file exists and the path is not a directory, let's return the file
	return f, nil
}


func main() {
	templatesDir = os.Getenv("TEMPLATES_DIR")
	staticAssetDir = os.Getenv("STATIC_ASSET_DIR")

	mux := http.NewServeMux()
	fs := http.FileServer(neuteredFileSystem{http.Dir(staticAssetDir)})
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/favicon.ico",faviconHandler)
	mux.Handle("/central",handleWithError(centralControllerHandler))
	mux.HandleFunc("/", mainPageHandler)

	var tlsCertPath = os.Getenv("TLS_CERT_PATH")
	var tlsKeyPath = os.Getenv("TLS_KEY_PATH")

	if len(tlsCertPath) > 0 && len(tlsKeyPath) > 0{
		err := http.ListenAndServeTLS(":443", tlsCertPath, tlsKeyPath, nil)
		if err != nil{
			panic(err)
		}
	}else{
		port := os.Getenv("PORT")
		if len(port) == 0{
			port = "8080"
		}
		fmt.Println("web service starting on port "+port)
		err := http.ListenAndServe(fmt.Sprintf(":%s",port),mux)
		if err != nil{
			panic(err)
		}
	}


}

type internalError struct{
	Error error //internal error, not to display
	Message string //message to display to the client
	Code int //return code
}
type clientError struct{
	Error error `json:"error"`
	AdditionalInfo string `json:"additional_info"`
}
type handleWithError func(http.ResponseWriter, *http.Request) *internalError
func (hwe handleWithError) ServeHTTP(w http.ResponseWriter, r *http.Request){
	if err := hwe(w, r); err != nil {
		clientError := clientError{Error:errors.New(err.Message)}

		w.Header().Set("content-type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(err.Code)
		if err.Error != nil{
			//TODO log better
			utils.Logger.Error(err)
		}
		encodeError := json.NewEncoder(w).Encode(clientError)
		dealWithEncodingError(w,encodeError)
	}
}
func dealWithEncodingError(w http.ResponseWriter, iEncodingError error){
	if iEncodingError != nil{
		//TODO log better
		_,_ = fmt.Fprintln(w, `{"error":"something seriously wrong happen on our side, we are sorry for the inconvenient"}`)
	}
}


func mainPageHandler(w http.ResponseWriter, r *http.Request){
	mainPage := filepath.Join(templatesDir,"index.tmpl")

	tmpl := template.Must(template.ParseFiles(mainPage))

	err := tmpl.Execute(w,nil)
	if err != nil{
		panic(err)
	}

}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r,path.Join(staticAssetDir,"images","favicon.ico"))
}
func centralControllerHandler(w http.ResponseWriter, r *http.Request) *internalError {

	if r.Method == http.MethodPost{
		action := r.FormValue("action")
		resourceName := r.FormValue("resourceName")
		resourceImdbID := r.FormValue("resourceImdbID")

		switch action {
		case "getResourceInfo":{
			resources, err := controller.GetResourceInfo(resourceName, resourceImdbID)
			if err != nil{
				return &internalError{Error:err,Message:"internal error",Code:http.StatusInternalServerError}
			}
			respondResourceInfo(w, resources)
		}
		default:
			return &internalError{Error: nil,Message:"action '"+action+"' not recognized",Code:http.StatusBadRequest}
		}

		return nil
	}else{
		return &internalError{Error: nil,Message:"expecting POST request to central, got "+r.Method,Code:http.StatusBadRequest}
	}
}


type resourceInfoResponse struct{
	Resources interface{} `json:"resources"`
}
func respondResourceInfo(w http.ResponseWriter, iResponseResources interface{}){

	w.Header().Set("content-type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	encodeError := json.NewEncoder(w).Encode(resourceInfoResponse{Resources:iResponseResources})
	dealWithEncodingError(w,encodeError)
}