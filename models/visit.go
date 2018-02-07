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

// Visit struct used for storing matches
type Visit struct {
	ID         int    `json:"id"`
	MatchID    int    `json:"match_id"`
	PlayerID   int    `json:"player_id"`
	FirstDart  *Dart  `json:"first_dart"`
	SecondDart *Dart  `json:"second_dart"`
	ThirdDart  *Dart  `json:"third_dart"`
	IsBust     bool   `json:"is_bust"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ValidateInput will verify the input does not containg any errors
func (visit Visit) ValidateInput() error {
	err := visit.FirstDart.ValidateInput()
	if err != nil {
		return err
	}
	err = visit.SecondDart.ValidateInput()
	if err != nil {
		return err
	}
	err = visit.ThirdDart.ValidateInput()
	if err != nil {
		return err
	}
	return nil
}

// SetModifiers will set IsBust and IsCheckoutAttempt modifiers for the given dartø
func (dart Dart) SetModifiers(currentScore int) {
	scoreAfterThrow := currentScore - dart.getScore()
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

// getScore will get the actual score of the dart (value * multiplier)
func (dart Dart) getScore() int {
	return int(dart.Value.Int64 * dart.Multiplier)
}

// SetVisitModifiers will set IsBust and IsCheckoutAttempt on all darts
func (visit Visit) SetVisitModifiers(currentScore int) {
	visit.FirstDart.SetModifiers(currentScore)
	currentScore = currentScore - visit.FirstDart.getScore()
	if currentScore < 2 {
		// If first dart is bust/checkout, then second/third dart doesn't count
		visit.SecondDart.Value.Valid = false
		visit.SecondDart.Multiplier = 1
		visit.ThirdDart.Value.Valid = false
		visit.ThirdDart.Multiplier = 1
	} else {
		visit.SecondDart.SetModifiers(currentScore)
		currentScore = currentScore - visit.SecondDart.getScore()
		if currentScore < 2 {
			// If second dart is bust, then third dart doesn't count
			visit.ThirdDart.Value.Valid = false
			visit.ThirdDart.Multiplier = 1
		} else {
			visit.ThirdDart.SetModifiers(currentScore)
			currentScore = currentScore - visit.ThirdDart.getScore()
		}
	}
}

// GetHitsMap will return a map where key is dart and value is count of single,double,triple hits
func GetHitsMap(visits []*Visit) (map[int64]*Hits, int) {
	hitsMap := make(map[int64]*Hits)
	// Populate the map with hits for each value (miss, 1-20, bull)
	for i := 0; i <= 20; i++ {
		hitsMap[int64(i)] = new(Hits)
	}
	hitsMap[25] = new(Hits)

	var dartsThrown int
	for _, visit := range visits {
		if visit.FirstDart.Value.Valid {
			hit := hitsMap[visit.FirstDart.Value.Int64]
			if visit.FirstDart.Multiplier == 1 {
				hit.Singles++
			}
			if visit.FirstDart.Multiplier == 2 {
				hit.Doubles++
			}
			if visit.FirstDart.Multiplier == 3 {
				hit.Triples++
			}
			dartsThrown++
		}
		if visit.SecondDart.Value.Valid {
			hit := hitsMap[visit.SecondDart.Value.Int64]
			if visit.SecondDart.Multiplier == 1 {
				hit.Singles++
			}
			if visit.SecondDart.Multiplier == 2 {
				hit.Doubles++
			}
			if visit.SecondDart.Multiplier == 3 {
				hit.Triples++
			}
			dartsThrown++
		}
		if visit.ThirdDart.Value.Valid {
			hit := hitsMap[visit.ThirdDart.Value.Int64]
			if visit.ThirdDart.Multiplier == 1 {
				hit.Singles++
			}
			if visit.ThirdDart.Multiplier == 2 {
				hit.Doubles++
			}
			if visit.ThirdDart.Multiplier == 3 {
				hit.Triples++
			}
			dartsThrown++
		}
	}
	return hitsMap, dartsThrown
}
