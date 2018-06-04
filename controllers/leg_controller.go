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

// GetLegsForMatch will return a list of all legs for the given match ID
func GetLegsForMatch(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	matchID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	legs, err := data.GetLegsForMatch(matchID)
	if err != nil {
		log.Println("Unable to get legs", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(legs)
}

// GetLeg will return a leg specified by the given id
func GetLeg(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	legID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	leg, err := data.GetLeg(legID)
	if err != nil {
		log.Println("Unable to get leg", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(leg)
}

// GetActiveLegs will return a list of all legs which are currently active
func GetActiveLegs(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	legs, err := data.GetActiveLegs()
	if err != nil {
		log.Println("Unable to get legs", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(legs)
}

// GetLegPlayers will return a leg specified by the given id
func GetLegPlayers(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	legID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	players, err := data.GetLegPlayers(legID)
	if err != nil {
		log.Println("Unable to get players for leg", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

// GetX01StatisticsForLeg will return X01 statistics for all players in the given leg
func GetX01StatisticsForLeg(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	legID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	leg, err := data.GetLeg(legID)
	if err != nil {
		log.Println("Unable to get leg")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	match, err := data.GetMatch(leg.MatchID)
	if err != nil {
		log.Println("Unable to get Match")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if match.MatchType.ID == models.SHOOTOUT {
		stats, err := data.GetShootoutStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else {
		stats, err := data.GetX01StatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	}
}

// ChangePlayerOrder will modify the order of players for the given leg
func ChangePlayerOrder(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	legID, err := strconv.Atoi(params["id"])
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

	err = data.ChangePlayerOrder(legID, orderMap)
	if err != nil {
		log.Println("Unable to change player order", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FinishLeg will finalize a leg
func FinishLeg(w http.ResponseWriter, r *http.Request) {
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

	err = data.FinishLeg(visit)
	if err != nil {
		log.Println("Unable to finalize leg", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteLeg will delete a leg
func DeleteLeg(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	legID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.DeleteLeg(legID)
	if err != nil {
		log.Println("Unable to delete leg", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UndoFinishLeg will undo a finalized leg
func UndoFinishLeg(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	legID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.UndoLegFinish(legID)
	if err != nil {
		log.Println("Unable to undo leg finish", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
