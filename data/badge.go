package data

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kcapp/api/models"
)

func GetBadges() ([]*models.Badge, error) {
	rows, err := models.DB.Query(`
		SELECT
			b.id,
			b.name,
			b.description,
			b.filename,
			b.secret,
			b.hidden,
			b.levels
		FROM badge b`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	badges := make([]*models.Badge, 0)
	for rows.Next() {
		badge := new(models.Badge)
		err := rows.Scan(&badge.ID, &badge.Name, &badge.Description, &badge.Filename, &badge.Secret, &badge.Hidden, &badge.Levels)
		if err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}
	return badges, nil
}

func GetBadge(badgeID int) (*models.Badge, error) {
	badge := new(models.Badge)
	err := models.DB.QueryRow(`
		SELECT
			b.id,
			b.name,
			b.description,
			b.filename,
			b.levels
		FROM badge b WHERE b.id = ?`, badgeID).Scan(&badge.ID, &badge.Name, &badge.Description, &badge.Filename, &badge.Levels)
	if err != nil {
		return nil, err
	}
	return badge, nil
}

func GetBadgesStatistics() ([]*models.BadgeStatistics, error) {
	players, err := GetPlayers()
	if err != nil {
		return nil, err
	}
	numPlayers := 0
	for _, player := range players {
		if !player.IsBot && !player.IsPlaceholder {
			numPlayers++
		}
	}
	rows, err := models.DB.Query(`
		SELECT
			b.id, p2b.level, p2b.value,
			COUNT(DISTINCT p2b.player_id) AS 'players_unlocked',
			MIN(p2b.created_at) AS 'first_unlock',
			GROUP_CONCAT(DISTINCT p2b.player_id ORDER BY p2b.created_at) AS 'players'
		FROM badge b
			LEFT JOIN player2badge p2b ON b.id = p2b.badge_id
		GROUP BY b.id, p2b.level`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	badges := make([]*models.BadgeStatistics, 0)
	for rows.Next() {
		var players []uint8
		statistic := new(models.BadgeStatistics)
		err := rows.Scan(&statistic.BadgeID, &statistic.Level, &statistic.Value, &statistic.UnlockedPlayers, &statistic.FirstUnlock, &players)
		if err != nil {
			return nil, err
		}
		playerStr := string(players)
		if playerStr != "" {
			playerIDs := strings.Split(playerStr, ",")
			for _, playerID := range playerIDs {
				id, err := strconv.Atoi(playerID)
				if err != nil {
					return nil, err
				}
				statistic.Players = append(statistic.Players, id)
			}
		}

		if statistic.UnlockedPlayers > 0 {
			statistic.UnlockedPercent = float32(statistic.UnlockedPlayers) / float32(numPlayers)
		}
		badges = append(badges, statistic)
	}
	return badges, nil
}

func GetBadgeStatistics(badgeID int) ([]*models.PlayerBadge, error) {
	rows, err := models.DB.Query(`
		SELECT
			b.id, b.name, b.description, b.filename,
			p2b.player_id, p2b.level, p2b.value, p2b.leg_id,
			p2b.match_id, p2b.tournament_id, p2b.opponent_player_id,
			p2b.visit_id,
			s.first_dart, IFNULL(s.first_dart_multiplier, 1),
			s.second_dart, IFNULL(s.second_dart_multiplier, 1),
			s.third_dart, IFNULL(s.third_dart_multiplier, 1),
			p2b.data,
			p2b.created_at
		FROM player2badge p2b
			LEFT JOIN badge b ON b.id = p2b.badge_id
			LEFT JOIN score s on s.id = p2b.visit_id
			LEFT JOIN player p on p.id = p2b.player_id
		WHERE b.id = ?
		ORDER BY level, created_at, p.first_name`, badgeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playerIds := make([]int, 0)
	badges := make([]*models.PlayerBadge, 0)
	for rows.Next() {
		badge := new(models.PlayerBadge)
		badge.Badge = new(models.Badge)
		darts := make([]*models.Dart, 3)
		darts[0] = new(models.Dart)
		darts[1] = new(models.Dart)
		darts[2] = new(models.Dart)

		err := rows.Scan(
			&badge.Badge.ID,
			&badge.Badge.Name,
			&badge.Badge.Description,
			&badge.Badge.Filename,
			&badge.PlayerID,
			&badge.Level,
			&badge.Value,
			&badge.LegID,
			&badge.MatchID,
			&badge.TournamentID,
			&badge.OpponentPlayerID,
			&badge.VisitID,
			&darts[0].Value, &darts[0].Multiplier,
			&darts[1].Value, &darts[1].Multiplier,
			&darts[2].Value, &darts[2].Multiplier,
			&badge.Data,
			&badge.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if badge.VisitID.Valid {
			badge.Darts = darts
		}
		badges = append(badges, badge)
		playerIds = append(playerIds, badge.PlayerID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(badges) > 0 {
		badge := badges[0]
		if badge.Level.Valid {
			stats, err := GetPlayerBadgeStatistics(playerIds, nil)
			if err != nil {
				return nil, err
			}

			for _, badge := range badges {
				badge.Statistics = stats[badge.PlayerID]
			}
		}
	}
	return badges, nil
}

func CheckLegForBadges(leg *models.Leg, statistics map[int]*models.PlayerBadgeStatistics) error {
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
		valid, playerID, visitID := badge.Validate(leg)
		if valid {
			if playerID != nil {
				err = addBadge(tx, badge.GetID(), nil, nil, *playerID, nil, &leg.ID, visitID, nil, leg.UpdatedAt)
				if err != nil {
					return err
				}
			} else {
				for _, playerID := range leg.Players {
					err = addBadge(tx, badge.GetID(), nil, nil, playerID, nil, &leg.ID, visitID, nil, leg.UpdatedAt)
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
			err = addBadge(tx, badge.GetID(), nil, nil, *playerID, nil, &leg.ID, nil, nil, leg.UpdatedAt)
			if err != nil {
				return err
			}
		}
	}

	for _, badge := range models.VisitBadgesLevel {
		for _, playerID := range leg.Players {
			stats := statistics[playerID]
			valid, level, visitID := badge.Validate(stats, leg.Visits)
			if valid {
				err = addBadge(tx, badge.GetID(), level, badge.Levels(), playerID, nil, &leg.ID, visitID, nil, leg.UpdatedAt)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, badge := range models.VisitBadges {
		for _, playerID := range leg.Players {
			valid, visitID := badge.Validate(playerID, leg.Visits)
			if valid {
				err = addBadge(tx, badge.GetID(), nil, nil, playerID, nil, &leg.ID, visitID, nil, leg.UpdatedAt)
				if err != nil {
					return err
				}
			}
		}
	}
	tx.Commit()

	return nil
}

func CheckMatchForBadges(match *models.Match) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}

	leaderboard, err := GetPlayersLastXLegsStatistics()
	if err != nil {
		return err
	}

	when := time.Now()
	if !match.EndTime.IsZero() {
		when = match.EndTime
	}
	for _, badge := range models.MatchBadges {
		valid, playerIDs := badge.Validate(match)
		if valid {
			if playerIDs != nil {
				for _, playerID := range playerIDs {
					err = addBadge(tx, badge.GetID(), nil, nil, playerID, &match.ID, nil, nil, nil, when)
					if err != nil {
						return err
					}
				}
			} else {
				for _, playerID := range match.Players {
					err = addBadge(tx, badge.GetID(), nil, nil, playerID, &match.ID, nil, nil, nil, when)
					if err != nil {
						return err
					}
				}
			}

		}
	}
	matchStatistics, err := GetX01StatisticsForMatch(match.ID)
	if err != nil {
		return err
	}
	for _, badge := range models.LeaderboardBadges {
		valid, playerID, opponentPlayerID := badge.Validate(match, matchStatistics, leaderboard)
		if valid {
			err = addBadge(tx, badge.GetID(), nil, nil, *playerID, &match.ID, nil, nil, opponentPlayerID, when)
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()

	return nil
}

func AddGlobalBadge(playerID int, badge models.GlobalBadge) error {
	return AddGlobalBadgeWithTime(playerID, badge, time.Now())
}

func AddGlobalBadgeWithTime(playerID int, badge models.GlobalBadge, when time.Time) error {
	_, err := models.DB.Exec("INSERT IGNORE INTO player2badge (player_id, badge_id, created_at) VALUES (?, ?, ?)",
		playerID, badge.GetID(), when)
	if err != nil {
		return err
	}
	log.Printf("Added global badge %d to player %d", badge.GetID(), playerID)
	return nil
}

func AddGlobalLevelBadge(playerID int, level int, badge models.GlobalLevelBadge) error {
	_, err := models.DB.Exec(`INSERT INTO player2badge (player_id, badge_id, level, value, created_at) VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE value=IF(?>level,?,value),level=?`,
		playerID, badge.GetID(), level, badge.Levels()[level-1], time.Now(), level, badge.Levels()[level-1], level)
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

func addBadge(tx *sql.Tx, badgeID int, level *int, levels []int, playerID int, matchID *int, legID *int, visitID *int, opponentPlayerID *int, when time.Time) error {
	if level != nil {
		_, err := tx.Exec(`INSERT INTO player2badge(player_id, badge_id, level, value, leg_id, visit_id,created_at) VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE leg_id=IF(?>level,?,leg_id), visit_id=IF(?>level,?,visit_id), created_at=IF(?>level,?,created_at), value=IF(?>level,?,value),level=?`,
			playerID, badgeID, level, levels[*level-1], legID, visitID, when, level, legID, level, visitID, level, when, level, levels[*level-1], level)
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("Added badge %d (level %d) to player %d on leg %d", badgeID, *level, playerID, *legID)
	} else {
		_, err := tx.Exec("INSERT IGNORE INTO player2badge (player_id, badge_id, match_id, leg_id, visit_id, opponent_player_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			playerID, badgeID, matchID, legID, visitID, opponentPlayerID, when)
		if err != nil {
			tx.Rollback()
			return err
		}
		if legID != nil {
			log.Printf("Added badge %d to player %d on leg %d", badgeID, playerID, *legID)
		} else if matchID != nil {
			log.Printf("Added badge %d to player %d on match %d", badgeID, playerID, *matchID)
		} else {
			log.Printf("Added badge %d to player %d", badgeID, playerID)

		}
	}
	return nil
}
