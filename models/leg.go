package models

import (
	"encoding/json"
	"math"
	"math/rand"
	"time"

	"github.com/guregu/null"
)

// Leg struct used for storing legs
type Leg struct {
	ID                 int                 `json:"id"`
	Endtime            null.String         `json:"end_time"`
	StartingScore      int                 `json:"starting_score"`
	IsFinished         bool                `json:"is_finished"`
	CurrentPlayerID    int                 `json:"current_player_id"`
	WinnerPlayerID     null.Int            `json:"winner_player_id"`
	CreatedAt          string              `json:"created_at"`
	UpdatedAt          string              `json:"updated_at"`
	BoardStreamURL     null.String         `json:"board_stream_url,omitempty"`
	MatchID            int                 `json:"match_id"`
	HasScores          bool                `json:"has_scores"`
	Players            []int               `json:"players,omitempty"`
	DartsThrown        int                 `json:"darts_thrown,omitempty"`
	Visits             []*Visit            `json:"visits"`
	Hits               map[int64]*Hits     `json:"hits,omitempty"`
	CheckoutStatistics *CheckoutStatistics `json:"checkout_statistics,omitempty"`
	Statistics         interface{}         `json:"statistics,omitempty"`
	Parameters         *LegParameters      `json:"parameters,omitempty"`
}

// LegParameters struct used for storing leg parameters
type LegParameters struct {
	LegID       int          `json:"leg_id,omitempty"`
	OutshotType *OutshotType `json:"outshot_type,omitempty"`
	Numbers     []int        `json:"numbers"`
	Hits        map[int]int  `json:"hits"`
}

// IsTicTacToeWinner will check if the given player has won a game of Tic Tac Toe
func (params LegParameters) IsTicTacToeWinner(playerID int) bool {
	hits := params.Hits
	numbers := params.Numbers

	for _, combo := range TicTacToeWinningCombos {
		if hits[numbers[combo[0]]] == playerID && hits[numbers[combo[1]]] == playerID && hits[numbers[combo[2]]] == playerID {
			return true
		}
	}
	return false
}

// GenerateTicTacToeNumbers will generate 9 unique numbers for a Tic-Tac-Toe board
func (params *LegParameters) GenerateTicTacToeNumbers(startingScore int) {
	rand.Seed(time.Now().UnixNano())

	bogey := []int{169, 168, 166, 165, 163, 162, 159}
	numbers := make([]int, 9)
	min := 21 + startingScore

	// Get 9 random numbers between the given range
	iteration := 1
	for i := range numbers {
		min = 21 + startingScore + ((startingScore / 4) * (iteration - 1))
		max := min + 10

		valid := true
		for valid {
			num := rand.Intn(max-min) + min
			// Make sure we don't select duplicates, and don't select bogey numbers
			if !containsInt(numbers, num) && !containsInt(bogey, num) {
				numbers[i] = num
				valid = false
				if i%3 == 0 {
					iteration++
				}
			}
		}
	}
	// Shuffle the numbers
	for i := range numbers {
		j := rand.Intn(i + 1)
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}

	// Make sure the middle number is the largest number
	max := 0
	idx := 0
	for i, e := range numbers {
		if i == 0 || e > max {
			max = e
			idx = i
		}
	}
	if idx != 4 {
		numbers[4], numbers[idx] = numbers[idx], numbers[4]
	}

	// The middle number should be more difficult, so we make sure it's odd, and increase it's value
	iteration = 0
	valid := true
	for valid {
		newMiddle := numbers[4] + iteration + 10 + rand.Intn(5)
		if newMiddle%2 == 0 {
			newMiddle++
		}
		if !containsInt(numbers, newMiddle) && !containsInt(bogey, newMiddle) {
			if newMiddle > 170 {
				// Numbers got too big, so reset counter
				iteration -= 10
			} else {
				numbers[4] = newMiddle
				break
			}
		}
		iteration++
	}
	params.Numbers = numbers
}

// IsTicTacToeDraw will check if the given parameters indicate a draw
func (params LegParameters) IsTicTacToeDraw() bool {
	hits := params.Hits
	numbers := params.Numbers

	draw := true
	for _, combo := range TicTacToeWinningCombos {
		num1 := numbers[combo[0]]
		num2 := numbers[combo[1]]
		num3 := numbers[combo[2]]

		// Check if keys exists
		_, exists1 := hits[num1]
		_, exists2 := hits[num2]
		_, exists3 := hits[num3]

		if (exists1 && exists2 && hits[num1] != hits[num2]) ||
			(exists1 && exists3 && hits[num1] != hits[num3]) ||
			(exists2 && exists3 && hits[num2] != hits[num3]) {
			// Two numbers are taken by two different players, which means this combo cannot be completed
			continue
		}
		draw = false
	}
	return draw
}

