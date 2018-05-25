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

// NewMatch will start a new match
func NewMatch(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	var matchInput models.Match
	err := json.NewDecoder(r.Body).Decode(&matchInput)
	if err != nil {
		log.Println("Unable to deserialize body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match, err := data.NewMatch(matchInput)
	if err != nil {
		log.Println("Unable to start new match", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(match)
}

// ContinueMatch will either return the current leg id or create a new leg
func ContinueMatch(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	match, err := data.ContinueMatch(id)
	if err != nil {
		log.Println("Unable to get match: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(match)
}

// GetMatches will return a list of all matches
func GetMatches(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	matches, err := data.GetMatches()
	if err != nil {
		log.Println("Unable to get matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}

// GetActiveMatches will return a list of active matches
func GetActiveMatches(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	matches, err := data.GetActiveMatches()
	if err != nil {
		log.Println("Unable to get active matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}

// GetMatchesLimit will return N matches from the given starting point
func GetMatchesLimit(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	start, err := strconv.Atoi(params["start"])
	if err != nil {
		log.Println("Invalid start parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(params["limit"])
	if err != nil {
		log.Println("Invalid limit parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	matches, err := data.GetMatchesLimit(start, limit)
	if err != nil {
		log.Println("Unable to get matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}

// GetMatch will reurn a the match with the given ID
func GetMatch(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	match, err := data.GetMatch(id)
	if err != nil {
		log.Println("Unable to get match: ", err)
		http.Error(w, "Unable to get match", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(match)
}

// GetX01StatisticsForMatch will return X01 statistics for all players in the given leg
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

// GetMatchesModes will return all match modes
func GetMatchesModes(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	modes, err := data.GetMatchModes()
	if err != nil {
		log.Println("Unable to get match modes", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(modes)
}

// GetMatchesTypes will return all match types
func GetMatchesTypes(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	types, err := data.GetMatchTypes()
	if err != nil {
		log.Println("Unable to get match types", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(types)
}
