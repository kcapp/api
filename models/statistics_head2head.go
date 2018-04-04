package models

// StatisticsHead2Head struct used for storing head to head statistics
type StatisticsHead2Head struct {
	LastGames           map[int][]*Game               `json:"last_games"`
	Head2HeadGames      []*Game                       `json:"head_to_head_games"`
	Head2HeadWins       map[int64]int                 `json:"head_to_head_wins"`
	Player301Statistics map[int]*StatisticsX01        `json:"player_301_statistics"`
	Player501Statistics map[int]*StatisticsX01        `json:"player_501_statistics"`
	PlayerVisits        map[int][]*Visit              `json:"player_visits"`
	PlayerCheckouts     map[int][]*CheckoutStatistics `json:"player_checkouts"`
}
