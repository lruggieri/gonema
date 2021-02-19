package injector

import (
	"errors"
	"github.com/lruggieri/gonema/pkg/util"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	connectionsPoolElements = 10
)

type injector struct {
	mainHost          string
	parametersToShoot []string
	tps               int

	started           bool
	injectionFinished int
}

func (i *injector) Run() {
	if i.started {
		util.Logger.Error("injector is already running")
	}
	i.started = true

	util.Logger.Info("Injector started for mainHost ", i.mainHost)

	for _, parameterToShoot := range i.parametersToShoot {
		_ = i.runConfiguration(i.mainHost, parameterToShoot)
	}
}
func (i *injector) runConfiguration(mainHost, parameters string) error {
	urlToCall := mainHost
	if len(parameters) > 0 {
		runeToCall := []rune(urlToCall)
		if runeToCall[len(runeToCall)-1] != rune('?') {
			urlToCall += "?"
		}
		urlToCall += parameters
	}

	//first, let's create a pool of connections to the server
	clientsPool := make([]http.Client, 0, connectionsPoolElements)
	for i := 0; i < connectionsPoolElements; i++ {
		clientsPool = append(clientsPool, http.Client{
			Transport: &http.Transport{
				DialContext:         (&net.Dialer{Timeout: 1 * time.Second}).DialContext,
				TLSHandshakeTimeout: 3 * time.Second,
			},
		})
	}
	util.Logger.Debug("connection pool created with ", connectionsPoolElements, " elements")

	util.Logger.Debug("performing first request...")
	request, _ := http.NewRequest(http.MethodGet, urlToCall, nil)
	//the first request (to fill the cache) should have a higher timeout
	firstClient := http.Client{
		Transport: &http.Transport{
			DialContext:         (&net.Dialer{Timeout: 20 * time.Second}).DialContext,
			TLSHandshakeTimeout: 3 * time.Second,
		},
	}
	firstResult := makeRequest(&firstClient, request)
	if firstResult.err != nil {
		return errors.New("problem during first request for '" + parameters + "': " + firstResult.err.Error())
	}

	util.Logger.Info("First request performed. Starting to shoot at " + strconv.Itoa(i.tps) + " tps")

	rate := time.Second / time.Duration(i.tps)
	throttle := time.Tick(rate)
	for {
		randomClient := clientsPool[util.GetRandomPositiveInt(len(clientsPool))]
		<-throttle
		go func(iClient http.Client) {
			res := makeRequest(&randomClient, request)
			util.Logger.Debug(res.statusCode)
		}(randomClient)
	}
}

func (i *injector) logStats() {}

func makeRequest(iClient *http.Client, iReq *http.Request) (oResult singleInjectionResult) {
	util.Logger.Debug("shooting request for ", iReq.URL)
	resp, err := iClient.Do(iReq)
	if err != nil {
		return singleInjectionResult{
			err: err,
		}
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body) //be sure to read all response to avoid connection to drop
	defer resp.Body.Close()

	if err != nil {
		return singleInjectionResult{
			err: err,
		}
	} else {
		return singleInjectionResult{
			timedOut:   false,
			statusCode: resp.StatusCode,
		}
	}
}

type singleInjectionResult struct {
	timedOut   bool
	statusCode int
	err        error
}

func NewInjector(iMainHost string, iParametersToShoot []string, iTps int) *injector {
	return &injector{
		mainHost:          iMainHost,
		parametersToShoot: iParametersToShoot,
		tps:               iTps,
	}
}
