package old

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/kcapp/api/util"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"

	"github.com/gorilla/mux"
	"github.com/jordic/goics"
)

// GetPlayerEloChangelog will return the elo changelog for the given player
func GetPlayerEloChangelog(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	start, err := strconv.Atoi(params["start"])
	if err != nil {
		log.Println("Invalid start parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(params["limit"])
	if err != nil {
		log.Println("Invalid limit parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	changelog, err := data.GetPlayerEloChangelog(id, start, limit)
	if err != nil {
		log.Println("Unable to get player elo changelog", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(changelog)
}

// GetPlayerMatchTypeStatistics will return statistics for the given player
func GetPlayerMatchTypeStatistics(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	matchType, err := strconv.Atoi(params["match_type"])
	if err != nil {
		log.Println("Invalid match type parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch matchType {
	case models.X01:
		stats, err := data.GetX01StatisticsForPlayer(id, models.X01)
		if err != nil {
			log.Println("Unable to get X01 statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.SHOOTOUT:
		stats, err := data.GetShootoutStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Cricket statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.X01HANDICAP:
		stats, err := data.GetX01StatisticsForPlayer(id, models.X01HANDICAP)
		if err != nil {
			log.Println("Unable to get X01 handicap statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.CRICKET:
		stats, err := data.GetCricketStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Cricket statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.DARTSATX:
		stats, err := data.GetDartsAtXStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Darts at X statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.AROUNDTHEWORLD:
		stats, err := data.GetAroundTheWorldStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Around The World Statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.SHANGHAI:
		stats, err := data.GetShanghaiStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Shanghai Statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.AROUNDTHECLOCK:
		stats, err := data.GetAroundTheClockStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Around the Clock Statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.TICTACTOE:
		stats, err := data.GetTicTacToeStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Tic Tac Toe Statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.BERMUDATRIANGLE:
		stats, err := data.GetBermudaTriangleStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Bermuda Triangle Statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.FOURTWENTY:
		stats, err := data.Get420StatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get 420 Statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.KILLBULL:
		stats, err := data.GetKillBullStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Kill Bull Statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.GOTCHA:
		stats, err := data.GetGotchaStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Gotcha Statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.JDCPRACTICE:
		stats, err := data.GetJDCPracticeStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get JDC Practice Statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.KNOCKOUT:
		stats, err := data.GetKnockoutStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Knockout Statistics for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
		return

	case models.SCAM:
		stats, err := data.GetScamStatisticsForPlayer(id)
		if err != nil {
			log.Println("Unable to get Scam Statistics for player", err)
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

// GetPlayerMatchTypeHistory will return history of match statistics for the given player
func GetPlayerMatchTypeHistory(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	matchType, err := strconv.Atoi(params["match_type"])
	if err != nil {
		log.Println("Invalid match type parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(params["limit"])
	if err != nil {
		log.Println("Invalid limit parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch matchType {
	case models.X01:
		legs, err := data.GetX01HistoryForPlayer(id, limit, models.X01)
		if err != nil {
			log.Println("Unable to get X01 history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.SHOOTOUT:
		legs, err := data.GetShootoutHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Shootout history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.X01HANDICAP:
		legs, err := data.GetX01HistoryForPlayer(id, limit, models.X01HANDICAP)
		if err != nil {
			log.Println("Unable to get X01 handicap history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.CRICKET:
		legs, err := data.GetCricketHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Cricket history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.DARTSATX:
		legs, err := data.GetDartsAtXHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Darts at X history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.AROUNDTHEWORLD:
		legs, err := data.GetAroundTheWorldHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Around The World history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.SHANGHAI:
		legs, err := data.GetShanghaiHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Shanghai history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.AROUNDTHECLOCK:
		legs, err := data.GetAroundTheClockHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Around the Clock history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.TICTACTOE:
		legs, err := data.GetTicTacToeHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Tic Tac Toe history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.BERMUDATRIANGLE:
		legs, err := data.GetBermudaTriangleHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Bermuda Triangle history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.FOURTWENTY:
		legs, err := data.Get420HistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get 420 history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.KILLBULL:
		legs, err := data.GetKillBullHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Kill Bull history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.GOTCHA:
		legs, err := data.GetGotchaHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Gotcha history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return
	case models.JDCPRACTICE:
		legs, err := data.GetJDCPracticeHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get JDC Practice history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.KNOCKOUT:
		legs, err := data.GetKnockoutHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Knockout history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	case models.SCAM:
		legs, err := data.GetScamHistoryForPlayer(id, limit)
		if err != nil {
			log.Println("Unable to get Scam history for player", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(legs)
		return

	default:
		log.Println("Unknown match type parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// AddPlayer will create a new player
func AddPlayer(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		log.Println("Unable to deserialize player json", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.AddPlayer(player)
	if err != nil {
		log.Println("Unable to add player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdatePlayer will update the given player
func UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var player models.Player
	err = json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		log.Println("Unable to deserialize player json", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.UpdatePlayer(id, player)
	if err != nil {
		log.Println("Unable to update player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetPlayerCheckouts will return all checkouts done by a player
func GetPlayerCheckouts(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	checkouts, err := data.GetPlayerCheckouts(id)
	if err != nil {
		log.Println("Unable to get player checkouts")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(checkouts)
}

// GetPlayerTournamentStandings will return all tournament standings for the given player
func GetPlayerTournamentStandings(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	standings, err := data.GetPlayerTournamentStandings(id)
	if err != nil {
		log.Println("Unable to get player tournament standings")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(standings)
}

// GetPlayerHeadToHead will return head to head statistics between the given players
func GetPlayerHeadToHead(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	player1, err := strconv.Atoi(params["player_1"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	player2, err := strconv.Atoi(params["player_2"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	head2head, err := data.GetPlayerHeadToHead(player1, player2)
	if err != nil {
		log.Println("Unable to get player head to head statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(head2head)
}

// SimulateMatch will return the result of a match between the two players
func SimulateMatch(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	player1, err := strconv.Atoi(params["player_1"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	player2, err := strconv.Atoi(params["player_2"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input struct {
		Player1Score int `json:"player1_score"`
		Player2Score int `json:"player2_score"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println("Unable to deserialize body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	elos, err := data.GetPlayersElo(player1, player2)
	if err != nil {
		log.Println("Unable to get player elos")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var output struct {
		Player1OldElo int `json:"player1_old_elo"`
		Player1NewElo int `json:"player1_new_elo"`
		Player2OldElo int `json:"player2_old_elo"`
		Player2NewElo int `json:"player2_new_elo"`
	}
	output.Player1OldElo = int(elos[0].TournamentElo.Int64)
	output.Player2OldElo = int(elos[1].TournamentElo.Int64)
	output.Player1NewElo, output.Player2NewElo = data.CalculateElo(output.Player1OldElo, elos[0].TournamentEloMatches, input.Player1Score,
		output.Player2OldElo, elos[1].TournamentEloMatches, input.Player2Score)

	json.NewEncoder(w).Encode(output)
}

// GetPlayerCalendar will return a calendar feed for all official matches for the given player
func GetPlayerCalendar(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	w.Header().Set("Content-type", "text/calendar")
	w.Header().Set("charset", "utf-8")
	w.Header().Set("Content-Disposition", "inline")
	w.Header().Set("filename", "kcapp-calendar.ics")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matches, err := data.GetPlayerOfficialMatches(id)
	if err != nil {
		log.Println("Unable to get official matches for player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	players, err := data.GetPlayers()
	if err != nil {
		log.Println("Unable to get players", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := models.Entries{}
	for _, match := range matches {
		if match.IsFinished {
			continue
		}
		entry := new(models.Entry)

		t := match.CreatedAt
		home := players[match.Players[0]]
		away := players[match.Players[1]]

		entry.DateStart = t
		entry.DateEnd = t.Add(time.Minute * time.Duration(30))
		entry.Summary = home.FirstName + " vs. " + away.FirstName
		location := "Dart Board"
		if match.Venue.Name.Valid {
			location = match.Venue.Name.String
		}
		entry.Location = location
		entry.Description = "Official Darts Match (" + strconv.Itoa(match.ID) + ") - " + home.FirstName + " vs. " + away.FirstName + " at " + location
		result = append(result, entry)
	}

	b := bytes.Buffer{}
	goics.NewICalEncode(&b).Encode(result)

	w.WriteHeader(http.StatusOK)
	w.Write(b.Bytes())
}

// GetRandomLegForPlayer will return a random leg for a given player and starting score
func GetRandomLegForPlayer(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	startingScore, err := strconv.Atoi(params["starting_score"])
	if err != nil {
		log.Println("Invalid starting score parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	visits, err := data.GetRandomLegForPlayer(id, startingScore)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Not enough data for player", http.StatusBadRequest)
			return
		}
		log.Println("Unable to get official matches for player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(visits)
}

// GetPlayerX01Statistics will return statistics for the given player
func GetPlayerX01Statistics(w http.ResponseWriter, r *http.Request) {
	util.SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := data.GetPlayerX01Statistics(id)
	if err != nil {
		log.Println("Unable to get player statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	visits, err := data.GetPlayerVisitCount(id)
	if err != nil {
		log.Println("Unable to get visits for player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stats.Visits = visits
	for _, v := range visits {
		stats.TotalVisits += v.Count
	}

	json.NewEncoder(w).Encode(stats)
}
