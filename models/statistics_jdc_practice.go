package models

import "github.com/guregu/null"

// StatisticsJDCPractice struct used for storing statistics for JDC Practice Routine
type StatisticsJDCPractice struct {
	ID             int        `json:"id"`
	LegID          int        `json:"leg_id"`
	PlayerID       int        `json:"player_id"`
	MatchesPlayed  int        `json:"matches_played"`
	MatchesWon     int        `json:"matches_won"`
	LegsPlayed     int        `json:"legs_played"`
	LegsWon        int        `json:"legs_won"`
	OfficeID       null.Int   `json:"office_id,omitempty"`
	DartsThrown    int        `json:"darts_thrown,omitempty"`
	Score          int        `json:"score"`
	MPR            null.Float `json:"mpr,omitempty"`
	ShanghaiCount  int        `json:"shanghai_count"`
	DoublesHitrate float64    `json:"doubles_hitrate"`
	HighestScore   int        `json:"highest_score"`
	// Values used only to calculate statistics
	Marks int64 `json:"-"`
}
