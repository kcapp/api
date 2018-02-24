package models

import (
	"github.com/guregu/null"
)

// Player struct used for storing players
type Player struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Nickname    null.String `json:"nickname,omitempty"`
	GamesPlayed int         `json:"games_played"`
	GamesWon    int         `json:"games_won"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at,omitempty"`
}
