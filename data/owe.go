package data

import (
	"errors"
	"log"

	"github.com/kcapp/api/models"
)

// GetOwes will return all current owes between players
func GetOwes() ([]*models.Owe, error) {
	rows, err := models.DB.Query(`
		SELECT 
			o.player_ower_id,
			o.player_owee_id,
			ot.id, ot.item,
			o. amount 
		FROM owes o 
		JOIN owe_type ot ON ot.id = o.owe_type_id
		WHERE o.amount > 0`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	owes := make([]*models.Owe, 0)
	for rows.Next() {
		o := new(models.Owe)
		o.OweType = new(models.OweType)
		err := rows.Scan(&o.PlayerOwerID, &o.PlayerOweeID, &o.OweType.ID, &o.OweType.Item, &o.Amount)
		if err != nil {
			return nil, err
		}
		owes = append(owes, o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return owes, nil
}

// RegisterPayback will register a payback between the given players
func RegisterPayback(owe models.Owe) error {
	stmt, err := models.DB.Prepare(`UPDATE owes SET amount = amount - ? WHERE player_ower_id = ? AND player_owee_id = ? and owe_type_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(owe.Amount, owe.PlayerOwerID, owe.PlayerOweeID, owe.OweType.ID)
	if err != nil {
		return err
	}

	updatedRows, err := res.RowsAffected()
	if updatedRows == 0 {
		return errors.New("No rows were updated when registering payback")
	}
	log.Printf("Player %d paid back %d %s to player %d", owe.PlayerOwerID, owe.Amount, owe.OweType.Item.String, owe.PlayerOweeID)
	return nil
}
