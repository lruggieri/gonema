package main

import (
	"fmt"
	"github.com/lruggieri/gonema/website/controller"
	"github.com/lruggieri/utils/netutil"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	if len(templatesDir) == 0{
		templatesDir = "website/templates"
	}
	staticAssetDir = os.Getenv("STATIC_ASSET_DIR")
	if len(staticAssetDir) == 0{
		staticAssetDir = "website/static"
	}

	mux := http.NewServeMux()
	fs := http.FileServer(neuteredFileSystem{http.Dir(staticAssetDir)})
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/favicon.ico",faviconHandler)
	mux.HandleFunc("/robots.txt",robotsHandler)
	mux.Handle("/central", netutil.HandleWithError(centralControllerHandler))
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

func mainPageHandler(w http.ResponseWriter, r *http.Request){
	mainPage := filepath.Join(templatesDir,"index.tmpl")

	tmpl := template.Must(template.ParseFiles(mainPage))

	err := tmpl.Execute(w,nil)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Website down, we are very sorry for the inconvenience."))
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r,path.Join(staticAssetDir,"images","favicon.ico"))
}
func robotsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r,path.Join("robots.txt"))
}
func centralControllerHandler(w http.ResponseWriter, r *http.Request) netutil.ResponseLayout {

	if r.Method == http.MethodPost{
		action := r.FormValue("action")

		switch action {
		case "getResourceInfo":{
			resourceName := r.FormValue("resourceName")
			resourceImdbID := r.FormValue("resourceImdbID")
			if len(resourceName) == 0 && len(resourceImdbID) == 0{
				return netutil.ResponseLayout{StatusCode:http.StatusBadRequest,Error:"no resource name or ID was given"}
			}
			resources, err := controller.GetResourceInfoFromOmdb(resourceName, resourceImdbID)
			if err != nil{
				return netutil.ResponseLayout{StatusCode:http.StatusInternalServerError,Error:err.Error(), IsInternalError:true}
			}
			return netutil.ResponseLayout{StatusCode:http.StatusOK,Response:resources}
		}
		case "getTorrents":{
			keyword := r.FormValue("keyword")
			tType := r.FormValue("type")
			if len(keyword) == 0{
				return netutil.ResponseLayout{StatusCode:http.StatusBadRequest,Error:"invalid keyword when fetching torrents"}
			}
			torrents,err := controller.GetTorrents(keyword,tType)
			if err != nil{
				return netutil.ResponseLayout{StatusCode:http.StatusInternalServerError,Error:err.Error(), IsInternalError:true}
			}
			return netutil.ResponseLayout{StatusCode:http.StatusOK,Response:torrents}
		}
		case "suggest":{
			resourceName := r.FormValue("resourceName")
			if len(resourceName) > 0{

				/*requestUrl := bytes.Buffer{}
				requestUrl.WriteString(os.Getenv("GONEMAES_API_HOST"))
				port := os.Getenv("GONEMAES_API_PORT")
				if len(port) > 0 {
					requestUrl.WriteString(":")
					requestUrl.WriteString(port)
				}
				requestUrl.WriteString("/complete?field=name&text=")
				requestUrl.WriteString(resourceName)
				requestUrl.WriteString("&size=10")
				resp, err := http.Get(requestUrl.String())
				if err != nil{
					return netutil.ResponseLayout{
						StatusCode:http.StatusInternalServerError,
						Error:"cannot get es_api. err: "+err.Error(),
						IsInternalError:true,
					}
				}
				defer resp.Body.Close()

				respBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil{
					return netutil.ResponseLayout{
						StatusCode:http.StatusInternalServerError,
						Error:"cannot read es_api response. err: "+err.Error(),
						IsInternalError:true,
					}
				}
				var respDecoded netutil.ResponseLayout
				err = json.Unmarshal(respBytes,&respDecoded)
				if err != nil{
					return netutil.ResponseLayout{
						StatusCode:http.StatusInternalServerError,
						Error:"cannot decode es_api response. err: "+err.Error(),
						IsInternalError:true,
					}
				}
				if len(respDecoded.Error) > 0{
					return netutil.ResponseLayout{
						StatusCode:http.StatusInternalServerError,
						Error:respDecoded.Error,
					}
				}
				return netutil.ResponseLayout{
					StatusCode:http.StatusOK,
					Response:respDecoded.Response,
				}*/
				resources, err := controller.GetResourceInfoFromOmdb(resourceName, "")
				if err != nil{
					return netutil.ResponseLayout{StatusCode:http.StatusInternalServerError,Error:err.Error(), IsInternalError:true}
				}
				return netutil.ResponseLayout{StatusCode:http.StatusOK,Response:resources}

			}else{
				return netutil.ResponseLayout{
					StatusCode:http.StatusBadRequest,
					Error:"invalid parameters",
				}
			}
		}
		case "getAggregations":{
			aggType := r.FormValue("aggType")
			resType := r.FormValue("resType")
			if len(aggType) == 0{
				return netutil.ResponseLayout{StatusCode:http.StatusBadRequest,Error:"invalid aggType when getting aggregations"}
			}
			if len(resType) == 0{
				return netutil.ResponseLayout{StatusCode:http.StatusBadRequest,Error:"invalid resType when getting aggregations"}
			}
			torrents,err := controller.GetAggregations(aggType,resType)
			if err != nil{
				return netutil.ResponseLayout{StatusCode:http.StatusInternalServerError,Error:err.Error(), IsInternalError:true}
			}
			return netutil.ResponseLayout{StatusCode:http.StatusOK,Response:torrents}
		}
		default:
			return netutil.ResponseLayout{
				StatusCode:http.StatusBadRequest,
				Error:"action '"+action+"' not recognized",
			}
		}
	}else{
		return netutil.ResponseLayout{
			StatusCode:http.StatusBadRequest,
			Error:"expecting POST request to central, got "+r.Method,
		}
	}
}