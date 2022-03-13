package data

import (
	"database/sql"
	"fmt"

	"github.com/guregu/null"
	"github.com/kcapp/api/models"
)

// GetDartsAtXStatistics will return statistics for all players active during the given period
func GetDartsAtXStatistics(from string, to string) ([]*models.StatisticsDartsAtX, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) as 'matches_played',
			COUNT(DISTINCT m2.id) as 'matches_won',
			COUNT(DISTINCT l.id) as 'legs_played',
			COUNT(DISTINCT l2.id) as 'legs_won',
			m.office_id AS 'office_id',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.singles) as 'singles',
			SUM(s.doubles) as 'doubles',
			SUM(s.triples) as 'triples',
			SUM(s.singles + s.doubles + s.triples) / (99 * COUNT(DISTINCT l.id)) as 'hit_rate',
			SUM(s.hits5) as 'hits5',
			SUM(s.hits6) as 'hits6',
			SUM(s.hits7) as 'hits7',
			SUM(s.hits8) as 'hits8',
			SUM(s.hits9) as 'hits9'
		FROM statistics_darts_at_x s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 5
		GROUP BY p.id, m.office_id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC, avg_score DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsDartsAtX, 0)
	for rows.Next() {
		s := new(models.StatisticsDartsAtX)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.OfficeID, &s.AvgScore,
			&s.Singles, &s.Doubles, &s.Triples, &s.HitRate, &s.Hits5, &s.Hits6, &s.Hits7, &s.Hits8, &s.Hits9)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetDartsAtXStatisticsForLeg will return statistics for all players in the given leg
func GetDartsAtXStatisticsForLeg(id int) ([]*models.StatisticsDartsAtX, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			s.score,
			s.singles,
			s.doubles,
			s.triples,
			s.hit_rate,
			s.hits5,
			s.hits6,
			s.hits7,
			s.hits8,
			s.hits9
		FROM statistics_darts_at_x s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE l.id = ? GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsDartsAtX, 0)
	for rows.Next() {
		s := new(models.StatisticsDartsAtX)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.Score, &s.Singles, &s.Doubles, &s.Triples, &s.HitRate, &s.Hits5, &s.Hits6, &s.Hits7, &s.Hits8, &s.Hits9)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetDartsAtXStatisticsForMatch will return statistics for all players in the given match
func GetDartsAtXStatisticsForMatch(id int) ([]*models.StatisticsDartsAtX, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.singles) as 'singles',
			SUM(s.doubles) as 'doubles',
			SUM(s.triples) as 'triples',
			SUM(s.singles + s.doubles + s.triples) / 99 * COUNT(DISTINCT l.id) as 'hit_rate',
			SUM(s.hits5) as 'hits5',
			SUM(s.hits6) as 'hits6',
			SUM(s.hits7) as 'hits7',
			SUM(s.hits8) as 'hits8',
			SUM(s.hits9) as 'hits9'
		FROM statistics_darts_at_x s
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

	stats := make([]*models.StatisticsDartsAtX, 0)
	for rows.Next() {
		s := new(models.StatisticsDartsAtX)
		err := rows.Scan(&s.PlayerID, &s.AvgScore, &s.Singles, &s.Doubles, &s.Triples, &s.HitRate, &s.Hits5, &s.Hits6, &s.Hits7, &s.Hits8, &s.Hits9)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetDartsAtXStatisticsForPlayer will return Darts at X statistics for the given player
func GetDartsAtXStatisticsForPlayer(id int) (*models.StatisticsDartsAtX, error) {
	s := new(models.StatisticsDartsAtX)
	err := models.DB.QueryRow(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) as 'matches_played',
			COUNT(DISTINCT m2.id) as 'matches_won',
			COUNT(DISTINCT l.id) as 'legs_played',
			COUNT(DISTINCT l2.id) as 'legs_won',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.singles) as 'singles',
			SUM(s.doubles) as 'doubles',
			SUM(s.triples) as 'triples',
			SUM(s.singles + s.doubles + s.triples) / (99 * COUNT(DISTINCT l.id)) as 'hit_rate',
			SUM(s.hits5) as 'hits5',
			SUM(s.hits6) as 'hits6',
			SUM(s.hits7) as 'hits7',
			SUM(s.hits8) as 'hits8',
			SUM(s.hits9) as 'hits9'
		FROM statistics_darts_at_x s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 5
		GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.AvgScore,
		&s.Singles, &s.Doubles, &s.Triples, &s.HitRate, &s.Hits5, &s.Hits6, &s.Hits7, &s.Hits8, &s.Hits9)
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.StatisticsDartsAtX), nil
		}
		return nil, err
	}

	rows, err := models.DB.Query(`
		SELECT
			l.starting_score,
			SUM(s.hit_rate) / COUNT(l.id) AS 'hit_rate'
		FROM statistics_darts_at_x s
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
		GROUP BY l.starting_score`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hitrates := make(map[int]float32)
	for i := 1; i <= 20; i++ {
		hitrates[i] = 0
	}
	hitrates[25] = 0
	s.Hitrates = hitrates

	for rows.Next() {
		var target int
		var hitrate float32
		err := rows.Scan(&target, &hitrate)
		if err != nil {
			return nil, err
		}
		s.Hitrates[target] = hitrate
	}
	return s, nil
}

