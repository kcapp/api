package models

import (
	"time"

	"github.com/guregu/null"
)

// Badge represents a badge model.
type Badge struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Filename    string `json:"filename"`
}

// PlayerBadge represents a Player2Badge model.
type PlayerBadge struct {
	Badge     *Badge    `json:"badge"`
	PlayerID  int       `json:"player_id"`
	LegID     null.Int  `json:"leg_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
