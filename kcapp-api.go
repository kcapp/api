package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/controllers"
	"github.com/kcapp/api/models"
)

// our main function
func main() {
	var configFileParam string

	if len(os.Args) > 1 {
		configFileParam = os.Args[1]
	}

	config, err := models.GetConfig(configFileParam)
	if err != nil {
		panic(err)
	}
	models.InitDB(config.GetMysqlConnectionString())

	router := mux.NewRouter()
	router.HandleFunc("/match", controllers.NewMatch).Methods("POST")
	router.HandleFunc("/match/active", controllers.GetActiveMatches).Methods("GET")
	router.HandleFunc("/match/types", controllers.GetMatchesTypes).Methods("GET")
	router.HandleFunc("/match/modes", controllers.GetMatchesModes).Methods("GET")
	router.HandleFunc("/match/{id}/continue", controllers.ContinueMatch).Methods("PUT")
	router.HandleFunc("/match", controllers.GetMatches).Methods("GET")
	router.HandleFunc("/match/{id}", controllers.GetMatch).Methods("GET")
	router.HandleFunc("/match/{id}/rematch", controllers.ReMatch).Methods("POST")
	router.HandleFunc("/match/{id}/statistics", controllers.GetX01StatisticsForMatch).Methods("GET")
	router.HandleFunc("/match/{id}/legs", controllers.GetLegsForMatch).Methods("GET")
	router.HandleFunc("/match/{start}/{limit}", controllers.GetMatchesLimit).Methods("GET")

	router.HandleFunc("/leg/active", controllers.GetActiveLegs).Methods("GET")
	router.HandleFunc("/leg/{id}", controllers.GetLeg).Methods("GET")
	router.HandleFunc("/leg/{id}", controllers.DeleteLeg).Methods("DELETE")
	router.HandleFunc("/leg/{id}/statistics", controllers.GetX01StatisticsForLeg).Methods("GET")
	router.HandleFunc("/leg/{id}/players", controllers.GetLegPlayers).Methods("GET")
	router.HandleFunc("/leg/{id}/order", controllers.ChangePlayerOrder).Methods("PUT")
	router.HandleFunc("/leg/{id}/finish", controllers.FinishLeg).Methods("PUT")
	router.HandleFunc("/leg/{id}/undo", controllers.UndoFinishLeg).Methods("PUT")

	router.HandleFunc("/visit", controllers.AddVisit).Methods("POST")
	router.HandleFunc("/visit/{id}/modify", controllers.ModifyVisit).Methods("PUT")
	router.HandleFunc("/visit/{id}", controllers.DeleteVisit).Methods("DELETE")
	router.HandleFunc("/visit/{leg_id}/last", controllers.DeleteLastVisit).Methods("DELETE")

	router.HandleFunc("/player", controllers.GetPlayers).Methods("GET")
	router.HandleFunc("/player/active", controllers.GetActivePlayers).Methods("GET")
	router.HandleFunc("/player/compare", controllers.GetPlayersX01Statistics).Methods("GET")
	router.HandleFunc("/player/{id}", controllers.GetPlayer).Methods("GET")
	router.HandleFunc("/player/{id}", controllers.UpdatePlayer).Methods("PUT")
	router.HandleFunc("/player/{id}/statistics", controllers.GetPlayerX01Statistics).Methods("GET")
	router.HandleFunc("/player/{id}/statistics/previous", controllers.GetPlayerX01PreviousStatistics).Methods("GET")
	router.HandleFunc("/player/{id}/progression", controllers.GetPlayerProgression).Methods("GET")
	router.HandleFunc("/player/{id}/checkouts", controllers.GetPlayerCheckouts).Methods("GET")
	router.HandleFunc("/player/{id}/tournament", controllers.GetPlayerTournamentStandings).Methods("GET")
	router.HandleFunc("/player/{player_1}/vs/{player_2}", controllers.GetPlayerHeadToHead).Methods("GET")
	router.HandleFunc("/player/{player_1}/vs/{player_2}/simulate", controllers.SimulateMatch).Methods("PUT")
	router.HandleFunc("/player", controllers.AddPlayer).Methods("POST")

	router.HandleFunc("/statistics/x01/{from}/{to}", controllers.GetX01Statistics).Methods("GET")
	router.HandleFunc("/statistics/shootout/{from}/{to}", controllers.GetShootoutStatistics).Methods("GET")

	router.HandleFunc("/owe", controllers.GetOwes).Methods("GET")
	router.HandleFunc("/owe/payback", controllers.RegisterPayback).Methods("PUT")

	router.HandleFunc("/owetype", controllers.GetOweTypes).Methods("GET")

	router.HandleFunc("/venue", controllers.GetVenues).Methods("GET")
	router.HandleFunc("/venue/{id}", controllers.GetVenue).Methods("GET")
	router.HandleFunc("/venue/{id}/spectate", controllers.SpectateVenue).Methods("GET")

	router.HandleFunc("/tournament", controllers.GetTournaments).Methods("GET")
	router.HandleFunc("/tournament/groups", controllers.GetTournamentGroups).Methods("GET")
	router.HandleFunc("/tournament/standings", controllers.GetTournamentStandings).Methods("GET")
	router.HandleFunc("/tournament/{id}", controllers.GetTournament).Methods("GET")
	router.HandleFunc("/tournament/{id}/matches", controllers.GetTournamentMatches).Methods("GET")
	router.HandleFunc("/tournament/{id}/overview", controllers.GetTournamentOverview).Methods("GET")
	router.HandleFunc("/tournament/{id}/statistics", controllers.GetTournamentStatistics).Methods("GET")

	log.Printf("Listening on port %d", config.APIConfig.Port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.APIConfig.Port), router))
}
