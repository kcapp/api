package data

import (
	"github.com/kcapp/api/models"
)

// GetShootoutStatistics will return statistics for all players active duing the given period
func GetShootoutStatistics(from string, to string) ([]*models.StatisticsShootout, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id),
			SUM(s.ppd) / COUNT(p.id) AS 'ppd',
			SUM(s.60s_plus),
			SUM(s.100s_plus),
			SUM(s.140s_plus),
			SUM(s.180s) AS '180s'
		FROM statistics_shootout s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND m.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 2
		GROUP BY p.id`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statsMap := make(map[int]*models.StatisticsShootout, 0)
	for rows.Next() {
		s := new(models.StatisticsShootout)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.PPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus, &s.Score180s)
		if err != nil {
			return nil, err
		}
		statsMap[s.PlayerID] = s
	}

	rows, err = models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			COUNT(m.winner_id) AS 'matches_won'
		FROM matches m
			JOIN player p ON p.id = m.winner_id
		WHERE m.updated_at >= ? AND m.updated_at < ?
		AND m.match_type_id = 2
		GROUP BY m.winner_id`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var playerID int
		var matchesWon int
		err := rows.Scan(&playerID, &matchesWon)
		if err != nil {
			return nil, err
		}
		statsMap[playerID].MatchesWon = matchesWon
	}

	stats := make([]*models.StatisticsShootout, 0)
	for _, s := range statsMap {
		stats = append(stats, s)
	}

	return stats, nil
}

// GetShootoutStatisticsForLeg will return statistics for all players in the given leg
func GetShootoutStatisticsForLeg(id int) ([]*models.StatisticsShootout, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id,
			p.id,
			SUM(s.ppd) / COUNT(p.id),
			SUM(60s_plus),
			SUM(100s_plus),
			SUM(140s_plus),
			SUM(180s) AS '180s'
		FROM statistics_shootout s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE l.id = ? GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsShootout, 0)
	for rows.Next() {
		s := new(models.StatisticsShootout)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.PPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus, &s.Score180s)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
