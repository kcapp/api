package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/data"
)

// GetX01Statistics will return X01 statistics for a given period
func GetX01Statistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)

	statistics, err := data.GetX01Statistics(params["from"], params["to"], 301, 501)
	if err != nil {
		log.Println("Unable to get X01 statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(statistics)
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

// GetGlobalStatistics will return some global statistics for all legs played
func GetGlobalStatistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)

	global, err := data.GetGlobalStatistics()
	if err != nil {
		log.Println("Unable to get global statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(global)
}

// GetOfficeStatistics will return statistics for the given office
func GetOfficeStatistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)

	id, err := strconv.Atoi(params["office_id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statistics, err := data.GetOfficeStatistics(id, params["from"], params["to"])
	if err != nil {
		log.Println("Unable to get statistics for office", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(statistics)
}
