package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/guregu/null"
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
)

// GetTournaments will return all tournaments
func GetTournaments(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	tournaments, err := data.GetTournaments()
	if err != nil {
		log.Println("Unable to get tournaments", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tournaments)
}

// AddTournamentGroup will add a new tournament group
func AddTournamentGroup(w http.ResponseWriter, r *http.Request) {
	var group models.TournamentGroup
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		log.Println("Unable to deserialize group json", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = data.AddTournamentGroup(group)
	if err != nil {
		log.Println("Unable to add group", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetTournamentGroups will return all tournaments
func GetTournamentGroups(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	groups, err := data.GetTournamentGroups()
	if err != nil {
		log.Println("Unable to get tournament groups", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(groups)
}

// GetTournament will return the given tournament
func GetTournament(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tournament, err := data.GetTournament(id)
	if err != nil {
		log.Println("Unable to get tournament", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tournament)
}

// GetCurrentTournament will return the current active tournament
func GetCurrentTournament(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	tournament, err := data.GetCurrentTournament()
	if err != nil {
		log.Println("Unable to get tournament", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tournament == nil {
		http.Error(w, "No current tournament for office", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(tournament)
}

// GetCurrentTournamentForOffice will return the current active tournament for a given office
func GetCurrentTournamentForOffice(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	officeID, err := strconv.Atoi(params["office_id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tournament, err := data.GetCurrentTournamentForOffice(officeID)
	if err != nil {
		log.Println("Unable to get tournament for office", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tournament == nil {
		http.Error(w, "No current tournament for office", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(tournament)
}

// GetTournamentsForOffice will return all tournaments for given office
func GetTournamentsForOffice(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	officeID, err := strconv.Atoi(params["office_id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tournament, err := data.GetTournamentsForOffice(officeID)
	if err != nil {
		log.Println("Unable to get tournaments for office", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tournament == nil {
		http.Error(w, "No tournaments for office", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(tournament)
}

func GetTournamentProbabilities(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	prob, err := data.GetTournamentProbabilities(id)
	if err != nil {
		log.Println("Unable to get tournament probabilities", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(prob)
}

// GetTournamentMatches will return all matches for the given tournament
func GetTournamentMatches(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	matches, err := data.GetTournamentMatches(id)
	if err != nil {
		log.Println("Unable to get tournament matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}

// GetTournamentMatchResults will reurn match results for the given ID
func GetTournamentMatchResults(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tournamentMatches, err := data.GetTournamentMatches(id)
	if err != nil {
		log.Println("Unable to get tournament matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	players, err := data.GetPlayers()
	if err != nil {
		log.Println("Unable to get players: ", err)
		http.Error(w, "Unable to get players", http.StatusBadRequest)
		return
	}

	output := make([]interface{}, 0)
	for _, matches := range tournamentMatches {
		for _, match := range matches {
			homePlayer := match.Players[0]
			awayPlayer := match.Players[1]
			homeWins := 0
			awayWins := 0
			for _, winnerID := range match.LegsWon {
				if homePlayer == winnerID {
					homeWins++
				} else if awayPlayer == winnerID {
					awayWins++
				}
			}

			output = append(output, struct {
				MatchID      int       `json:"match_id"`
				IsFinished   bool      `json:"is_finished"`
				IsWalkover   bool      `json:"is_walkover"`
				MatchTime    time.Time `json:"scheduled_time"`
				WinnerID     null.Int  `json:"winner_id"`
				HomeScore    int       `json:"home_score"`
				HomePlayerID int       `json:"home_player_id"`
				HomePlayer   string    `json:"home_player_name"`
				AwayScore    int       `json:"away_score"`
				AWayPlayerID int       `json:"away_player_id"`
				AwayPlayer   string    `json:"away_player_name"`
			}{
				match.ID,
				match.IsFinished,
				match.IsWalkover,
				match.CreatedAt,
				match.WinnerID,
				homeWins,
				homePlayer,
				players[homePlayer].GetName(),
				awayWins,
				awayPlayer,
				players[awayPlayer].GetName(),
			})
		}

	}

	json.NewEncoder(w).Encode(output)
}

// GetTournamentOverview will return statistics for the given tournament
func GetTournamentOverview(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	stats, err := data.GetTournamentOverview(id)
	if err != nil {
		log.Println("Unable to get tournament overview", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetTournamentStatistics will return statistics for the given tournament
func GetTournamentStatistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	stats, err := data.GetTournamentStatistics(id)
	if err != nil {
		log.Println("Unable to get tournament statistics", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetNextTournamentMatch will return the next tournament match
func GetNextTournamentMatch(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	match, err := data.GetNextTournamentMatch(id)
	if err != nil {
		log.Println("Unable to get next tournament match", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if match == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	json.NewEncoder(w).Encode(match)
}

// GetTournamentStandings will return statistics for the given tournament
func GetTournamentStandings(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	stats, err := data.GetTournamentStandings()
	if err != nil {
		log.Println("Unable to get tournament standings", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// NewTournament will create a new tournament
func NewTournament(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	var tournamentInput models.Tournament
	err := json.NewDecoder(r.Body).Decode(&tournamentInput)
	if err != nil {
		log.Println("Unable to deserialize body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tournament, err := data.NewTournament(tournamentInput)
	if err != nil {
		log.Println("Unable to create new tournament", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tournament)
}

// GetTournamentPlayerMatches will return all matches for the given tournament and player
func GetTournamentPlayerMatches(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	playerID, err := strconv.Atoi(params["player_id"])
	if err != nil {
		log.Println("Invalid player id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matches, err := data.GetTournamentMatchesForPlayer(id, playerID)
	if err != nil {
		log.Println("Unable to get official matches for player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}
