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

// GetPlayers will return a map containing all players
func GetPlayers(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	players, err := data.GetPlayers()
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
	player, err := data.GetPlayer(id)
	if err != nil {
		log.Println("Unable to get player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(player)
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

	stats, err := data.GetPlayerStatistics(id)
	if err != nil {
		log.Println("Unable to get player statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	visits, err := data.GetPlayerVisitCount(id)
	if err != nil {
		log.Println("Unable to get visits for player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stats.Visits = visits
	for _, v := range visits {
		stats.TotalVisits += v.Count
	}

	json.NewEncoder(w).Encode(stats)
}

// GetPlayersStatistics will return statistics for the given players
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

	stats, err := data.GetPlayersStatistics(ids)
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
	err := data.AddPlayer(player)
	if err != nil {
		log.Println("Unable to add player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetPlayerProgression will return statistics for the given player
func GetPlayerProgression(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
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

// GetPlayerCheckouts will return all checkouts done by a player
func GetPlayerCheckouts(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	checkouts, err := data.GetPlayerCheckouts(id)
	if err != nil {
		log.Println("Unable to get player checkouts")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(checkouts)
}

// GetPlayerHeadToHead will return head to head statistics between the given players
func GetPlayerHeadToHead(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	player1, err := strconv.Atoi(params["player_1"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	player2, err := strconv.Atoi(params["player_2"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	head2head, err := data.GetPlayerHeadToHead(player1, player2)
	if err != nil {
		log.Println("Unable to get player head to head statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(head2head)
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
