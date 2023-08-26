package data

import (
	"database/sql"
	"log"

	"github.com/kcapp/api/models"
)

func GetBadges() ([]*models.Badge, error) {
	rows, err := models.DB.Query(`
		SELECT
			b.id,
			b.name,
			b.description,
			b.filename
		FROM badge b`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	badges := make([]*models.Badge, 0)
	for rows.Next() {
		badge := new(models.Badge)
		err := rows.Scan(
			&badge.ID,
			&badge.Name,
			&badge.Description,
			&badge.Filename,
		)
		if err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}
	return badges, nil
}

func CheckLegForBadges(leg *models.Leg) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}
	for _, badge := range models.LegBadges {
		valid, playerID := badge.Validate(leg)
		if valid {
			if playerID != nil {
				addLegBadge(tx, *playerID, leg.ID, badge)
				if err != nil {
					return err
				}
			} else {
				for _, playerID := range leg.Players {
					addLegBadge(tx, playerID, leg.ID, badge)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	tx.Commit()

	return nil
}

func addLegBadge(tx *sql.Tx, playerID int, legID int, badge models.LegBadge) error {
	_, err := tx.Exec("INSERT INTO player2badge (player_id, badge_id, leg_id) VALUES (?, ?, ?)",
		playerID, badge.GetID(), legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Printf("Added badge %d to player %d on leg %d", badge.GetID(), playerID, legID)
	return nil
}
