package models

// PlayerElo struct used for storing elo information
type PlayerElo struct {
	PlayerID             int `json:"id"`
	CurrentElo           int `json:"current_elo"`
	CurrentEloMatches    int `json:"current_elo_matches"`
	TournamentElo        int `json:"tournament_elo"`
	TournamentEloMatches int `json:"tournament_elo_matches"`
}
