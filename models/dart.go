package models

import (
	"errors"
	"fmt"

	"github.com/guregu/null"
)

const (
	// SINGLE const representing single
	SINGLE = 1
	// DOUBLE const representing double
	DOUBLE = 2
	// TRIPLE const representing triple
	TRIPLE = 3
)

const BULLSEYE = 25

var (
	// CRICKETDARTS var holding darts aimed at in a game of Cricket
	CRICKETDARTS = []int{15, 16, 17, 18, 19, 20, 25}
)

// Dart struct used for storing darts
type Dart struct {
	Value      null.Int `json:"value"`
	Multiplier int64    `json:"multiplier"`
}

// IsBust will check if the given dart is a bust
func (dart *Dart) IsBust(currentScore int, outshotTypeId int) bool {
	scoreAfterThrow := currentScore - dart.GetScore()
	if scoreAfterThrow == 0 {
		if outshotTypeId == OUTSHOTANY ||
			(outshotTypeId == OUTSHOTDOUBLE && dart.IsDouble()) ||
			(outshotTypeId == OUTSHOTMASTER && (dart.IsDouble() || dart.IsTriple())) {
			return false
		}
	}
	if outshotTypeId == OUTSHOTANY {
		if scoreAfterThrow < 1 {
			return true
		}
	} else {
		if scoreAfterThrow < 2 {
			return true
		}
	}

	// If the throw is not a bust, make sure the dart is valid
	if !dart.Value.Valid {
		dart.Value = null.IntFrom(0)
	}
	return false
}

// IsBustAbove will check if the given dart puts score above the given target
func (dart *Dart) IsBustAbove(currentScore int, targetScore int) bool {
	scoreAfterThrow := currentScore + dart.GetScore()
	if scoreAfterThrow > targetScore {
		return true
	}

	// If the throw is not a bust, make sure the dart is valid
	if !dart.Value.Valid {
		dart.Value = null.IntFrom(0)
	}
	return false
}

// ValidateInput will verify that the dart contains valid values
func (dart *Dart) ValidateInput() error {
	if dart.Value.Int64 < 0 {
		return errors.New("value cannot be less than 0")
	} else if dart.Value.Int64 > 25 || (dart.Value.Int64 > 20 && dart.Value.Int64 < 25) {
		return errors.New("value has to be less than 21 (or 25 (bull))")
	} else if dart.Multiplier > 3 || dart.Multiplier < 1 {
		return errors.New("multiplier has to be one of 1 (single), 2 (douhle), 3 (triple)")
	}

	// Make sure multiplier is 1 on miss
	if !dart.Value.Valid || dart.Value.Int64 == 0 {
		dart.Multiplier = 1
	}
	return nil
}

// GetScore will get the actual score of the dart (value * multiplier)
func (dart Dart) GetScore() int {
	return int(dart.Value.Int64 * dart.Multiplier)
}

// GetBermudaTriangleScore will get the Bermuda Triangle score for the given dart on target
func (dart Dart) GetBermudaTriangleScore(target Target) int {
	if (target.Value == -1 || target.Value == dart.ValueRaw()) && contains(target.multipliers, dart.Multiplier) {
		if target.score > 0 {
			return target.score
		}
		return dart.GetScore()
	}
	return 0
}

// Get420Score will get the 420 score for the given dart on target
func (dart Dart) Get420Score(target Target) int {
	if dart.Multiplier == 2 && dart.ValueRaw() == target.Value {
		return dart.GetScore()
	}
	return 0
}

// GetJDCPracticeScore will get the JDC Practice score for the given dart on target
func (dart Dart) GetJDCPracticeScore(target Target) int {
	if target.Value == dart.ValueRaw() && contains(target.multipliers, dart.Multiplier) {
		return dart.GetScore()
	}
	return 0
}

// IsCheckoutAttempt checks if this dart was a checkout attempt
func (dart Dart) IsCheckoutAttempt(currentScore int, dartNum int, outshotTypeId int) bool {
	if !dart.Value.Valid {
		// Dart was not actually thrown, player busted/checked out already
		return false
	}

	if outshotTypeId == OUTSHOTANY {
		if ((currentScore <= 20 || currentScore == 25) ||
			((currentScore == 50 || currentScore <= 40) && currentScore%2 == 0) ||
			(currentScore <= 60 && currentScore%3 == 0)) && currentScore > 0 {
			return true
		}
	} else if outshotTypeId == OUTSHOTMASTER {
		if currentScore <= 60 && currentScore%3 == 0 && currentScore > 1 {
			return true
		}
	}
	if currentScore-dart.GetScore() == 0 && dart.IsDouble() {
		// Actual checkout
		return true
	} else if (dartNum == 3 && currentScore == 50) || (currentScore <= 40 && currentScore%2 == 0 && currentScore > 1) {
		// Checkout attempt (bull only counts if it was on the third dart)
		return true
	}
	return false
}

