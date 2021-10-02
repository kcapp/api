package data

import (
	"database/sql"
	"log"

	"github.com/guregu/null"
	"github.com/kcapp/api/models"
)

// GetJDCPracticeStatistics will return statistics for all players active during the given period
func GetJDCPracticeStatistics(from string, to string) ([]*models.StatisticsJDCPractice, error) {
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

	stats := make([]*models.StatisticsJDCPractice, 0)
	for rows.Next() {
		s := new(models.StatisticsJDCPractice)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.OfficeID, &s.DartsThrown,
			&s.Score, &s.MPR, &s.ShanghaiCount, &s.DoublesHitrate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetJDCPracticeStatisticsForLeg will return statistics for all players in the given leg
func GetJDCPracticeStatisticsForLeg(id int) ([]*models.StatisticsJDCPractice, error) {
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

	stats := make([]*models.StatisticsJDCPractice, 0)
	for rows.Next() {
		s := new(models.StatisticsJDCPractice)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.ShanghaiCount, &s.DoublesHitrate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetJDCPracticeStatisticsForMatch will return statistics for all players in the given match
func GetJDCPracticeStatisticsForMatch(id int) ([]*models.StatisticsJDCPractice, error) {
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

	stats := make([]*models.StatisticsJDCPractice, 0)
	for rows.Next() {
		s := new(models.StatisticsJDCPractice)
		err := rows.Scan(&s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.ShanghaiCount, &s.DoublesHitrate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetJDCPracticeStatisticsForPlayer will return JDC Practice statistics for the given player
func GetJDCPracticeStatisticsForPlayer(id int) (*models.StatisticsJDCPractice, error) {
	s := new(models.StatisticsJDCPractice)
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
			return new(models.StatisticsJDCPractice), nil
		}
		return nil, err
	}
	return s, nil
}

// GetJDCPracticeHistoryForPlayer will return history of JDC Practice statistics for the given player
func GetJDCPracticeHistoryForPlayer(id int, limit int) ([]*models.Leg, error) {
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
		s := new(models.StatisticsJDCPractice)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.ShanghaiCount, &s.DoublesHitrate)
		if err != nil {
			return nil, err
		}
		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// CalculateJDCPracticeStatistics will generate JDC Practice statistics for the given leg
func CalculateJDCPracticeStatistics(legID int) (map[int]*models.StatisticsJDCPractice, error) {
	leg, err := GetLeg(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}

	statisticsMap := make(map[int]*models.StatisticsJDCPractice)
	for _, player := range players {
		stats := new(models.StatisticsJDCPractice)
		stats.PlayerID = player.PlayerID
		stats.Score = player.CurrentScore
		statisticsMap[player.PlayerID] = stats
		stats.MPR = null.FloatFrom(0)
	}

	round := 0
	for i, visit := range leg.Visits {
		if i > 0 && i%len(players) == 0 {
			round++
		}
		stats := statisticsMap[visit.PlayerID]

		first := visit.FirstDart
		second := visit.SecondDart
		third := visit.ThirdDart

		target := models.TargetsJDCPractice[round]
		if target.Values == nil {
			if first.ValueRaw() == target.Value {
				stats.Marks += first.Multiplier
			}
			if second.ValueRaw() == target.Value {
				stats.Marks += second.Multiplier
			}
			if third.ValueRaw() == target.Value {
				stats.Marks += third.Multiplier
			}
			if first.ValueRaw() == target.Value && visit.IsShanghai() {
				stats.ShanghaiCount++
			}
		} else {
			values := target.Values

			if first.IsDouble() && first.ValueRaw() == values[0] {
				stats.DoublesHitrate++
			}
			if second.IsDouble() && second.ValueRaw() == values[1] {
				stats.DoublesHitrate++
			}
			if third.IsDouble() && third.ValueRaw() == values[2] {
				stats.DoublesHitrate++
			}
		}

		stats.DartsThrown = visit.DartsThrown
	}

	for _, stats := range statisticsMap {
		stats.MPR = null.FloatFrom(float64(stats.Marks) / 12.0)
		stats.DoublesHitrate = stats.DoublesHitrate / 21.0
	}
	return statisticsMap, nil
}

// ReCalculateJDCPracticeStatistics will recaulcate statistics for JDC Practice legs
func ReCalculateJDCPracticeStatistics() (map[int]map[int]*models.StatisticsJDCPractice, error) {
	legs, err := GetLegsOfType(models.JDCPRACTICE, true)
	if err != nil {
		return nil, err
	}

	s := make(map[int]map[int]*models.StatisticsJDCPractice)
	for _, leg := range legs {
		stats, err := CalculateJDCPracticeStatistics(leg.ID)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			log.Printf(`UPDATE statistics_jdc_practice SET darts_thrown = %d, score = %d, mpr = %f,
				shanghai_count = %d, doubles_hitrate = %f WHERE leg_id = %d AND player_id = %d;`,
				stat.DartsThrown, stat.Score, stat.MPR.Float64, stat.ShanghaiCount, stat.DoublesHitrate, leg.ID, playerID)
		}
		s[leg.ID] = stats
	}

	return s, err
}
