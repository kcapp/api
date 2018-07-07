package models

import (
	"github.com/guregu/null"
)

const (
	// X01 constant representing MatchType 1
	X01 = 1
	// SHOOTOUT constant representing MatchType 2
	SHOOTOUT = 2
	// X01HANDICAP constant representing MatchType 3
	X01HANDICAP = 3
)

// Match struct used for storing matches
type Match struct {
	ID              int              `json:"id"`
	CurrentLegID    null.Int         `json:"current_leg_id"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
	EndTime         string           `json:"end_time,omitempty"`
	MatchType       *MatchType       `json:"match_type"`
	MatchMode       *MatchMode       `json:"match_mode"`
	WinnerID        null.Int         `json:"winner_id"`
	IsFinished      bool             `json:"is_finished"`
	IsAbandoned     bool             `json:"is_abandoned"`
	IsWalkover      bool             `json:"is_walkover"`
	OweTypeID       null.Int         `json:"owe_type_id"`
	VenueID         null.Int         `json:"venue_id"`
	Venue           *Venue           `json:"venue"`
	OweType         *OweType         `json:"owe_type,omitempty"`
	TournamentID    null.Int         `json:"tournament_id,omitempty"`
	Tournament      *MatchTournament `json:"tournament,omitempty"`
	Players         []int            `json:"players"`
	Legs            []*Leg           `json:"legs,omitempty"`
	PlayerHandicaps map[int]int      `json:"player_handicaps,omitempty"`
	FirstThrow      null.String      `json:"first_throw_time,omitempty"`
	LastThrow       null.String      `json:"last_throw_time,omitempty"`
}

// MatchType struct used for storing match types
type MatchType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// MatchMode struct used for storing match modes
type MatchMode struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	ShortName    string   `json:"short_name"`
	WinsRequired int      `json:"wins_required"`
	LegsRequired null.Int `json:"legs_required"`
}

// Venue struct used for storing venues
type Venue struct {
	ID          null.Int    `json:"id"`
	Name        null.String `json:"name"`
	Description null.String `json:"description"`
}

// MatchTournament struct for storing tournament information
type MatchTournament struct {
	TournamentID        null.Int    `json:"tournament_id"`
	TournamentName      null.String `json:"tournament_name"`
	TournamentGroupID   null.Int    `json:"tournament_group_id"`
	TournamentGroupName null.String `json:"tournament_group_name"`
}
