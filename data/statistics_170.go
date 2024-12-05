package data

import (
	"database/sql"
	"fmt"

	"github.com/guregu/null"
	"github.com/kcapp/api/models"
)

// Get170Statistics will return statistics for all players active during the given period
func Get170Statistics(from string, to string) ([]*models.Statistics170, error) {
	rows, err := models.DB.Query(`
			SELECT
				p.id,
				COUNT(DISTINCT m.id) AS 'matches_played',
				COUNT(DISTINCT m2.id) AS 'matches_won',
				COUNT(DISTINCT l.id) AS 'legs_played',
				COUNT(DISTINCT l2.id) AS 'legs_won',
				m.office_id AS 'office_id',
				SUM(s.points),
				IF(s.darts_thrown = 0, 0, SUM(s.ppd_score) / SUM(s.darts_thrown)),
				IF(s.darts_thrown = 0, 0, SUM(s.ppd_score) / SUM(s.darts_thrown) * 3),
				SUM(s.rounds),
				COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100,
				SUM(s.checkout_completed),
				SUM(s.checkout_attempts),
				MAX(s.highest_checkout),
				SUM(s.darts_thrown),
				SUM(s.checkout_9_darts),
				SUM(s.checkout_8_darts),
				SUM(s.checkout_7_darts),
				SUM(s.checkout_6_darts),
				SUM(s.checkout_5_darts),
				SUM(s.checkout_4_darts),
				SUM(s.checkout_3_darts)
			FROM statistics_170 s
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

	stats := make([]*models.Statistics170, 0)
	for rows.Next() {
		s := new(models.Statistics170)
		var darts9, darts8, darts7, darts6, darts5, darts4, darts3 int

		checkoutDarts := make(map[int]int, 0)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.Points, &s.PPD, &s.ThreeDartAvg, &s.Rounds, &s.CheckoutPercentage, &s.CheckoutCompleted,
			&s.CheckoutAttempts, &s.HighestCheckout, &s.DartsThrown, &darts9, &darts8, &darts7, &darts6, &darts5, &darts4, &darts3)
		if err != nil {
			return nil, err
		}
		checkoutDarts[9] = darts9
		checkoutDarts[8] = darts8
		checkoutDarts[7] = darts7
		checkoutDarts[6] = darts6
		checkoutDarts[5] = darts5
		checkoutDarts[4] = darts4
		checkoutDarts[3] = darts3
		s.CheckoutDarts = checkoutDarts
		stats = append(stats, s)
	}
	return stats, nil
}

// Get170StatisticsForLeg will return statistics for all players in the given leg
func Get170StatisticsForLeg(id int) ([]*models.Statistics170, error) {
	rows, err := models.DB.Query(`
			SELECT
				l.id,
				p.id,
				s.points,
				s.ppd,
				s.ppd_score / s.darts_thrown * 3,
				s.rounds,
				s.checkout_percentage,
				s.checkout_completed,
				s.checkout_attempts,
				s.highest_checkout,
				s.darts_thrown,
				s.checkout_9_darts,
				s.checkout_8_darts,
				s.checkout_7_darts,
				s.checkout_6_darts,
				s.checkout_5_darts,
				s.checkout_4_darts,
				s.checkout_3_darts
			FROM statistics_170 s
				JOIN player p ON p.id = s.player_id
				JOIN leg l ON l.id = s.leg_id
				JOIN player2leg p2l on l.id = p2l.leg_id AND p.id = p2l.player_id
			WHERE l.id = ? GROUP BY p.id ORDER BY p2l.order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.Statistics170, 0)
	for rows.Next() {
		s := new(models.Statistics170)
		var darts9, darts8, darts7, darts6, darts5, darts4, darts3 int

		checkoutDarts := make(map[int]int, 0)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.Points, &s.PPD, &s.ThreeDartAvg, &s.Rounds, &s.CheckoutPercentage, &s.CheckoutCompleted,
			&s.CheckoutAttempts, &s.HighestCheckout, &s.DartsThrown, &darts9, &darts8, &darts7, &darts6, &darts5, &darts4, &darts3)
		if err != nil {
			return nil, err
		}
		checkoutDarts[9] = darts9
		checkoutDarts[8] = darts8
		checkoutDarts[7] = darts7
		checkoutDarts[6] = darts6
		checkoutDarts[5] = darts5
		checkoutDarts[4] = darts4
		checkoutDarts[3] = darts3
		s.CheckoutDarts = checkoutDarts
		stats = append(stats, s)
	}
	return stats, nil
}

