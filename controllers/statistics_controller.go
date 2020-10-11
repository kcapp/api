package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
)

// GetStatistics will return statistics for the given match type
func GetStatistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	matchType, err := strconv.Atoi(params["match_type"])
	if err != nil {
		log.Println("Invalid match type parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch matchType {
	case models.X01:
		statistics, err := data.GetX01Statistics(params["from"], params["to"], 301, 501)
		if err != nil {
			log.Println("Unable to get X01 statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(statistics)
		return

	case models.SHOOTOUT:
		stats, err := data.GetShootoutStatistics(params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get Shootout statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.CRICKET:
		stats, err := data.GetCricketStatistics(params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get Cricket statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.DARTSATX:
		stats, err := data.GetDartsAtXStatistics(params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get Darts At X statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.AROUNDTHEWORLD:
		stats, err := data.GetAroundTheWorldStatistics(params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get Around The World statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.SHANGHAI:
		stats, err := data.GetShanghaiStatistics(params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get Shanghai statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.AROUNDTHECLOCK:
		stats, err := data.GetAroundTheClockStatistics(params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get Around The Clock statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.TICTACTOE:
		stats, err := data.GetTicTacToeStatistics(params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get Tic Tac Toe statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.BERMUDATRIANGLE:
		stats, err := data.GetBermudaTriangleStatistics(params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get Bermuda Triangle statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.FOURTWENTY:
		stats, err := data.Get420Statistics(params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get 420 Statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.KILLBULL:
		stats, err := data.GetKillBullStatistics(params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get Kill Bull Statistics", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	default:
		log.Println("Unknown match type parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// GetGlobalStatistics will return some global statistics for all matches
func GetGlobalStatistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)

	global, err := data.GetGlobalStatistics()
	if err != nil {
		log.Println("Unable to get global statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(global)
}

// GetGlobalStatisticsFnc will return global fish and chips counter
func GetGlobalStatisticsFnc(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)

	global, err := data.GetGlobalStatisticsFnc()
	if err != nil {
		log.Println("Unable to get global fish and chips statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(global)
}

// GetOfficeStatistics will return statistics for the given office
func GetOfficeStatistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)

	id, err := strconv.Atoi(params["office_id"])
	if err != nil {
		statistics, err := data.GetOfficeStatistics(params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get statistics for office", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(statistics)
	} else {
		statistics, err := data.GetOfficeStatisticsForOffice(id, params["from"], params["to"])
		if err != nil {
			log.Println("Unable to get statistics for office", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(statistics)
	}
}

// GetDartStatistics will return dart statistics for all players
func GetDartStatistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)

	dart, err := strconv.Atoi(params["dart"])
	if err != nil {
		log.Println("Unable to get dart", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	statistics, err := data.GetDartStatistics(dart)
	if err != nil {
		log.Println("Unable to get dart statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(statistics)

}
