package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kcapp/api/models"

	"github.com/kcapp/api/controllers"

	"github.com/gorilla/mux"
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
	router.HandleFunc("/game", controllers.NewGame).Methods("POST")
	router.HandleFunc("/game/types", controllers.GetGamesTypes).Methods("GET")
	router.HandleFunc("/game/modes", controllers.GetGamesModes).Methods("GET")
	router.HandleFunc("/game/{id}/continue", controllers.ContinueGame).Methods("PUT")
	router.HandleFunc("/game", controllers.GetGames).Methods("GET")
	router.HandleFunc("/game/{id}", controllers.GetGame).Methods("GET")
	router.HandleFunc("/game/{id}/statistics", controllers.GetX01StatisticsForGame).Methods("GET")
	router.HandleFunc("/game/{id}/matches", controllers.GetMatchesForGame).Methods("GET")

	router.HandleFunc("/match/active", controllers.GetActiveMatches).Methods("GET")
	router.HandleFunc("/match/{id}", controllers.GetMatch).Methods("GET")
	router.HandleFunc("/match/{id}", controllers.DeleteMatch).Methods("DELETE")
	router.HandleFunc("/match/{id}/statistics", controllers.GetX01StatisticsForMatch).Methods("GET")
	router.HandleFunc("/match/{id}/players", controllers.GetMatchPlayers).Methods("GET")
	router.HandleFunc("/match/{id}/order", controllers.ChangePlayerOrder).Methods("PUT")
	router.HandleFunc("/match/{id}/finish", controllers.FinishMatch).Methods("PUT")

	router.HandleFunc("/visit", controllers.AddVisit).Methods("POST")
	router.HandleFunc("/visit/{id}/modify", controllers.ModifyVisit).Methods("PUT")
	router.HandleFunc("/visit/{id}", controllers.DeleteVisit).Methods("DELETE")

	router.HandleFunc("/player", controllers.GetPlayers).Methods("GET")
	router.HandleFunc("/player/compare", controllers.GetPlayersStatistics).Methods("GET")
	router.HandleFunc("/player/{id}", controllers.GetPlayer).Methods("GET")
	router.HandleFunc("/player/{id}/statistics", controllers.GetPlayerStatistics).Methods("GET")
	router.HandleFunc("/player/{id}/progression", controllers.GetPlayerProgression).Methods("GET")
	router.HandleFunc("/player", controllers.AddPlayer).Methods("POST")

	router.HandleFunc("/statistics/x01/{from}/{to}", controllers.GetX01Statistics).Methods("GET")
	router.HandleFunc("/statistics/shootout/{from}/{to}", controllers.GetShootoutStatistics).Methods("GET")

	router.HandleFunc("/owe", controllers.GetOwes).Methods("GET")
	router.HandleFunc("/owe/payback", controllers.RegisterPayback).Methods("PUT")

	router.HandleFunc("/owetype", controllers.GetOweTypes).Methods("GET")

	log.Printf("Listening on port %d", config.APIConfig.Port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.APIConfig.Port), router))
}
