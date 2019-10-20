package controllers

import (
	"encoding/json"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"
	"bytes"

	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"

	"github.com/gorilla/mux"
	"github.com/jordic/goics"
)

// GetPlayers will return a map containing all players
func GetPlayers(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	players, err := data.GetPlayers()
	if err != nil {
		log.Println("Unable to get players", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

// GetActivePlayers will return a map containing all active players
func GetActivePlayers(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	players, err := data.GetActivePlayers()
	if err != nil {
		log.Println("Unable to get active players", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

// GetPlayer will return a player with the given ID
func GetPlayer(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	player, err := data.GetPlayer(id)
	if err != nil {
		log.Println("Unable to get player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(player)
}

// GetPlayerX01Statistics will return statistics for the given player
func GetPlayerX01Statistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
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

// GetPlayerX01PreviousStatistics will return statistics for the given player
func GetPlayerX01PreviousStatistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := data.GetPlayerX01PreviousStatistics(id)
	if err != nil {
		log.Println("Unable to get player statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

// GetPlayersX01Statistics will return statistics for the given players
func GetPlayersX01Statistics(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := r.URL.Query()["id"]
	if params == nil {
		log.Println("No players specified to compare")
		http.Error(w, "No players specified to compare", http.StatusBadRequest)
		return
	}
	ids, err := sliceAtoi(params)
	if err != nil {
		log.Println("Unable to convert params to int")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := data.GetPlayersX01Statistics(ids)
	if err != nil {
		log.Println("Unable to get players statistics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// AddPlayer will create a new player
func AddPlayer(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	_ = json.NewDecoder(r.Body).Decode(&player)
	err := data.AddPlayer(player)
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
	_ = json.NewDecoder(r.Body).Decode(&player)
	err = data.UpdatePlayer(id, player)
	if err != nil {
		log.Println("Unable to update player", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetPlayerProgression will return statistics for the given player
func GetPlayerProgression(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := data.GetPlayerProgression(id)
	if err != nil {
		log.Println("Unable to get player progression", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// GetPlayerCheckouts will return all checkouts done by a player
func GetPlayerCheckouts(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
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
	SetHeaders(w)
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
	SetHeaders(w)
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
	SetHeaders(w)
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
	SetHeaders(w)
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

		t, err := time.Parse("2006-01-02 15:04:05", match.CreatedAt)
		if err != nil {
			log.Println("Unable to parse time")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
		entry.Description =  "Official Darts Match (" + strconv.Itoa(match.ID) + ") - " + home.FirstName + " vs. " + away.FirstName + " at " + location
		result = append(result, entry)
	}

	b := bytes.Buffer{}
	goics.NewICalEncode(&b).Encode(result)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(b.String()))
}

// GetRandomLegForPlayer will return a random leg for a given player and starting score
func GetRandomLegForPlayer(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)

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

func sliceAtoi(sa []string) ([]int, error) {
	si := make([]int, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.Atoi(a)
		if err != nil {
			return si, err
		}
		si = append(si, i)
	}
	return si, nil
}
