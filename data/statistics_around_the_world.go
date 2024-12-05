package data

import (
	"database/sql"
	"fmt"

	"github.com/guregu/null"
	"github.com/kcapp/api/models"
)

// GetAroundTheWorldStatistics will return statistics for all players active during the given period
func GetAroundTheWorldStatistics(from string, to string) ([]*models.StatisticsAroundThe, error) {
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
		FROM statistics_around_the s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
			AND m.match_type_id = 6
		GROUP BY p.id, m.office_id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsAroundThe, 0)
	for rows.Next() {
		s := new(models.StatisticsAroundThe)
		h := make([]*float64, 26)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.OfficeID, &s.DartsThrown,
			&s.Score, &s.MPR, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10],
			&h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20], &h[25])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		hitrates[25] = *h[25]
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// GetAroundTheWorldStatisticsForLeg will return statistics for all players in the given leg
func GetAroundTheWorldStatisticsForLeg(id int) ([]*models.StatisticsAroundThe, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			s.darts_thrown,
			s.score,
			s.mpr,
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
		FROM statistics_around_the s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE l.id = ? GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsAroundThe, 0)
	for rows.Next() {
		s := new(models.StatisticsAroundThe)
		h := make([]*float64, 26)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.TotalHitRate,
			&h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10], &h[11],
			&h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20], &h[25])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		hitrates[25] = *h[25]
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// GetAroundTheWorldStatisticsForMatch will return statistics for all players in the given match
func GetAroundTheWorldStatisticsForMatch(id int) ([]*models.StatisticsAroundThe, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			SUM(s.darts_thrown) as 'darts_thrown',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.mpr) / COUNT(DISTINCT l.id) as 'mpr',
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
			SUM(s.hit_rate_bull) / COUNT(l.id) as 'hit_rate_25'
		FROM statistics_around_the s
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

	stats := make([]*models.StatisticsAroundThe, 0)
	for rows.Next() {
		s := new(models.StatisticsAroundThe)
		h := make([]*float64, 26)
		err := rows.Scan(&s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.TotalHitRate,
			&h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10], &h[11],
			&h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20], &h[25])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		hitrates[25] = *h[25]
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// GetAroundTheWorldStatisticsForPlayer will return Around the World statistics for the given player
func GetAroundTheWorldStatisticsForPlayer(id int) (*models.StatisticsAroundThe, error) {
	s := new(models.StatisticsAroundThe)
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
			SUM(s.mpr) / COUNT(DISTINCT l.id) as 'mpr',
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
		FROM statistics_around_the s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
			AND m.match_type_id = 6
		GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.DartsThrown,
		&s.Score, &s.MPR, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10],
		&h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20], &h[25])
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.StatisticsAroundThe), nil
		}
		return nil, err
	}
	hitrates := make(map[int]float64)
	for i := 1; i <= 20; i++ {
		hitrates[i] = *h[i]
	}
	hitrates[25] = *h[25]
	s.Hitrates = hitrates
	return s, nil
}

// GetAroundTheWorldHistoryForPlayer will return history of Around the World statistics for the given player
func GetAroundTheWorldHistoryForPlayer(id int, start int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.AROUNDTHEWORLD, id, start, limit, false)
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
		FROM statistics_around_the s
			LEFT JOIN player p ON p.id = s.player_id
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 6
		ORDER BY l.id DESC
		LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.StatisticsAroundThe)
		h := make([]*float64, 26)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.TotalHitRate,
			&h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10], &h[11],
			&h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20], &h[25])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		hitrates[25] = *h[25]
		s.Hitrates = hitrates

		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// GetShanghaiStatistics will return statistics for all players active during the given period
