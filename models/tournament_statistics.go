package models

import "github.com/guregu/null"

// TournamentStatistics struct for storing tournament statistics
type TournamentStatistics struct {
	GeneralStatistics  *TournamentGeneralStatistics `json:"general_statistics"`
	HighestCheckout    []*BestStatistic             `json:"checkout_highest,omitempty"`
	BestThreeDartAvg   []*BestStatisticFloat        `json:"best_three_dart_avg,omitempty"`
	BestFirstNineAvg   []*BestStatisticFloat        `json:"best_first_nine_avg,omitempty"`
	Best301DartsThrown []*BestStatistic             `json:"best_301_darts_thrown,omitempty"`
	Best501DartsThrown []*BestStatistic             `json:"best_501_darts_thrown,omitempty"`
	Best701DartsThrown []*BestStatistic             `json:"best_701_darts_thrown,omitempty"`
}

// TournamentGeneralStatistics struct for storing tournament statistics
type TournamentGeneralStatistics struct {
	Score60sPlus        int `json:"scores_60s_plus"`
	Score100sPlus       int `json:"scores_100s_plus"`
	Score140sPlus       int `json:"scores_140s_plus"`
	Score180s           int `json:"scores_180s"`
	ScoreFishNChips     int `json:"scores_fish_n_chips"`
	ScoreBullseye       int `json:"scores_bullseye"`
	ScoreDoubleBullseye int `json:"scores_double_bullseye"`
	D1Checkouts         int `json:"checkout_d1"`
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
	ManualOrder           null.Int         `json:"manual_order"`
}
