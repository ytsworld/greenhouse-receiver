package receiver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	greenhouse "github.com/ytsworld/greenhouse-client/pkg"
)

var (
	valid = greenhouse.Data{
		UnixTimestampUTC:       (time.Now()).Unix(),
		Success:                true,
		Temperature:            20.1,
		Humidity:               76,
		SoilMoistureResistance: 845,
	}
)

func TestHealthCheckHandler(t *testing.T) {
	stringPayload, err := json.Marshal(&valid)
	if err != nil {
		t.Fatal(err)
	}
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/api/v1/greenhouse", bytes.NewReader(stringPayload))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(EntryPoint)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	expected := `{"success": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",rr.Body.String(), expected)
	}
}
