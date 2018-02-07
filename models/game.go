package models

import (
	"github.com/guregu/null"
)

// Game struct used for storing games
type Game struct {
	ID             int       `json:"id"`
	IsFinished     bool      `json:"is_finished"`
	CurrentMatchID null.Int  `json:"current_match_id"`
	GameType       *GameType `json:"game_type"`
	WinnerID       null.Int  `json:"winner_id"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
	OweTypeID      null.Int  `json:"owe_type_id"`
	OweType        *OweType  `json:"owe_type,omitempty"`
	Players        []int     `json:"players"`
	Matches        []*Match  `json:"matches,omitempty"`
}

// GameType struct used for storing game types
type GameType struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	ShortName       string   `json:"short_name"`
	WinsRequired    int      `json:"wins_required"`
	MatchesRequired null.Int `json:"matches_required"`
}
