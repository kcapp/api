package models

import "github.com/guregu/null"

// PlayerElo struct used for storing elo information
type PlayerElo struct {
	PlayerID             int      `json:"player_id"`
	CurrentElo           int      `json:"current_elo"`
	CurrentEloMatches    int      `json:"current_elo_matches"`
	CurrentEloNew        int      `json:"current_elo_new,omitempty"`
	TournamentElo        null.Int `json:"tournament_elo,omitempty"`
	TournamentEloMatches int      `json:"tournament_elo_matches"`
	TournamentEloNew     null.Int `json:"tournament_elo_new,omitempty"`
	WinProbability       float64  `json:"win_probability,omitempty"`
}
