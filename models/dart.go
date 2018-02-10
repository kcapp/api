package models

import (
	"errors"

	"github.com/guregu/null"
)

// Dart struct used for storing darts
type Dart struct {
	Value             null.Int `json:"value"`
	Multiplier        int64    `json:"multiplier"`
	IsCheckoutAttempt bool     `json:"is_checkout"`
	IsBust            bool     `json:"is_bust,omitempty"`
}

// SetModifiers will set IsBust and IsCheckoutAttempt modifiers for the given dartø
func (dart Dart) SetModifiers(currentScore int) {
	scoreAfterThrow := currentScore - dart.GetScore()
	if scoreAfterThrow == 0 && dart.Multiplier == 2 {
		// Check if this throw was a checkout
		dart.IsBust = false
		dart.IsCheckoutAttempt = true
	} else if scoreAfterThrow < 2 {
		dart.IsBust = true
		dart.IsCheckoutAttempt = false
	} else {
		dart.IsBust = false
		if currentScore == 50 || (currentScore <= 40 && currentScore%2 == 0) {
			dart.IsCheckoutAttempt = true
		} else {
			dart.IsCheckoutAttempt = false
		}
	}
}

// ValidateInput will verify that the dart contains valid values
func (dart Dart) ValidateInput() error {
	if dart.Value.Int64 < 0 {
		return errors.New("Value cannot be less than 0")
	} else if dart.Value.Int64 > 25 || (dart.Value.Int64 > 20 && dart.Value.Int64 < 25) {
		return errors.New("Value has to be 20 or less (or 25 (bull))")
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

// IsCheckout checks if this dart was a checkout attempt
func (dart Dart) IsCheckout(currentScore int) bool {
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
