package controllers

import (
	"net/http"
)

// Healthcheck will return OK
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("☄ HTTP status code returned!"))
}
