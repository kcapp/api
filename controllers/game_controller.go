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

// NewGame will start a new game
func NewGame(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	var gameInput models.Game
	err := json.NewDecoder(r.Body).Decode(&gameInput)
	if err != nil {
		log.Println("Unable to deserialize body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game, err := data.NewGame(gameInput)
	if err != nil {
		log.Println("Unable to start new game", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(game)
}

// ContinueGame will either return the current match id or create a new match
func ContinueGame(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	game, err := data.ContinueGame(id)
	if err != nil {
		log.Println("Unable to get game: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(game)
}

// GetGames will return a list of all games
func GetGames(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	games, err := data.GetGames()
	if err != nil {
		log.Println("Unable to get games", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(games)
}

// GetGame will reurn a the game with the given ID
func GetGame(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	game, err := data.GetGame(id)
	if err != nil {
		log.Println("Unable to get game: ", err)
		http.Error(w, "Unable to get game", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(game)
}

// GetX01StatisticsForGame will return X01 statistics for all players in the given match
func GetX01StatisticsForGame(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	gameID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := data.GetX01StatisticsForGame(gameID)
	if err != nil {
		log.Println("Unable to get statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetGamesModes will return all game modes
func GetGamesModes(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	modes, err := data.GetGameModes()
	if err != nil {
		log.Println("Unable to get game modes", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(modes)
}

// GetGamesTypes will return all game types
func GetGamesTypes(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	types, err := data.GetGameTypes()
	if err != nil {
		log.Println("Unable to get game types", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(types)
}
