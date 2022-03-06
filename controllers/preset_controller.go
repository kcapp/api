package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"

	"github.com/gorilla/mux"
)

// AddPreset will create a new preset
func AddPreset(w http.ResponseWriter, r *http.Request) {
	var preset models.MatchPreset
	err := json.NewDecoder(r.Body).Decode(&preset)
	if err != nil {
		log.Println("Unable to deserialize preset json", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.AddPreset(preset)
	if err != nil {
		log.Println("Unable to add preset", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetPresets will return a list of all presets
func GetPresets(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	players, err := data.GetPresets()
	if err != nil {
		log.Println("Unable to get presets", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

// GetPreset will return a preset with the given ID
func GetPreset(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	preset, err := data.GetPreset(id)
	if err != nil {
		log.Println("Unable to get preset", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(preset)
}

// UpdatePreset will update the given preset
func UpdatePreset(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var preset models.MatchPreset
	err = json.NewDecoder(r.Body).Decode(&preset)
	if err != nil {
		log.Println("Unable to deserialize preset json", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.UpdatePreset(id, preset)
	if err != nil {
		log.Println("Unable to update preset", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeletePreset will delete a preset with the given ID
func DeletePreset(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = data.DeletePreset(id)
	if err != nil {
		log.Println("Unable to delete preset", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
