package models

import "github.com/guregu/null"

// DefaultOptions struct used for storing default options
type DefaultOptions struct {
	MatchType                    *MatchType   `json:"match_type"`
	MatchMode                    *MatchMode   `json:"match_mode"`
	StartingScore                int          `json:"starting_score"`
	MaxRounds                    null.Int     `json:"max_rounds,omitempty"`
	OutshotType                  *OutshotType `json:"outshot_type,omitempty"`
	LeaderboardLastLegsCount     int          `json:"leaderboard_last_legs_count"`
	LeaderboardActivePeriodWeeks int          `json:"leaderboard_active_period_weeks"`
}
