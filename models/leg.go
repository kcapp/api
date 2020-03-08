package models

import (
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
	Hits            map[int]int64    `json:"hits"`
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
		}
	}
}

// AddVisitHits adds hits from the visit to the player
func (p2l *Player2Leg) AddVisitHits(visit Visit) {
	addInMap(p2l.Hits, int(visit.FirstDart.Value.Int64), visit.FirstDart.Multiplier)
	if visit.SecondDart.Value.Valid {
		addInMap(p2l.Hits, int(visit.SecondDart.Value.Int64), visit.SecondDart.Multiplier)
	}
	if visit.ThirdDart.Value.Valid {
		addInMap(p2l.Hits, int(visit.ThirdDart.Value.Int64), visit.ThirdDart.Multiplier)
	}

}

// addInMap will add the given value to the key if it doesn't exist, or add the value to the existing value
func addInMap(m map[int]int64, k int, v int64) {
	if _, ok := m[k]; !ok {
		m[k] = v
	} else {
		m[k] += v
	}
}
