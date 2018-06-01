package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/data"
)

// GetTournaments will return all tournaments
func GetTournaments(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	tournaments, err := data.GetTournaments()
	if err != nil {
		log.Println("Unable to get tournaments", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tournaments)
}

// GetTournamentGroups will return all tournaments
func GetTournamentGroups(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	groups, err := data.GetTournamentGroups()
	if err != nil {
		log.Println("Unable to get tournament groups", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(groups)
}

// GetTournament will return the given tournament
func GetTournament(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tournament, err := data.GetTournament(id)
	if err != nil {
		log.Println("Unable to get tournament", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tournament)
}

// GetTournamentMatches will return all matches for the given tournament
func GetTournamentMatches(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	matches, err := data.GetTournamentMatches(id)
	if err != nil {
		log.Println("Unable to get tournament matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}

// GetTournamentOverview will return statistics for the given tournament
func GetTournamentOverview(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	stats, err := data.GetTournamentOverview(id)
	if err != nil {
		log.Println("Unable to get tournament overview", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetTournamentOverview will return statistics for the given tournament
func GetTournamentStatistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	stats, err := data.GetTournamentStatistics(id)
	if err != nil {
		log.Println("Unable to get tournament statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}
