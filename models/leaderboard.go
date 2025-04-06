package models

import "github.com/guregu/null"

// MatchTypeLeaderboard struct to hold leaderboard data for a specific match type
type MatchTypeLeaderboard struct {
	MatchTypeID int      `json:"match_type_id"`
	PlayerID    int      `json:"player_id"`
	LegID       int      `json:"leg_id"`
	DartsThrown null.Int `json:"darts_thrown"`
	Score       int      `json:"score"`
}
