package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/kcapp/api/models"

	"github.com/kcapp/api/controllers"

	"github.com/gorilla/mux"
)

// our main function
func main() {

	dbConfig, err := models.GetConfig()
	if err != nil {
		log.Fatal(errors.New("Could not load db config yaml file"))
	}

	connectionString, err := dbConfig.GetMysqlConnectionString()
	if err != nil {
		log.Fatal(errors.New("Malformed connection string"))
	}

	models.InitDB(connectionString)

	router := mux.NewRouter()
	router.HandleFunc("/game", controllers.NewGame).Methods("POST")
	router.HandleFunc("/game/{id}/continue", controllers.ContinueGame).Methods("PUT")
	router.HandleFunc("/game", controllers.GetGames).Methods("GET")
	router.HandleFunc("/game/{id}", controllers.GetGame).Methods("GET")
	router.HandleFunc("/game/{id}/statistics", controllers.GetX01StatisticsForGame).Methods("GET")
	router.HandleFunc("/game/{id}/matches", controllers.GetMatchesForGame).Methods("GET")

	router.HandleFunc("/match/{id}", controllers.GetMatch).Methods("GET")
	router.HandleFunc("/match/{id}/statistics", controllers.GetX01StatisticsForMatch).Methods("GET")
	router.HandleFunc("/match/{id}/players", controllers.GetMatchPlayers).Methods("GET")
	router.HandleFunc("/match/{id}/order", controllers.ChangePlayerOrder).Methods("PUT")

	router.HandleFunc("/visit", controllers.AddVisit).Methods("POST")
	router.HandleFunc("/visit/{id}/modify", controllers.ModifyVisit).Methods("PUT")
	router.HandleFunc("/visit/{id}", controllers.DeleteVisit).Methods("DELETE")

	router.HandleFunc("/player", controllers.GetPlayers).Methods("GET")
	router.HandleFunc("/player/{id}", controllers.GetPlayer).Methods("GET")
	router.HandleFunc("/player/{id}/statistics", controllers.GetPlayerStatistics).Methods("GET")
	router.HandleFunc("/player/compare/", controllers.GetPlayersStatistics).Methods("GET")
	router.HandleFunc("/player", controllers.AddPlayer).Methods("POST")

	router.HandleFunc("/statistics/x01/{from}/{to}", controllers.GetX01Statistics).Methods("GET")

	router.HandleFunc("/owe", controllers.GetOwes).Methods("GET")
	router.HandleFunc("/owe/payback", controllers.RegisterPayback).Methods("PUT")

	apiPort, err := dbConfig.GetAPIPort()
	if err != nil {
		log.Fatal(errors.New("Could not read API port"))
	}

	log.Println(http.ListenAndServe(fmt.Sprintf(":%s", apiPort), router))
}
