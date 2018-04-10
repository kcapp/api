package models

import (
	"github.com/guregu/null"
)

const (
	// X01 constant representing GameType 1
	X01 = 1
	// SHOOTOUT constant representing GameType 2
	SHOOTOUT = 2
	// X01HANDICAP constant representing GameType 3
	X01HANDICAP = 3
)

// Game struct used for storing games
type Game struct {
	ID              int         `json:"id"`
	IsFinished      bool        `json:"is_finished"`
	CurrentMatchID  null.Int    `json:"current_match_id"`
	WinnerID        null.Int    `json:"winner_id"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
	EndTime         string      `json:"end_time,omitempty"`
	GameType        *GameType   `json:"game_type"`
	GameMode        *GameMode   `json:"game_mode"`
	OweTypeID       null.Int    `json:"owe_type_id"`
	OweType         *OweType    `json:"owe_type,omitempty"`
	Players         []int       `json:"players"`
	Matches         []*Match    `json:"matches,omitempty"`
	PlayerHandicaps map[int]int `json:"player_handicaps,omitempty"`
}

// GameType struct used for storing game types
type GameType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GameMode struct used for stroing game modes
type GameMode struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	ShortName       string   `json:"short_name"`
	WinsRequired    int      `json:"wins_required"`
	MatchesRequired null.Int `json:"matches_required"`
}