// GetString will return a string representing the dart of the fomat "<multiplier>-<value>"
func (dart Dart) GetString() string {
	if dart.Value.Valid {
		return fmt.Sprintf("%d-%d", dart.Multiplier, dart.Value.Int64)
	}
	return fmt.Sprintf("%d-NULL", dart.Multiplier)
}

func (dart Dart) String() string {
	if !dart.Value.Valid {
		return ""
	}
	if dart.Multiplier == TRIPLE {
		return fmt.Sprintf("T%d", dart.ValueRaw())
	} else if dart.Multiplier == DOUBLE {
		return fmt.Sprintf("D%d", dart.ValueRaw())
	}
	return fmt.Sprintf("%d", dart.ValueRaw())
}

// NewDart will return a new dart with the given settings
func NewDart(value null.Int, multipler int64) *Dart {
	return &Dart{Value: value, Multiplier: multipler}
}

// IsSingle will check if this dart multipler was a single
func (dart Dart) IsSingle() bool {
	return dart.Multiplier == SINGLE
}

// IsDouble will check if this dart multipler was a double
func (dart Dart) IsDouble() bool {
	return dart.Multiplier == DOUBLE
}

// IsTriple will check if this dart multipler was a triple
func (dart Dart) IsTriple() bool {
	return dart.Multiplier == TRIPLE
}

// IsBull will check if this dart was a bullseye (single or double)
func (dart Dart) IsBull() bool {
	return dart.ValueRaw() == 25
}

// IsMiss will check if this dart was a miss
func (dart Dart) IsMiss() bool {
	return dart.ValueRaw() == 0
}

// IsCricketMiss will check if this dart was a miss on cricket numbers
func (dart Dart) IsCricketMiss() bool {
	for _, num := range CRICKETDARTS {
		if dart.ValueRaw() == num {
			return false
		}
	}
	return true
}

// ValueRaw will return the value of the dart or 0 if invalid
func (dart Dart) ValueRaw() int {
	if dart.Value.Valid {
		return int(dart.Value.Int64)
	}
	return 0
}

func (dart Dart) IsValue(values []int) bool {
	for _, value := range values {
		if dart.ValueRaw() == value {
			return true
		}
	}
	return false
}

// GetMarksHit will return the number of marks hit by the given darts, accounting for numbers requiring less than 3 hits to close
// If the number is still open by other players hits = multiplier, otherwise hits = multiplier - prev_hits
func (dart *Dart) GetMarksHit(hits map[int]int64, open bool) int64 {
	marks := int64(0)

	val := dart.ValueRaw()
	multiplier := dart.Multiplier
	if _, ok := hits[val]; !ok {
		hits[val] = multiplier
		marks += multiplier
	} else {
		if !open && hits[val]+multiplier > 3 {
			marks += multiplier - hits[val]
		} else {
			marks += multiplier
		}
		hits[val] += multiplier
	}
	return marks
}

// CalculateCricketScore will calculate the score for each player for the given dart
func (dart *Dart) CalculateCricketScore(playerID int, scores map[int]*Player2Leg) int {
	if !dart.Value.Valid {
		return 0
	}
	if !dart.IsHit(CRICKETDARTS) {
		return 0
	}

	score := int(dart.Value.Int64)
	hitsMap := scores[playerID].Hits
	if _, ok := hitsMap[score]; !ok {
		hitsMap[score] = new(Hits)
	}
	hits := hitsMap[score].Total
	hitsMap[score].Total += int(dart.Multiplier)
	multiplier := hitsMap[score].Total - hits
	if hits < 3 {
		multiplier = hitsMap[score].Total - 3
	}
	points := int(dart.Value.Int64) * multiplier

	pointsGiven := false
	if hitsMap[score].Total > 3 {
		for id, p2l := range scores {
			if id == playerID {
				continue
			}
			if val, ok := p2l.Hits[score]; ok {
				if val.Total < 3 {
					p2l.CurrentScore += points
					pointsGiven = true
				}
			} else {
				p2l.CurrentScore += points
				pointsGiven = true
			}
		}
	}
	if points < 0 || !pointsGiven {
		points = 0
	}
	return points
}

// IsHit will return true if the given dart hit one of the supplied targets
func (dart *Dart) IsHit(targets []int) bool {
	score := int(dart.Value.Int64)
	return containsInt(targets, score)
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func removeInt(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
