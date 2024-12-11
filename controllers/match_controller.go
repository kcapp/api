package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"

	"github.com/gorilla/mux"
)

// NewMatch will start a new match
func NewMatch(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	var matchInput models.Match
	err := json.NewDecoder(r.Body).Decode(&matchInput)
	if err != nil {
		log.Println("Unable to deserialize body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match, err := data.NewMatch(matchInput)
	if err != nil {
		switch t := err.(type) {
		default:
			log.Println("Unable to start new match", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		case *models.MatchConfigError:
			log.Println("Unable to start new match", t)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	json.NewEncoder(w).Encode(match)
}

// ReMatch will start a new match with same settings as the given match ID
func ReMatch(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	match, err := data.GetMatch(id)
	if err != nil {
		log.Println("Unable to get match: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reverse order of players for rematch
	for i := len(match.Players)/2 - 1; i >= 0; i-- {
		opp := len(match.Players) - 1 - i
		match.Players[i], match.Players[opp] = match.Players[opp], match.Players[i]
	}
	players, err := data.GetPlayersScore(int(match.CurrentLegID.Int64))
	if err != nil {
		log.Println("Unable to get players in match: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	match.BotPlayerConfig = make(map[int]*models.BotConfig)
	for _, player := range players {
		if player.BotConfig != nil {
			match.BotPlayerConfig[player.PlayerID] = player.BotConfig
		}
	}
	match.CreatedAt = time.Now().UTC()
	match, err = data.NewMatch(*match)
	if err != nil {
		log.Println("Unable to rematch: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(match)
}

// GetMatches will return a list of all matches
func GetMatches(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	matches, err := data.GetMatches()
	if err != nil {
		log.Println("Unable to get matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}

// GetActiveMatches will return a list of active matches
func GetActiveMatches(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	since, err := strconv.Atoi(r.URL.Query().Get("since"))
	if err != nil {
		since = 2
	}
	matches, err := data.GetActiveMatches(since)
	if err != nil {
		log.Println("Unable to get active matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}

func GetMatchProbabilities(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	prob, err := data.GetMatchProbabilities(id)
	if err != nil {
		log.Println("Unable to get match probabilities", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(prob)
}

// GetMatchesLimit will return N matches from the given starting point
func GetMatchesLimit(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
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
	matches, err := data.GetMatchesLimit(start, limit)
	if err != nil {
		log.Println("Unable to get matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	count, err := data.GetMatchesCount()
	if err != nil {
		log.Println("Unable to get matches", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(count))
	json.NewEncoder(w).Encode(matches)
}

// GetMatch will reurn a the match with the given ID
func GetMatch(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	match, err := data.GetMatch(id)
	if err != nil {
		log.Println("Unable to get match: ", err)
		http.Error(w, "Unable to get match", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(match)
}

// SetScore will set the score of a given match
func SetScore(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input models.MatchResult
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println("Unable to deserialize body", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match, err := data.SetScore(id, input)
	if err != nil {
		log.Println("Unable to set score for match: ", err)
		http.Error(w, "Unable to set score for match", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(match)
}

// GetMatchMetadata will return metadata for the given match
func GetMatchMetadata(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metadata, err := data.GetMatchMetadata(id)
	if err != nil {
		log.Println("Unable to get match metadata: ", err)
		http.Error(w, "Unable to get match metadata", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(metadata)
}

// GetMatchMetadataForTournament will return metadata for all matches in a tournament
func GetMatchMetadataForTournament(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	tournamentID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metadata, err := data.GetMatchMetadataForTournament(tournamentID)
	if err != nil {
		log.Println("Unable to get match metadata for tournament: ", err)
		http.Error(w, "Unable to get match metadata for tournament", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(metadata)
}

// GetStatisticsForMatch will return statistics for all players in the given match
func GetStatisticsForMatch(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	params := mux.Vars(r)
	matchID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Invalid id parameter")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match, err := data.GetMatch(matchID)
	if err != nil {
		log.Printf("Unable to get Match %d", matchID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if match.MatchType.ID == models.SHOOTOUT {
		stats, err := data.GetShootoutStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get shootout statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.CRICKET {
		stats, err := data.GetCricketStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get cricket statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.DARTSATX {
		stats, err := data.GetDartsAtXStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get darts at x statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.AROUNDTHECLOCK {
		stats, err := data.GetAroundTheClockStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get around the clock statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.AROUNDTHEWORLD {
		stats, err := data.GetAroundTheWorldStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get around the world statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.SHANGHAI {
		stats, err := data.GetShanghaiStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get shanghai statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.TICTACTOE {
		stats, err := data.GetTicTacToeStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get tic tac toe statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.BERMUDATRIANGLE {
		stats, err := data.GetBermudaTriangleStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get bermuda triangle statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.FOURTWENTY {
		stats, err := data.Get420StatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get 420 statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.KILLBULL {
		stats, err := data.GetKillBullStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get Kill Bull statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.GOTCHA {
		stats, err := data.GetGotchaStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get Gotcha statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.JDCPRACTICE {
		stats, err := data.GetJDCPracticeStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get JDC Practice statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.KNOCKOUT {
		stats, err := data.GetKnockoutStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get Knockout statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.SCAM {
		stats, err := data.GetScamStatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get Scam statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else if match.MatchType.ID == models.ONESEVENTY {
		stats, err := data.Get170StatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get 170 statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	} else {
		stats, err := data.GetX01StatisticsForMatch(matchID)
		if err != nil {
			log.Printf("Unable to get x01 statistics for match %d: %s", matchID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(stats)
	}
}

// GetMatchesModes will return all match modes
func GetMatchesModes(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	modes, err := data.GetMatchModes()
	if err != nil {
		log.Println("Unable to get match modes", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(modes)
}

// GetMatchesTypes will return all match types
func GetMatchesTypes(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	types, err := data.GetMatchTypes()
	if err != nil {
		log.Println("Unable to get match types", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(types)
}

// GetOutshotTypes will return all outshot types
func GetOutshotTypes(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	types, err := data.GetOutshotTypes()
	if err != nil {
		log.Println("Unable to get outshot types", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(types)
}
