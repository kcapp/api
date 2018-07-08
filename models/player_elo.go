package models

// PlayerElo struct used for storing elo information
type PlayerElo struct {
	PlayerID             int     `json:"player_id"`
	CurrentElo           int     `json:"current_elo"`
	CurrentEloMatches    int     `json:"current_elo_matches"`
	CurrentEloNew        int     `json:"-"`
	TournamentElo        int     `json:"tournament_elo"`
	TournamentEloMatches int     `json:"tournament_elo_matches"`
	TournamentEloNew     int     `json:"-"`
	WinProbability       float64 `json:"win_probability,omitempty"`
}
