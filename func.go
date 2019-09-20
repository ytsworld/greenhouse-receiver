package hello

import (
	"fmt"
	"net/http"
)

// EntryPoint Will be provided as function to execute for Google Cloud functions
func EntryPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
