package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kcapp/api/data"
)

// GetDefaultOptions will return the default options
func GetDefaultOptions(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	offices, err := data.GetDefaultOptions()
	if err != nil {
		log.Println("Unable to get default options", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(offices)
}
