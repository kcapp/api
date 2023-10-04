package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/data"
)

func GetBadges(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	badges, err := data.GetBadges()
	if err != nil {
		log.Println("Unable to get badges")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(badges)
}

func GetBadgesStatistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	badges, err := data.GetBadgesStatistics()
	if err != nil {
		log.Println("Unable to get badge statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(badges)
}

func GetBadgeStatistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	badges, err := data.GetBadgeStatistics(id)
	if err != nil {
		log.Println("Unable to get badge statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(badges)
}
