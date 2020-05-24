package data

import (
	"github.com/guregu/null"
	"github.com/kcapp/api/models"
)

// GetAroundTheClockStatistics will return statistics for all players active during the given period
func GetAroundTheClockStatistics(from string, to string) ([]*models.StatisticsAroundThe, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			SUM(s.darts_thrown) as 'darts_thrown',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			MAX(s.longest_streak) as 'longest_streak',
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
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 8
		GROUP BY p.id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsAroundThe, 0)
	for rows.Next() {
		s := new(models.StatisticsAroundThe)
		h := make([]*float32, 26)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.DartsThrown,
			&s.Score, &s.LongestStreak, &s.TotalHitRate, &h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9],
			&h[10], &h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20], &h[25])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float32)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		hitrates[25] = *h[25]
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// GetAroundTheClockStatisticsForLeg will return statistics for all players in the given leg
func GetAroundTheClockStatisticsForLeg(id int) ([]*models.StatisticsAroundThe, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			s.darts_thrown,
			s.score,
			s.longest_streak,
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
		h := make([]*float32, 26)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.LongestStreak, &s.TotalHitRate,
			&h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9], &h[10], &h[11],
			&h[12], &h[13], &h[14], &h[15], &h[16], &h[17], &h[18], &h[19], &h[20], &h[25])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float32)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		hitrates[25] = *h[25]
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// GetAroundTheClockStatisticsForMatch will return statistics for all players in the given match
func GetAroundTheClockStatisticsForMatch(id int) ([]*models.StatisticsAroundThe, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			SUM(s.darts_thrown) as 'darts_thrown',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			CAST(SUM(s.longest_streak) / COUNT(l.id) AS SIGNED) as 'avg_longest_streak',
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
		h := make([]*float32, 26)
		err := rows.Scan(&s.PlayerID, &s.DartsThrown, &s.Score, &s.LongestStreak, &s.TotalHitRate,
			&h[1], &h[2], &h[3], &h[4], &h[5], &h[6], &h[7], &h[8], &h[9],
			&h[10], &h[11], &h[12], &h[13], &h[14], &h[15], &h[16], &h[17],
			&h[18], &h[19], &h[20], &h[25])
		if err != nil {
			return nil, err
		}
		hitrates := make(map[int]float32)
		for i := 1; i <= 20; i++ {
			hitrates[i] = *h[i]
		}
		hitrates[25] = *h[25]
		s.Hitrates = hitrates
		stats = append(stats, s)
	}
	return stats, nil
}

// CalculateAroundTheClockStatistics will generate around the clock statistics for the given leg
func CalculateAroundTheClockStatistics(legID int) (map[int]*models.StatisticsAroundThe, error) {
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
		stats.Score = 1
		statisticsMap[player.PlayerID] = stats
		stats.Hitrates = make(map[int]float32)
		for i := 1; i <= 21; i++ {
			stats.Hitrates[i] = 0
		}
	}

	for _, visit := range leg.Visits {
		stats := statisticsMap[visit.PlayerID]

		currentScore := stats.Score
		currentScore = isHit(stats, currentScore, visit.FirstDart)
		currentScore = isHit(stats, currentScore, visit.SecondDart)
		currentScore = isHit(stats, currentScore, visit.ThirdDart)

		score := visit.CalculateAroundTheClockScore(stats.Score)
		stats.Score += score
		stats.DartsThrown = visit.DartsThrown
	}

	for _, stats := range statisticsMap {
		stats.Hitrates[25] = stats.Hitrates[21]
		delete(stats.Hitrates, 21)

		for i := 1; i <= 20; i++ {
			if stats.Score <= i {
				stats.Hitrates[i] = 0
			}
			stats.TotalHitRate += stats.Hitrates[i]
		}
		stats.TotalHitRate += stats.Hitrates[25]
		stats.TotalHitRate = stats.TotalHitRate / 21

		if stats.CurrentStreak > stats.LongestStreak.ValueOrZero() {
			stats.LongestStreak = null.IntFrom(stats.CurrentStreak)
		}
	}
	return statisticsMap, nil
}

func isHit(stats *models.StatisticsAroundThe, currentScore int, dart *models.Dart) int {
	target := currentScore + 1
	if (dart.ValueRaw() == target && dart.IsSingle()) || (target == 21 && dart.IsBull()) {
		stats.Hitrates[target] = 1 / (1 + stats.Hitrates[target])
		target++
		stats.CurrentStreak++
	} else {
		stats.Hitrates[target]++
		if stats.CurrentStreak > stats.LongestStreak.ValueOrZero() {
			stats.LongestStreak = null.IntFrom(stats.CurrentStreak)
		}
		stats.CurrentStreak = 0
	}
	return target
}
