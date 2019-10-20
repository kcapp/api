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
	ID              int                `json:"id"`
	CurrentLegID    null.Int           `json:"current_leg_id"`
	CreatedAt       string             `json:"created_at"`
	UpdatedAt       string             `json:"updated_at"`
	EndTime         string             `json:"end_time,omitempty"`
	MatchType       *MatchType         `json:"match_type"`
	MatchMode       *MatchMode         `json:"match_mode"`
	WinnerID        null.Int           `json:"winner_id"`
	IsFinished      bool               `json:"is_finished"`
	IsAbandoned     bool               `json:"is_abandoned"`
	IsWalkover      bool               `json:"is_walkover"`
	OfficeID        null.Int           `json:"office_id,omitempty"`
	OweTypeID       null.Int           `json:"owe_type_id"`
	VenueID         null.Int           `json:"venue_id"`
	IsPractice      bool               `json:"is_practice"`
	Venue           *Venue             `json:"venue"`
	OweType         *OweType           `json:"owe_type,omitempty"`
	TournamentID    null.Int           `json:"tournament_id,omitempty"`
	Tournament      *MatchTournament   `json:"tournament,omitempty"`
	Players         []int              `json:"players"`
	Legs            []*Leg             `json:"legs,omitempty"`
	PlayerHandicaps map[int]int        `json:"player_handicaps,omitempty"`
	BotPlayerConfig map[int]*BotConfig `json:"bot_player_config,omitempty"`
	FirstThrow      null.String        `json:"first_throw_time,omitempty"`
	LastThrow       null.String        `json:"last_throw_time,omitempty"`
	EloChange       map[int]*PlayerElo `json:"elo_change,omitempty"`
	LegsWon         []int              `json:"legs_won,omitempty"`
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

// MatchTournament struct for storing tournament information
type MatchTournament struct {
	TournamentID        null.Int    `json:"tournament_id"`
	TournamentName      null.String `json:"tournament_name"`
	TournamentGroupID   null.Int    `json:"tournament_group_id"`
	TournamentGroupName null.String `json:"tournament_group_name"`
	OfficeID            null.Int    `json:"office_id"`
}

// MatchMetadata struct used for storing metadata about matches
type MatchMetadata struct {
	ID                   int              `json:"id"`
	MatchID              int              `json:"match_id"`
	OrderOfPlay          int              `json:"order_of_play"`
	TournamentGroup      *TournamentGroup `json:"tournament_group"`
	HomePlayer           int              `json:"player_home"`
	AwayPlayer           int              `json:"player_away"`
	MatchDisplayname     string           `json:"match_displayname"`
	Elimination          bool             `json:"elimination"`
	Trophy               bool             `json:"trophy"`
	Promotion            bool             `json:"promotion"`
	SemiFinal            bool             `json:"semi_final"`
	GrandFinal           bool             `json:"grand_final"`
	WinnerOutcomeMatchID null.Int         `json:"winner_outcome_match_id"`
	IsWinnerOutcomeHome  bool             `json:"is_winner_outcome_home"`
	LooserOutcomeMatchID null.Int         `json:"looser_outcome_match_id"`
	IsLooserOutcomeHome  bool             `json:"is_looser_outcome_home"`
	WinnerOutcome        null.String      `json:"winner_outcome"`
	LooserOutcome        null.String      `json:"looser_outcome"`
}
