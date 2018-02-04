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
	Value      null.Int `json:"value"`
	Multiplier int64    `json:"multiplier"`
	IsCheckout bool     `json:"is_checkout"`
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
			&first.Value, &first.Multiplier, &first.IsCheckout,
			&second.Value, &second.Multiplier, &second.IsCheckout,
			&third.Value, &third.Multiplier, &third.IsCheckout,
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
			&first.Value, &first.Multiplier, &first.IsCheckout,
			&second.Value, &second.Multiplier, &second.IsCheckout,
			&third.Value, &third.Multiplier, &third.IsCheckout,
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
		&v.FirstDart.Value, &v.FirstDart.Multiplier, &v.FirstDart.IsCheckout,
		&v.SecondDart.Value, &v.SecondDart.Multiplier, &v.SecondDart.IsCheckout,
		&v.ThirdDart.Value, &v.ThirdDart.Multiplier, &v.ThirdDart.IsCheckout,
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
