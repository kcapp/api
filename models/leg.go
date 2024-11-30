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
	Endtime            null.Time           `json:"end_time"`
	StartingScore      int                 `json:"starting_score"`
	IsFinished         bool                `json:"is_finished"`
	CurrentPlayerID    int                 `json:"current_player_id"`
	WinnerPlayerID     null.Int            `json:"winner_player_id"`
	LegType            *MatchType          `json:"leg_type"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
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
	LegID         int          `json:"leg_id,omitempty"`
	OutshotType   *OutshotType `json:"outshot_type,omitempty"`
	Numbers       []int        `json:"numbers"`
	Hits          map[int]int  `json:"hits"`
	StartingLives null.Int     `json:"starting_lives,omitempty"`
	PointsToWin   null.Int     `json:"points_to_win,omitempty"`
	MaxRounds     null.Int     `json:"max_rounds,omitempty"`
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

	// Get 9 random numbers between the given range
	iteration := 1
	for i := range numbers {
		min := 21 + startingScore + ((startingScore / 4) * (iteration - 1))
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
		StartTime          null.Time           `json:"start_time,omitempty"`
		Endtime            null.Time           `json:"end_time"`
		StartingScore      int                 `json:"starting_score"`
		IsFinished         bool                `json:"is_finished"`
		CurrentPlayerID    int                 `json:"current_player_id"`
		WinnerPlayerID     null.Int            `json:"winner_player_id"`
		LegType            *MatchType          `json:"leg_type"`
		CreatedAt          time.Time           `json:"created_at"`
		UpdatedAt          time.Time           `json:"updated_at"`
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
	if leg.LegType.ID == ONESEVENTY {
		round = len(leg.Visits)/len(leg.Players)/3 + 1
	}

	startTime := null.TimeFrom(leg.CreatedAt)
	endTime := leg.Endtime
	if leg.Visits != nil && len(leg.Visits) > 0 {
		startTime = null.TimeFrom(leg.Visits[0].CreatedAt)
		endTime = null.TimeFrom(leg.GetLastVisit().CreatedAt)
	}

	return json.Marshal(legJSON{
		ID:                 leg.ID,
		StartTime:          startTime,
		Endtime:            endTime,
		StartingScore:      leg.StartingScore,
		IsFinished:         leg.IsFinished,
		CurrentPlayerID:    leg.CurrentPlayerID,
		WinnerPlayerID:     leg.WinnerPlayerID,
		LegType:            leg.LegType,
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

// Player2Leg struct used for storing players in a leg
type Player2Leg struct {
	LegID           int              `json:"leg_id"`
	PlayerID        int              `json:"player_id"`
	PlayerName      string           `json:"player_name"`
	Order           int              `json:"order"`
	CurrentScore    int              `json:"current_score"`
	StartingScore   int              `json:"starting_score"`
	IsCurrentPlayer bool             `json:"is_current_player"`
	Wins            int              `json:"wins,omitempty"`
	VisitStatistics *VisitStatistics `json:"visit_statistics,omitempty"`
	Handicap        null.Int         `json:"handicap,omitempty"`
	Lives           null.Int         `json:"lives,omitempty"`
	Modifiers       *PlayerModifiers `json:"modifiers,omitempty"`
	Player          *Player          `json:"player,omitempty"`
	BotConfig       *BotConfig       `json:"bot_config,omitempty"`
	Hits            HitsMap          `json:"hits"`
	DartsThrown     int              `json:"darts_thrown"`
	IsStopper       null.Bool        `json:"is_stopper,omitempty"`
	IsScorer        null.Bool        `json:"is_scorer,omitempty"`
	CurrentPoints   null.Int         `json:"current_points"`
}

type HitsMap map[int]*Hits

// Contains will check if the map contains all the given values
func (m HitsMap) Contains(modifier int, values ...int) bool {
	for _, v := range values {
		if hits, ok := m[v]; ok {
			count := 0
			if modifier == SINGLE {
				count += hits.Singles
			} else if modifier == DOUBLE {
				count += hits.Doubles
			} else if modifier == TRIPLE {
				count += hits.Triples
			} else {
				count += hits.Total
			}
			if count < 1 {
				// No hits on the given modifier
				return false
			}
		} else {
			// No hits on the number at all
			return false
		}
	}
	return true
}

// Contains will check if the map contains all the given values
func (m HitsMap) GetHits(value int, modifier int) int {
	hits := 0
	if val, ok := m[value]; ok {
		if modifier == SINGLE {
			hits = val.Singles
		} else if modifier == DOUBLE {
			hits = val.Doubles
		} else if modifier == TRIPLE {
			hits = val.Triples
		} else {
			hits = val.Total
		}
	}
	return hits
}

// Add will add the given dart to the HitsMap, inserting it if needed or incrementing the existing value
func (m HitsMap) Add(d *Dart) {
	if _, ok := m[d.ValueRaw()]; !ok {
		m[d.ValueRaw()] = new(Hits)
	}
	hits := m[d.ValueRaw()]
	hits.Add(d)
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
	Score60PlusCounter  int `json:"score_60_plus_counter"`
	Score100PlusCounter int `json:"score_100_plus_counter"`
	Score140PlusCounter int `json:"score_140_plus_counter"`
	Score180Counter     int `json:"score_180_counter"`
}

// PlayerModifiers struct used for storing visit modifiers for a player
type PlayerModifiers struct {
	IsViliusVisit  bool `json:"is_vilius_visit"`
	IsBeerMatch    bool `json:"is_beer_match"`
	IsFishAndChips bool `json:"is_fish_and_chips"`
	IsScore60Plus  bool `json:"is_score_60_plus"`
	IsScore100Plus bool `json:"is_score_100_plus"`
	IsScore140Plus bool `json:"is_score_140_plus"`
	IsScore180     bool `json:"is_score_180"`
}

// AddVisitStatistics adds information about the given visit
func (p2l *Player2Leg) AddVisitStatistics(leg Leg) {
	p2l.VisitStatistics = new(VisitStatistics)
	if p2l.Player.IsBot {
		// Don't add visit statistics for bots
		return
	}
	for _, visit := range leg.Visits {
		if visit.PlayerID == p2l.PlayerID {
			if visit.IsFishAndChips() {
				p2l.VisitStatistics.FishAndChipsCounter++
			} else if visit.IsViliusVisit() {
				p2l.VisitStatistics.ViliusVisitCounter++
			} else if visit.IsScore60Plus() {
				p2l.VisitStatistics.Score60PlusCounter++
			} else if visit.IsScore100Plus() {
				p2l.VisitStatistics.Score100PlusCounter++
			} else if visit.IsScore140Plus() {
				p2l.VisitStatistics.Score140PlusCounter++
			} else if visit.IsScore180() {
				p2l.VisitStatistics.Score180Counter++
			}
			if leg.LegType.ID != ONESEVENTY {
				p2l.DartsThrown = visit.DartsThrown
			}
		}
	}
}

// IsOut will check if the given player is out of the current match
func (player *Player2Leg) IsOut(matchType int, visit Visit) bool {
	if matchType == KNOCKOUT {
		// If player has less than 1 life, and is not the current player
		return player.Lives.Int64 < 1 && player.PlayerID != visit.PlayerID
	}
	// For all other types players are never out
	return false
}

// SetStopper will mark the player as a stopper in SCAM match type
func (p *Player2Leg) SetStopper() {
	p.IsStopper = null.BoolFrom(true)
	p.IsScorer = null.BoolFrom(false)
}

// SetScorer will mark the player as a scorer in SCAM match type
func (p *Player2Leg) SetScorer() {
	p.IsStopper = null.BoolFrom(false)
	p.IsScorer = null.BoolFrom(true)

}

// DecorateVisitsScam will add information about stopper/scorer to each visit
func DecorateVisitsScam(players map[int]*Player2Leg, visits []*Visit) {
	stopperOrder := 1
	for _, player := range players {
		if player.Order == stopperOrder {
			player.SetStopper()
		} else {
			player.SetScorer()
		}
		player.Hits = make(HitsMap)
	}

	hits := make(HitsMap)
	for _, visit := range visits {
		player := players[visit.PlayerID]
		if player.IsStopper.Bool {
			hits.Add(visit.FirstDart)
			hits.Add(visit.SecondDart)
			hits.Add(visit.ThirdDart)
			player.Hits = hits

			visit.IsStopper = null.BoolFrom(true)
			if hits.Contains(SINGLE, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20) {
				stopperOrder++
				for _, player := range players {
					if player.Order == stopperOrder {
						player.SetStopper()
					} else {
						player.SetScorer()
					}
				}
				hits = make(HitsMap)
			}
		}
	}
}

// GetLastVisit returns the last visit of the leg.
func (leg Leg) GetLastVisit() *Visit {
	return leg.Visits[len(leg.Visits)-1]
}

// IsX01 returns true if this leg is a X01 leg
func (leg Leg) IsX01() bool {
	return leg.LegType.ID == X01 || leg.LegType.ID == X01HANDICAP
}

// GetFirstHitDart will return the first (non-Miss) dart for the given player
func (leg Leg) GetFirstHitDart(playerID int) *Dart {
	for _, visit := range leg.Visits {
		if visit.PlayerID != playerID {
			continue
		}
		for _, dart := range visit.GetDarts() {
			if !dart.IsMiss() {
				return &dart
			}
		}
	}
	return nil
}

func (leg Leg) IsLegCheckout() bool {
	if !leg.IsX01() {
		return false
	}
	lastVisit := leg.Visits[len(leg.Visits)-1]
	return lastVisit.IsCheckout
}
