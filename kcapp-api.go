package main

import (
	"fmt"
	"log"
	"math"
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
	router.HandleFunc("/player/{id}/progression", controllers.GetPlayerProgression).Methods("GET")
	router.HandleFunc("/player/{id}/checkouts", controllers.GetPlayerCheckouts).Methods("GET")
	router.HandleFunc("/player/{id}/tournament", controllers.GetPlayerTournamentStandings).Methods("GET")
	router.HandleFunc("/player/{player_1}/vs/{player_2}", controllers.GetPlayerHeadToHead).Methods("GET")
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
	router.HandleFunc("/tournament/{id}", controllers.GetTournament).Methods("GET")
	router.HandleFunc("/tournament/{id}/matches", controllers.GetTournamentMatches).Methods("GET")
	router.HandleFunc("/tournament/{id}/overview", controllers.GetTournamentOverview).Methods("GET")
	router.HandleFunc("/tournament/{id}/statistics", controllers.GetTournamentStatistics).Methods("GET")
	router.HandleFunc("/tournament/groups", controllers.GetTournamentGroups).Methods("GET")

	p1Elo, p2Elo := CalculateElo(1774, 240, 1, 1545, 125, 0)
	log.Printf("P1 = %d, P2 = %d", p1Elo, p2Elo)

	log.Printf("Listening on port %d", config.APIConfig.Port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.APIConfig.Port), router))
}

// CalculateElo will calculate the Elo for each player based on the given information. Returned value is new Elo for player1 and player2 respectively
func CalculateElo(player1Elo int, player1Matches int, player1Score int, player2Elo int, player2Matches int, player2Score int) (int, int) {
	if player1Matches == 0 {
		player1Matches = 1
	}
	if player2Matches == 0 {
		player2Matches = 1
	}

	// P1 = Winner
	// P2 = Looser
	// PD = Points Difference
	// Multiplier = ln(abs(PD)+1) * (2.2 / ((P1(old)-P2(old)) * 0.001 + 2.2))
	// Elo Winner = P1(old) + 800/num_matches * (1 - 1/(1 + 10 ^ (P2(old) - P1(old) / 400) ) )
	// Elo Looser = P2(old) + 800/num_matches * (0 - 1/(1 + 10 ^ (P2(old) - P1(old) / 400) ) )

	if player1Score > player2Score {
		multiplier := math.Log(math.Abs(float64(player1Score-player2Score))+1) * (2.2 / ((float64(player1Elo-player2Elo))*0.001 + 2.2))
		player1Elo, player2Elo = calculateElo(player1Elo, player1Matches, player2Elo, player2Matches, multiplier, false)
	} else if player1Score < player2Score {
		multiplier := math.Log(math.Abs(float64(player1Score-player2Score))+1) * (2.2 / ((float64(player2Elo-player1Elo))*0.001 + 2.2))
		player2Elo, player1Elo = calculateElo(player2Elo, player2Matches, player1Elo, player1Matches, multiplier, false)
	} else {
		player1Elo, player2Elo = calculateElo(player1Elo, player1Matches, player2Elo, player2Matches, 1.0, true)
	}
	return player1Elo, player2Elo
}

func calculateElo(winnerElo int, winnerMatches int, looserElo int, looserMatches int, multiplier float64, isDraw bool) (int, int) {
	constant := 800.0

	Wwinner := 1.0
	Wlooser := 0.0
	if isDraw {
		Wwinner = 0.5
		Wlooser = 0.5
	}
	changeWinner := int((constant / float64(winnerMatches) * (Wwinner - (1 / (1 + math.Pow(10, float64(looserElo-winnerElo)/400))))) * multiplier)
	calculatedWinner := winnerElo + changeWinner

	changeLooser := int((constant / float64(looserMatches) * (Wlooser - (1 / (1 + math.Pow(10, float64(winnerElo-looserElo)/400))))) * multiplier)
	calculatedLooser := looserElo + changeLooser

	return calculatedWinner, calculatedLooser
}
