package models

import "github.com/guregu/null"

// StatisticsDartsAtX struct used for storing statistics for cricket
type StatisticsDartsAtX struct {
	ID            int      `json:"id,omitempty"`
	LegID         int      `json:"leg_id,omitempty"`
	PlayerID      int      `json:"player_id,omitempty"`
	MatchesPlayed int      `json:"matches_played,omitempty"`
	MatchesWon    int      `json:"matches_won,omitempty"`
	LegsPlayed    int      `json:"legs_played,omitempty"`
	LegsWon       int      `json:"legs_won,omitempty"`
	AvgScore      int      `json:"avg_score"`
	Score         null.Int `json:"score,omitempty"`
	Singles       int      `json:"singles"`
	Doubles       int      `json:"doubles"`
	Triples       int      `json:"triples"`
	HitRate       float32  `json:"hit_rate"`
}
