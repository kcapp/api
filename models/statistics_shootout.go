package models

import "github.com/guregu/null"

// StatisticsShootout struct used for storing statistics for shootout legs
type StatisticsShootout struct {
	ID            int             `json:"id,omitempty"`
	LegID         int             `json:"leg_id,omitempty"`
	PlayerID      int             `json:"player_id,omitempty"`
	MatchesPlayed int             `json:"matches_played,omitempty"`
	MatchesWon    int             `json:"matches_won,omitempty"`
	LegsPlayed    int             `json:"legs_played,omitempty"`
	LegsWon       int             `json:"legs_won"`
	OfficeID      null.Int        `json:"office_id,omitempty"`
	Score         int             `json:"score"`
	PPD           float32         `json:"ppd"`
	DartsThrown   null.Int        `json:"darts_thrown"`
	Score60sPlus  int             `json:"scores_60s_plus"`
	Score100sPlus int             `json:"scores_100s_plus"`
	Score140sPlus int             `json:"scores_140s_plus"`
	Score180s     int             `json:"scores_180s"`
	Hits          map[int64]*Hits `json:"hits,omitempty"`
	HighestScore  int             `json:"highest_score"`
}
