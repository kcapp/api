package models

// TournamentStatistics struct for stroring tournament statistics
type TournamentStatistics struct {
	HighestCheckout    []*BestStatistic      `json:"checkout_highest,omitempty"`
	BestPPD            []*BestStatisticFloat `json:"best_ppd,omitempty"`
	BestFirstNinePPD   []*BestStatisticFloat `json:"best_first_nine_ppd,omitempty"`
	Best301DartsThrown []*BestStatistic      `json:"best_301_darts_thrown,omitempty"`
	Best501DartsThrown []*BestStatistic      `json:"best_501_darts_thrown,omitempty"`
	Best701DartsThrown []*BestStatistic      `json:"best_701_darts_thrown,omitempty"`
}

// TournamentOverview stuct for storing tournament overview
type TournamentOverview struct {
	Tournament            *Tournament      `json:"tournament"`
	Group                 *TournamentGroup `json:"tournament_group"`
	PlayerID              int              `json:"player_id"`
	Played                int              `json:"played"`
	MatchesWon            int              `json:"matches_won"`
	MatchesDraw           int              `json:"matches_draw"`
	MatchesLost           int              `json:"matches_lost"`
	LegsFor               int              `json:"legs_for"`
	LegsAgainst           int              `json:"legs_against"`
	LegsDifference        int              `json:"legs_difference"`
	Points                int              `json:"points"`
	PPD                   float32          `json:"ppd"`
	FirstNinePPD          float32          `json:"first_nine_ppd"`
	ThreeDartAvg          float32          `json:"three_dart_avg"`
	FirstNineThreeDartAvg float32          `json:"first_nine_three_dart_avg"`
	CheckoutAttempts      int              `json:"checkout_attempts"`
	CheckoutPercentage    float32          `json:"checkout_percentage"`
	Score60sPlus          int              `json:"scores_60s_plus"`
	Score100sPlus         int              `json:"scores_100s_plus"`
	Score140sPlus         int              `json:"scores_140s_plus"`
	Score180s             int              `json:"scores_180s"`
	Accuracy20            float32          `json:"accuracy_20"`
	Accuracy19            float32          `json:"accuracy_19"`
	AccuracyOverall       float32          `json:"accuracy_overall"`
	IsPromoted            bool             `json:"is_promoted"`
	IsRelegated           bool             `json:"is_relegated"`
	IsWinner              bool             `json:"is_winner"`
}
