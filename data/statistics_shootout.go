package data

import (
	"github.com/kcapp/api/models"
)

// GetShootoutStatistics will return statistics for all players active duing the given period
func GetShootoutStatistics(from string, to string) ([]*models.StatisticsShootout, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			SUM(s.ppd) / COUNT(p.id) AS 'ppd',
			SUM(s.60s_plus),
			SUM(s.100s_plus),
			SUM(s.140s_plus),
			SUM(s.180s) AS '180s'
		FROM statistics_shootout s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND m.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 2
		GROUP BY p.id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC, ppd DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsShootout, 0)
	for rows.Next() {
		s := new(models.StatisticsShootout)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.PPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus, &s.Score180s)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetShootoutStatisticsForMatch will return statistics for the given match
func GetShootoutStatisticsForMatch(matchID int) ([]*models.StatisticsShootout, error) {
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
		WHERE m.id = ?
			AND m.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 2
		GROUP BY p.id`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsShootout, 0)
	for rows.Next() {
		s := new(models.StatisticsShootout)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.PPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus, &s.Score180s)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetShootoutStatisticsForLeg will return statistics for all players in the given leg
func GetShootoutStatisticsForLeg(id int) ([]*models.StatisticsShootout, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			ppd,
			60s_plus,
			100s_plus,
			140s_plus,
			180s
		FROM statistics_shootout s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
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

// GetShootoutStatisticsForPlayer will return Shootout statistics for the given player
func GetShootoutStatisticsForPlayer(id int) (*models.StatisticsShootout, error) {
	s := new(models.StatisticsShootout)
	err := models.DB.QueryRow(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			SUM(s.ppd) / COUNT(p.id) AS 'ppd',
			SUM(s.60s_plus),
			SUM(s.100s_plus),
			SUM(s.140s_plus),
			SUM(s.180s) AS '180s'
		FROM statistics_shootout s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 2
		GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.PPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus, &s.Score180s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetShootoutHistoryForPlayer will return history of Shootout statistics for the given player
func GetShootoutHistoryForPlayer(id int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.SHOOTOUT, false)
	if err != nil {
		return nil, err
	}
	m := make(map[int]*models.Leg)
	for _, leg := range legs {
		m[leg.ID] = leg
	}

	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			ppd,
			60s_plus,
			100s_plus,
			140s_plus,
			180s
		FROM statistics_shootout s
			LEFT JOIN player p ON p.id = s.player_id
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 2
		ORDER BY l.id DESC
		LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.StatisticsShootout)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.PPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus, &s.Score180s)
		if err != nil {
			return nil, err
		}
		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// CalculateShootoutStatistics will generate shootout statistics for the given leg
func CalculateShootoutStatistics(legID int) (map[int]*models.StatisticsShootout, error) {
	visits, err := GetLegVisits(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}
	statisticsMap := make(map[int]*models.StatisticsShootout)
	for _, player := range players {
		stats := new(models.StatisticsShootout)
		statisticsMap[player.PlayerID] = stats
	}

	for _, visit := range visits {
		stats := statisticsMap[visit.PlayerID]

		visitScore := visit.GetScore()
		stats.PPD += float32(visitScore)

		if visitScore >= 60 && visitScore < 100 {
			stats.Score60sPlus++
		} else if visitScore >= 100 && visitScore < 140 {
			stats.Score100sPlus++
		} else if visitScore >= 140 && visitScore < 180 {
			stats.Score140sPlus++
		} else if visitScore == 180 {
			stats.Score180s++
		}
	}

	for _, stats := range statisticsMap {
		stats.PPD = stats.PPD / float32(9)
	}
	return statisticsMap, nil
}
