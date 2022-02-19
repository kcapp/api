package controllers_v2

import (
	"encoding/json"
	"log"
	"net/http"

	data_v2 "github.com/kcapp/api/data/v2"
)

// GetPlayers will return a map containing all players
func GetPlayers(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	players, err := data_v2.GetPlayers()
	if err != nil {
		log.Println("Unable to get players", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}
