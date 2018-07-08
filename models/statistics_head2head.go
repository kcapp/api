package models

// StatisticsHead2Head struct used for storing head to head statistics
type StatisticsHead2Head struct {
	LastMatches         map[int][]*Match              `json:"last_matches"`
	Head2HeadMatches    []*Match                      `json:"head_to_head_matches"`
	Head2HeadWins       map[int64]int                 `json:"head_to_head_wins"`
	Player301Statistics map[int]*StatisticsX01        `json:"player_301_statistics"`
	Player501Statistics map[int]*StatisticsX01        `json:"player_501_statistics"`
	PlayerVisits        map[int][]*Visit              `json:"player_visits"`
	PlayerCheckouts     map[int][]*CheckoutStatistics `json:"player_checkouts"`
	PlayerElos          map[int]*PlayerElo            `json:"player_elo"`
}
