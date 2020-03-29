package models

import (
	"encoding/json"
	"math"

	"github.com/guregu/null"
)

// Leg struct used for storing legs
type Leg struct {
	ID                 int                 `json:"id"`
	Endtime            null.String         `json:"end_time"`
	StartingScore      int                 `json:"starting_score"`
	IsFinished         bool                `json:"is_finished"`
	CurrentPlayerID    int                 `json:"current_player_id"`
	WinnerPlayerID     null.Int            `json:"winner_player_id"`
	CreatedAt          string              `json:"created_at"`
	UpdatedAt          string              `json:"updated_at"`
	BoardStreamURL     null.String         `json:"board_stream_url,omitempty"`
	MatchID            int                 `json:"match_id"`
	HasScores          bool                `json:"has_scores"`
	Players            []int               `json:"players,omitempty"`
	DartsThrown        int                 `json:"darts_thrown,omitempty"`
	Visits             []*Visit            `json:"visits"`
	Hits               map[int64]*Hits     `json:"hits,omitempty"`
	CheckoutStatistics *CheckoutStatistics `json:"checkout_statistics,omitempty"`
}

// MarshalJSON will marshall the given object to JSON
func (leg Leg) MarshalJSON() ([]byte, error) {
	// Use a type to get consistnt order of JSON key-value pairs.
	type legJSON struct {
		ID                 int                 `json:"id"`
		Endtime            null.String         `json:"end_time"`
		StartingScore      int                 `json:"starting_score"`
		IsFinished         bool                `json:"is_finished"`
		CurrentPlayerID    int                 `json:"current_player_id"`
		WinnerPlayerID     null.Int            `json:"winner_player_id"`
		CreatedAt          string              `json:"created_at"`
		UpdatedAt          string              `json:"updated_at"`
		BoardStreamURL     null.String         `json:"board_stream_url,omitempty"`
		MatchID            int                 `json:"match_id"`
		HasScores          bool                `json:"has_scores"`
		Round              int                 `json:"round"`
		Players            []int               `json:"players,omitempty"`
		DartsThrown        int                 `json:"darts_thrown,omitempty"`
		Visits             []*Visit            `json:"visits"`
		Hits               map[int64]*Hits     `json:"hits,omitempty"`
		CheckoutStatistics *CheckoutStatistics `json:"checkout_statistics,omitempty"`
	}
	round := int(math.Floor(float64(len(leg.Visits))/float64(len(leg.Players))) + 1)

	return json.Marshal(legJSON{
		ID:                 leg.ID,
		Endtime:            leg.Endtime,
		StartingScore:      leg.StartingScore,
		IsFinished:         leg.IsFinished,
		CurrentPlayerID:    leg.CurrentPlayerID,
		WinnerPlayerID:     leg.WinnerPlayerID,
		CreatedAt:          leg.CreatedAt,
		UpdatedAt:          leg.UpdatedAt,
		BoardStreamURL:     leg.BoardStreamURL,
		MatchID:            leg.MatchID,
		HasScores:          leg.HasScores,
		Round:              round,
		Players:            leg.Players,
		DartsThrown:        leg.DartsThrown,
		Visits:             leg.Visits,
		Hits:               leg.Hits,
		CheckoutStatistics: leg.CheckoutStatistics,
	})
}

// Player2Leg struct used for stroring players in a leg
type Player2Leg struct {
	LegID           int              `json:"leg_id"`
	PlayerID        int              `json:"player_id"`
	PlayerName      string           `json:"player_name"`
	Order           int              `json:"order"`
	CurrentScore    int              `json:"current_score"`
	IsCurrentPlayer bool             `json:"is_current_player"`
	Wins            int              `json:"wins,omitempty"`
	VisitStatistics *VisitStatistics `json:"visit_statistics,omitempty"`
	Handicap        null.Int         `json:"handicap,omitempty"`
	Modifiers       *PlayerModifiers `json:"modifiers,omitempty"`
	Player          *Player          `json:"player,omitempty"`
	BotConfig       *BotConfig       `json:"bot_config,omitempty"`
	Hits            map[int]*Hits    `json:"hits"`
	DartsThrown     int              `json:"darts_thrown,omitempty"`
}

// BotConfig struct used for storing bot configuration
type BotConfig struct {
	PlayerID null.Int `json:"player_id"`
	Skill    null.Int `json:"skill_level"`
}

// VisitStatistics tells about the
type VisitStatistics struct {
	FishAndChipsCounter int `json:"fish_and_chips_counter"`
	ViliusVisitCounter  int `json:"vilius_visit_counter"`
}

// PlayerModifiers struct used for storing visit modifiers for a player
type PlayerModifiers struct {
	IsViliusVisit  bool `json:"is_vilius_visit"`
	IsBeerMatch    bool `json:"is_beer_match"`
	IsFishAndChips bool `json:"is_fish_and_chips"`
}

// AddVisitStatistics adds information about
func (p2l *Player2Leg) AddVisitStatistics(leg Leg) {
	p2l.VisitStatistics = new(VisitStatistics)
	for _, visit := range leg.Visits {
		if visit.PlayerID == p2l.PlayerID {
			if visit.IsFishAndChips() {
				p2l.VisitStatistics.FishAndChipsCounter++
			}
			if visit.IsViliusVisit() {
				p2l.VisitStatistics.ViliusVisitCounter++
			}
			p2l.DartsThrown = visit.DartsThrown
		}
	}
}
