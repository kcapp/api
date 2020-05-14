package models

import "github.com/guregu/null"

// StatisticsAroundThe struct used for storing statistics for around the world/clock/shanghai
type StatisticsAroundThe struct {
	ID            int             `json:"id,omitempty"`
	LegID         int             `json:"leg_id,omitempty"`
	PlayerID      int             `json:"player_id,omitempty"`
	MatchesPlayed int             `json:"matches_played,omitempty"`
	MatchesWon    int             `json:"matches_won,omitempty"`
	LegsPlayed    int             `json:"legs_played,omitempty"`
	LegsWon       int             `json:"legs_won,omitempty"`
	DartsThrown   int             `json:"darts_thrown,omitempty"`
	Score         int             `json:"score,omitempty"`
	Shanghai      null.Int        `json:"shanghai,omitempty"`
	MPR           null.Float      `json:"mpr,omitempty"`
	LongestStreak null.Int        `json:"longest_streak,omitempty"`
	TotalHitRate  float32         `json:"total_hit_rate"`
	Hitrates      map[int]float32 `json:"hitrates,omitempty"`

	// Values used only to calculate statistics
	Marks         int64 `json:"-"`
	CurrentStreak int64 `json:"-"`
}