func GetShanghaiStatistics(from string, to string) ([]*models.StatisticsAroundThe, error) {
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
			SUM(s.total_hit_rate) / COUNT(l.id) as 'total_hit_rate',
			IFNULL(SUM(s.hit_rate_1) / SUM(IF(shanghai < 1, 0, 1)), 0) as 'hit_rate_1',
			IFNULL(SUM(s.hit_rate_2) / SUM(IF(shanghai < 2, 0, 1)), 0) as 'hit_rate_2',
			IFNULL(SUM(s.hit_rate_3) / SUM(IF(shanghai < 3, 0, 1)), 0) as 'hit_rate_3',
			IFNULL(SUM(s.hit_rate_4) / SUM(IF(shanghai < 4, 0, 1)), 0) as 'hit_rate_4',
			IFNULL(SUM(s.hit_rate_5) / SUM(IF(shanghai < 5, 0, 1)), 0) as 'hit_rate_5',
			IFNULL(SUM(s.hit_rate_6) / SUM(IF(shanghai < 6, 0, 1)), 0) as 'hit_rate_6',
			IFNULL(SUM(s.hit_rate_7) / SUM(IF(shanghai < 7, 0, 1)), 0) as 'hit_rate_7',
			IFNULL(SUM(s.hit_rate_8) / SUM(IF(shanghai < 8, 0, 1)), 0) as 'hit_rate_8',
			IFNULL(SUM(s.hit_rate_9) / SUM(IF(shanghai < 9, 0, 1)), 0) as 'hit_rate_9',
			IFNULL(SUM(s.hit_rate_10) / SUM(IF(shanghai < 10, 0, 1)), 0) as 'hit_rate_10',
			IFNULL(SUM(s.hit_rate_11) / SUM(IF(shanghai < 11, 0, 1)), 0) as 'hit_rate_11',
			IFNULL(SUM(s.hit_rate_12) / SUM(IF(shanghai < 12, 0, 1)), 0) as 'hit_rate_12',
			IFNULL(SUM(s.hit_rate_13) / SUM(IF(shanghai < 13, 0, 1)), 0) as 'hit_rate_13',
			IFNULL(SUM(s.hit_rate_14) / SUM(IF(shanghai < 14, 0, 1)), 0) as 'hit_rate_14',
			IFNULL(SUM(s.hit_rate_15) / SUM(IF(shanghai < 15, 0, 1)), 0) as 'hit_rate_15',
			IFNULL(SUM(s.hit_rate_16) / SUM(IF(shanghai < 16, 0, 1)), 0) as 'hit_rate_16',
			IFNULL(SUM(s.hit_rate_17) / SUM(IF(shanghai < 17, 0, 1)), 0) as 'hit_rate_17',
			IFNULL(SUM(s.hit_rate_18) / SUM(IF(shanghai < 18, 0, 1)), 0) as 'hit_rate_18',
			IFNULL(SUM(s.hit_rate_19) / SUM(IF(shanghai < 19, 0, 1)), 0) as 'hit_rate_19',
			IFNULL(SUM(s.hit_rate_20) / SUM(IF(shanghai < 20, 0, 1)), 0) as 'hit_rate_20'
		FROM statistics_around_the s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
			AND m.match_type_id = 7
		GROUP BY p.id, m.office_id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsAroundThe, 0)
	for rows.Next() {
		s := new(models.StatisticsAroundThe)
		h := make([]*null.Float, 21)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.OfficeID,
			&s.DartsThrown, &s.Score, &s.MPR, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7],
			&h[8], &h[9], &h[10], &h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrate := h[i]
			if hitrate.Valid {
				hitrates[i] = hitrate.Float64
			} else {
				hitrates[i] = 0
			}
		}
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// GetShanghaiStatisticsForLeg will return statistics for all players in the given leg
func GetShanghaiStatisticsForLeg(id int) ([]*models.StatisticsAroundThe, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			s.darts_thrown,
			s.score,
			s.shanghai,
			s.mpr,
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
			s.hit_rate_20
		FROM statistics_around_the s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE l.id = ? GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsAroundThe, 0)
	for rows.Next() {
		s := new(models.StatisticsAroundThe)
		h := make([]*float64, 21)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.Shanghai,
			&s.MPR, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7],
			&h[8], &h[9], &h[10], &h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17],
			&h[18], &h[19], &h[20])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// GetShanghaiStatisticsForMatch will return statistics for all players in the given match
