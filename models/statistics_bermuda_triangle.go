package models

import "github.com/guregu/null"

// StatisticsBermudaTriangle struct used for storing statistics for Bermuda Triangle
type StatisticsBermudaTriangle struct {
	ID                  int             `json:"id"`
	LegID               int             `json:"leg_id"`
	PlayerID            int             `json:"player_id"`
	MatchesPlayed       int             `json:"matches_played"`
	MatchesWon          int             `json:"matches_won"`
	LegsPlayed          int             `json:"legs_played"`
	LegsWon             int             `json:"legs_won"`
	OfficeID            null.Int        `json:"office_id,omitempty"`
	DartsThrown         int             `json:"darts_thrown,omitempty"`
	Score               int             `json:"score,omitempty"`
	TotalMarks          int             `json:"total_marks,omitempty"`
	MPR                 float64         `json:"mpr,omitempty"`
	HighestScoreReached int             `json:"highest_score_reached,omitempty"`
	TotalHitRate        float64         `json:"total_hit_rate"`
	Hitrates            map[int]float64 `json:"hitrates,omitempty"`
	HitCount            int             `json:"hit_count,omitempty"`
}
