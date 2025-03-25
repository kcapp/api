package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/controllers"
	controllers_v2 "github.com/kcapp/api/controllers/v2"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API",
	Run: func(cmd *cobra.Command, args []string) {
		models.InitDB(models.GetMysqlConnectionString())

		router := mux.NewRouter()
		router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
			w.WriteHeader(http.StatusNoContent)
			return
		})

		router.HandleFunc("/health", controllers.Healthcheck).Methods("HEAD")

		router.HandleFunc("/match", controllers.NewMatch).Methods("POST")
		router.HandleFunc("/match/active", controllers.GetActiveMatches).Methods("GET")
		router.HandleFunc("/match/types", controllers.GetMatchesTypes).Methods("GET")
		router.HandleFunc("/match/modes", controllers.GetMatchesModes).Methods("GET")
		router.HandleFunc("/match/outshot", controllers.GetOutshotTypes).Methods("GET")
		router.HandleFunc("/match", controllers.GetMatches).Methods("GET")
		router.HandleFunc("/match/{id}", controllers.GetMatch).Methods("GET")
		router.HandleFunc("/match/{id}", controllers.UpdateMatch).Methods("PUT")
		router.HandleFunc("/match/{id}/score", controllers.SetScore).Methods("PUT")
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
		router.HandleFunc("/leg/{id}/warmup", controllers.StartWarmup).Methods("PUT")
		router.HandleFunc("/leg/{id}/undo", controllers.UndoFinishLeg).Methods("PUT")
		router.HandleFunc("/leg/{id}/finish", controllers.FinishLeg).Methods("PUT")

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
		router.HandleFunc("/player/{id}/hits", controllers.GetPlayerHits).Methods("PUT")
		router.HandleFunc("/player/{id}/statistics/previous", controllers.GetPlayerX01PreviousStatistics).Methods("GET")
		router.HandleFunc("/player/{id}/progression", controllers.GetPlayerProgression).Methods("GET")
		router.HandleFunc("/player/{id}/checkouts", controllers.GetPlayerCheckouts).Methods("GET")
		router.HandleFunc("/player/{id}/tournament", controllers.GetPlayerTournamentStandings).Methods("GET")
		router.HandleFunc("/player/{id}/badges", controllers.GetPlayerBadges).Methods("GET")
		router.HandleFunc("/player/{id}/elo/{start}/{limit}", controllers.GetPlayerEloChangelog).Methods("GET")
		router.HandleFunc("/player/{player_1}/vs/{player_2}", controllers.GetPlayerHeadToHead).Methods("GET")
		router.HandleFunc("/player/{player_1}/vs/{player_2}/simulate", controllers.SimulateMatch).Methods("PUT")
		router.HandleFunc("/player", controllers.AddPlayer).Methods("POST")
		router.HandleFunc("/player/{id}/calendar", controllers.GetPlayerCalendar).Methods("GET")
		router.HandleFunc("/player/{id}/random/{starting_score}", controllers.GetRandomLegForPlayer).Methods("GET")
		router.HandleFunc("/player/{id}/statistics/{match_type}", controllers.GetPlayerMatchTypeStatistics).Methods("GET")
		router.HandleFunc("/player/{id}/statistics/{match_type}/history/{limit}", controllers.GetPlayerMatchTypeHistory).Methods("GET")

		// v2
		router.HandleFunc("/players", controllers_v2.GetPlayers).Methods("GET")

		router.HandleFunc("/preset", controllers.AddPreset).Methods("POST")
		router.HandleFunc("/preset", controllers.GetPresets).Methods("GET")
		router.HandleFunc("/preset/{id}", controllers.GetPreset).Methods("GET")
		router.HandleFunc("/preset/{id}", controllers.UpdatePreset).Methods("PUT")
		router.HandleFunc("/preset/{id}", controllers.DeletePreset).Methods("DELETE")

		router.HandleFunc("/option/default", controllers.GetDefaultOptions).Methods("GET")

		router.HandleFunc("/statistics/global", controllers.GetGlobalStatistics).Methods("GET")
		router.HandleFunc("/statistics/global/fnc", controllers.GetGlobalStatisticsFnc).Methods("GET")
		router.HandleFunc("/statistics/office/{from}/{to}", controllers.GetOfficeStatistics).Methods("GET")
		router.HandleFunc("/statistics/office/{office_id}/{from}/{to}", controllers.GetOfficeStatistics).Methods("GET")
		router.HandleFunc("/statistics/{dart}/hits", controllers.GetDartStatistics).Methods("GET")
		router.HandleFunc("/statistics/x01/player/{legs}", controllers.GetPlayersLastXLegsStatistics).Methods("GET")
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
		router.HandleFunc("/venue/{id}/players", controllers.GetRecentPlayers).Methods("GET")
		router.HandleFunc("/venue/{id}/matches", controllers.GetActiveVenueMatches).Methods("GET")

		router.HandleFunc("/tournament", controllers.NewTournament).Methods("POST")
		router.HandleFunc("/tournament/generate", controllers.GenerateTournament).Methods("POST")
		router.HandleFunc("/tournament/generate/playoffs/{id}", controllers.GeneratePlayoffsTournament).Methods("POST")
		router.HandleFunc("/tournament", controllers.GetTournaments).Methods("GET")
		router.HandleFunc("/tournament/current", controllers.GetCurrentTournament).Methods("GET")
		router.HandleFunc("/tournament/current/{office_id}", controllers.GetCurrentTournamentForOffice).Methods("GET")
		router.HandleFunc("/tournament/office/{office_id}", controllers.GetTournamentsForOffice).Methods("GET")
		router.HandleFunc("/tournament/groups", controllers.AddTournamentGroup).Methods("POST")
		router.HandleFunc("/tournament/groups", controllers.GetTournamentGroups).Methods("GET")
		router.HandleFunc("/tournament/standings", controllers.GetTournamentStandings).Methods("GET")
		router.HandleFunc("/tournament/preset", controllers.GetTournamentPresets).Methods("GET")
		router.HandleFunc("/tournament/preset/{id}", controllers.GetTournamentPreset).Methods("GET")
		router.HandleFunc("/tournament/{id}", controllers.GetTournament).Methods("GET")
		router.HandleFunc("/tournament/{id}/player", controllers.AddPlayerToTournament).Methods("POST")
		router.HandleFunc("/tournament/{id}/player/{player_id}", controllers.GetTournamentPlayerMatches).Methods("GET")
		router.HandleFunc("/tournament/{id}/matches", controllers.GetTournamentMatches).Methods("GET")
		router.HandleFunc("/tournament/{id}/matches/result", controllers.GetTournamentMatchResults).Methods("GET")
		router.HandleFunc("/tournament/{id}/metadata", controllers.GetMatchMetadataForTournament).Methods("GET")
		router.HandleFunc("/tournament/{id}/overview", controllers.GetTournamentOverview).Methods("GET")
		router.HandleFunc("/tournament/{id}/statistics", controllers.GetTournamentStatistics).Methods("GET")
		router.HandleFunc("/tournament/match/{id}/next", controllers.GetNextTournamentMatch).Methods("GET")
		router.HandleFunc("/tournament/{id}/probabilities", controllers.GetTournamentProbabilities).Methods("GET")
		router.HandleFunc("/tournament/match/{id}/probabilities", controllers.GetMatchProbabilities).Methods("GET")

		router.HandleFunc("/badge", controllers.GetBadges).Methods("GET")
		router.HandleFunc("/badge/statistics", controllers.GetBadgesStatistics).Methods("GET")
		router.HandleFunc("/badge/{id}", controllers.GetBadge).Methods("GET")
		router.HandleFunc("/badge/{id}/statistics", controllers.GetBadgeStatistics).Methods("GET")

		port := viper.GetInt("api.port")
		log.Printf("Listening on port %d", port)
		log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), router))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
