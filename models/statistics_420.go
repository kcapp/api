package models

import "github.com/guregu/null"

// Statistics420 struct used for storing statistics for 420
type Statistics420 struct {
	ID            int             `json:"id"`
	LegID         int             `json:"leg_id"`
	PlayerID      int             `json:"player_id"`
	MatchesPlayed int             `json:"matches_played"`
	MatchesWon    int             `json:"matches_won"`
	LegsPlayed    int             `json:"legs_played"`
	LegsWon       int             `json:"legs_won"`
	OfficeID      null.Int        `json:"office_id,omitempty"`
	Score         int             `json:"score,omitempty"`
	TotalHitRate  float64         `json:"total_hit_rate"`
	Hitrates      map[int]float64 `json:"hitrates,omitempty"`
}
