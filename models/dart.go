package models

import (
	"errors"

	"github.com/guregu/null"
)

// Dart struct used for storing darts
type Dart struct {
	Value      null.Int `json:"value"`
	Multiplier int64    `json:"multiplier"`
}

// IsBust will check if the given dart is a bust
func (dart *Dart) IsBust(currentScore int) bool {
	scoreAfterThrow := currentScore - dart.GetScore()
	if scoreAfterThrow == 0 && dart.Multiplier == 2 {
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
func (dart Dart) IsCheckoutAttempt(currentScore int) bool {
	if !dart.Value.Valid {
		// Dart was not actually thrown, player busted/checked out already
		return false
	}
	if currentScore-dart.GetScore() == 0 && dart.Multiplier == 2 {
		// Actual checkout
		return true
	} else if currentScore == 50 || (currentScore <= 40 && currentScore%2 == 0) {
		// Checkout attempt
		return true
	}
	return false
}
