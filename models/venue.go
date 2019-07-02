package models

import "github.com/guregu/null"

// Venue struct used for storing venues
type Venue struct {
	ID          null.Int     `json:"id"`
	Name        null.String  `json:"name"`
	OfficeID    null.Int     `json:"office_id"`
	Description null.String  `json:"description"`
	Config      *VenueConfig `json:"config,omitempty"`
}

// VenueConfig struct used for storing venue configuration
type VenueConfig struct {
	VenueID                int         `json:"id"`
	HasDualMonitor         bool        `json:"has_dual_monitor"`
	HasLEDLights           bool        `json:"has_led_lights"`
	HasSmartboard          bool        `json:"has_smartboard"`
	SmartboardUUID         null.String `json:"smartboard_uuid,omitempty"`
	SmartboardButtonNumber null.Int    `json:"smartboard_button_number,omitempty"`
}
