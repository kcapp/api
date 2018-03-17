package models

import (
	"github.com/guregu/null"
)

// Match struct used for storing matches
type Match struct {
	ID              int             `json:"id"`
	Endtime         null.String     `json:"end_time"`
	StartingScore   int             `json:"starting_score"`
	IsFinished      bool            `json:"is_finished"`
	CurrentPlayerID int             `json:"current_player_id"`
	WinnerPlayerID  null.Int        `json:"winner_player_id"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
	GameID          int             `json:"game_id"`
	Players         []int           `json:"players,omitempty"`
	DartsThrown     int             `json:"darts_thrown,omitempty"`
	Visits          []*Visit        `json:"visits"`
	Hits            map[int64]*Hits `json:"hits,omitempty"`
}

// Player2Match struct used for stroring players in a match
type Player2Match struct {
	MatchID         int              `json:"match_id"`
	PlayerID        int              `json:"player_id"`
	Order           int              `json:"order"`
	CurrentScore    int              `json:"current_score"`
	IsCurrentPlayer bool             `json:"is_current_player"`
	Wins            int              `json:"wins,omitempty"`
	Modifiers       *PlayerModifiers `json:"modifiers,omitempty"`
}

// PlayerModifiers struct used for storing visit modifiers for a player
type PlayerModifiers struct {
	IsViliusVisit bool `json:"is_vilius_visit"`
}
