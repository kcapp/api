package data

import (
	"log"

	"github.com/kcapp/api/models"
)

// GetScameStatistics will return statistics for all players active during the given period
/*
func GetScameStatistics(from string, to string) ([]*models.StatisticsScam, error) {
	rows, err := models.DB.Query(`
			SELECT
				p.id,
				COUNT(DISTINCT m.id) AS 'matches_played',
				COUNT(DISTINCT m2.id) AS 'matches_won',
				COUNT(DISTINCT l.id) AS 'legs_played',
				COUNT(DISTINCT l2.id) AS 'legs_won',
				m.office_id AS 'office_id',
				SUM(s.darts_thrown) as 'darts_thrown',
				CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
				SUM(s.mpr) / COUNT(DISTINCT l.id) as 'mpr',
				SUM(s.shanghai_count) as 'shanghai_count',
				SUM(s.doubles_hitrate) / COUNT(l.id) as 'doubles_hitrate'
			FROM statistics_jdc_practice s
				JOIN player p ON p.id = s.player_id
				JOIN leg l ON l.id = s.leg_id
				JOIN matches m ON m.id = l.match_id
				LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
				LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
			WHERE m.updated_at >= ? AND m.updated_at < ?
				AND l.is_finished = 1 AND m.is_abandoned = 0
				AND m.match_type_id = 14
			GROUP BY p.id, m.office_id
			ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsScam, 0)
	for rows.Next() {
		s := new(models.StatisticsScam)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.OfficeID, &s.DartsThrown,
			&s.Score, &s.MPR, &s.ShanghaiCount, &s.DoublesHitrate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetScamStatisticsForLeg will return statistics for all players in the given leg
func GetScamStatisticsForLeg(id int) ([]*models.StatisticsScam, error) {
	rows, err := models.DB.Query(`
			SELECT
				l.id,
				p.id,
				s.darts_thrown,
				s.score,
				s.mpr,
				s.shanghai_count,
				s.doubles_hitrate
			FROM statistics_jdc_practice s
				JOIN player p ON p.id = s.player_id
				JOIN leg l ON l.id = s.leg_id
				JOIN player2leg p2l on l.id = p2l.leg_id AND p.id = p2l.player_id
			WHERE l.id = ? GROUP BY p.id ORDER BY p2l.order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsScam, 0)
	for rows.Next() {
		s := new(models.StatisticsScam)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.ShanghaiCount, &s.DoublesHitrate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetScamStatisticsForMatch will return statistics for all players in the given match
func GetScamStatisticsForMatch(id int) ([]*models.StatisticsScam, error) {
	rows, err := models.DB.Query(`
			SELECT
				p.id,
				SUM(s.darts_thrown) as 'darts_thrown',
				CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
				SUM(s.mpr) / COUNT(DISTINCT l.id) as 'mpr',
				SUM(s.shanghai_count) as 'shanghai_count',
				SUM(s.doubles_hitrate) / COUNT(l.id) as 'doubles_hitrate'
			FROM statistics_jdc_practice s
				JOIN player p ON p.id = s.player_id
				JOIN leg l ON l.id = s.leg_id
				JOIN matches m ON m.id = l.match_id
				JOIN player2leg p2l ON p2l.leg_id = l.id AND p2l.player_id = s.player_id
			WHERE m.id = ?
			GROUP BY p.id
			ORDER BY p2l.order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsScam, 0)
	for rows.Next() {
		s := new(models.StatisticsScam)
		err := rows.Scan(&s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.ShanghaiCount, &s.DoublesHitrate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetScamStatisticsForPlayer will return Scam statistics for the given player
func GetScamStatisticsForPlayer(id int) (*models.StatisticsScam, error) {
	s := new(models.StatisticsScam)
	err := models.DB.QueryRow(`
			SELECT
				p.id,
				COUNT(DISTINCT m.id) AS 'matches_played',
				COUNT(DISTINCT m2.id) AS 'matches_won',
				COUNT(DISTINCT l.id) AS 'legs_played',
				COUNT(DISTINCT l2.id) AS 'legs_won',
				SUM(s.darts_thrown) as 'darts_thrown',
				CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
				SUM(s.mpr) / COUNT(DISTINCT l.id) as 'mpr',
				SUM(s.shanghai_count) as 'shanghai_count',
				SUM(s.doubles_hitrate) / COUNT(l.id) as 'doubles_hitrate'
			FROM statistics_jdc_practice s
				JOIN player p ON p.id = s.player_id
				JOIN leg l ON l.id = s.leg_id
				JOIN matches m ON m.id = l.match_id
				LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
				LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
			WHERE s.player_id = ?
				AND l.is_finished = 1 AND m.is_abandoned = 0
				AND m.match_type_id = 14
			GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.DartsThrown,
		&s.Score, &s.MPR, &s.ShanghaiCount, &s.DoublesHitrate)
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.StatisticsScam), nil
		}
		return nil, err
	}
	return s, nil
}

// GetScamHistoryForPlayer will return history of Scam statistics for the given player
func GetScamHistoryForPlayer(id int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.JDCPRACTICE, false)
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
				s.darts_thrown,
				s.score,
				s.mpr,
				s.shanghai_count,
				s.doubles_hitrate
			FROM statistics_jdc_practice s
				LEFT JOIN player p ON p.id = s.player_id
				LEFT JOIN leg l ON l.id = s.leg_id
				LEFT JOIN matches m ON m.id = l.match_id
			WHERE s.player_id = ?
				AND l.is_finished = 1 AND m.is_abandoned = 0
				AND m.match_type_id = 14
			ORDER BY l.id DESC
			LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.StatisticsScam)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.ShanghaiCount, &s.DoublesHitrate)
		if err != nil {
			return nil, err
		}
		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}*/

