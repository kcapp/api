package data

import (
	"database/sql"
	"log"
	"time"

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

func CheckLegForBadges(leg *models.Leg, statistics map[int]*models.BadgeStatistics) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}
	playersMap, err := GetPlayersScore(leg.ID)
	if err != nil {
		return err
	}
	players := make([]*models.Player2Leg, 0, len(playersMap))
	for _, value := range playersMap {
		players = append(players, value)
	}

	for _, badge := range models.LegBadges {
		valid, playerID := badge.Validate(leg)
		if valid {
			if playerID != nil {
				err = addLegBadge(tx, *playerID, leg.ID, badge, leg.UpdatedAt)
				if err != nil {
					return err
				}
			} else {
				for _, playerID := range leg.Players {
					err = addLegBadge(tx, playerID, leg.ID, badge, leg.UpdatedAt)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	for _, badge := range models.LegPlayerBadges {
		valid, playerID := badge.Validate(leg, players)
		if valid {
			err = addLegPlayerBadge(tx, *playerID, leg.ID, badge, leg.UpdatedAt)
			if err != nil {
				return err
			}
		}
	}

	for _, badge := range models.VisitBadges {
		for _, playerID := range leg.Players {
			stats := statistics[playerID]
			valid, level := badge.Validate(stats, leg.Visits)
			if valid {
				err = addVisitBadge(tx, playerID, *level, leg.ID, badge, leg.UpdatedAt)
				if err != nil {
					return err
				}
			}
		}
	}
	tx.Commit()

	return nil
}

func AddBadge(playerID int, badge models.GlobalBadge) error {
	_, err := models.DB.Exec("INSERT IGNORE INTO player2badge (player_id, badge_id, created_at) VALUES (?, ?, ?)",
		playerID, badge.GetID(), time.Now())
	if err != nil {
		return err
	}
	log.Printf("Added global badge %d to player %d", badge.GetID(), playerID)
	return nil
}

func AddTournamentBadge(playerID int, tournamentID int, badge models.GlobalBadge, when time.Time) error {
	_, err := models.DB.Exec("INSERT IGNORE INTO player2badge (player_id, badge_id, tournament_id, created_at) VALUES (?, ?, ?, ?)",
		playerID, badge.GetID(), tournamentID, when)
	if err != nil {
		return err
	}
	log.Printf("Added tournament badge %d to player %d", badge.GetID(), playerID)
	return nil
}

func addLegBadge(tx *sql.Tx, playerID int, legID int, badge models.LegBadge, when time.Time) error {
	_, err := tx.Exec("INSERT IGNORE INTO player2badge (player_id, badge_id, leg_id, created_at) VALUES (?, ?, ?, ?)",
		playerID, badge.GetID(), legID, when)
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Printf("Added leg badge %d to player %d on leg %d", badge.GetID(), playerID, legID)
	return nil
}

func addLegPlayerBadge(tx *sql.Tx, playerID int, legID int, badge models.LegPlayerBadge, when time.Time) error {
	_, err := tx.Exec("INSERT IGNORE INTO player2badge (player_id, badge_id, leg_id, created_at) VALUES (?, ?, ?, ?)",
		playerID, badge.GetID(), legID, when)
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Printf("Added leg player badge %d to player %d on leg %d", badge.GetID(), playerID, legID)
	return nil
}

func addVisitBadge(tx *sql.Tx, playerID int, level int, legID int, badge models.VisitBadge, when time.Time) error {
	_, err := tx.Exec(`INSERT INTO player2badge (player_id, badge_id, level, value, leg_id, created_at) VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE leg_id=IF(?>level,?,leg_id), created_at=IF(?>level,?,created_at), value=IF(?>level,?,value),level=?`,
		playerID, badge.GetID(), level, badge.Levels()[level-1], legID, when, level, legID, level, when, level, badge.Levels()[level-1], level)
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Printf("Added visit badge %d (level %d) to player %d on leg %d", badge.GetID(), level, playerID, legID)
	return nil
}