// MarshalJSON will marshall the given object to JSON
func (leg Leg) MarshalJSON() ([]byte, error) {
	// Use a type to get consistent order of JSON key-value pairs.
	type legJSON struct {
		ID                 int                 `json:"id"`
		Endtime            null.String         `json:"end_time"`
		StartingScore      int                 `json:"starting_score"`
		IsFinished         bool                `json:"is_finished"`
		CurrentPlayerID    int                 `json:"current_player_id"`
		WinnerPlayerID     null.Int            `json:"winner_player_id"`
		CreatedAt          string              `json:"created_at"`
		UpdatedAt          string              `json:"updated_at"`
		BoardStreamURL     null.String         `json:"board_stream_url,omitempty"`
		MatchID            int                 `json:"match_id"`
		HasScores          bool                `json:"has_scores"`
		Round              int                 `json:"round"`
		Players            []int               `json:"players,omitempty"`
		DartsThrown        int                 `json:"darts_thrown,omitempty"`
		Visits             []*Visit            `json:"visits"`
		Hits               map[int64]*Hits     `json:"hits,omitempty"`
		CheckoutStatistics *CheckoutStatistics `json:"checkout_statistics,omitempty"`
		Statistics         interface{}         `json:"statistics,omitempty"`
		Parameters         *LegParameters      `json:"parameters,omitempty"`
	}
	round := int(math.Floor(float64(len(leg.Visits))/float64(len(leg.Players))) + 1)

	return json.Marshal(legJSON{
		ID:                 leg.ID,
		Endtime:            leg.Endtime,
		StartingScore:      leg.StartingScore,
		IsFinished:         leg.IsFinished,
		CurrentPlayerID:    leg.CurrentPlayerID,
		WinnerPlayerID:     leg.WinnerPlayerID,
		CreatedAt:          leg.CreatedAt,
		UpdatedAt:          leg.UpdatedAt,
		BoardStreamURL:     leg.BoardStreamURL,
		MatchID:            leg.MatchID,
		HasScores:          leg.HasScores,
		Round:              round,
		Players:            leg.Players,
		DartsThrown:        leg.DartsThrown,
		Visits:             leg.Visits,
		Hits:               leg.Hits,
		CheckoutStatistics: leg.CheckoutStatistics,
		Statistics:         leg.Statistics,
		Parameters:         leg.Parameters,
	})
}

// Player2Leg struct used for stroring players in a leg
type Player2Leg struct {
	LegID           int              `json:"leg_id"`
	PlayerID        int              `json:"player_id"`
	PlayerName      string           `json:"player_name"`
	Order           int              `json:"order"`
	CurrentScore    int              `json:"current_score"`
	StartingScore   int              `json:"-"`
	IsCurrentPlayer bool             `json:"is_current_player"`
	Wins            int              `json:"wins,omitempty"`
	VisitStatistics *VisitStatistics `json:"visit_statistics,omitempty"`
	Handicap        null.Int         `json:"handicap,omitempty"`
	Modifiers       *PlayerModifiers `json:"modifiers,omitempty"`
	Player          *Player          `json:"player,omitempty"`
	BotConfig       *BotConfig       `json:"bot_config,omitempty"`
	Hits            map[int]*Hits    `json:"hits"`
	DartsThrown     int              `json:"darts_thrown,omitempty"`
}

// BotConfig struct used for storing bot configuration
type BotConfig struct {
	PlayerID null.Int `json:"player_id"`
	Skill    null.Int `json:"skill_level"`
}

// VisitStatistics tells about the
type VisitStatistics struct {
	FishAndChipsCounter int `json:"fish_and_chips_counter"`
	ViliusVisitCounter  int `json:"vilius_visit_counter"`
}

// PlayerModifiers struct used for storing visit modifiers for a player
type PlayerModifiers struct {
	IsViliusVisit  bool `json:"is_vilius_visit"`
	IsBeerMatch    bool `json:"is_beer_match"`
	IsFishAndChips bool `json:"is_fish_and_chips"`
}

// AddVisitStatistics adds information about
func (p2l *Player2Leg) AddVisitStatistics(leg Leg) {
	p2l.VisitStatistics = new(VisitStatistics)
	for _, visit := range leg.Visits {
		if visit.PlayerID == p2l.PlayerID {
			if visit.IsFishAndChips() {
				p2l.VisitStatistics.FishAndChipsCounter++
			}
			if visit.IsViliusVisit() {
				p2l.VisitStatistics.ViliusVisitCounter++
			}
			p2l.DartsThrown = visit.DartsThrown
		}
	}
}
