package models

// TournamentStatistics stuct for storing tournament  statistics
type TournamentStatistics struct {
	Tournament         *Tournament      `json:"tournament"`
	Group              *TournamentGroup `json:"tournament_group"`
	PlayerID           int              `json:"player_id"`
	Played             int              `json:"played"`
	MatchesWon         int              `json:"matches_won"`
	MatchesDraw        int              `json:"matches_draw"`
	MatchesLost        int              `json:"matches_lost"`
	LegsFor            int              `json:"legs_for"`
	LegsAgainst        int              `json:"legs_against"`
	LegsDifference     int              `json:"legs_difference"`
	Points             int              `json:"points"`
	PPD                float32          `json:"ppd"`
	FirstNinePPD       float32          `json:"first_nine_ppd"`
	CheckoutPercentage float32          `json:"checkout_percentage"`
	Score60sPlus       int              `json:"scores_60s_plus"`
	Score100sPlus      int              `json:"scores_100s_plus"`
	Score140sPlus      int              `json:"scores_140s_plus"`
	Score180s          int              `json:"scores_180s"`
	Accuracy20         float32          `json:"accuracy_20"`
	Accuracy19         float32          `json:"accuracy_19"`
	AccuracyOverall    float32          `json:"accuracy_overall"`
}
