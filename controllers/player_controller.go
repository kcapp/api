package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/kcapp/api/util"
	"log"
	"net/http"
	"strconv"
)

// GetPlayers will return a map containing all players
func GetPlayers(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	players, err := data.GetPlayers()
	if err != nil {
		log.Println("Unable to get players", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

// GetActivePlayers will return a map containing all active players
func GetActivePlayers(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	players, err := data.GetActivePlayers()
	if err != nil {
		log.Println("Unable to get active players", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

// GetPlayer will return a player with the given ID
func GetPlayer(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	player, err := data.GetPlayer(id)
	if err != nil {
		log.Println("Unable to get player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(player)
}

// GetPlayersX01Statistics will return statistics for the given players
func GetPlayersX01Statistics(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	id1, err := strconv.Atoi(params["id1"])
	if err != nil {
		log.Println("Invalid player1 id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id2, err := strconv.Atoi(params["id2"])
	if err != nil {
		log.Println("Invalid player 2 id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ids []int
	ids = append(ids, id1, id2)

	if err != nil {
		log.Println("Unable to convert params to int")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := data.GetPlayersX01Statistics(ids)
	if err != nil {
		log.Println("Unable to get players statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetPlayerStatistics will return statistics for the given player
func GetPlayerStatistics(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	statistics := new(models.PlayerStatistics)

	x01, err := data.GetPlayerX01Statistics(id)
	if err != nil {
		log.Println("Unable to get player x01 statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	statistics.X01 = x01
	json.NewEncoder(w).Encode(statistics)
}

// GetPlayerX01PreviousStatistics will return statistics for the given player
func GetPlayerX01PreviousStatistics(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := data.GetPlayerX01PreviousStatistics(id)
	if err != nil {
		log.Println("Unable to get player statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

// GetPlayerProgression will return statistics for the given player
func GetPlayerProgression(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := data.GetPlayerProgression(id)
	if err != nil {
		log.Println("Unable to get player progression", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}
