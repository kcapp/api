package models

import "github.com/guregu/null"

// StatisticsKillBull struct used for storing statistics for kill bull
type StatisticsKillBull struct {
	ID            int      `json:"id,omitempty"`
	MatchesPlayed int      `json:"matches_played"`
	MatchesWon    int      `json:"matches_won"`
	LegsPlayed    int      `json:"legs_played"`
	LegsWon       int      `json:"legs_won"`
	OfficeID      null.Int `json:"office_id,omitempty"`
	LegID         int      `json:"leg_id,omitempty"`
	PlayerID      int      `json:"player_id,omitempty"`
	DartsThrown   int      `json:"darts_thrown"`
	Score         int      `json:"score"`
	Marks3        int      `json:"marks_3"`
	Marks4        int      `json:"marks_4"`
	Marks5        int      `json:"marks_5"`
	Marks6        int      `json:"marks_6"`
	LongestStreak int      `json:"longest_streak"`
	TimesBusted   int      `json:"times_busted"`
	TotalHitRate  float64  `json:"total_hit_rate"`

	CurrentStreak int `json:"-"`
}