// Get170StatisticsForMatch will return statistics for all players in the given match
func Get170StatisticsForMatch(id int) ([]*models.Statistics170, error) {
	rows, err := models.DB.Query(`
			SELECT
				p.id,
				SUM(s.points),
				IF(s.darts_thrown = 0, 0, SUM(s.ppd_score) / SUM(s.darts_thrown)),
				IF(s.darts_thrown = 0, 0, SUM(s.ppd_score) / SUM(s.darts_thrown) * 3),
				SUM(s.rounds),
				SUM(s.checkout_completed) / SUM(s.checkout_attempts) * 100,
				SUM(s.checkout_completed),
				SUM(s.checkout_attempts),
				MAX(s.highest_checkout),
				SUM(s.darts_thrown),
				SUM(s.checkout_9_darts),
				SUM(s.checkout_8_darts),
				SUM(s.checkout_7_darts),
				SUM(s.checkout_6_darts),
				SUM(s.checkout_5_darts),
				SUM(s.checkout_4_darts),
				SUM(s.checkout_3_darts)
			FROM statistics_170 s
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

	stats := make([]*models.Statistics170, 0)
	for rows.Next() {
		s := new(models.Statistics170)
		var darts9, darts8, darts7, darts6, darts5, darts4, darts3 int

		checkoutDarts := make(map[int]int, 0)
		err := rows.Scan(&s.PlayerID, &s.Points, &s.PPD, &s.ThreeDartAvg, &s.Rounds, &s.CheckoutPercentage, &s.CheckoutCompleted,
			&s.CheckoutAttempts, &s.HighestCheckout, &s.DartsThrown, &darts9, &darts8, &darts7, &darts6, &darts5, &darts4, &darts3)
		if err != nil {
			return nil, err
		}
		checkoutDarts[9] = darts9
		checkoutDarts[8] = darts8
		checkoutDarts[7] = darts7
		checkoutDarts[6] = darts6
		checkoutDarts[5] = darts5
		checkoutDarts[4] = darts4
		checkoutDarts[3] = darts3
		s.CheckoutDarts = checkoutDarts
		stats = append(stats, s)
	}
	return stats, nil
}

// Get170StatisticsForPlayer will return Scam statistics for the given player
func Get170StatisticsForPlayer(id int) (*models.Statistics170, error) {
	s := new(models.Statistics170)
	var darts9, darts8, darts7, darts6, darts5, darts4, darts3 int
	err := models.DB.QueryRow(`
			SELECT
				p.id,
				SUM(s.points),
				SUM(s.ppd_score) / SUM(s.darts_thrown),
				SUM(s.rounds),
				SUM(s.checkout_attempts) / COUNT(DISTINCT l2.id),
				SUM(s.checkout_attempts),
				MAX(s.highest_checkout),
				SUM(s.darts_thrown),
				SUM(s.checkout_9_darts),
				SUM(s.checkout_8_darts),
				SUM(s.checkout_7_darts),
				SUM(s.checkout_6_darts),
				SUM(s.checkout_5_darts),
				SUM(s.checkout_4_darts),
				SUM(s.checkout_3_darts)
			FROM statistics_170 s
				JOIN player p ON p.id = s.player_id
				JOIN leg l ON l.id = s.leg_id
				JOIN matches m ON m.id = l.match_id
				LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
				LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
			WHERE s.player_id = ?
				AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
				AND m.match_type_id = 17
			GROUP BY p.id`, id).Scan(&s.LegID, &s.PlayerID, &s.Points, &s.PPD, &s.Rounds, &s.CheckoutPercentage, &s.CheckoutAttempts,
		&s.HighestCheckout, &s.DartsThrown, &darts9, &darts8, &darts7, &darts6, &darts5, &darts4, &darts3)
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.Statistics170), nil
		}
		return nil, err
	}
	checkoutDarts := make(map[int]int, 0)
	checkoutDarts[9] = darts9
	checkoutDarts[8] = darts8
	checkoutDarts[7] = darts7
	checkoutDarts[6] = darts6
	checkoutDarts[5] = darts5
	checkoutDarts[4] = darts4
	checkoutDarts[3] = darts3
	s.CheckoutDarts = checkoutDarts
	return s, nil
}

// Get170HistoryForPlayer will return history of Scam statistics for the given player
func Get170HistoryForPlayer(id int, start int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.ONESEVENTY, id, start, limit, false)
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
		player.DartsThrown += 3

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

		if !visit.IsBust {
			player.CurrentScore -= visit.GetScore()
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
		if stats.PPDScore == 0 {
			stats.PPD = 0
		} else {
			stats.PPD = float32(stats.PPDScore) / float32(stats.DartsThrown)
		}
		if stats.CheckoutCompleted > 0 {
			stats.CheckoutPercentage = null.FloatFrom(float64(stats.CheckoutCompleted) / float64(stats.CheckoutAttempts) * 100.0)
		}
		stats.ThreeDartAvg = stats.PPD * 3
		stats.Rounds -= 1
	}
	return statisticsMap, nil
}

// ReCalculate170Statistics will recaulcate statistics for 170 legs
func ReCalculate170Statistics(legs []int) ([]string, error) {
	queries := make([]string, 0)
	for _, legID := range legs {
		stats, err := Calculate170Statistics(legID)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			query := fmt.Sprintf(`UPDATE statistics_170 SET points = %d, ppd = %f, ppd_score = %d, rounds = %d, checkout_attempts = %d, 
			checkout_completed = %d, darts_thrown = %d,  checkout_9_darts = %d, checkout_8_darts = %d, 
			checkout_7_darts = %d, checkout_6_darts = %d, checkout_5_darts = %d, checkout_4_darts = %d, checkout_3_darts = %d`,
				stat.Points, stat.PPD, stat.PPDScore, stat.Rounds, stat.CheckoutAttempts, stat.CheckoutCompleted,
				stat.DartsThrown, stat.CheckoutDarts[9], stat.CheckoutDarts[8], stat.CheckoutDarts[7],
				stat.CheckoutDarts[6], stat.CheckoutDarts[5], stat.CheckoutDarts[4], stat.CheckoutDarts[3])

			if stat.CheckoutPercentage.Valid {
				query += fmt.Sprintf(", checkout_percentage = %f", stat.CheckoutPercentage.Float64)
			}
			if stat.HighestCheckout.Valid {
				query += fmt.Sprintf(", highest_checkout = %d", stat.HighestCheckout.Int64)
			}
			query += fmt.Sprintf(" WHERE leg_id = %d AND player_id = %d;", legID, playerID)
			queries = append(queries, query)
		}
	}

	return queries, nil
}
