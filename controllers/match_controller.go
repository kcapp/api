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

// GetMatchesForGame will return a list of all matches for the given game ID
func GetMatchesForGame(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	gameID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matches, err := data.GetMatchesForGame(gameID)
	if err != nil {
		log.Println("Unable to get matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}

// GetMatch will return a match specified by the given id
func GetMatch(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	matchID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match, err := data.GetMatch(matchID)
	if err != nil {
		log.Println("Unable to get match", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(match)
}

// GetActiveMatches will return a list of all matches which are currently active
func GetActiveMatches(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	matches, err := data.GetActiveMatches()
	if err != nil {
		log.Println("Unable to get matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}

// GetMatchPlayers will return a match specified by the given id
func GetMatchPlayers(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	matchID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	players, err := data.GetMatchPlayers(matchID)
	if err != nil {
		log.Println("Unable to get players for match", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

// GetX01StatisticsForMatch will return X01 statistics for all players in the given match
func GetX01StatisticsForMatch(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	matchID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := data.GetX01StatisticsForMatch(matchID)
	if err != nil {
		log.Println("Unable to get statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// ChangePlayerOrder will modify the order of players for the given match
func ChangePlayerOrder(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	matchID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderMap := make(map[string]int)
	err = json.NewDecoder(r.Body).Decode(&orderMap)
	if err != nil {
		log.Println("Unable to deserialize body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.ChangePlayerOrder(matchID, orderMap)
	if err != nil {
		log.Println("Unable to change player order", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FinishMatch will finalize a match
func FinishMatch(w http.ResponseWriter, r *http.Request) {
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

	_, err = data.FinishMatch(visit)
	if err != nil {
		log.Println("Unable to finalize match", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
