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
		log.Printf("[%d] Unable to get players for leg: %s", legID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

// GetStatisticsForLeg will return statistics for all players in the given leg
func GetStatisticsForLeg(w http.ResponseWriter, r *http.Request) {
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

	matchType := match.MatchType.ID
	if leg.LegType != nil {
		matchType = leg.LegType.ID
	}
	if matchType == models.SHOOTOUT {
		stats, err := data.GetShootoutStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get shootout statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.CRICKET {
		stats, err := data.GetCricketStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get cricket statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.DARTSATX {
		stats, err := data.GetDartsAtXStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get Darts At X statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.AROUNDTHECLOCK {
		stats, err := data.GetAroundTheClockStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get Around the Clock statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.AROUNDTHEWORLD {
		stats, err := data.GetAroundTheWorldStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get Around the World statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.SHANGHAI {
		stats, err := data.GetShanghaiStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get Shanghai statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.TICTACTOE {
		stats, err := data.GetTicTacToeStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get Tic Tac Toe statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.BERMUDATRIANGLE {
		stats, err := data.GetBermudaTriangleStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get Bermuda Triangle statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.FOURTWENTY {
		stats, err := data.Get420StatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get 420 statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.KILLBULL {
		stats, err := data.GetKillBullStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get Kill Bull statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.GOTCHA {
		stats, err := data.GetGotchaStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get Gotcha statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.JDCPRACTICE {
		stats, err := data.GetJDCPracticeStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get JDC Practice statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.KNOCKOUT {
		stats, err := data.GetKnockoutStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get Knockout statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.SCAM {
		stats, err := data.GetScamStatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get Scam statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if matchType == models.ONESEVENTY {
		stats, err := data.Get170StatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get 170 statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else {
		stats, err := data.GetX01StatisticsForLeg(legID)
		if err != nil {
			log.Println("Unable to get x01 statistics", err)
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
		log.Println("Unable to deserialize order body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.ChangePlayerOrder(legID, orderMap)
	if err != nil {
		log.Println("Unable to change player order", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err := data.GetLegPlayers(legID)
	if err != nil {
		log.Printf("[%d] Unable to get players for leg: %s", legID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

// StartWarmup will set the leg as warm up
func StartWarmup(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	legID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	venue := new(models.Venue)
	err = json.NewDecoder(r.Body).Decode(&venue)
	if err != nil {
		log.Println("Unable to deserialize venue body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.StartWarmup(legID, int(venue.ID.Int64))
	if err != nil {
		log.Println("Unable to start warmup", err)
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
