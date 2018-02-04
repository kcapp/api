package models

import "github.com/guregu/null"

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

// Dart struct used for storing darts
type Dart struct {
	Value             null.Int `json:"value"`
	Multiplier        int64    `json:"multiplier"`
	IsCheckoutAttempt bool     `json:"is_checkout"`
	IsBust            bool     `json:"is_bust,omitempty"`
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

	setVisitModifiers(currentScore, visit.FirstDart, visit.SecondDart, visit.ThirdDart)
	visit.IsBust = visit.FirstDart.IsBust || visit.SecondDart.IsBust || visit.ThirdDart.IsBust

	stmt, err := db.Prepare(`
		INSERT INTO score(
			match_id, player_id,
			first_dart, first_dart_multiplier, is_checkout_first,
			second_dart, second_dart_multiplier, is_checkout_second,
			third_dart, third_dart_multiplier, is_checkout_third,
			is_bust, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(visit.MatchID, visit.PlayerID,
		visit.FirstDart.Value, visit.FirstDart.Multiplier, visit.FirstDart.IsCheckoutAttempt,
		visit.SecondDart.Value, visit.SecondDart.Multiplier, visit.SecondDart.IsCheckoutAttempt,
		visit.ThirdDart.Value, visit.ThirdDart.Multiplier, visit.ThirdDart.IsCheckoutAttempt,
		visit.IsBust)
	if err != nil {
		return err
	}

	// TODO set current player
	stmt, err = db.Prepare(`UPDATE ` + "`match`" + ` SET current_player_id = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return nil
}

func getPlayerScore(playerID int, matchID int) (int, error) {
	var currentScore int
	err := db.QueryRow(`
		SELECT m.starting_score - SUM(first_dart * first_dart_multiplier + second_dart * second_dart_multiplier + third_dart * third_dart_multiplier)
		FROM score s LEFT JOIN `+"`match`"+` m ON m.id = s.match_id
		WHERE player_id = ? AND match_id = ? AND is_bust = 0`, playerID, matchID).Scan(&currentScore)
	if err != nil {
		return 0, err
	}
	return currentScore, nil
}

func setVisitModifiers(currentScore int, firstDart *Dart, secondDart *Dart, thirdDart *Dart) {
	// Make sure to write multipliers as 1 (single) on miss
	if !firstDart.Value.Valid || firstDart.Value.Int64 == 0 {
		firstDart.Multiplier = 1
	}
	if !secondDart.Value.Valid || secondDart.Value.Int64 == 0 {
		secondDart.Multiplier = 1
	}
	if !thirdDart.Value.Valid || thirdDart.Value.Int64 == 0 {
		thirdDart.Multiplier = 1
	}

	firstDart.IsBust = isBust(firstDart, currentScore)
	if !firstDart.IsBust {
		firstDart.IsCheckoutAttempt = isCheckoutAttempt(currentScore)
	}
	currentScore = currentScore - int((firstDart.Value.Int64 * firstDart.Multiplier))

	if !firstDart.IsBust && !isCheckout(firstDart, currentScore) {
		secondDart.IsBust = isBust(secondDart, currentScore)
		if !secondDart.IsBust {
			secondDart.IsCheckoutAttempt = isCheckoutAttempt(currentScore)
		}
		currentScore = currentScore - int((secondDart.Value.Int64 * secondDart.Multiplier))

		if !secondDart.IsBust && !isCheckout(secondDart, currentScore) {
			thirdDart.IsBust = isBust(thirdDart, currentScore)
			if !thirdDart.IsBust {
				thirdDart.IsCheckoutAttempt = isCheckoutAttempt(currentScore)
			}
			currentScore = currentScore - int((thirdDart.Value.Int64 * thirdDart.Multiplier))
		}
	}
}

func isBust(dart *Dart, currentScore int) bool {
	scoreAfterThrow := currentScore - int((dart.Value.Int64 * dart.Multiplier))
	if isCheckout(dart, currentScore) {
		return false
	} else if scoreAfterThrow < 2 {
		return true
	}
	return false
}

func isCheckoutAttempt(currentScore int) bool {
	return currentScore == 50 || (currentScore <= 40 && currentScore%2 == 0)
}

func isCheckout(dart *Dart, currentScore int) bool {
	return currentScore-int((dart.Value.Int64*dart.Multiplier)) == 0 && dart.Multiplier == 2
}
