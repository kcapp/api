package data

import (
	"database/sql"
	"log"

	"github.com/kcapp/api/models"
)

// GetKillBullStatistics will return statistics for all players active during the given period
func GetKillBullStatistics(from string, to string) ([]*models.StatisticsKillBull, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			m.office_id AS 'office_id',
			SUM(s.darts_thrown) as 'darts_thrown',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.marks3) as 'marks3',
			SUM(s.marks4) as 'marks4',
			SUM(s.marks5) as 'marks5',
			SUM(s.marks6) as 'marks6',
			MAX(s.longest_streak) as 'longest_streak',
			SUM(s.times_busted) as 'times_busted',
			SUM(s.total_hit_rate) / COUNT(l.id) as 'total_hit_rate'
		FROM statistics_kill_bull s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 12
		GROUP BY p.id, m.office_id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsKillBull, 0)
	for rows.Next() {
		s := new(models.StatisticsKillBull)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.OfficeID, &s.DartsThrown, &s.Score, &s.Marks3, &s.Marks4, &s.Marks5, &s.Marks6, &s.LongestStreak, &s.TimesBusted, &s.TotalHitRate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetKillBullStatisticsForLeg will return statistics for all players in the given leg
func GetKillBullStatisticsForLeg(id int) ([]*models.StatisticsKillBull, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			s.darts_thrown,
			s.score,
			s.marks3,
			s.marks4,
			s.marks5,
			s.marks6,
			s.longest_streak,
			s.times_busted,
			s.total_hit_rate
		FROM statistics_kill_bull s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE l.id = ? GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsKillBull, 0)
	for rows.Next() {
		s := new(models.StatisticsKillBull)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.Marks3, &s.Marks4, &s.Marks5, &s.Marks6, &s.LongestStreak, &s.TimesBusted, &s.TotalHitRate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetKillBullStatisticsForMatch will return statistics for all players in the given match
func GetKillBullStatisticsForMatch(id int) ([]*models.StatisticsKillBull, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			SUM(s.darts_thrown) as 'darts_thrown',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.marks3) as 'marks3',
			SUM(s.marks4) as 'marks4',
			SUM(s.marks5) as 'marks5',
			SUM(s.marks6) as 'marks6',
			MAX(s.longest_streak) as 'longest_streak',
			SUM(s.times_busted) as 'times_busted',
			SUM(s.total_hit_rate) / COUNT(l.id) as 'total_hit_rate'
		FROM statistics_kill_bull s
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

	stats := make([]*models.StatisticsKillBull, 0)
	for rows.Next() {
		s := new(models.StatisticsKillBull)
		err := rows.Scan(&s.PlayerID, &s.DartsThrown, &s.Score, &s.Marks3, &s.Marks4, &s.Marks5, &s.Marks6, &s.LongestStreak, &s.TimesBusted, &s.TotalHitRate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetKillBullStatisticsForPlayer will return AtW statistics for the given player
func GetKillBullStatisticsForPlayer(id int) (*models.StatisticsKillBull, error) {
	s := new(models.StatisticsKillBull)
	err := models.DB.QueryRow(`
		SELECT
			p.id,
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			SUM(s.darts_thrown) as 'darts_thrown',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.marks3) as 'marks3',
			SUM(s.marks4) as 'marks4',
			SUM(s.marks5) as 'marks5',
			SUM(s.marks6) as 'marks6',
			MAX(s.longest_streak) as 'longest_streak',
			SUM(s.times_busted) as 'times_busted',
			SUM(s.total_hit_rate) / COUNT(l.id) as 'total_hit_rate'
		FROM statistics_kill_bull s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 12
		GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.DartsThrown, &s.Score, &s.Marks3, &s.Marks4, &s.Marks5, &s.Marks6,
		&s.LongestStreak, &s.TimesBusted, &s.TotalHitRate)
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.StatisticsKillBull), nil
		}
		return nil, err
	}
	return s, nil
}

// GetKillBullHistoryForPlayer will return history of AtW statistics for the given player
func GetKillBullHistoryForPlayer(id int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.KILLBULL, false)
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
			s.marks3,
			s.marks4,
			s.marks5,
			s.marks6,
			s.longest_streak,
			s.times_busted,
			s.total_hit_rate
		FROM statistics_kill_bull s
			LEFT JOIN player p ON p.id = s.player_id
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 12
		ORDER BY l.id DESC
		LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.StatisticsKillBull)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.Marks3, &s.Marks4, &s.Marks5, &s.Marks6, &s.LongestStreak, &s.TimesBusted, &s.TotalHitRate)
		if err != nil {
			return nil, err
		}
		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// CalculateKillBullStatistics will generate around the clock statistics for the given leg
func CalculateKillBullStatistics(legID int) (map[int]*models.StatisticsKillBull, error) {
	leg, err := GetLeg(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}

	statisticsMap := make(map[int]*models.StatisticsKillBull)
	for _, player := range players {
		stats := new(models.StatisticsKillBull)
		stats.PlayerID = player.PlayerID
		stats.Score = player.StartingScore
		statisticsMap[player.PlayerID] = stats
	}

	for _, visit := range leg.Visits {
		stats := statisticsMap[visit.PlayerID]
		player := players[visit.PlayerID]

		score := visit.CalculateKillBullScore()
		if score == 0 {
			if stats.Score < player.StartingScore {
				stats.TimesBusted++
			}
			stats.Score = player.StartingScore
		} else {
			stats.Score -= score
			if stats.Score < 0 {
				stats.Score = 0
			}
		}
		stats.DartsThrown = visit.DartsThrown

		hits := 0
		marks := 0
		if visit.FirstDart.IsBull() {
			hits++
			marks += int(visit.FirstDart.Multiplier)
		}
		if visit.SecondDart.IsBull() {
			hits++
			marks += int(visit.SecondDart.Multiplier)
		}
		if visit.ThirdDart.IsBull() {
			hits++
			marks += int(visit.ThirdDart.Multiplier)
		}
		if marks > 0 {
			stats.CurrentStreak++
		} else {
			if stats.CurrentStreak > stats.LongestStreak {
				stats.LongestStreak = stats.CurrentStreak
			}
			stats.CurrentStreak = 0
		}

		if marks == 3 {
			stats.Marks3++
		} else if marks == 4 {
			stats.Marks4++
		} else if marks == 5 {
			stats.Marks5++
		} else if marks == 6 {
			stats.Marks6++
		}
		stats.TotalHitRate += float64(hits)
	}

	for _, stats := range statisticsMap {
		if stats.CurrentStreak > stats.LongestStreak {
			stats.LongestStreak = stats.CurrentStreak
		}
		stats.TotalHitRate = float64(stats.TotalHitRate) / float64(stats.DartsThrown)

	}
	return statisticsMap, nil
}

// ReCalculateKillBullStatistics will recaulcate statistics for Around the Clock legs
func ReCalculateKillBullStatistics() (map[int]map[int]*models.StatisticsKillBull, error) {
	legs, err := GetLegsOfType(models.FOURTWENTY, true)
	if err != nil {
		return nil, err
	}

	s := make(map[int]map[int]*models.StatisticsKillBull)
	for _, leg := range legs {
		stats, err := CalculateKillBullStatistics(leg.ID)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			log.Printf(`UPDATE statistics_kill_bull SET darts_thrown = %d, score = %d, marks3 = %d, marks4 = %d, marks5 = %d, marks6 = %d, longest_streak = %d, times_busted = %d, total_hit_rate = %f
			WHERE leg_id = %d AND player_id = %d;`, stat.DartsThrown, stat.Score, stat.Marks3, stat.Marks4, stat.Marks5, stat.Marks6, stat.LongestStreak, stat.TimesBusted, stat.TotalHitRate, leg.ID, playerID)
		}
		s[leg.ID] = stats
	}

	return s, err
}
