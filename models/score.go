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

// Dart struct used for storing darts
type Dart struct {
	Value             null.Int `json:"value"`
	Multiplier        int64    `json:"multiplier"`
	IsCheckoutAttempt bool     `json:"is_checkout"`
	IsBust            bool     `json:"is_bust,omitempty"`
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

// GetPlayerVisits will return all visits for a given player
func GetPlayerVisits(id int) ([]*Visit, error) {
	rows, err := db.Query(`
		SELECT
			id, match_id, player_id, 
			first_dart, first_dart_multiplier, is_checkout_first,
			second_dart, second_dart_multiplier, is_checkout_second,
			third_dart, third_dart_multiplier, is_checkout_third,
			is_bust,
			created_at,
			updated_at
		FROM score s
		WHERE player_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	visits := make([]*Visit, 0)
	for rows.Next() {
		v := new(Visit)
		first := new(Dart)
		second := new(Dart)
		third := new(Dart)
		err := rows.Scan(&v.ID, &v.MatchID, &v.PlayerID,
			&first.Value, &first.Multiplier, &first.IsCheckoutAttempt,
			&second.Value, &second.Multiplier, &second.IsCheckoutAttempt,
			&third.Value, &third.Multiplier, &third.IsCheckoutAttempt,
			&v.IsBust, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}
		v.FirstDart = first
		v.SecondDart = second
		v.ThirdDart = third
		visits = append(visits, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return visits, nil
}

// GetMatchVisits will return all visits for a given match
func GetMatchVisits(id int) ([]*Visit, error) {
	rows, err := db.Query(`
		SELECT
			id, match_id, player_id, 
			first_dart, first_dart_multiplier, is_checkout_first,
			second_dart, second_dart_multiplier, is_checkout_second,
			third_dart, third_dart_multiplier, is_checkout_third,
			is_bust,
			created_at,
			updated_at
		FROM score s
		WHERE match_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	visits := make([]*Visit, 0)
	for rows.Next() {
		v := new(Visit)
		first := new(Dart)
		second := new(Dart)
		third := new(Dart)
		err := rows.Scan(&v.ID, &v.MatchID, &v.PlayerID,
			&first.Value, &first.Multiplier, &first.IsCheckoutAttempt,
			&second.Value, &second.Multiplier, &second.IsCheckoutAttempt,
			&third.Value, &third.Multiplier, &third.IsCheckoutAttempt,
			&v.IsBust, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}
		v.FirstDart = first
		v.SecondDart = second
		v.ThirdDart = third
		visits = append(visits, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return visits, nil
}

// GetVisit will return the visit with the given ID
func GetVisit(id int) (*Visit, error) {
	v := new(Visit)
	v.FirstDart = new(Dart)
	v.SecondDart = new(Dart)
	v.ThirdDart = new(Dart)
	err := db.QueryRow(`
		SELECT
			id, match_id, player_id, 
			first_dart, first_dart_multiplier, is_checkout_first,
			second_dart, second_dart_multiplier, is_checkout_second,
			third_dart, third_dart_multiplier, is_checkout_third,
			is_bust,
			created_at,
			updated_at
		FROM score s
		WHERE s.id = ?`, id).Scan(&v.ID, &v.MatchID, &v.PlayerID,
		&v.FirstDart.Value, &v.FirstDart.Multiplier, &v.FirstDart.IsCheckoutAttempt,
		&v.SecondDart.Value, &v.SecondDart.Multiplier, &v.SecondDart.IsCheckoutAttempt,
		&v.ThirdDart.Value, &v.ThirdDart.Multiplier, &v.ThirdDart.IsCheckoutAttempt,
		&v.IsBust, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// ModifyVisit modify the scores of a visit
func ModifyVisit(visit Visit) error {
	// FIXME: We need to check if this is a checkout/bust
	stmt, err := db.Prepare(`
		UPDATE score SET 
    		first_dart = ?,
    		first_dart_multiplier = ?,
    		second_dart = ?,
    		second_dart_multiplier = ?,
    		third_dart = ?,
		    third_dart_multiplier = ?,
			updated_at = NOW()
		WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(visit.FirstDart.Value, visit.FirstDart.Multiplier, visit.SecondDart.Value, visit.SecondDart.Multiplier,
		visit.ThirdDart.Value, visit.ThirdDart.Multiplier, visit.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteVisit will delete the visit for the given ID
func DeleteVisit(id int) error {
	visit, err := GetVisit(id)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Delete the visit
	_, err = tx.Exec("DELETE FROM score WHERE id = ?", id)
	if err != nil {
		return err
	}
	// Set current player to the player of the last visit
	_, err = tx.Exec("UPDATE `match` SET current_player_id = ? WHERE id = ?", visit.PlayerID, visit.MatchID)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// AddVisit will write thegiven visit to database
func AddVisit(visit Visit) error {
	currentScore, err := getPlayerScore(visit.PlayerID, visit.MatchID)
	if err != nil {
		return err
	}

	// TODO Don't allow to save score for same player twice in a row
	// Only allow saving score for match.current_player_id ?

	// Set visit modifiers
	setVisitModifiers(currentScore, visit.FirstDart, visit.SecondDart, visit.ThirdDart)
	visit.IsBust = visit.FirstDart.IsBust || visit.SecondDart.IsBust || visit.ThirdDart.IsBust

	// Determine who the next player will be
	players, err := GetMatchPlayers(visit.MatchID)
	if err != nil {
		return err
	}

	currentPlayerOrder := 1
	order := make(map[int]int)
	for _, player := range players {
		if player.PlayerID == visit.PlayerID {
			currentPlayerOrder = player.Order
		}
		order[player.Order] = player.PlayerID
	}
	nextPlayerId := order[(currentPlayerOrder%len(players))+1]

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT INTO score(
			match_id, player_id,
			first_dart, first_dart_multiplier, is_checkout_first,
			second_dart, second_dart_multiplier, is_checkout_second,
			third_dart, third_dart_multiplier, is_checkout_third,
			is_bust, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`, visit.MatchID, visit.PlayerID,
		visit.FirstDart.Value, visit.FirstDart.Multiplier, visit.FirstDart.IsCheckoutAttempt,
		visit.SecondDart.Value, visit.SecondDart.Multiplier, visit.SecondDart.IsCheckoutAttempt,
		visit.ThirdDart.Value, visit.ThirdDart.Multiplier, visit.ThirdDart.IsCheckoutAttempt,
		visit.IsBust)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE `+"`match`"+` SET current_player_id = ? WHERE id = ?`, nextPlayerId, visit.MatchID)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func getPlayerScore(playerID int, matchID int) (int, error) {
	var currentScore int
	err := db.QueryRow(`
		SELECT m.starting_score - IFNULL(SUM(first_dart * first_dart_multiplier + second_dart * second_dart_multiplier + third_dart * third_dart_multiplier), 0)
		FROM score s LEFT JOIN `+"`match`"+` m ON m.id = s.match_id
		WHERE player_id = ? AND m.id = ? AND is_bust = 0`, playerID, matchID).Scan(&currentScore)
	if err != nil {
		return 0, err
	}
	return currentScore, nil
}

func setVisitModifiers(currentScore int, firstDart *Dart, secondDart *Dart, thirdDart *Dart) {
	firstDart.SetModifiers(currentScore)
	currentScore = currentScore - firstDart.getScore()
	if currentScore < 2 {
		// If first dart is bust/checkout, then second/third dart doesn't count
		secondDart.Value.Valid = false
		secondDart.Multiplier = 1
		thirdDart.Value.Valid = false
		thirdDart.Multiplier = 1
	} else {
		secondDart.SetModifiers(currentScore)
		currentScore = currentScore - secondDart.getScore()
		if currentScore < 2 {
			// If second dart is bust, then third dart doesn't count
			thirdDart.Value.Valid = false
			thirdDart.Multiplier = 1
		} else {
			thirdDart.SetModifiers(currentScore)
			currentScore = currentScore - thirdDart.getScore()
		}
	}
}
