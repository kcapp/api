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
	router.HandleFunc("/health", controllers.Healthcheck).Methods("HEAD")

	router.HandleFunc("/match", controllers.NewMatch).Methods("POST")
	router.HandleFunc("/match/active", controllers.GetActiveMatches).Methods("GET")
	router.HandleFunc("/match/types", controllers.GetMatchesTypes).Methods("GET")
	router.HandleFunc("/match/modes", controllers.GetMatchesModes).Methods("GET")
	router.HandleFunc("/match/outshot", controllers.GetOutshotTypes).Methods("GET")
	router.HandleFunc("/match/{id}/continue", controllers.ContinueMatch).Methods("PUT")
	router.HandleFunc("/match", controllers.GetMatches).Methods("GET")
	router.HandleFunc("/match/{id}", controllers.GetMatch).Methods("GET")
	router.HandleFunc("/match/{id}/metadata", controllers.GetMatchMetadata).Methods("GET")
	router.HandleFunc("/match/{id}/rematch", controllers.ReMatch).Methods("POST")
	router.HandleFunc("/match/{id}/statistics", controllers.GetStatisticsForMatch).Methods("GET")
	router.HandleFunc("/match/{id}/legs", controllers.GetLegsForMatch).Methods("GET")
	router.HandleFunc("/match/{start}/{limit}", controllers.GetMatchesLimit).Methods("GET")

	router.HandleFunc("/leg/active", controllers.GetActiveLegs).Methods("GET")
	router.HandleFunc("/leg/{id}", controllers.GetLeg).Methods("GET")
	router.HandleFunc("/leg/{id}", controllers.DeleteLeg).Methods("DELETE")
	router.HandleFunc("/leg/{id}/statistics", controllers.GetStatisticsForLeg).Methods("GET")
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
	router.HandleFunc("/player/{id}/statistics", controllers.GetPlayerStatistics).Methods("GET")

	router.HandleFunc("/player/{id}/statistics/previous", controllers.GetPlayerX01PreviousStatistics).Methods("GET")
	router.HandleFunc("/player/{id}/progression", controllers.GetPlayerProgression).Methods("GET")
	router.HandleFunc("/player/{id}/checkouts", controllers.GetPlayerCheckouts).Methods("GET")
	router.HandleFunc("/player/{id}/tournament", controllers.GetPlayerTournamentStandings).Methods("GET")
	router.HandleFunc("/player/{id}/elo/{start}/{limit}", controllers.GetPlayerEloChangelog).Methods("GET")
	router.HandleFunc("/player/{player_1}/vs/{player_2}", controllers.GetPlayerHeadToHead).Methods("GET")
	router.HandleFunc("/player/{player_1}/vs/{player_2}/simulate", controllers.SimulateMatch).Methods("PUT")
	router.HandleFunc("/player", controllers.AddPlayer).Methods("POST")
	router.HandleFunc("/player/{id}/calendar", controllers.GetPlayerCalendar).Methods("GET")
	router.HandleFunc("/player/{id}/random/{starting_score}", controllers.GetRandomLegForPlayer).Methods("GET")
	router.HandleFunc("/player/{id}/statistics/{match_type}", controllers.GetPlayerMatchTypeStatistics).Methods("GET")
	router.HandleFunc("/player/{id}/statistics/{match_type}/history/{limit}", controllers.GetPlayerMatchTypeHistory).Methods("GET")

	router.HandleFunc("/statistics/global", controllers.GetGlobalStatistics).Methods("GET")
	router.HandleFunc("/statistics/global/fnc", controllers.GetGlobalStatisticsFnc).Methods("GET")
	router.HandleFunc("/statistics/office/{from}/{to}", controllers.GetOfficeStatistics).Methods("GET")
	router.HandleFunc("/statistics/office/{office_id}/{from}/{to}", controllers.GetOfficeStatistics).Methods("GET")
	router.HandleFunc("/statistics/{dart}/hits", controllers.GetDartStatistics).Methods("GET")
	router.HandleFunc("/statistics/{match_type}/{from}/{to}", controllers.GetStatistics).Methods("GET")

	router.HandleFunc("/owe", controllers.GetOwes).Methods("GET")
	router.HandleFunc("/owe/payback", controllers.RegisterPayback).Methods("PUT")

	router.HandleFunc("/owetype", controllers.GetOweTypes).Methods("GET")

	router.HandleFunc("/office", controllers.AddOffice).Methods("POST")
	router.HandleFunc("/office/{id}", controllers.UpdateOffice).Methods("PUT")
	router.HandleFunc("/office", controllers.GetOffices).Methods("GET")

	router.HandleFunc("/venue", controllers.AddVenue).Methods("POST")
	router.HandleFunc("/venue/{id}", controllers.UpdateVenue).Methods("PUT")
	router.HandleFunc("/venue", controllers.GetVenues).Methods("GET")
	router.HandleFunc("/venue/{id}", controllers.GetVenue).Methods("GET")
	router.HandleFunc("/venue/{id}/config", controllers.GetVenueConfiguration).Methods("GET")
	router.HandleFunc("/venue/{id}/spectate", controllers.SpectateVenue).Methods("GET")

	router.HandleFunc("/tournament", controllers.NewTournament).Methods("POST")
	router.HandleFunc("/tournament", controllers.GetTournaments).Methods("GET")
	router.HandleFunc("/tournament/current", controllers.GetCurrentTournament).Methods("GET")
	router.HandleFunc("/tournament/current/{office_id}", controllers.GetCurrentTournamentForOffice).Methods("GET")
	router.HandleFunc("/tournament/groups", controllers.GetTournamentGroups).Methods("GET")
	router.HandleFunc("/tournament/standings", controllers.GetTournamentStandings).Methods("GET")
	router.HandleFunc("/tournament/{id}", controllers.GetTournament).Methods("GET")
	router.HandleFunc("/tournament/{id}/player/{player_id}", controllers.GetTournamentPlayerMatches).Methods("GET")
	router.HandleFunc("/tournament/{id}/matches", controllers.GetTournamentMatches).Methods("GET")
	router.HandleFunc("/tournament/{id}/metadata", controllers.GetMatchMetadataForTournament).Methods("GET")
	router.HandleFunc("/tournament/{id}/overview", controllers.GetTournamentOverview).Methods("GET")
	router.HandleFunc("/tournament/{id}/statistics", controllers.GetTournamentStatistics).Methods("GET")

	log.Printf("Listening on port %d", config.APIConfig.Port)
	log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.APIConfig.Port), router))
}
