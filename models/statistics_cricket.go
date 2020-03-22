package models

import "github.com/guregu/null"

// StatisticsCricket struct used for storing statistics for cricket
type StatisticsCricket struct {
	ID             int      `json:"id,omitempty"`
	LegID          int      `json:"leg_id,omitempty"`
	PlayerID       int      `json:"player_id,omitempty"`
	MatchesPlayed  int      `json:"matches_played,omitempty"`
	MatchesWon     int      `json:"matches_won,omitempty"`
	LegsPlayed     int      `json:"legs_played,omitempty"`
	LegsWon        int      `json:"legs_won,omitempty"`
	TotalMarks     int      `json:"total_marks"`
	Rounds         int      `json:"rounds"`
	Score          null.Int `json:"score"`
	FirstNineMarks int      `json:"first_nine_marks"`
	MPR            float32  `json:"mpr"`
	FirstNineMPR   float32  `json:"first_nine_mpr"`
	Marks5         int      `json:"marks_5"`
	Marks6         int      `json:"marks_6"`
	Marks7         int      `json:"marks_7"`
	Marks8         int      `json:"marks_8"`
	Marks9         int      `json:"marks_9"`
}
