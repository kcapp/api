package data

import (
	"log"

	"github.com/kcapp/api/models"
)

// AddVisit will write thegiven visit to database
func AddVisit(visit models.Visit) error {
	currentScore, err := GetPlayerScore(visit.PlayerID, visit.MatchID)
	if err != nil {
		return err
	}

	// TODO Don't allow to save score for same player twice in a row
	// Only allow saving score for match.current_player_id ?

	visit.SetIsBust(currentScore)

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
	nextPlayerID := order[(currentPlayerOrder%len(players))+1]

	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT INTO score(
			match_id, player_id,
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier,
			is_bust, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		visit.MatchID, visit.PlayerID,
		visit.FirstDart.Value, visit.FirstDart.Multiplier,
		visit.SecondDart.Value, visit.SecondDart.Multiplier,
		visit.ThirdDart.Value, visit.ThirdDart.Multiplier,
		visit.IsBust)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE `+"`match`"+` SET current_player_id = ? WHERE id = ?`, nextPlayerID, visit.MatchID)
	if err != nil {
		return err
	}
	tx.Commit()

	log.Printf("[%d] Added score for player %d, (%d-%d, %d-%d, %d-%d, %t)", visit.MatchID, visit.PlayerID, visit.FirstDart.Value.Int64,
		visit.FirstDart.Multiplier, visit.SecondDart.Value.Int64, visit.SecondDart.Multiplier, visit.ThirdDart.Value.Int64, visit.ThirdDart.Multiplier,
		visit.IsBust)

	return nil
}

// ModifyVisit modify the scores of a visit
func ModifyVisit(visit models.Visit) error {
	// FIXME: We need to check if this is a checkout/bust
	stmt, err := models.DB.Prepare(`
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
	log.Printf("[%d] Modified score %d, throws: (%d-%d, %d-%d, %d-%d)", visit.MatchID, visit.ID, visit.FirstDart.Value.Int64,
		visit.FirstDart.Multiplier, visit.SecondDart.Value.Int64, visit.SecondDart.Multiplier, visit.ThirdDart.Value.Int64, visit.ThirdDart.Multiplier)

	return nil
}

// DeleteVisit will delete the visit for the given ID
func DeleteVisit(id int) error {
	visit, err := GetVisit(id)
	if err != nil {
		return err
	}
	tx, err := models.DB.Begin()
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

	log.Printf("[%d] Deleted visit %d", visit.MatchID, visit.ID)
	return nil
}

// GetPlayerVisits will return all visits for a given player
func GetPlayerVisits(id int) ([]*models.Visit, error) {
	rows, err := models.DB.Query(`
		SELECT
			id, match_id, player_id, 
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier,
			is_bust,
			created_at,
			updated_at
		FROM score s
		WHERE player_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	visits := make([]*models.Visit, 0)
	for rows.Next() {
		v := new(models.Visit)
		v.FirstDart = new(models.Dart)
		v.SecondDart = new(models.Dart)
		v.ThirdDart = new(models.Dart)
		err := rows.Scan(&v.ID, &v.MatchID, &v.PlayerID,
			&v.FirstDart.Value, &v.FirstDart.Multiplier,
			&v.SecondDart.Value, &v.SecondDart.Multiplier,
			&v.ThirdDart.Value, &v.ThirdDart.Multiplier,
			&v.IsBust, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}
		visits = append(visits, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return visits, nil
}

// GetMatchVisits will return all visits for a given match
func GetMatchVisits(id int) ([]*models.Visit, error) {
	rows, err := models.DB.Query(`
		SELECT
			id, match_id, player_id, 
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier,
			is_bust,
			created_at,
			updated_at
		FROM score s
		WHERE match_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	visits := make([]*models.Visit, 0)
	for rows.Next() {
		v := new(models.Visit)
		v.FirstDart = new(models.Dart)
		v.SecondDart = new(models.Dart)
		v.ThirdDart = new(models.Dart)
		err := rows.Scan(&v.ID, &v.MatchID, &v.PlayerID,
			&v.FirstDart.Value, &v.FirstDart.Multiplier,
			&v.SecondDart.Value, &v.SecondDart.Multiplier,
			&v.ThirdDart.Value, &v.ThirdDart.Multiplier,
			&v.IsBust, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}
		visits = append(visits, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return visits, nil
}

// GetVisit will return the visit with the given ID
func GetVisit(id int) (*models.Visit, error) {
	v := new(models.Visit)
	v.FirstDart = new(models.Dart)
	v.SecondDart = new(models.Dart)
	v.ThirdDart = new(models.Dart)
	err := models.DB.QueryRow(`
		SELECT
			id, match_id, player_id, 
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier,
			is_bust,
			created_at,
			updated_at
		FROM score s
		WHERE s.id = ?`, id).Scan(&v.ID, &v.MatchID, &v.PlayerID,
		&v.FirstDart.Value, &v.FirstDart.Multiplier,
		&v.SecondDart.Value, &v.SecondDart.Multiplier,
		&v.ThirdDart.Value, &v.ThirdDart.Multiplier,
		&v.IsBust, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return v, nil
}
