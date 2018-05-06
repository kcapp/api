package models

import (
	"github.com/guregu/null"
)

// Player struct used for storing players
type Player struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Nickname      null.String `json:"nickname,omitempty"`
	MatchesPlayed int         `json:"matches_played"`
	MatchesWon    int         `json:"matches_won"`
	LegsPlayed    int         `json:"legs_played"`
	LegsWon       int         `json:"legs_won"`
	Color         null.String `json:"color,omitempty"`
	ProfilePicURL null.String `json:"profile_pic_url,omitempty"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at,omitempty"`
}
