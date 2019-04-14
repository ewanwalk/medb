package utils

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Vars
// obtain all the query params applicable to the request
// this merges mux & standard lib
func Vars(r *http.Request) map[string]string {

	vars := make(map[string]string)

	for key, values := range r.URL.Query() {
		vars[key] = values[0]
	}

	// mux takes priority and will overwrite any existing collisions
	for key, val := range mux.Vars(r) {
		vars[key] = val
	}

	return vars
}

func QueryParams(r *http.Request) map[string]string {
	vars := make(map[string]string)

	for key, values := range r.URL.Query() {
		vars[key] = values[0]
	}

	return vars
}