func GetShanghaiStatisticsForMatch(id int) ([]*models.StatisticsAroundThe, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			SUM(s.darts_thrown) as 'darts_thrown',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.mpr) / COUNT(DISTINCT l.id) as 'mpr',
			SUM(s.total_hit_rate) / COUNT(l.id) as 'total_hit_rate',
			MAX(shanghai) as 'shanghai',
			IFNULL(SUM(s.hit_rate_1) / SUM(IF(shanghai < 1, 0, 1)), 0) as 'hit_rate_1',
			IFNULL(SUM(s.hit_rate_2) / SUM(IF(shanghai < 2, 0, 1)), 0) as 'hit_rate_2',
			IFNULL(SUM(s.hit_rate_3) / SUM(IF(shanghai < 3, 0, 1)), 0) as 'hit_rate_3',
			IFNULL(SUM(s.hit_rate_4) / SUM(IF(shanghai < 4, 0, 1)), 0) as 'hit_rate_4',
			IFNULL(SUM(s.hit_rate_5) / SUM(IF(shanghai < 5, 0, 1)), 0) as 'hit_rate_5',
			IFNULL(SUM(s.hit_rate_6) / SUM(IF(shanghai < 6, 0, 1)), 0) as 'hit_rate_6',
			IFNULL(SUM(s.hit_rate_7) / SUM(IF(shanghai < 7, 0, 1)), 0) as 'hit_rate_7',
			IFNULL(SUM(s.hit_rate_8) / SUM(IF(shanghai < 8, 0, 1)), 0) as 'hit_rate_8',
			IFNULL(SUM(s.hit_rate_9) / SUM(IF(shanghai < 9, 0, 1)), 0) as 'hit_rate_9',
			IFNULL(SUM(s.hit_rate_10) / SUM(IF(shanghai < 10, 0, 1)), 0) as 'hit_rate_10',
			IFNULL(SUM(s.hit_rate_11) / SUM(IF(shanghai < 11, 0, 1)), 0) as 'hit_rate_11',
			IFNULL(SUM(s.hit_rate_12) / SUM(IF(shanghai < 12, 0, 1)), 0) as 'hit_rate_12',
			IFNULL(SUM(s.hit_rate_13) / SUM(IF(shanghai < 13, 0, 1)), 0) as 'hit_rate_13',
			IFNULL(SUM(s.hit_rate_14) / SUM(IF(shanghai < 14, 0, 1)), 0) as 'hit_rate_14',
			IFNULL(SUM(s.hit_rate_15) / SUM(IF(shanghai < 15, 0, 1)), 0) as 'hit_rate_15',
			IFNULL(SUM(s.hit_rate_16) / SUM(IF(shanghai < 16, 0, 1)), 0) as 'hit_rate_16',
			IFNULL(SUM(s.hit_rate_17) / SUM(IF(shanghai < 17, 0, 1)), 0) as 'hit_rate_17',
			IFNULL(SUM(s.hit_rate_18) / SUM(IF(shanghai < 18, 0, 1)), 0) as 'hit_rate_18',
			IFNULL(SUM(s.hit_rate_19) / SUM(IF(shanghai < 19, 0, 1)), 0) as 'hit_rate_19',
			IFNULL(SUM(s.hit_rate_20) / SUM(IF(shanghai < 20, 0, 1)), 0) as 'hit_rate_20'
		FROM statistics_around_the s
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

	stats := make([]*models.StatisticsAroundThe, 0)
	for rows.Next() {
		s := new(models.StatisticsAroundThe)
		h := make([]*null.Float, 21)
		err := rows.Scan(&s.PlayerID, &s.DartsThrown, &s.Score, &s.MPR, &s.TotalHitRate, &s.Shanghai,
			&h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10], &h[11],
			&h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrate := h[i]
			if hitrate.Valid {
				hitrates[i] = hitrate.Float64
			} else {
				hitrates[i] = 0
			}
		}
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// GetShanghaiStatisticsForPlayer will return Shanghai statistics for the given player
func GetShanghaiStatisticsForPlayer(id int) (*models.StatisticsAroundThe, error) {
	s := new(models.StatisticsAroundThe)
	h := make([]null.Float, 26)
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
			SUM(s.total_hit_rate) / COUNT(l.id) as 'total_hit_rate',
			SUM(s.hit_rate_1) / SUM(IF(shanghai < 1, 0, 1)) as 'hit_rate_1',
			SUM(s.hit_rate_2) / SUM(IF(shanghai < 2, 0, 1)) as 'hit_rate_2',
			SUM(s.hit_rate_3) / SUM(IF(shanghai < 3, 0, 1)) as 'hit_rate_3',
			SUM(s.hit_rate_4) / SUM(IF(shanghai < 4, 0, 1)) as 'hit_rate_4',
			SUM(s.hit_rate_5) / SUM(IF(shanghai < 5, 0, 1)) as 'hit_rate_5',
			SUM(s.hit_rate_6) / SUM(IF(shanghai < 6, 0, 1)) as 'hit_rate_6',
			SUM(s.hit_rate_7) / SUM(IF(shanghai < 7, 0, 1)) as 'hit_rate_7',
			SUM(s.hit_rate_8) / SUM(IF(shanghai < 8, 0, 1)) as 'hit_rate_8',
			SUM(s.hit_rate_9) / SUM(IF(shanghai < 9, 0, 1)) as 'hit_rate_9',
			SUM(s.hit_rate_10) / SUM(IF(shanghai < 10, 0, 1)) as 'hit_rate_10',
			SUM(s.hit_rate_11) / SUM(IF(shanghai < 11, 0, 1)) as 'hit_rate_11',
			SUM(s.hit_rate_12) / SUM(IF(shanghai < 12, 0, 1)) as 'hit_rate_12',
			SUM(s.hit_rate_13) / SUM(IF(shanghai < 13, 0, 1)) as 'hit_rate_13',
			SUM(s.hit_rate_14) / SUM(IF(shanghai < 14, 0, 1)) as 'hit_rate_14',
			SUM(s.hit_rate_15) / SUM(IF(shanghai < 15, 0, 1)) as 'hit_rate_15',
			SUM(s.hit_rate_16) / SUM(IF(shanghai < 16, 0, 1)) as 'hit_rate_16',
			SUM(s.hit_rate_17) / SUM(IF(shanghai < 17, 0, 1)) as 'hit_rate_17',
			SUM(s.hit_rate_18) / SUM(IF(shanghai < 18, 0, 1)) as 'hit_rate_18',
			SUM(s.hit_rate_19) / SUM(IF(shanghai < 19, 0, 1)) as 'hit_rate_19',
			SUM(s.hit_rate_20) / SUM(IF(shanghai < 20, 0, 1)) as 'hit_rate_20'
		FROM statistics_around_the s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
			AND m.match_type_id = 7
		GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.DartsThrown,
		&s.Score, &s.MPR, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10],
		&h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20])
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.StatisticsAroundThe), nil
		}
		return nil, err
	}
	hitrates := make(map[int]float64)
	for i := 1; i <= 20; i++ {
		hitrate := h[i]
		if hitrate.Valid {
			hitrates[i] = hitrate.Float64
		} else {
			hitrates[i] = 0
		}
	}
	s.Hitrates = hitrates
	return s, nil
}

