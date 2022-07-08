package models

import "github.com/guregu/null"

// StatisticsJDCPractice struct used for storing statistics for JDC Practice Routine
type StatisticsScam struct {
	ID                 int      `json:"id"`
	LegID              int      `json:"leg_id"`
	PlayerID           int      `json:"player_id"`
	MatchesPlayed      int      `json:"matches_played"`
	MatchesWon         int      `json:"matches_won"`
	LegsPlayed         int      `json:"legs_played"`
	LegsWon            int      `json:"legs_won"`
	OfficeID           null.Int `json:"office_id,omitempty"`
	DartsThrownStopper int      `json:"darts_thrown_stopper"`
	DartsThrownScorer  int      `json:"darts_thrown_scorer"`
	MPR                float32  `json:"mpr"`
	PPD                float32  `json:"ppd"`
	PPDScore           int      `json:"-"`
	ThreeDartAvg       float32  `json:"three_dart_avg"`
	Score              int      `json:"score"`
}
