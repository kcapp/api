package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kcapp/api/controllers"
	controllers_v2 "github.com/kcapp/api/controllers/v2"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API",
	Run: func(cmd *cobra.Command, args []string) {
		configFileParam, _ := cmd.Flags().GetString("config")
		config, err := models.GetConfig(configFileParam)
		if err != nil {
			panic(err)
		}
		models.InitDB(config.GetMysqlConnectionString())
		GetLongestWinStreak()

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
		router.HandleFunc("/match/{id}", controllers.SetScore).Methods("PUT")
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
		//router.HandleFunc("/tournament/preset", controllers.AddTournamentPreset).Methods("POST")
		router.HandleFunc("/tournament/preset", controllers.GetTournamentPresets).Methods("GET")
		router.HandleFunc("/tournament/preset/{id}", controllers.GetTournamentPreset).Methods("GET")
		//router.HandleFunc("/tournament/preset/{id}", controllers.UpdateTournamentPreset).Methods("PUT")
		router.HandleFunc("/tournament/{id}", controllers.GetTournament).Methods("GET")
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

		log.Printf("Listening on port %d", config.APIConfig.Port)
		log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.APIConfig.Port), router))
	},
}

func GetLongestWinStreak() {
	// Establish a connection to the MySQL database
	db, err := sql.Open("mysql", "developer:abcd1234@tcp(localhost:3306)/kcapp")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Set maximum number of open and idle connections
	db.SetMaxOpenConns(10) // Adjust the value based on your MySQL server configuration
	db.SetMaxIdleConns(5)  // Adjust the value based on your MySQL server configuration

	// Set a maximum connection lifetime
	db.SetConnMaxLifetime(5 * time.Minute) // Adjust the value based on your application requirements

	// Get a list of all players
	players := []int{1, 2, 3, 4, 5, 6, 7, 8, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 35, 36, 37, 38, 39, 40, 42, 44, 46, 52, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 66, 67, 68, 69, 70, 76, 79, 82, 83, 84, 90, 97, 99, 102, 103, 106, 107, 108, 109, 110, 111, 112, 113, 114, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 135, 139, 141, 142, 143, 144, 148, 151, 152, 157, 159, 160, 163, 165, 168, 169, 175, 176, 177, 178, 179, 180, 181, 183, 184, 185, 187, 188, 191, 192, 194, 196, 198, 199, 200, 202, 206, 207, 208, 209, 210, 213, 217, 219, 221, 222, 223, 224, 225, 226, 229, 236, 237, 238, 239, 240, 250, 252, 253, 254, 255, 256, 257, 258, 260, 261, 262, 264, 276, 277, 278, 281, 286, 287, 290, 292, 296, 297, 298, 303, 304, 306, 307, 308, 311, 313, 314, 318, 322, 331, 332, 333, 334, 343, 344, 345, 347, 350, 361, 363, 364, 365, 366, 368, 371, 376, 380, 385, 387, 389, 401, 402, 409, 411, 412, 413, 414, 417, 418, 420, 422, 428, 429, 432, 435} // Replace with your actual list of player IDs

	// Initialize variables to track the longest win streak
	longestStreak := 0
	playerWithLongestStreak := 0

	// Iterate over each player
	for _, playerID := range players {
		// Dynamically update the query
		query := fmt.Sprintf(`
            SELECT
                winner_id,
                MAX(streak) AS longest_streak
            FROM (
                SELECT
                    winner_id,
                    @streak := IF(@prev_winner = winner_id, @streak + 1, 1) AS streak,
                    @prev_winner := winner_id
                FROM
                    (SELECT @streak := 0, @prev_winner := NULL) AS vars,
                    (SELECT * FROM matches m WHERE id IN (SELECT DISTINCT match_id FROM player2leg WHERE player_id = %d) AND is_finished = 1 AND match_type_id = 1 ORDER BY updated_at) AS m
            ) AS streaks
			WHERE winner_id = %d
            GROUP BY winner_id
            ORDER BY longest_streak DESC
        `, playerID, playerID)

		// Execute the query
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Fetch the result
		var winnerID, streak int
		if rows.Next() {
			err := rows.Scan(&winnerID, &streak)
			if err != nil {
				log.Fatal(err)
			}

			// Check if the player has a longer streak than the current longest
			if streak > longestStreak {
				longestStreak = streak
				playerWithLongestStreak = winnerID
			}
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Player ID = %d, Streak = %d\n", winnerID, streak)
	}
	fmt.Printf("Player with the longest win streak: Player ID = %d, Streak = %d\n", playerWithLongestStreak, longestStreak)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
