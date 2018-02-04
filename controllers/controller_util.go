package controllers

import "net/http"

// SetHeaders will set the default headers used by all requests
func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
