package main

import (
	"kcapp-api/controllers"
	"kcapp-api/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// our main function
func main() {
	models.InitDB("developer:abcd1234@(10.12.100.66:3306)/cakeapp_dev")

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

	router.HandleFunc("/match/{id}/modify", controllers.ModifyMatchVisit).Methods("PUT")

	router.HandleFunc("/player", controllers.GetPlayers).Methods("GET")
	router.HandleFunc("/player/{id}", controllers.GetPlayer).Methods("GET")
	router.HandleFunc("/player/{id}/statistics", controllers.GetPlayerStatistics).Methods("GET")
	router.HandleFunc("/player/compare/", controllers.GetPlayersStatistics).Methods("GET")
	router.HandleFunc("/player", controllers.AddPlayer).Methods("POST")

	router.HandleFunc("/statistics/x01/{from}/{to}", controllers.GetX01Statistics).Methods("GET")

	router.HandleFunc("/owe", controllers.GetOwes).Methods("GET")
	router.HandleFunc("/owe/payback", controllers.RegisterPayback).Methods("PUT")

	log.Println(http.ListenAndServe(":8001", router))
}
