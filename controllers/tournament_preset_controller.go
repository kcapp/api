package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/data"
)

/*
// AddTournamentPreset will create a new tournament preset
func AddTournamentPreset(w http.ResponseWriter, r *http.Request) {
	var preset models.TournamentPreset
	err := json.NewDecoder(r.Body).Decode(&preset)
	if err != nil {
		log.Println("Unable to deserialize preset json", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.AddTournamentPreset(preset)
	if err != nil {
		log.Println("Unable to add preset", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}*/

// GetTournamentPresets will return a list of all presets
func GetTournamentPresets(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	players, err := data.GetTournamentPresets()
	if err != nil {
		log.Println("Unable to get presets", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

// GetTournamentPreset will return a preset with the given ID
func GetTournamentPreset(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	preset, err := data.GetTournamentPreset(id)
	if err != nil {
		log.Println("Unable to get preset", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(preset)
}

/*
// UpdateTournamentPreset will update the given preset
func UpdateTournamentPreset(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var preset models.TournamentPreset
	err = json.NewDecoder(r.Body).Decode(&preset)
	if err != nil {
		log.Println("Unable to deserialize preset json", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.UpdateTournamentPreset(id, preset)
	if err != nil {
		log.Println("Unable to update preset", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
*/
