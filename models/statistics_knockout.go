package models

import "github.com/guregu/null"

// StatisticsKnockout struct used for storing statistics for Knockout
type StatisticsKnockout struct {
	ID            int      `json:"id"`
	LegID         int      `json:"leg_id"`
	PlayerID      int      `json:"player_id"`
	MatchesPlayed int      `json:"matches_played"`
	MatchesWon    int      `json:"matches_won"`
	LegsPlayed    int      `json:"legs_played"`
	LegsWon       int      `json:"legs_won"`
	OfficeID      null.Int `json:"office_id,omitempty"`
	DartsThrown   int      `json:"darts_thrown,omitempty"`
	AvgScore      float64  `json:"avg_score"`
	LivesLost     int      `json:"lives_lost"`
	LivesTaken    int      `json:"lives_taken"`
	FinalPosition int      `json:"final_position"`
}
