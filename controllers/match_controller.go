package controllers

import (
	"encoding/json"
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/util"
	"log"
	"net/http"
)

// GetMatchTypes will return all match types
func GetMatchTypes(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	types, err := data.GetMatchTypes()
	if err != nil {
		log.Println("Unable to get match types", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(types)
}

// GetMatchModes will return all match modes
func GetMatchModes(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	modes, err := data.GetMatchModes()
	if err != nil {
		log.Println("Unable to get match modes", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(modes)
}