// GetShanghaiHistoryForPlayer will return history of Shanghai statistics for the given player
func GetShanghaiHistoryForPlayer(id int, start int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.SHANGHAI, id, start, limit, false)
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
			s.shanghai,
			s.mpr,
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
			s.hit_rate_20
		FROM statistics_around_the s
			LEFT JOIN player p ON p.id = s.player_id
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 7
		ORDER BY l.id DESC
		LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.StatisticsAroundThe)
		h := make([]*float64, 26)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.Shanghai,
			&s.MPR, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7],
			&h[8], &h[9], &h[10], &h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17],
			&h[18], &h[19], &h[20])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float64)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		s.Hitrates = hitrates

		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// CalculateAroundTheWorldStatistics will generate Around the World statistics for the given leg
func CalculateAroundTheWorldStatistics(legID int, matchType int) (map[int]*models.StatisticsAroundThe, error) {
	leg, err := GetLeg(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}

	statisticsMap := make(map[int]*models.StatisticsAroundThe)
	for _, player := range players {
		stats := new(models.StatisticsAroundThe)
		stats.PlayerID = player.PlayerID
		stats.Score = 0
		statisticsMap[player.PlayerID] = stats
		stats.Hitrates = make(map[int]float64)
		for i := 1; i <= 21; i++ {
			stats.Hitrates[i] = 0
		}
		stats.MPR = null.FloatFrom(0)
	}

	round := 1
	shanghai := 0
	for i, visit := range leg.Visits {
		if i > 0 && i%len(players) == 0 {
			round++
		}
		stats := statisticsMap[visit.PlayerID]

		if visit.FirstDart.ValueRaw() == round || (round == 21 && visit.FirstDart.IsBull()) {
			stats.Hitrates[round]++
			stats.Marks += visit.FirstDart.Multiplier
		}

		if visit.SecondDart.ValueRaw() == round || (round == 21 && visit.SecondDart.IsBull()) {
			stats.Hitrates[round]++
			stats.Marks += visit.SecondDart.Multiplier
		}

		if visit.ThirdDart.ValueRaw() == round || (round == 21 && visit.ThirdDart.IsBull()) {
			stats.Hitrates[round]++
			stats.Marks += visit.ThirdDart.Multiplier
		}

		if matchType == models.SHANGHAI && visit.IsShanghai() && visit.FirstDart.ValueRaw() == round {
			stats.Shanghai = null.IntFrom(int64(visit.FirstDart.ValueRaw()))
			shanghai = visit.FirstDart.ValueRaw()
		}

		score := visit.CalculateAroundTheWorldScore(round)
		stats.Score += score
		stats.DartsThrown = visit.DartsThrown
	}

	for _, stats := range statisticsMap {
		totalHitRate := float64(0)
		for i := 1; i <= 20; i++ {
			totalHitRate += stats.Hitrates[i]
			stats.Hitrates[i] = stats.Hitrates[i] / 3
		}
		totalHitRate += stats.Hitrates[25]
		stats.Hitrates[25] = stats.Hitrates[21] / 3
		delete(stats.Hitrates, 21)

		if shanghai > 0 {
			stats.TotalHitRate = totalHitRate / float64(shanghai*3)
		} else {
			stats.TotalHitRate = totalHitRate / float64(round*3)
		}
		stats.MPR = null.FloatFrom(float64(stats.Marks) / float64(round))
	}
	return statisticsMap, nil
}

