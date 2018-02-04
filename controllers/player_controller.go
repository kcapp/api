package controllers

import (
	"encoding/json"
	"kcapp-api/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetPlayers will return a map containing all players
func GetPlayers(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	players, err := models.GetPlayers()
	if err != nil {
		log.Println("Unable to get players", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

// GetPlayer will return a player with the given ID
func GetPlayer(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(id)
	// TODO
}

// GetPlayerStatistics will return statistics for the given player
func GetPlayerStatistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := models.GetPlayerStatistics(id)
	if err != nil {
		log.Println("Unable to get player statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetPlayersStatistics will return statistics for the given player
func GetPlayersStatistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)

	params := r.URL.Query()["id"]
	if params == nil {
		log.Println("No players specified to compare")
		http.Error(w, "No players specified to compare", http.StatusBadRequest)
		return
	}
	ids, err := sliceAtoi(params)
	if err != nil {
		log.Println("Unable to convert params to int")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := models.GetPlayersStatistics(ids)
	if err != nil {
		log.Println("Unable to get players statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// AddPlayer will create a new player
func AddPlayer(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	_ = json.NewDecoder(r.Body).Decode(&player)
	err := models.AddPlayer(player)
	if err != nil {
		log.Println("Unable to add player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sliceAtoi(sa []string) ([]int, error) {
	si := make([]int, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.Atoi(a)
		if err != nil {
			return si, err
		}
		si = append(si, i)
	}
	return si, nil
}
