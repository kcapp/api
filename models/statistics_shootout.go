package models

// StatisticsShootout struct used for storing statistics for shootout matches
type StatisticsShootout struct {
	ID            int             `json:"id,omitempty"`
	MatchID       int             `json:"match_id,omitempty"`
	PlayerID      int             `json:"player_id,omitempty"`
	PPD           float32         `json:"ppd"`
	Score60sPlus  int             `json:"scores_60s_plus"`
	Score100sPlus int             `json:"scores_100s_plus"`
	Score140sPlus int             `json:"scores_140s_plus"`
	Score180s     int             `json:"scores_180s"`
	Hits          map[int64]*Hits `json:"hits,omitempty"`
}
