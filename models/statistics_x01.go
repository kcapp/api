package models

import (
	"github.com/guregu/null"
)

// StatisticsX01 struct used for storing statistics
type StatisticsX01 struct {
	ID                 int             `json:"id,omitempty"`
	MatchID            int             `json:"match_id,omitempty"`
	PlayerID           int             `json:"player_id,omitempty"`
	PPD                float32         `json:"ppd"`
	FirstNinePPD       float32         `json:"first_nine_ppd"`
	CheckoutPercentage float32         `json:"checkout_percentage"`
	DartsThrown        int             `json:"darts_thrown"`
	Score60sPlus       int             `json:"scores_60s_plus"`
	Score100sPlus      int             `json:"scores_100s_plus"`
	Score140sPlus      int             `json:"scores_140s_plus"`
	Score180s          int             `json:"scores_180s"`
	Accuracy20         null.Float      `json:"accuracy_20"`
	Accuracy19         null.Float      `json:"accuracy_19"`
	AccuracyOverall    null.Float      `json:"accuracy_overall"`
	Hits               map[int64]*Hits `json:"hits,omitempty"`
	GamesPlayed        int             `json:"games_played,omitempty"`
	GamesWon           int             `json:"games_won,omitempty"`
	BestPPD            float32         `json:"best_ppd,omitempty"`
	BestFirstNinePPD   float32         `json:"best_first_nine_ppd,omitempty"`
	Best301            int             `json:"best_301,omitempty"`
	Best501            int             `json:"best_501,omitempty"`
	HighestCheckout    int             `json:"highest_checkout,omitempty"`
	StartingScore      null.Int        `json:"-"`
}

// Hits sturct used to store summary of hits for players/matches
type Hits struct {
	Singles int `json:"1,omitempty"`
	Doubles int `json:"2,omitempty"`
	Triples int `json:"3,omitempty"`
}
