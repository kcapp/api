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

// Dart struct used for storing darts
type Dart struct {
	Value      null.Int `json:"value"`
	Multiplier int64    `json:"multiplier"`
}

// IsBust will check if the given dart is a bust
func (dart *Dart) IsBust(currentScore int) bool {
	scoreAfterThrow := currentScore - dart.GetScore()
	if scoreAfterThrow == 0 && dart.IsDouble() {
		return false
	} else if scoreAfterThrow < 2 {
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
		return errors.New("Value cannot be less than 0")
	} else if dart.Value.Int64 > 25 || (dart.Value.Int64 > 20 && dart.Value.Int64 < 25) {
		return errors.New("Value has to be less than 21 (or 25 (bull))")
	} else if dart.Multiplier > 3 || dart.Multiplier < 1 {
		return errors.New("Multiplier has to be one of 1 (single), 2 (douhle), 3 (triple)")
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

// IsCheckoutAttempt checks if this dart was a checkout attempt
func (dart Dart) IsCheckoutAttempt(currentScore int, num int) bool {
	if !dart.Value.Valid {
		// Dart was not actually thrown, player busted/checked out already
		return false
	}
	if currentScore-dart.GetScore() == 0 && dart.IsDouble() {
		// Actual checkout
		return true
	} else if (num == 3 && currentScore == 50) || (currentScore <= 40 && currentScore%2 == 0 && currentScore > 1) {
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

// ValueRaw will return the value of the dart or 0 if invalid
func (dart Dart) ValueRaw() int {
	if dart.Value.Valid {
		return int(dart.Value.Int64)
	}
	return 0
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
