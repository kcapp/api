package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/data"
)

// GetVenues will return all venues
func GetVenues(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	venues, err := data.GetVenues()
	if err != nil {
		log.Println("Unable to get venues", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(venues)
}

// GetVenue will return the given venue
func GetVenue(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	venue, err := data.GetVenue(id)
	if err != nil {
		log.Println("Unable to get venue", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(venue)
}

// GetVenueConfiguration will return the configuration for the given venue
func GetVenueConfiguration(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	config, err := data.GetVenueConfiguration(id)
	if err != nil {
		log.Println("Unable to get venue configuration", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(config)
}

// SpectateVenue will spectate the current match active at a given venue
func SpectateVenue(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	matches, err := data.SpectateVenue(id)
	if err != nil {
		log.Println("Unable to spectate venue", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}
