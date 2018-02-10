package models

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

// SetVisitModifiers will set IsBust and IsCheckoutAttempt on all darts
func (visit Visit) SetVisitModifiers(currentScore int) {
	visit.FirstDart.SetModifiers(currentScore)
	currentScore = currentScore - visit.FirstDart.GetScore()
	if currentScore < 2 {
		// If first dart is bust/checkout, then second/third dart doesn't count
		visit.SecondDart.Value.Valid = false
		visit.SecondDart.Multiplier = 1
		visit.ThirdDart.Value.Valid = false
		visit.ThirdDart.Multiplier = 1
	} else {
		visit.SecondDart.SetModifiers(currentScore)
		currentScore = currentScore - visit.SecondDart.GetScore()
		if currentScore < 2 {
			// If second dart is bust, then third dart doesn't count
			visit.ThirdDart.Value.Valid = false
			visit.ThirdDart.Multiplier = 1
		} else {
			visit.ThirdDart.SetModifiers(currentScore)
			currentScore = currentScore - visit.ThirdDart.GetScore()
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

// GetScore will return the total points scored during the given visit
func (visit Visit) GetScore() int {
	return visit.FirstDart.GetScore() + visit.SecondDart.GetScore() + visit.ThirdDart.GetScore()
}
