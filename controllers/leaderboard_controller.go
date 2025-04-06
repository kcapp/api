package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kcapp/api/data"
)

// GetMatchTypeLeaderboard will return leaderboard for each match type
func GetMatchTypeLeaderboard(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	leaderboard, err := data.GetMatchTypeLeaderboard()
	if err != nil {
		log.Println("Unable to get match type leaderboards", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(leaderboard)
}
