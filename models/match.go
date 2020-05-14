package models

import (
	"encoding/json"
	"strconv"

	"github.com/guregu/null"
)

const (
	// X01 constant representing type 1
	X01 = 1
	// SHOOTOUT constant representing type 2
	SHOOTOUT = 2
	// X01HANDICAP constant representing type 3
	X01HANDICAP = 3
	// CRICKET constant representing type 4
	CRICKET = 4
	// DARTSATX constant representing type 5
	DARTSATX = 5
	// AROUNDTHEWORLD constant representing type 6
	AROUNDTHEWORLD = 6
	// SHANGHAI constant representing type 7
	SHANGHAI = 7
	// AROUNDTHECLOCK constant representing type 9
	AROUNDTHECLOCK = 8
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

// MarshalJSON will marshall the given object to JSON
func (match Match) MarshalJSON() ([]byte, error) {
	// Use a type to get consistnt order of JSON key-value pairs.
	type matchJSON struct {
		ID               int                `json:"id"`
		CurrentLegID     null.Int           `json:"current_leg_id"`
		CreatedAt        string             `json:"created_at"`
		UpdatedAt        string             `json:"updated_at"`
		EndTime          string             `json:"end_time,omitempty"`
		MatchType        *MatchType         `json:"match_type"`
		MatchMode        *MatchMode         `json:"match_mode"`
		WinnerID         null.Int           `json:"winner_id"`
		IsFinished       bool               `json:"is_finished"`
		IsAbandoned      bool               `json:"is_abandoned"`
		IsWalkover       bool               `json:"is_walkover"`
		OfficeID         null.Int           `json:"office_id,omitempty"`
		OweTypeID        null.Int           `json:"owe_type_id"`
		VenueID          null.Int           `json:"venue_id"`
		IsPractice       bool               `json:"is_practice"`
		Venue            *Venue             `json:"venue"`
		OweType          *OweType           `json:"owe_type,omitempty"`
		TournamentID     null.Int           `json:"tournament_id,omitempty"`
		Tournament       *MatchTournament   `json:"tournament,omitempty"`
		Players          []int              `json:"players"`
		Legs             []*Leg             `json:"legs,omitempty"`
		CurrentLegNumber string             `json:"current_leg_num"`
		PlayerHandicaps  map[int]int        `json:"player_handicaps,omitempty"`
		BotPlayerConfig  map[int]*BotConfig `json:"bot_player_config,omitempty"`
		FirstThrow       null.String        `json:"first_throw_time,omitempty"`
		LastThrow        null.String        `json:"last_throw_time,omitempty"`
		EloChange        map[int]*PlayerElo `json:"elo_change,omitempty"`
		LegsWon          []int              `json:"legs_won,omitempty"`
	}
	legPostfix := [4]string{"st", "nd", "rd", "th"}
	idx := ((len(match.Legs)+90)%100-10)%10 - 1
	if idx < 0 {
		idx = 0
	} else if idx > 3 {
		idx = 3
	}

	legNum := strconv.Itoa(len(match.Legs)) + legPostfix[idx]
	return json.Marshal(matchJSON{
		ID:               match.ID,
		CurrentLegID:     match.CurrentLegID,
		CreatedAt:        match.CreatedAt,
		UpdatedAt:        match.UpdatedAt,
		EndTime:          match.EndTime,
		MatchType:        match.MatchType,
		MatchMode:        match.MatchMode,
		WinnerID:         match.WinnerID,
		IsFinished:       match.IsFinished,
		IsAbandoned:      match.IsAbandoned,
		IsWalkover:       match.IsWalkover,
		OfficeID:         match.OfficeID,
		OweTypeID:        match.OweTypeID,
		VenueID:          match.VenueID,
		IsPractice:       match.IsPractice,
		Venue:            match.Venue,
		OweType:          match.OweType,
		TournamentID:     match.TournamentID,
		Tournament:       match.Tournament,
		Players:          match.Players,
		Legs:             match.Legs,
		CurrentLegNumber: legNum,
		PlayerHandicaps:  match.PlayerHandicaps,
		BotPlayerConfig:  match.BotPlayerConfig,
		FirstThrow:       match.FirstThrow,
		LastThrow:        match.LastThrow,
		EloChange:        match.EloChange,
		LegsWon:          match.LegsWon,
	})
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
