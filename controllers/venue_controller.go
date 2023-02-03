package controllers

import (
	"encoding/json"
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/util"
	"log"
	"net/http"
)

// GetVenues will return all venues
func GetVenues(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	venues, err := data.GetVenues()
	if err != nil {
		log.Println("Unable to get venues", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(venues)
}
