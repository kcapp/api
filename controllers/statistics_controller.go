package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/data"
)

// GetX01Statistics will return X01 statistics for a given period
func GetX01Statistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)
	stats, err := data.GetX01Statistics(params["from"], params["to"])
	if err != nil {
		log.Println("Unable to get X01 statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetShootoutStatistics will return Shootout statistics for a given period
func GetShootoutStatistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)
	stats, err := data.GetShootoutStatistics(params["from"], params["to"])
	if err != nil {
		log.Println("Unable to get Shootout statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}
