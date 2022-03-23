package models

import "github.com/guregu/null"

// Probability struct used for storing matches
type Probability struct {
	ID                         int             `json:"id"`
	CreatedAt                  string          `json:"created_at"`
	UpdatedAt                  string          `json:"updated_at"`
	IsFinished                 bool            `json:"is_finished"`
	IsAbandoned                bool            `json:"is_abandoned"`
	IsWalkover                 bool            `json:"is_walkover"`
	WinnerID                   null.Int        `json:"winner_id"`
	Players                    []int           `json:"players"`
	Elos                       map[int]int     `json:"player_elo"`
	PlayerWinningProbabilities map[int]float64 `json:"player_winning_probabilities"`
	PlayerOdds                 map[int]float32 `json:"player_odds"`
}
