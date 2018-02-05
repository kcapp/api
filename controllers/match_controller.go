package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

	matches, err := models.GetMatchesForGame(gameID)
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

	match, err := models.GetMatch(matchID)
	if err != nil {
		log.Println("Unable to get match", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(match)
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

	playersMap, err := models.GetMatchPlayers(matchID)
	if err != nil {
		log.Println("Unable to get players for match", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	players, err := models.GetPlayersScore(matchID, playersMap)
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

	stats, err := models.GetX01StatisticsForMatch(matchID)
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

	err = models.ChangePlayerOrder(matchID, orderMap)
	if err != nil {
		log.Println("Unable to change player order", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
