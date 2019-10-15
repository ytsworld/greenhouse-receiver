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

// EntryPoint is provided as function to execute when a new authorized request is received
func EntryPoint(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" && r.URL.Path == "/api/v1/greenhouse" {

		rawData, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Warnf("Error reading request body. %s", err)
			handleError(http.StatusServiceUnavailable, defaultServerSideMessage, w)
			return
		}
		var data = greenhouse.Data{}
		err = json.Unmarshal(rawData, &data)
		if err != nil {
			handleError(http.StatusBadRequest, "Json payload is not valid", w)
			return
		}

		log.Infof("Got data: %+v", data)

		if !data.Success {
			log.Infof("Greenhouse client had error: %s", data.Message)
			handleSuccess(w)
			return
		}

		err = persistAll(data)
		if err != nil {
			log.Warnf("error while persisting data: %s", err)
			handleError(http.StatusInternalServerError, "Error while persisting data", w)
			return
		}

		handleSuccess(w)
		return

	}

	// Unkown request path / method combination
	log.Infof("Page not found for %s - %s", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))

}

func handleSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("{\"success\": true}"))
}

func handleError(statusCode int, message string, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
