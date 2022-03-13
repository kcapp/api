package data

import (
	"database/sql"
	"fmt"

	"github.com/kcapp/api/models"
)

// GetBermudaTriangleStatistics will return statistics for all players active during the given period
func GetBermudaTriangleStatistics(from string, to string) ([]*models.StatisticsBermudaTriangle, error) {
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
			SUM(s.total_marks) / (COUNT(DISTINCT l.id) * 13) as 'mpr',
			MAX(s.highest_score_reached) as 'highest_score_reached',
			SUM(s.total_hit_rate) / COUNT(l.id) as 'total_hit_rate',
			SUM(s.hit_rate_1) / COUNT(l.id) as 'hit_rate_1',
			SUM(s.hit_rate_2) / COUNT(l.id) as 'hit_rate_2',
			SUM(s.hit_rate_3) / COUNT(l.id) as 'hit_rate_3',
			SUM(s.hit_rate_4) / COUNT(l.id) as 'hit_rate_4',
			SUM(s.hit_rate_5) / COUNT(l.id) as 'hit_rate_5',
			SUM(s.hit_rate_6) / COUNT(l.id) as 'hit_rate_6',
			SUM(s.hit_rate_7) / COUNT(l.id) as 'hit_rate_7',
			SUM(s.hit_rate_8) / COUNT(l.id) as 'hit_rate_8',
			SUM(s.hit_rate_9) / COUNT(l.id) as 'hit_rate_9',
			SUM(s.hit_rate_10) / COUNT(l.id) as 'hit_rate_10',
			SUM(s.hit_rate_11) / COUNT(l.id) as 'hit_rate_11',
			SUM(s.hit_rate_12) / COUNT(l.id) as 'hit_rate_12',
			SUM(s.hit_rate_13) / COUNT(l.id) as 'hit_rate_13',
			CAST(SUM(s.hit_count) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_hit_count'
		FROM statistics_bermuda_triangle s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 10
		GROUP BY p.id, m.office_id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsBermudaTriangle, 0)
	for rows.Next() {
		s := new(models.StatisticsBermudaTriangle)
		h := make([]*float64, 14)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.OfficeID, &s.DartsThrown,
			&s.Score, &s.MPR, &s.HighestScoreReached, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8],
			&h[9], &h[10], &h[11], &h[12], &h[13], &s.HitCount)
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 13; i++ {
			hitrates[i] = *h[i]
		}
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// GetBermudaTriangleStatisticsForLeg will return statistics for all players in the given leg
func GetBermudaTriangleStatisticsForLeg(id int) ([]*models.StatisticsBermudaTriangle, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			s.darts_thrown,
			s.score,
			s.mpr,
			s.highest_score_reached,
			s.total_hit_rate,
			s.hit_rate_1,
			s.hit_rate_2,
			s.hit_rate_3,
			s.hit_rate_4,
			s.hit_rate_5,
			s.hit_rate_6,
			s.hit_rate_7,
			s.hit_rate_8,
			s.hit_rate_9,
			s.hit_rate_10,
			s.hit_rate_11,
			s.hit_rate_12,
			s.hit_rate_13,
			s.hit_count
		FROM statistics_bermuda_triangle s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE l.id = ? GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsBermudaTriangle, 0)
	for rows.Next() {
		s := new(models.StatisticsBermudaTriangle)
		h := make([]*float64, 14)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.HighestScoreReached, &s.TotalHitRate,
			&h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10], &h[11], &h[12], &h[13], &s.HitCount)
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 13; i++ {
			hitrates[i] = *h[i]
		}
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// GetBermudaTriangleStatisticsForMatch will return statistics for all players in the given match
func GetBermudaTriangleStatisticsForMatch(id int) ([]*models.StatisticsBermudaTriangle, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			SUM(s.darts_thrown) as 'darts_thrown',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.total_marks) / (COUNT(DISTINCT l.id) * 13) as 'mpr',
			MAX(s.highest_score_reached) as 'highest_score_reached',
			SUM(s.total_hit_rate) / COUNT(l.id) as 'total_hit_rate',
			SUM(s.hit_rate_1) / COUNT(l.id) as 'hit_rate_1',
			SUM(s.hit_rate_2) / COUNT(l.id) as 'hit_rate_2',
			SUM(s.hit_rate_3) / COUNT(l.id) as 'hit_rate_3',
			SUM(s.hit_rate_4) / COUNT(l.id) as 'hit_rate_4',
			SUM(s.hit_rate_5) / COUNT(l.id) as 'hit_rate_5',
			SUM(s.hit_rate_6) / COUNT(l.id) as 'hit_rate_6',
			SUM(s.hit_rate_7) / COUNT(l.id) as 'hit_rate_7',
			SUM(s.hit_rate_8) / COUNT(l.id) as 'hit_rate_8',
			SUM(s.hit_rate_9) / COUNT(l.id) as 'hit_rate_9',
			SUM(s.hit_rate_10) / COUNT(l.id) as 'hit_rate_10',
			SUM(s.hit_rate_11) / COUNT(l.id) as 'hit_rate_11',
			SUM(s.hit_rate_12) / COUNT(l.id) as 'hit_rate_12',
			SUM(s.hit_rate_13) / COUNT(l.id) as 'hit_rate_13',
			CAST(SUM(s.hit_count) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_hit_count'
		FROM statistics_bermuda_triangle s
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

	stats := make([]*models.StatisticsBermudaTriangle, 0)
	for rows.Next() {
		s := new(models.StatisticsBermudaTriangle)
		h := make([]*float64, 14)
		err := rows.Scan(&s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.HighestScoreReached, &s.TotalHitRate,
			&h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10], &h[11], &h[12], &h[13], &s.HitCount)
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 13; i++ {
			hitrates[i] = *h[i]
		}
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// GetBermudaTriangleStatisticsForPlayer will return Bermuda Triangle statistics for the given player
func GetBermudaTriangleStatisticsForPlayer(id int) (*models.StatisticsBermudaTriangle, error) {
	s := new(models.StatisticsBermudaTriangle)
	h := make([]*float64, 26)
	err := models.DB.QueryRow(`
		SELECT
			p.id,
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			SUM(s.darts_thrown) as 'darts_thrown',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.total_marks) / (COUNT(DISTINCT l.id) * 13) as 'mpr',
			MAX(s.highest_score_reached) as 'highest_score_reached',
			SUM(s.total_hit_rate) / COUNT(l.id) as 'total_hit_rate',
			SUM(s.hit_rate_1) / COUNT(l.id) as 'hit_rate_1',
			SUM(s.hit_rate_2) / COUNT(l.id) as 'hit_rate_2',
			SUM(s.hit_rate_3) / COUNT(l.id) as 'hit_rate_3',
			SUM(s.hit_rate_4) / COUNT(l.id) as 'hit_rate_4',
			SUM(s.hit_rate_5) / COUNT(l.id) as 'hit_rate_5',
			SUM(s.hit_rate_6) / COUNT(l.id) as 'hit_rate_6',
			SUM(s.hit_rate_7) / COUNT(l.id) as 'hit_rate_7',
			SUM(s.hit_rate_8) / COUNT(l.id) as 'hit_rate_8',
			SUM(s.hit_rate_9) / COUNT(l.id) as 'hit_rate_9',
			SUM(s.hit_rate_10) / COUNT(l.id) as 'hit_rate_10',
			SUM(s.hit_rate_11) / COUNT(l.id) as 'hit_rate_11',
			SUM(s.hit_rate_12) / COUNT(l.id) as 'hit_rate_12',
			SUM(s.hit_rate_13) / COUNT(l.id) as 'hit_rate_13',
			CAST(SUM(s.hit_count) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_hit_count'
		FROM statistics_bermuda_triangle s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 10
		GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.DartsThrown,
		&s.Score, &s.MPR, &s.HighestScoreReached, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8],
		&h[9], &h[10], &h[11], &h[12], &h[13], &s.HitCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.StatisticsBermudaTriangle), nil
		}
		return nil, err
	}
	hitrates := make(map[int]float64)
	for i := 1; i <= 13; i++ {
		hitrates[i] = *h[i]
	}
	s.Hitrates = hitrates
	return s, nil
}

// GetBermudaTriangleHistoryForPlayer will return history of Bermuda Triangle statistics for the given player
func GetBermudaTriangleHistoryForPlayer(id int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.BERMUDATRIANGLE, false)
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
			s.highest_score_reached,
			s.total_hit_rate,
			s.hit_rate_1,
			s.hit_rate_2,
			s.hit_rate_3,
			s.hit_rate_4,
			s.hit_rate_5,
			s.hit_rate_6,
			s.hit_rate_7,
			s.hit_rate_8,
			s.hit_rate_9,
			s.hit_rate_10,
			s.hit_rate_11,
			s.hit_rate_12,
			s.hit_rate_13,
			s.hit_count
		FROM statistics_bermuda_triangle s
			LEFT JOIN player p ON p.id = s.player_id
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 10
		ORDER BY l.id DESC
		LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.StatisticsBermudaTriangle)
		h := make([]*float64, 14)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.HighestScoreReached, &s.TotalHitRate,
			&h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10], &h[11], &h[12], &h[13], &s.HitCount)
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 13; i++ {
			hitrates[i] = *h[i]
		}
		s.Hitrates = hitrates

		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// CalculateBermudaTriangleStatistics will generate Bermuda Triangle statistics for the given leg
func CalculateBermudaTriangleStatistics(legID int) (map[int]*models.StatisticsBermudaTriangle, error) {
	leg, err := GetLeg(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}

	statisticsMap := make(map[int]*models.StatisticsBermudaTriangle)
	for _, player := range players {
		stats := new(models.StatisticsBermudaTriangle)
		stats.PlayerID = player.PlayerID
		stats.Score = 0
		stats.Hitrates = make(map[int]float64)
		for i := 0; i < 13; i++ {
			stats.Hitrates[i] = 0
		}

		statisticsMap[player.PlayerID] = stats
	}

	round := 0
	for i, visit := range leg.Visits {
		if i > 0 && i%len(players) == 0 {
			round++
		}
		stats := statisticsMap[visit.PlayerID]

		target := models.TargetsBermudaTriangle[round]
		score := visit.CalculateBermudaTriangleScore(round)
		if score == 0 {
			stats.Score = stats.Score / 2
		} else {
			stats.Score += score
		}

		marks := 0
		hits := 0
		if visit.FirstDart.GetBermudaTriangleScore(target) > 0 {
			marks += int(visit.FirstDart.Multiplier)
			hits++
		}
		if visit.SecondDart.GetBermudaTriangleScore(target) > 0 {
			marks += int(visit.SecondDart.Multiplier)
			hits++
		}
		if visit.ThirdDart.GetBermudaTriangleScore(target) > 0 {
			marks += int(visit.ThirdDart.Multiplier)
			hits++
		}
		stats.Hitrates[round] = float64(hits) / 3.0

		stats.TotalMarks += marks
		stats.HitCount += hits
		stats.DartsThrown = visit.DartsThrown

		if stats.Score > stats.HighestScoreReached {
			stats.HighestScoreReached = stats.Score
		}
	}

	for _, stats := range statisticsMap {
		stats.MPR = float64(stats.TotalMarks) / 13.0
		stats.TotalHitRate = float64(stats.HitCount) / 39.0

	}
	return statisticsMap, nil
}

// RecalculateBermudaTriangleStatistics will recaulcate statistics for Bermuda Triangle legs
func RecalculateBermudaTriangleStatistics(legs []int) ([]string, error) {
	queries := make([]string, 0)
	for _, legID := range legs {
		stats, err := CalculateBermudaTriangleStatistics(legID)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			queries = append(queries, fmt.Sprintf(`UPDATE statistics_bermuda_triangle SET darts_thrown = %d, score = %d, mpr = %f, total_marks = %d, highest_score_reached = 22%d4, total_hit_rate = %f, hit_rate_1 = %f, hit_rate_2 = %f, hit_rate_3 = %f, hit_rate_4 = %f, hit_rate_5 = %f, hit_rate_6 = %f, hit_rate_7 = %f, hit_rate_8 = %f, hit_rate_9 = %f, hit_rate_10 = %f, hit_rate_11 = %f, hit_rate_12 = %f, hit_rate_13 = %f, hit_count = %d WHERE leg_id = %d AND player_id = %d;`,
				stat.DartsThrown, stat.Score, stat.MPR, stat.TotalMarks, stat.HighestScoreReached, stat.TotalHitRate, stat.Hitrates[0], stat.Hitrates[1], stat.Hitrates[2], stat.Hitrates[3], stat.Hitrates[4],
				stat.Hitrates[5], stat.Hitrates[6], stat.Hitrates[7], stat.Hitrates[8], stat.Hitrates[9], stat.Hitrates[10], stat.Hitrates[11], stat.Hitrates[12], stat.HitCount, legID, playerID))
		}
	}
	return queries, nil
}
