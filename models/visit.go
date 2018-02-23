package models

import (
	"errors"

	"github.com/guregu/null"
)

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
	if visit.FirstDart == nil {
		return errors.New("First dart cannot be null")
	}
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

// SetIsBust will set IsBust for the given visit
func (visit *Visit) SetIsBust(currentScore int) {
	isBust := false
	isBust = visit.FirstDart.IsBust(currentScore)
	currentScore = currentScore - visit.FirstDart.GetScore()
	if !isBust && currentScore > 0 {
		isBust = visit.SecondDart.IsBust(currentScore)
		currentScore = currentScore - visit.SecondDart.GetScore()
		if !isBust && currentScore > 0 {
			isBust = visit.ThirdDart.IsBust(currentScore)
		} else {
			// Invalidate third dart if second was bust
			visit.ThirdDart.Value = null.IntFromPtr(nil)
		}
	} else {
		// Invalidate second/third dart if first was bust
		visit.SecondDart.Value = null.IntFromPtr(nil)
		visit.ThirdDart.Value = null.IntFromPtr(nil)
	}

	if !isBust && currentScore > 0 {
		// If this visit was not a bust, make sure that darts are set
		// as 0 (miss) instead of 'nil' (not thrown)
		if !visit.FirstDart.Value.Valid {
			visit.FirstDart.Value = null.IntFrom(0)
		}
		if !visit.SecondDart.Value.Valid {
			visit.SecondDart.Value = null.IntFrom(0)
		}
		if !visit.ThirdDart.Value.Valid {
			visit.ThirdDart.Value = null.IntFrom(0)
		}
	}

	visit.IsBust = isBust
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

// GetScore will return the total points scored during the given visit
func (visit Visit) GetScore() int {
	return visit.FirstDart.GetScore() + visit.SecondDart.GetScore() + visit.ThirdDart.GetScore()
}
