package models

import "github.com/guregu/null"

// Venue struct used for storing venues
type Venue struct {
	ID          null.Int    `json:"id"`
	Name        null.String `json:"name"`
	OfficeID    null.Int    `json:"office_id"`
	Description null.String `json:"description"`
}
