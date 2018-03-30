package models

// StatisticsHead2Head struct used for storing head to head statistics
type StatisticsHead2Head struct {
	LastGames           map[int][]*Game               `json:"last_games"`
	LastGamesHeadToHead []*Game                       `json:"last_games_head_to_head"`
	PlayerStatistics    map[int]*StatisticsX01        `json:"player_statistics"`
	PlayerVisits        map[int][]*Visit              `json:"player_visits"`
	PlayerCheckouts     map[int][]*CheckoutStatistics `json:"player_checkouts"`
}
