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

// PlayerEloChangelogs struct used for storing elo changelog information
type PlayerEloChangelogs struct {
	Total     int                   `json:"total"`
	Changelog []*PlayerEloChangelog `json:"changelog"`
}

// PlayerEloChangelog struct used for storing elo changelog information
type PlayerEloChangelog struct {
	ID         int        `json:"id"`
	MatchID    int        `json:"match_id"`
	FinishedAt string     `json:"finished_at"`
	MatchMode  string     `json:"match_mode"`
	MatchType  string     `json:"match_type"`
	IsOfficial bool       `json:"is_official"`
	WinnerID   null.Int   `json:"winner_id"`
	HomePlayer *PlayerElo `json:"home_player"`
	AwayPlayer *PlayerElo `json:"away_player"`
}
