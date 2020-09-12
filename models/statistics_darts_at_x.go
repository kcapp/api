package models

import "github.com/guregu/null"

// StatisticsDartsAtX struct used for storing statistics for cricket
type StatisticsDartsAtX struct {
	ID            int             `json:"id,omitempty"`
	LegID         int             `json:"leg_id,omitempty"`
	PlayerID      int             `json:"player_id,omitempty"`
	MatchesPlayed int             `json:"matches_played,omitempty"`
	MatchesWon    int             `json:"matches_won,omitempty"`
	LegsPlayed    int             `json:"legs_played,omitempty"`
	LegsWon       int             `json:"legs_won,omitempty"`
	OfficeID      null.Int        `json:"office_id,omitempty"`
	AvgScore      int             `json:"avg_score"`
	Score         null.Int        `json:"score,omitempty"`
	Singles       int             `json:"singles"`
	Doubles       int             `json:"doubles"`
	Triples       int             `json:"triples"`
	HitRate       float32         `json:"hit_rate"`
	Hits5         int             `json:"hits_5"`
	Hits6         int             `json:"hits_6"`
	Hits7         int             `json:"hits_7"`
	Hits8         int             `json:"hits_8"`
	Hits9         int             `json:"hits_9"`
	Hitrates      map[int]float32 `json:"hitrates,omitempty"`
}
