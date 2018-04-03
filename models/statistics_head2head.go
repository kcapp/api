package models

// StatisticsHead2Head struct used for storing head to head statistics
type StatisticsHead2Head struct {
	LastGames           map[int][]*Game               `json:"last_games"`
	LastGamesHeadToHead []*Game                       `json:"last_games_head_to_head"`
	Player301Statistics map[int]*StatisticsX01        `json:"player_301_statistics"`
	Player501Statistics map[int]*StatisticsX01        `json:"player_501_statistics"`
	PlayerVisits        map[int][]*Visit              `json:"player_visits"`
	PlayerCheckouts     map[int][]*CheckoutStatistics `json:"player_checkouts"`
}
