package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kcapp/api/models"

	"github.com/gorilla/mux"
)

// GetX01Statistics will return X01 statistics for a given period
func GetX01Statistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)
	stats, err := models.GetX01Statistics(params["from"], params["to"])
	if err != nil {
		log.Println("Unable to get statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}
