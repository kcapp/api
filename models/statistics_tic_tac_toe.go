package models

import "github.com/guregu/null"

var (
	// TicTacToeWinningCombos contains all winning combincations for tic tac toe
	TicTacToeWinningCombos = [8][3]int{
		// Horizontally
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},

		// Diagonally
		{0, 4, 8},
		{2, 4, 6},

		// Vertically
		{0, 3, 6},
		{1, 4, 7},
		{2, 5, 8},
	}
)

// StatisticsTicTacToe struct used for storing statistics for tic tac toe
type StatisticsTicTacToe struct {
	ID            int      `json:"id,omitempty"`
	LegID         int      `json:"leg_id,omitempty"`
	PlayerID      int      `json:"player_id,omitempty"`
	MatchesPlayed int      `json:"matches_played"`
	MatchesWon    int      `json:"matches_won"`
	LegsPlayed    int      `json:"legs_played"`
	LegsWon       int      `json:"legs_won"`
	OfficeID      null.Int `json:"office_id,omitempty"`
	DartsThrown   int      `json:"darts_thrown"`
	Score         int      `json:"score" `
	NumbersClosed int      `json:"numbers_closed"`
	HighestClosed int      `json:"highest_closed"`
}
