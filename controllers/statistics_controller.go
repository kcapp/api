package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/data"
)

// GetX01Statistics will return X01 statistics for a given period
func GetX01Statistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)

	statistics, err := data.GetX01Statistics(params["from"], params["to"], 301, 501)
	if err != nil {
		log.Println("Unable to get X01 statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(statistics)
}

// GetShootoutStatistics will return Shootout statistics for a given period
func GetShootoutStatistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)
	stats, err := data.GetShootoutStatistics(params["from"], params["to"])
	if err != nil {
		log.Println("Unable to get Shootout statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetCricketStatistics will return Cricket statistics for a given period
func GetCricketStatistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)
	stats, err := data.GetCricketStatistics(params["from"], params["to"])
	if err != nil {
		log.Println("Unable to get Cricket statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetDartsAtXStatistics will return Cricket statistics for a given period
func GetDartsAtXStatistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)
	stats, err := data.GetDartsAtXStatistics(params["from"], params["to"])
	if err != nil {
		log.Println("Unable to get Darts At X statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetAroundTheClockStatistics will return Around The Clock statistics for a given period
func GetAroundTheClockStatistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)
	stats, err := data.GetAroundTheClockStatistics(params["from"], params["to"])
	if err != nil {
		log.Println("Unable to get Around The Clock statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetAroundTheWorldStatistics will return Around The World statistics for a given period
func GetAroundTheWorldStatistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)
	stats, err := data.GetAroundTheWorldStatistics(params["from"], params["to"])
	if err != nil {
		log.Println("Unable to get Around The World statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetShanghaiStatistics will return Shanghai statistics for a given period
func GetShanghaiStatistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)
	stats, err := data.GetShanghaiStatistics(params["from"], params["to"])
	if err != nil {
		log.Println("Unable to get Shanghai statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetTicTacToeStatistics will return Tic Tac Toe statistics for a given period
func GetTicTacToeStatistics(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	SetHeaders(w)
	stats, err := data.GetTicTacToeStatistics(params["from"], params["to"])
	if err != nil {
		log.Println("Unable to get Tic Tac Toe statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
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
