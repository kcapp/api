package models

import (
	"github.com/guregu/null"
)

// Owe struct used for storing owes
type Owe struct {
	PlayerOwerID int      `json:"player_ower_id"`
	PlayerOweeID int      `json:"player_owee_id"`
	OweType      *OweType `json:"owe_type"`
	Amount       int      `json:"amount"`
}

// OweType struct used for storing owe types
type OweType struct {
	ID   null.Int    `json:"id"`
	Item null.String `json:"item"`
}
