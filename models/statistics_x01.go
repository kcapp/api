package models

import (
	"time"

	"github.com/guregu/null"
)

// BestStatistic struct used for storing a value and leg where the statistic was achieved
type BestStatistic struct {
	Value    int `json:"value"`
	LegID    int `json:"leg_id"`
	PlayerID int `json:"player_id"`
}

// BestStatisticFloat struct used for storing a value and leg where the statistic was achieved
type BestStatisticFloat struct {
	Value    float32 `json:"value"`
	LegID    int     `json:"leg_id"`
	PlayerID int     `json:"player_id"`
}

// StatisticsX01 struct used for storing statistics
type StatisticsX01 struct {
	ID                    int                 `json:"id,omitempty"`
	LegID                 int                 `json:"leg_id,omitempty"`
	PlayerID              int                 `json:"player_id,omitempty"`
	WinnerID              int                 `json:"winner_id,omitempty"`
	PPD                   float32             `json:"ppd"`
	PPDScore              int                 `json:"-"`
	FirstNinePPD          float32             `json:"first_nine_ppd"`
	FirstNinePPDScore     int                 `json:"-"`
	ThreeDartAvg          float32             `json:"three_dart_avg"`
	FirstNineThreeDartAvg float32             `json:"first_nine_three_dart_avg"`
	CheckoutPercentage    null.Float          `json:"checkout_percentage"`
	CheckoutAttempts      int                 `json:"checkout_attempts,omitempty"`
	Checkout              null.Int            `json:"checkout,omitempty"`
	DartsThrown           int                 `json:"darts_thrown,omitempty"`
	TotalVisits           int                 `json:"total_visits,omitempty"`
	Score60sPlus          int                 `json:"scores_60s_plus"`
	Score100sPlus         int                 `json:"scores_100s_plus"`
	Score140sPlus         int                 `json:"scores_140s_plus"`
	Score180s             int                 `json:"scores_180s"`
	Accuracy20            null.Float          `json:"accuracy_20"`
	Accuracy19            null.Float          `json:"accuracy_19"`
	AccuracyOverall       null.Float          `json:"accuracy_overall"`
	AccuracyStatistics    *AccuracyStatistics `json:"accuracy,omitempty"`
	Visits                []*Visit            `json:"visits,omitempty"`
	Hits                  map[int64]*Hits     `json:"hits,omitempty"`
	MatchesPlayed         int                 `json:"matches_played"`
	MatchesWon            int                 `json:"matches_won"`
	LegsPlayed            int                 `json:"legs_played"`
	LegsWon               int                 `json:"legs_won"`
	OfficeID              null.Int            `json:"office_id,omitempty"`
	BestThreeDartAvg      *BestStatisticFloat `json:"best_three_dart_avg,omitempty"`
	BestFirstNineAvg      *BestStatisticFloat `json:"best_first_nine_avg,omitempty"`
	Best301               *BestStatistic      `json:"best_301,omitempty"`
	Best501               *BestStatistic      `json:"best_501,omitempty"`
	Best701               *BestStatistic      `json:"best_701,omitempty"`
	HighestCheckout       *BestStatistic      `json:"highest_checkout,omitempty"`
	StartingScore         null.Int            `json:"-"`
	LastPlayedLeg         time.Time           `json:"last_played_leg,omitempty"`
}

// PlayerX01Progression struct used for storing player statistics in a bucket
type PlayerX01Progression struct {
	PlayerID     int            `json:"player_id"`
	Bucket       int            `json:"bucket"`
	FirstLegID   int            `json:"first_leg_id"`
	LastLegID    int            `json:"last_leg_id"`
	LegsInBucket int            `json:"legs_in_bucket"`
	StartDate    time.Time      `json:"start_date"`
	EndDate      time.Time      `json:"end_date"`
	Statistics   *StatisticsX01 `json:"statistics"`
}

// GlobalStatistics struct used for storing global statistics
type GlobalStatistics struct {
	FishNChips     int `json:"fish_n_chips"`
	Matches        int `json:"matches"`
	Legs           int `json:"legs"`
	Visits         int `json:"visits"`
	Darts          int `json:"darts"`
	Points         int `json:"points"`
	PointsBusted   int `json:"points_busted"`
	Score180s      int `json:"score_180s"`
	ScoreBullseyes int `json:"score_bullseyes"`
}

// CheckoutStatistics stuct used for storing detailed checkout statistics
type CheckoutStatistics struct {
	Checkout         int         `json:"checkout"`
	Count            int         `json:"count"`
	Completed        bool        `json:"completed"`
	Visits           []*Visit    `json:"visits,omitempty"`
	CheckoutAttempts map[int]int `json:"checkout_attempts,omitempty"`
}

// Hits struct used to store summary of hits for players/legs
type Hits struct {
	Singles int `json:"1,omitempty"`
	Doubles int `json:"2,omitempty"`
	Triples int `json:"3,omitempty"`
	Total   int `json:"total,omitempty"`
}

// Add will add the given dart to the hits map
func (h *Hits) Add(dart *Dart) {
	if dart.IsTriple() {
		h.Triples++
	} else if dart.IsDouble() {
		h.Doubles++
	} else {
		h.Singles++
	}
	h.Total += int(dart.Multiplier)

}

// OfficeStatistics struct used for storing statistics for a office
type OfficeStatistics struct {
	PlayerID int      `json:"player_id,omitempty"`
	LegID    int      `json:"leg_id,omitempty"`
	OfficeID null.Int `json:"office_id,omitempty"`
	Checkout int      `json:"checkout,omitempty"`
	Darts    []*Dart  `json:"darts"`
}