// RecalculateAroundTheWorldStatistics will recaulcate statistics for Around the World legs
func RecalculateAroundTheWorldStatistics(legs []int) ([]string, error) {
	queries := make([]string, 0)
	for _, legID := range legs {
		stats, err := CalculateAroundTheWorldStatistics(legID, models.AROUNDTHEWORLD)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			queries = append(queries, fmt.Sprintf(`UPDATE statistics_around_the SET darts_thrown = %d, score = %d, mpr = %f, total_hit_rate = %f, hit_rate_1 = %f, hit_rate_2 = %f, hit_rate_3 = %f, hit_rate_4 = %f, hit_rate_5 = %f, hit_rate_6 = %f, hit_rate_7 = %f, hit_rate_8 = %f, hit_rate_9 = %f, hit_rate_10 = %f, hit_rate_11 = %f, hit_rate_12 = %f, hit_rate_13 = %f, hit_rate_14 = %f, hit_rate_15 = %f, hit_rate_16 = %f, hit_rate_17 = %f, hit_rate_18 = %f, hit_rate_19 = %f, hit_rate_20 = %f, hit_rate_bull = %f WHERE leg_id = %d AND player_id = %d;`,
				stat.DartsThrown, stat.Score, stat.MPR.Float64, stat.TotalHitRate, stat.Hitrates[1], stat.Hitrates[2], stat.Hitrates[3], stat.Hitrates[4], stat.Hitrates[5],
				stat.Hitrates[6], stat.Hitrates[7], stat.Hitrates[8], stat.Hitrates[9], stat.Hitrates[10], stat.Hitrates[11], stat.Hitrates[12], stat.Hitrates[13],
				stat.Hitrates[14], stat.Hitrates[15], stat.Hitrates[16], stat.Hitrates[17], stat.Hitrates[18], stat.Hitrates[19], stat.Hitrates[20], stat.Hitrates[25],
				legID, playerID))
		}
	}
	return queries, nil
}

// RecalculateShanghaiStatistics will recaulcate statistics for Shanghai legs
func RecalculateShanghaiStatistics(legs []int) ([]string, error) {
	queries := make([]string, 0)
	for _, legID := range legs {
		stats, err := CalculateAroundTheWorldStatistics(legID, models.SHANGHAI)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			query := fmt.Sprintf(`UPDATE statistics_around_the SET darts_thrown = %d, score = %d, mpr = %f, total_hit_rate = %f, hit_rate_1 = %f, hit_rate_2 = %f, hit_rate_3 = %f, hit_rate_4 = %f, hit_rate_5 = %f, hit_rate_6 = %f, hit_rate_7 = %f, hit_rate_8 = %f, hit_rate_9 = %f, hit_rate_10 = %f, hit_rate_11 = %f, hit_rate_12 = %f, hit_rate_13 = %f, hit_rate_14 = %f, hit_rate_15 = %f, hit_rate_16 = %f, hit_rate_17 = %f, hit_rate_18 = %f, hit_rate_19 = %f, hit_rate_20 = %f`,
				stat.DartsThrown, stat.Score, stat.MPR.Float64, stat.TotalHitRate, stat.Hitrates[1], stat.Hitrates[2], stat.Hitrates[3], stat.Hitrates[4], stat.Hitrates[5],
				stat.Hitrates[6], stat.Hitrates[7], stat.Hitrates[8], stat.Hitrates[9], stat.Hitrates[10], stat.Hitrates[11], stat.Hitrates[12], stat.Hitrates[13],
				stat.Hitrates[14], stat.Hitrates[15], stat.Hitrates[16], stat.Hitrates[17], stat.Hitrates[18], stat.Hitrates[19], stat.Hitrates[20])

			if stat.Shanghai.Valid {
				query += fmt.Sprintf(`, shanghai = %d`, stat.Shanghai.Int64)
			}
			query += fmt.Sprintf(" WHERE leg_id = %d AND player_id = %d;", legID, playerID)
			queries = append(queries, query)
		}
	}
	return queries, nil
}
