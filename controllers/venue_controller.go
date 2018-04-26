package controllers

import (
	"encoding/json"
	"log"
	"net/http"

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
