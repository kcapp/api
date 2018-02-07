package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
)

// GetOwes will return a list of all games
func GetOwes(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	owes, err := data.GetOwes()
	if err != nil {
		log.Println("Unable to get owes", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(owes)
}

// RegisterPayback will register a payback between the given players
func RegisterPayback(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	var owe models.Owe
	err := json.NewDecoder(r.Body).Decode(&owe)
	if err != nil {
		log.Println("Unable to deserialize body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.RegisterPayback(owe)
	if err != nil {
		log.Println("Unable to register payback", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
