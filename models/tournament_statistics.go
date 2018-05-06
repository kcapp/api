package models

// TournamentStatistics stuct for storing tournament  statistics
type TournamentStatistics struct {
	Tournament     *Tournament      `json:"tournament"`
	Group          *TournamentGroup `json:"tournament_group"`
	PlayerID       int              `json:"player_id"`
	Played         int              `json:"played"`
	MatchesWon     int              `json:"matches_won"`
	MatchesDraw    int              `json:"matches_draw"`
	MatchesLost    int              `json:"matches_lost"`
	LegsFor        int              `json:"legs_for"`
	LegsAgainst    int              `json:"legs_against"`
	LegsDifference int              `json:"legs_difference"`
	Points         int              `json:"points"`
}
