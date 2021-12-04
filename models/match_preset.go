package models

import "github.com/guregu/null"

// MatchPreset struct used for storing a match preset
type MatchPreset struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	MatchType     *MatchType  `json:"match_type"`
	MatchMode     *MatchMode  `json:"match_mode"`
	StartingScore null.Int    `json:"starting_score"`
	SmartcardUID  null.String `json:"smartcard_uid,omitempty"`
	Description   null.String `json:"description,omitempty"`
}
