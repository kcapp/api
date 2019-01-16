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

// AddVisit will add the visit to the database
func AddVisit(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	var visit models.Visit
	err := json.NewDecoder(r.Body).Decode(&visit)
	if err != nil {
		log.Println("Unable to deserialize body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = visit.ValidateInput()
	if err != nil {
		log.Println("Invalid visit", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	insertedVisit, err := data.AddVisit(visit)
	if err != nil {
		log.Println("Unable to add visit", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(insertedVisit)
}

// ModifyVisit will modify the scores of the given visit
func ModifyVisit(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	var visit models.Visit
	err := json.NewDecoder(r.Body).Decode(&visit)
	if err != nil {
		log.Println("Unable to deserialize body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.ModifyVisit(visit)
	if err != nil {
		log.Println("Unable to modify visit", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteVisit will delete the given visit
func DeleteVisit(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = data.DeleteVisit(id)
	if err != nil {
		log.Println("Unable to delete visit: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteLastVisit will delete the last visit for a given leg
func DeleteLastVisit(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	legID, err := strconv.Atoi(params["leg_id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = data.DeleteLastVisit(legID)
	if err != nil {
		log.Println("Unable to delete visit: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
