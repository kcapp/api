package models

import (
	"github.com/guregu/null"
)

// Statistics170 struct used for storing statistics
type Statistics170 struct {
	ID                 int         `json:"id"`
	LegID              int         `json:"leg_id"`
	PlayerID           int         `json:"player_id"`
	MatchesPlayed      int         `json:"matches_played"`
	MatchesWon         int         `json:"matches_won"`
	LegsPlayed         int         `json:"legs_played"`
	LegsWon            int         `json:"legs_won"`
	Points             int         `json:"points"`
	PPD                float32     `json:"ppd"`
	PPDScore           int         `json:"-"`
	ThreeDartAvg       float32     `json:"three_dart_avg"`
	Rounds             int         `json:"rounds"`
	CheckoutPercentage null.Float  `json:"checkout_percentage"`
	CheckoutAttempts   int         `json:"checkout_attempts"`
	CheckoutCompleted  int         `json:"checkout_completed"`
	HighestCheckout    null.Int    `json:"highest_checkout,omitempty"`
	DartsThrown        int         `json:"darts_thrown"`
	CheckoutDarts      map[int]int `json:"checkout_darts"`
}