// CalculateScamStatistics will generate Scam statistics for the given leg
func CalculateScamStatistics(legID int) (map[int]*models.StatisticsScam, error) {
	leg, err := GetLeg(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}

	statisticsMap := make(map[int]*models.StatisticsScam)
	for _, player := range players {
		player.Hits = make(models.HitsMap) // Reset HitsMap to accurately calculate MPR

		stats := new(models.StatisticsScam)
		stats.PlayerID = player.PlayerID
		stats.Score = player.CurrentScore
		statisticsMap[player.PlayerID] = stats
	}
	models.DecorateVisitsScam(players, leg.Visits)

	for _, visit := range leg.Visits {
		player := players[visit.PlayerID]
		stats := statisticsMap[visit.PlayerID]

		if visit.IsStopper.Bool {
			stats.DartsThrownStopper += visit.GetDartsThrown()

			hits := player.Hits
			if hits.Contains(visit.FirstDart.ValueRaw()) {
				stats.MPR += 1
			}
			hits.Add(visit.FirstDart)

			if hits.Contains(visit.SecondDart.ValueRaw()) {
				stats.MPR += 1
			}
			hits.Add(visit.SecondDart)

			if hits.Contains(visit.ThirdDart.ValueRaw()) {
				stats.MPR += 1
			}
			hits.Add(visit.ThirdDart)
		} else {
			stats.DartsThrownScorer += 3
		}
	}

	for _, stats := range statisticsMap {
		stats.PPD = float32(stats.Score) / float32(stats.DartsThrownScorer)
		stats.ThreeDartAvg = stats.PPD * 3
		stats.MPR = stats.MPR / float32((stats.DartsThrownStopper / 3))
	}
	return statisticsMap, nil
}

// ReCalculateScamStatistics will recaulcate statistics for Scam legs
func ReCalculateScamStatistics() (map[int]map[int]*models.StatisticsScam, error) {
	legs, err := GetLegsOfType(models.SCAM, true)
	if err != nil {
		return nil, err
	}

	s := make(map[int]map[int]*models.StatisticsScam)
	for _, leg := range legs {
		stats, err := CalculateScamStatistics(leg.ID)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			log.Printf(`UPDATE statistics_scam SET darts_thrown_stopper = %d, darts_thrown_scorer = %d, mpr = %f, score = %d, WHERE leg_id = %d AND player_id = %d;`,
				stat.DartsThrownStopper, stat.DartsThrownScorer, stat.MPR, stat.Score, leg.ID, playerID)
		}
		s[leg.ID] = stats
	}

	return s, err
}
