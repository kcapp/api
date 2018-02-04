package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/models"
)

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

	err = models.ModifyVisit(visit)
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
	err = models.DeleteVisit(id)
	if err != nil {
		log.Println("Unable to delete visit: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
