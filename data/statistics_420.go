package data

import (
	"database/sql"
	"fmt"

	"github.com/kcapp/api/models"
)

// Get420Statistics will return statistics for all players active during the given period
func Get420Statistics(from string, to string) ([]*models.Statistics420, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			m.office_id AS 'office_id',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
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
			SUM(s.hit_rate_14) / COUNT(l.id) as 'hit_rate_14',
			SUM(s.hit_rate_15) / COUNT(l.id) as 'hit_rate_15',
			SUM(s.hit_rate_16) / COUNT(l.id) as 'hit_rate_16',
			SUM(s.hit_rate_17) / COUNT(l.id) as 'hit_rate_17',
			SUM(s.hit_rate_18) / COUNT(l.id) as 'hit_rate_18',
			SUM(s.hit_rate_19) / COUNT(l.id) as 'hit_rate_19',
			SUM(s.hit_rate_20) / COUNT(l.id) as 'hit_rate_20',
			SUM(s.hit_rate_bull) / COUNT(l.id) as 'hit_rate_bull'
		FROM statistics_420 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
			AND m.match_type_id = 11
		GROUP BY p.id, m.office_id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.Statistics420, 0)
	for rows.Next() {
		s := new(models.Statistics420)
		h := make([]*float64, 22)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.OfficeID, &s.Score, &s.TotalHitRate,
			&h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10], &h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17],
			&h[18], &h[19], &h[20], &h[21])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		hitrates[25] = *h[21]
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// Get420StatisticsForLeg will return statistics for all players in the given leg
func Get420StatisticsForLeg(id int) ([]*models.Statistics420, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			s.score,
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
			s.hit_rate_14,
			s.hit_rate_15,
			s.hit_rate_16,
			s.hit_rate_17,
			s.hit_rate_18,
			s.hit_rate_19,
			s.hit_rate_20,
			s.hit_rate_bull
		FROM statistics_420 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE l.id = ? GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.Statistics420, 0)
	for rows.Next() {
		s := new(models.Statistics420)
		h := make([]*float64, 22)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.Score, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9],
			&h[10], &h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20], &h[21])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		hitrates[25] = *h[21]
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// Get420StatisticsForMatch will return statistics for all players in the given match
func Get420StatisticsForMatch(id int) ([]*models.Statistics420, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
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
			SUM(s.hit_rate_14) / COUNT(l.id) as 'hit_rate_14',
			SUM(s.hit_rate_15) / COUNT(l.id) as 'hit_rate_15',
			SUM(s.hit_rate_16) / COUNT(l.id) as 'hit_rate_16',
			SUM(s.hit_rate_17) / COUNT(l.id) as 'hit_rate_17',
			SUM(s.hit_rate_18) / COUNT(l.id) as 'hit_rate_18',
			SUM(s.hit_rate_19) / COUNT(l.id) as 'hit_rate_19',
			SUM(s.hit_rate_20) / COUNT(l.id) as 'hit_rate_20',
			SUM(s.hit_rate_bull) / COUNT(l.id) as 'hit_rate_bull'
		FROM statistics_420 s
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

	stats := make([]*models.Statistics420, 0)
	for rows.Next() {
		s := new(models.Statistics420)
		h := make([]*float64, 22)
		err := rows.Scan(&s.PlayerID, &s.Score, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9],
			&h[10], &h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20], &h[21])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		hitrates[25] = *h[21]
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// Get420StatisticsForPlayer will return 420 statistics for the given player
func Get420StatisticsForPlayer(id int) (*models.Statistics420, error) {
	s := new(models.Statistics420)
	h := make([]*float64, 22)
	err := models.DB.QueryRow(`
		SELECT
			p.id,
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
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
			SUM(s.hit_rate_14) / COUNT(l.id) as 'hit_rate_14',
			SUM(s.hit_rate_15) / COUNT(l.id) as 'hit_rate_15',
			SUM(s.hit_rate_16) / COUNT(l.id) as 'hit_rate_16',
			SUM(s.hit_rate_17) / COUNT(l.id) as 'hit_rate_17',
			SUM(s.hit_rate_18) / COUNT(l.id) as 'hit_rate_18',
			SUM(s.hit_rate_19) / COUNT(l.id) as 'hit_rate_19',
			SUM(s.hit_rate_20) / COUNT(l.id) as 'hit_rate_20',
			SUM(s.hit_rate_bull) / COUNT(l.id) as 'hit_rate_bull'
		FROM statistics_420 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
			AND m.match_type_id = 11
		GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.Score,
		&s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10], &h[11], &h[12], &h[13],
		&h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20], &h[21])
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.Statistics420), nil
		}
		return nil, err
	}
	hitrates := make(map[int]float64)
	for i := 1; i <= 20; i++ {
		hitrates[i] = *h[i]
	}
	hitrates[25] = *h[21]
	s.Hitrates = hitrates
	return s, nil
}

// Get420HistoryForPlayer will return history of 420 statistics for the given player
func Get420HistoryForPlayer(id int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.FOURTWENTY, false)
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
			s.score,
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
			s.hit_rate_14,
			s.hit_rate_15,
			s.hit_rate_16,
			s.hit_rate_17,
			s.hit_rate_18,
			s.hit_rate_19,
			s.hit_rate_20,
			s.hit_rate_bull
		FROM statistics_420 s
			LEFT JOIN player p ON p.id = s.player_id
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 11
		ORDER BY l.id DESC
		LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.Statistics420)
		h := make([]*float64, 22)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.Score, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5],
			&h[6], &h[7], &h[8], &h[9], &h[10], &h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18],
			&h[19], &h[20], &h[21])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		hitrates[25] = *h[21]
		s.Hitrates = hitrates

		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// Calculate420Statistics will generate 420 statistics for the given leg
func Calculate420Statistics(legID int) (map[int]*models.Statistics420, error) {
	leg, err := GetLeg(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}

	statisticsMap := make(map[int]*models.Statistics420)
	for _, player := range players {
		stats := new(models.Statistics420)
		stats.PlayerID = player.PlayerID
		stats.Score = 420
		stats.Hitrates = make(map[int]float64)
		for i := 1; i <= 20; i++ {
			stats.Hitrates[i] = 0
		}
		stats.Hitrates[25] = 0

		statisticsMap[player.PlayerID] = stats
	}

	round := 0
	for i, visit := range leg.Visits {
		if i > 0 && i%len(players) == 0 {
			round++
		}
		stats := statisticsMap[visit.PlayerID]
		target := models.Targets420[round]

		score := visit.Calculate420Score(round)
		stats.Score -= score

		hits := 0.0
		if visit.FirstDart.Get420Score(target) > 0 {
			hits++
		}
		if visit.SecondDart.Get420Score(target) > 0 {
			hits++
		}
		if visit.ThirdDart.Get420Score(target) > 0 {
			hits++
		}
		stats.Hitrates[target.Value] = float64(hits) / 3.0
		stats.TotalHitRate += hits
	}

	for _, stats := range statisticsMap {
		stats.TotalHitRate = float64(stats.TotalHitRate) / 60.0

	}
	return statisticsMap, nil
}

// Recalculate420Statistics will recaulcate statistics for 420 legs
func Recalculate420Statistics(legs []int) ([]string, error) {
	queries := make([]string, 0)
	for _, legID := range legs {
		stats, err := Calculate420Statistics(legID)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			queries = append(queries, fmt.Sprintf(`UPDATE statistics_420 SET score = %d, total_hit_rate = %f, hit_rate_1 = %f, hit_rate_2 = %f, hit_rate_3 = %f, hit_rate_4 = %f, hit_rate_5 = %f, hit_rate_6 = %f, hit_rate_7 = %f, hit_rate_8 = %f, hit_rate_9 = %f, hit_rate_10 = %f, hit_rate_11 = %f, hit_rate_12 = %f, hit_rate_13 = %f, hit_rate_14 = %f, hit_rate_15 = %f, hit_rate_16 = %f, hit_rate_17 = %f, hit_rate_18 = %f, hit_rate_19 = %f, hit_rate_20 = %f, hit_rate_bull = %f WHERE leg_id = %d AND player_id = %d;`,
				stat.Score, stat.TotalHitRate, stat.Hitrates[1], stat.Hitrates[2], stat.Hitrates[3], stat.Hitrates[4], stat.Hitrates[5], stat.Hitrates[6], stat.Hitrates[7], stat.Hitrates[8], stat.Hitrates[9], stat.Hitrates[10],
				stat.Hitrates[11], stat.Hitrates[12], stat.Hitrates[13], stat.Hitrates[14], stat.Hitrates[15], stat.Hitrates[16], stat.Hitrates[17], stat.Hitrates[18], stat.Hitrates[19], stat.Hitrates[20], stat.Hitrates[25], legID, playerID))
		}
	}
	return queries, nil
}
