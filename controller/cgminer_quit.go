package controller

import (
	"time"
	"encoding/json"
	"net/http"
	"log"
	"github.com/blockassets/cgminer_client"
	)

// Implements Controller interface
type CGQuitCtrl struct{}

func (c CGQuitCtrl) build() *Controller {
	return &Controller{
		Methods: []string{http.MethodGet},
		Path:    "/cgminer/quit",
		Handler: c.makeHandler(),
	}
}

func (c CGQuitCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("CGMiner Quit Requested")

			bamStat := BAMStatus{"OK", nil}
			httpStat := http.StatusOK

			err := cgmQuit()
			if  err != nil {
				httpStat = http.StatusBadGateway
				bamStat = BAMStatus{"Error", err}
			}
			w.WriteHeader(httpStat)
			resp, _ := json.Marshal(bamStat)
			w.Write(resp)
		})
}

func cgmQuit() error {
	clnt := cgminer_client.New( "localhost", 4028, 5*time.Second)
	return clnt.Quit()

}