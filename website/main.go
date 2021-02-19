package main

import (
	"fmt"
	"github.com/NYTimes/gziphandler"
	uCache "github.com/lruggieri/gonema/pkg/util/cache"
	"github.com/lruggieri/gonema/website/controller"
	"github.com/lruggieri/utils/netutil"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

var (
	templatesDir   string
	staticAssetDir string

	cache *uCache.Cache
)

const (
	cacheRootAggregation = uCache.CacheElementRoot("aggregations")
)

//neuteredFileSystem is used to prevent directory listing of static assets
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
	if len(templatesDir) == 0 {
		templatesDir = "website/templates"
	}
	staticAssetDir = os.Getenv("STATIC_ASSET_DIR")
	if len(staticAssetDir) == 0 {
		staticAssetDir = "website/static"
	}

	mux := http.NewServeMux()
	fs := http.FileServer(neuteredFileSystem{http.Dir(staticAssetDir)})
	mux.Handle("/static/", http.StripPrefix("/static/", gziphandler.GzipHandler(fs)))
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/robots.txt", robotsHandler)
	mux.Handle("/central", gziphandler.GzipHandler(netutil.HandleWithError(centralControllerHandler)))
	mux.Handle("/", gziphandler.GzipHandler(http.HandlerFunc(mainPageHandler)))

	cache = uCache.NewCache()
	cache.SetNewRootElementDuration(cacheRootAggregation, 10*time.Minute)

	var tlsCertPath = os.Getenv("TLS_CERT_PATH")
	var tlsKeyPath = os.Getenv("TLS_KEY_PATH")

	if len(tlsCertPath) > 0 && len(tlsKeyPath) > 0 {
		err := http.ListenAndServeTLS(":443", tlsCertPath, tlsKeyPath, nil)
		if err != nil {
			panic(err)
		}
	} else {
		port := os.Getenv("PORT")
		if len(port) == 0 {
			port = "8080"
		}
		fmt.Println("web service starting on port " + port)
		err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
		if err != nil {
			panic(err)
		}
	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	mainPage := filepath.Join(templatesDir, "index.tmpl")
	tmpl := template.Must(template.ParseFiles(mainPage))

	err := tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Website down, we are very sorry for the inconvenience."))
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(staticAssetDir, "images", "favicon.ico"))
}
func robotsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("robots.txt"))
}
func centralControllerHandler(w http.ResponseWriter, r *http.Request) netutil.ResponseLayout {

	if r.Method == http.MethodPost {
		action := r.FormValue("action")

		switch action {
		case "getResourceInfo":
			{
				resourceName := r.FormValue("resourceName")
				resourceImdbID := r.FormValue("resourceImdbID")
				if len(resourceName) == 0 && len(resourceImdbID) == 0 {
					return netutil.ResponseLayout{StatusCode: http.StatusBadRequest, Error: "no resource name or ID was given"}
				}
				resources, err := controller.GetResourceInfoFromOmdb(resourceName, resourceImdbID)
				if err != nil {
					return netutil.ResponseLayout{StatusCode: http.StatusInternalServerError, Error: err.Error(), IsInternalError: true}
				}
				return netutil.ResponseLayout{StatusCode: http.StatusOK, Response: resources}
			}
		case "getTorrents":
			{
				keyword := r.FormValue("keyword")
				imdbID := r.FormValue("imdbID")
				tType := r.FormValue("type")
				if len(keyword) == 0 {
					return netutil.ResponseLayout{StatusCode: http.StatusBadRequest, Error: "invalid keyword when fetching torrents"}
				}
				torrents, err := controller.GetTorrents(keyword, imdbID, tType)
				if err != nil {
					return netutil.ResponseLayout{StatusCode: http.StatusInternalServerError, Error: err.Error(), IsInternalError: true}
				}
				return netutil.ResponseLayout{StatusCode: http.StatusOK, Response: torrents}
			}
		case "suggest":
			{
				resourceName := r.FormValue("resourceName")
				if len(resourceName) > 0 {
					resources, err := controller.GetResourceInfoFromOmdb(resourceName, "")
					if err != nil {
						return netutil.ResponseLayout{StatusCode: http.StatusInternalServerError, Error: err.Error(), IsInternalError: true}
					}
					return netutil.ResponseLayout{StatusCode: http.StatusOK, Response: resources}

				} else {
					return netutil.ResponseLayout{
						StatusCode: http.StatusBadRequest,
						Error:      "invalid parameters",
					}
				}
			}
		case "getAggregations":
			{

				aggregationKey := uCache.CacheElementKey(r.Form.Encode())

				cachedValue := cache.Fetch(cacheRootAggregation, aggregationKey)
				if cachedValue != nil {
					return netutil.ResponseLayout{StatusCode: http.StatusOK, Response: cachedValue}
				} else {
					aggType := r.FormValue("aggType")
					resType := r.FormValue("resType")
					if len(aggType) == 0 {
						return netutil.ResponseLayout{StatusCode: http.StatusBadRequest, Error: "invalid aggType when getting aggregations"}
					}
					if len(resType) == 0 {
						return netutil.ResponseLayout{StatusCode: http.StatusBadRequest, Error: "invalid resType when getting aggregations"}
					}
					torrents, err := controller.GetAggregations(aggType, resType)
					if err != nil {
						return netutil.ResponseLayout{StatusCode: http.StatusInternalServerError, Error: err.Error(), IsInternalError: true}
					}
					cache.Insert(cacheRootAggregation, aggregationKey, torrents)
					return netutil.ResponseLayout{StatusCode: http.StatusOK, Response: torrents}
				}
			}
		default:
			return netutil.ResponseLayout{
				StatusCode: http.StatusBadRequest,
				Error:      "action '" + action + "' not recognized",
			}
		}
	} else {
		return netutil.ResponseLayout{
			StatusCode: http.StatusBadRequest,
			Error:      "expecting POST request to central, got " + r.Method,
		}
	}
}
