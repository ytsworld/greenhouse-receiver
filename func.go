package receiver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	greenhouse "github.com/ytsworld/greenhouse-client/pkg"
)

var (
	defaultServerSideMessage = "Unexpected error occured on server side, please try again later"
)

// EntryPoint Will be provided as function to execute for Google Cloud functions
func EntryPoint(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" && r.URL.Path == "/api/v1/greenhouse" {

		rawData, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Warnf("Error reading request body. %s", err)
			handleError(http.StatusServiceUnavailable, defaultServerSideMessage, w)
		}
		var data = greenhouse.Data{}
		err = json.Unmarshal(rawData, &data)
		if err != nil {
			handleError(http.StatusBadRequest, "Json payload is not valid", w)
		}

		log.Infof("Got data: %+v", data)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("{\"success\": true}"))

	} else {
		log.Infof("Page not found for %s - %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}

}

func handleError(statusCode int, message string, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
