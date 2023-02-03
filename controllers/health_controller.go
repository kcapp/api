package controllers

import (
	"github.com/kcapp/api/util"
	"net/http"
)

// Healthcheck will return OK
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("☄ HTTP status code returned!"))
}
