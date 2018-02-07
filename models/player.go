package models

import (
	"github.com/guregu/null"
)

// Player struct used for storing players
type Player struct {
	ID           int         `json:"id"`
	Name         string      `json:"name"`
	Nickname     null.String `json:"nickname,omitempty"`
	GamesPlayed  int         `json:"games_played"`
	GamesWon     int         `json:"games_won"`
	PPD          float32     `json:"ppd,omitempty"`
	FirstNinePPD float32     `json:"first_nine_ppd,omitempty"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at,omitempty"`
}
