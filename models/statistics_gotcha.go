package models

import "github.com/guregu/null"

// StatisticsGotcha struct used for storing statistics for Gotcha
type StatisticsGotcha struct {
	ID            int      `json:"id"`
	LegID         int      `json:"leg_id"`
	PlayerID      int      `json:"player_id"`
	MatchesPlayed int      `json:"matches_played"`
	MatchesWon    int      `json:"matches_won"`
	LegsPlayed    int      `json:"legs_played"`
	LegsWon       int      `json:"legs_won"`
	OfficeID      null.Int `json:"office_id,omitempty"`
	DartsThrown   int      `json:"darts_thrown"`
	HighestScore  int      `json:"highest_score"`
	TimesReset    int      `json:"times_reset"`
	OthersReset   int      `json:"others_reset"`
	Score         int      `json:"score,omitempty"`
}
