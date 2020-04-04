package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
)

// AddOffice will create a new office
func AddOffice(w http.ResponseWriter, r *http.Request) {
	var office models.Office
	err := json.NewDecoder(r.Body).Decode(&office)
	if err != nil {
		log.Println("Unable to deserialize office json", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.AddOffice(office)
	if err != nil {
		log.Println("Unable to add office", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateOffice will update the given office
func UpdateOffice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var office models.Office
	err = json.NewDecoder(r.Body).Decode(&office)
	if err != nil {
		log.Println("Unable to deserialize office json", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.UpdateOffice(id, office)
	if err != nil {
		log.Println("Unable to update office", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetOffices will return all offices
func GetOffices(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	offices, err := data.GetOffices()
	if err != nil {
		log.Println("Unable to get offices", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(offices)
}

// GetOffice will return a office with the given ID
func GetOffice(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	office, err := data.GetOffice(id)
	if err != nil {
		log.Println("Unable to get office", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(office)
}
