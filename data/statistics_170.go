package data

import (
	"database/sql"
	"fmt"

	"github.com/guregu/null"
	"github.com/kcapp/api/models"
)

// Get170Statistics will return statistics for all players active during the given period
func Get170Statistics(from string, to string) ([]*models.StatisticsScam, error) {
	rows, err := models.DB.Query(`
			SELECT
				p.id,
				COUNT(DISTINCT m.id) AS 'matches_played',
				COUNT(DISTINCT m2.id) AS 'matches_won',
				COUNT(DISTINCT l.id) AS 'legs_played',
				COUNT(DISTINCT l2.id) AS 'legs_won',
				m.office_id AS 'office_id',
				SUM(s.darts_thrown_stopper) as 'darts_thrown_stopper',
				SUM(s.darts_thrown_scorer) as 'darts_thrown_scorer',
				SUM(s.score) / SUM(s.darts_thrown_scorer) as 'ppd',
				SUM(s.score) / SUM(s.darts_thrown_scorer) * 3 as 'three_dart_avg',
				CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
				(20 * COUNT(DISTINCT l.id)) / SUM(darts_thrown_stopper) * 3 as 'mpr'
			FROM statistics_scam s
				JOIN player p ON p.id = s.player_id
				JOIN leg l ON l.id = s.leg_id
				JOIN matches m ON m.id = l.match_id
				LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
				LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
			WHERE m.updated_at >= ? AND m.updated_at < ?
				AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
				AND m.match_type_id = 17
			GROUP BY p.id, m.office_id
			ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsScam, 0)
	for rows.Next() {
		s := new(models.StatisticsScam)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.OfficeID, &s.DartsThrownStopper, &s.DartsThrownScorer,
			&s.PPD, &s.ThreeDartAvg, &s.Score, &s.MPR)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// Get170StatisticsForLeg will return statistics for all players in the given leg
func Get170StatisticsForLeg(id int) ([]*models.StatisticsScam, error) {
	rows, err := models.DB.Query(`
			SELECT
				l.id,
				p.id,
				s.darts_thrown_scorer,
				s.darts_thrown_stopper,
				s.score,
				s.mpr,
				s.ppd,
				s.ppd / 3 as 'three_dart_avg'
			FROM statistics_scam s
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
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrownScorer, &s.DartsThrownStopper, &s.Score, &s.MPR, &s.PPD, &s.ThreeDartAvg)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// Get170StatisticsForMatch will return statistics for all players in the given match
func Get170StatisticsForMatch(id int) ([]*models.StatisticsScam, error) {
	rows, err := models.DB.Query(`
			SELECT
				p.id,
				SUM(s.darts_thrown_scorer) as 'darts_thrown_scorer',
				SUM(s.darts_thrown_stopper) as 'darts_thrown_stopper',
				CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
				SUM(s.score) / SUM(s.darts_thrown_scorer) as 'ppd',
				SUM(s.score) / SUM(s.darts_thrown_scorer) * 3 as 'three_dart_avg',
				(20 * COUNT(DISTINCT l.id)) / SUM(darts_thrown_stopper) * 3 as 'mpr'
			FROM statistics_scam s
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
		err := rows.Scan(&s.PlayerID, &s.DartsThrownScorer, &s.DartsThrownStopper, &s.Score, &s.PPD, &s.ThreeDartAvg, &s.MPR)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// Get170StatisticsForPlayer will return Scam statistics for the given player
func Get170StatisticsForPlayer(id int) (*models.StatisticsScam, error) {
	s := new(models.StatisticsScam)
	err := models.DB.QueryRow(`
			SELECT
				p.id,
				COUNT(DISTINCT m.id) AS 'matches_played',
				COUNT(DISTINCT m2.id) AS 'matches_won',
				COUNT(DISTINCT l.id) AS 'legs_played',
				COUNT(DISTINCT l2.id) AS 'legs_won',
				SUM(s.darts_thrown_scorer) as 'darts_thrown_scorer',
				SUM(s.darts_thrown_stopper) as 'darts_thrown_stopper',
				CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
				SUM(darts_thrown_stopper) / 20 * COUNT(DISTINCT l.id) * 3 as 'mpr',
				SUM(s.score) / SUM(s.darts_thrown_scorer) as 'ppd',
				SUM(s.score) / SUM(s.darts_thrown_scorer) * 3 as 'three_dart_avg'
			FROM statistics_scam s
				JOIN player p ON p.id = s.player_id
				JOIN leg l ON l.id = s.leg_id
				JOIN matches m ON m.id = l.match_id
				LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
				LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
			WHERE s.player_id = ?
				AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
				AND m.match_type_id = 16
			GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.DartsThrownScorer, &s.DartsThrownStopper,
		&s.Score, &s.MPR, &s.PPD, &s.ThreeDartAvg)
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.StatisticsScam), nil
		}
		return nil, err
	}
	return s, nil
}