// GetDartsAtXHistoryForPlayer will return history of Darts at X statistics for the given player
func GetDartsAtXHistoryForPlayer(id int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.DARTSATX, false)
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
			s.singles,
			s.doubles,
			s.triples,
			s.hit_rate,
			s.hits5,
			s.hits6,
			s.hits7,
			s.hits8,
			s.hits9
		FROM statistics_darts_at_x s
			LEFT JOIN player p ON p.id = s.player_id
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 5
		ORDER BY l.id DESC
		LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.StatisticsDartsAtX)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.Score, &s.Singles, &s.Doubles, &s.Triples, &s.HitRate, &s.Hits5, &s.Hits6, &s.Hits7, &s.Hits8, &s.Hits9)
		if err != nil {
			return nil, err
		}
		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// CalculateDartsAtXStatistics will generate statistics for the given leg
func CalculateDartsAtXStatistics(legID int) (map[int]*models.StatisticsDartsAtX, error) {
	visits, err := GetLegVisits(legID)
	if err != nil {
		return nil, err
	}

	leg, err := GetLeg(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}
	statisticsMap := make(map[int]*models.StatisticsDartsAtX)
	playerHitsMap := make(map[int]map[int]int64)
	for _, player := range players {
		stats := new(models.StatisticsDartsAtX)
		stats.PlayerID = player.PlayerID
		stats.Score = null.IntFrom(int64(player.CurrentScore))
		statisticsMap[player.PlayerID] = stats
		playerHitsMap[player.PlayerID] = make(map[int]int64)
	}

	number := leg.StartingScore
	for i := 0; i < len(visits); i++ {
		visit := visits[i]
		stats := statisticsMap[visit.PlayerID]

		hits := addDart(number, visit.FirstDart, stats)
		hits += addDart(number, visit.SecondDart, stats)
		hits += addDart(number, visit.ThirdDart, stats)
		switch hits {
		case 5:
			stats.Hits5++
		case 6:
			stats.Hits6++
		case 7:
			stats.Hits7++
		case 8:
			stats.Hits8++
		case 9:
			stats.Hits9++
		}
	}
	for _, stat := range statisticsMap {
		stat.HitRate = float32(stat.Singles+stat.Doubles+stat.Triples) / 99
	}
	return statisticsMap, nil
}

// RecalculateDartsAtXStatistics will recaulcate statistics for Darts at X legs
func RecalculateDartsAtXStatistics(legs []int) ([]string, error) {
	queries := make([]string, 0)
	for _, legID := range legs {
		stats, err := CalculateDartsAtXStatistics(legID)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			queries = append(queries, fmt.Sprintf(`UPDATE statistics_darts_at_x SET score = %d, singles = %d, doubles = %d, triples = %d, hit_rate = %f, hits5 = %d, hits6 = %d, hits7 = %d, hits8 = %d, hits9 = %d WHERE leg_id = %d AND player_id = %d;`,
				stat.Score.Int64, stat.Singles, stat.Doubles, stat.Triples, stat.HitRate, stat.Hits5, stat.Hits6, stat.Hits7, stat.Hits8, stat.Hits9, legID, playerID))
		}
	}
	return queries, nil
}

func addDart(number int, dart *models.Dart, stats *models.StatisticsDartsAtX) int {
	if dart.ValueRaw() == number {
		if dart.IsTriple() {
			stats.Triples++
		} else if dart.IsDouble() {
			stats.Doubles++
		} else {
			stats.Singles++
		}
		return int(dart.Multiplier)
	}
	return 0
}
