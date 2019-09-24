package hello

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	greenhouse "github.com/ytsworld/greenhouse-client/pkg"
)

// EntryPoint Will be provided as function to execute for Google Cloud functions
func EntryPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")

	if r.Method == "POST" && r.URL.Path == "/api/v1/greenhouse" {

		rawData, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			handleError(err, w)
		}
		var data greenhouse.Data
		err = json.Unmarshal(rawData, &data)
		if err != nil {
			handleError(err, w)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("nope"))
	}

}

func handleError(err error, w http.ResponseWriter) {
	// TODO LOG
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("Unexpected error occured on server side, please try again later"))
}