// Get170HistoryForPlayer will return history of Scam statistics for the given player
func Get170HistoryForPlayer(id int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.ONESEVENTY, false)
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
			s.points,
			s.ppd,
			s.ppd_score,
			s.rounds,
			s.checkout_percentage,
			s.checkout_attempts,
			s.checkout_completed,
			s.highest_checkout,
			s.checkout_9_darts,
			s.checkout_8_darts,
			s.checkout_7_darts,
			s.checkout_6_darts,
			s.checkout_5_darts,
			s.checkout_4_darts,
			s.checkout_3_darts
		FROM kcapp.statistics_170 s
			LEFT JOIN player p ON p.id = s.player_id
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 17
		ORDER BY l.id DESC
		LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.Statistics170)
		s.CheckoutDarts = make(map[int]int)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.Points, &s.PPD, &s.PPDScore, &s.Rounds, &s.CheckoutPercentage,
			&s.CheckoutAttempts, &s.CheckoutCompleted, &s.HighestCheckout, s.CheckoutDarts[9], s.CheckoutDarts[8],
			s.CheckoutDarts[7], s.CheckoutDarts[6], s.CheckoutDarts[5], s.CheckoutDarts[4], s.CheckoutDarts[3])
		if err != nil {
			return nil, err
		}
		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// Calculate170Statistics will generate 170 statistics for the given leg
func Calculate170Statistics(legID int) (map[int]*models.Statistics170, error) {
	leg, err := GetLeg(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}

	statisticsMap := make(map[int]*models.Statistics170)
	for _, player := range players {
		stats := new(models.Statistics170)
		stats.PlayerID = player.PlayerID
		stats.Points = int(player.CurrentPoints.Int64)
		stats.Rounds = len(leg.Visits)/len(leg.Players)/3 + 1
		stats.HighestCheckout = null.IntFrom(0)
		stats.CheckoutDarts = make(map[int]int)
		for i := 3; i <= 0; i++ {
			stats.CheckoutDarts[i] = 0
		}

		player.CurrentScore = 170
		player.DartsThrown = 0
		player.CurrentPoints = null.IntFrom(0)
		statisticsMap[player.PlayerID] = stats
	}
	round := 1
	for i, visit := range leg.Visits {
		if i > 0 && i%len(players) == 0 {
			round++
		}
		stats := statisticsMap[visit.PlayerID]
		player := players[visit.PlayerID]

		if !visit.IsBust {
			stats.PPDScore += visit.GetScore()
		}
		stats.DartsThrown = visit.DartsThrown

		stats.DartsThrown += 3
		if !visit.IsBust {
			player.CurrentScore -= visit.GetScore()
		}

		// TODO checkout attempts
		currentScore := player.CurrentScore
		if visit.FirstDart.IsCheckoutAttempt(currentScore, 1, models.OUTSHOTDOUBLE) {
			stats.CheckoutAttempts++
		}
		currentScore -= visit.FirstDart.GetScore()
		if visit.SecondDart.IsCheckoutAttempt(currentScore, 2, models.OUTSHOTDOUBLE) {
			stats.CheckoutAttempts++
		}
		currentScore -= visit.SecondDart.GetScore()
		if visit.ThirdDart.IsCheckoutAttempt(currentScore, 3, models.OUTSHOTDOUBLE) {
			stats.CheckoutAttempts++
		}

		if player.CurrentScore == 0 && visit.GetLastDart().IsDouble() {
			stats.CheckoutCompleted++
			stats.CheckoutDarts[player.DartsThrown]++
			if int(stats.HighestCheckout.Int64) < visit.GetScore() {
				stats.HighestCheckout = null.IntFrom(int64(visit.GetScore()))
			}
			player.CurrentScore = 170
			player.DartsThrown = 0
		} else if round != 1 && player.DartsThrown%9 == 0 {
			// 9 Darts have been thrown, reset
			player.CurrentScore = 170
			player.DartsThrown = 0
		}
	}

	for _, stats := range statisticsMap {
		stats.PPD = float32(stats.PPDScore) / float32(stats.DartsThrown)
		if stats.CheckoutAttempts > 0 {
			stats.CheckoutPercentage = null.FloatFrom(100 / float64(stats.CheckoutAttempts))
		}
		stats.ThreeDartAvg = stats.PPD * 3
	}
	return statisticsMap, nil
}

// ReCalculate170Statistics will recaulcate statistics for 170 legs
func ReCalculate170Statistics(legs []int) ([]string, error) {
	queries := make([]string, 0)
	for _, legID := range legs {
		stats, err := CalculateScamStatistics(legID)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			queries = append(queries, fmt.Sprintf(`UPDATE statistics_170 SET darts_thrown_stopper = %d, darts_thrown_scorer = %d, mpr = %f, score = %d, WHERE leg_id = %d AND player_id = %d;`,
				stat.DartsThrownStopper, stat.DartsThrownScorer, stat.MPR, stat.Score, legID, playerID))
		}
	}

	return queries, nil
}
