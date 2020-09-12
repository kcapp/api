package models

import "github.com/guregu/null"

// StatisticsAroundThe struct used for storing statistics for around the world/clock/shanghai
type StatisticsAroundThe struct {
	ID            int             `json:"id"`
	LegID         int             `json:"leg_id"`
	PlayerID      int             `json:"player_id"`
	MatchesPlayed int             `json:"matches_played"`
	MatchesWon    int             `json:"matches_won"`
	LegsPlayed    int             `json:"legs_played"`
	LegsWon       int             `json:"legs_won"`
	OfficeID      null.Int        `json:"office_id,omitempty"`
	DartsThrown   int             `json:"darts_thrown,omitempty"`
	Score         int             `json:"score,omitempty"`
	Shanghai      null.Int        `json:"shanghai,omitempty"`
	MPR           null.Float      `json:"mpr,omitempty"`
	LongestStreak null.Int        `json:"longest_streak,omitempty"`
	TotalHitRate  float64         `json:"total_hit_rate"`
	Hitrates      map[int]float64 `json:"hitrates,omitempty"`

	// Values used only to calculate statistics
	Marks         int64 `json:"-"`
	CurrentStreak int64 `json:"-"`
}
