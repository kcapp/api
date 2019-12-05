package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/kcapp/api/models"
)

// GetX01Statistics will return statistics for all players active duing the given period
func GetX01Statistics(from string, to string, startingScores ...int) ([]*models.StatisticsX01, error) {
	q, args, err := sqlx.In(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT m.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			SUM(s.ppd_score) / SUM(s.darts_thrown) AS 'ppd',
			SUM(s.first_nine_ppd) / (COUNT(p.id)) AS 'first_nine_ppd',
			(SUM(s.ppd_score) / SUM(s.darts_thrown)) * 3 AS 'three_dart_avg',
			SUM(s.first_nine_ppd) / COUNT(p.id) * 3 AS 'first_nine_three_dart_avg',
			SUM(60s_plus) AS '60s_plus',
			SUM(100s_plus) AS '100s_plus',
			SUM(140s_plus) AS '140s_plus',
			SUM(180s) AS '180s',
			SUM(accuracy_20) / COUNT(accuracy_20) AS 'accuracy_20s',
			SUM(accuracy_19) / COUNT(accuracy_19) AS 'accuracy_19s',
			SUM(overall_accuracy) / COUNT(overall_accuracy) AS 'accuracy_overall',
			COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100 AS 'checkout_percentage'
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND l.starting_score IN (?)
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_practice = 0
			AND m.match_type_id = 1
		GROUP BY p.id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC,
			(COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100) DESC`, from, to, startingScores)
	if err != nil {
		return nil, err
	}
	rows, err := models.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsX01, 0)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.PPD, &s.FirstNinePPD, &s.ThreeDartAvg,
			&s.FirstNineThreeDartAvg, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus, &s.Score180s, &s.Accuracy20, &s.Accuracy19,
			&s.AccuracyOverall, &s.CheckoutPercentage)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetX01StatisticsForLeg will return statistics for all players in the given leg
func GetX01StatisticsForLeg(id int) ([]*models.StatisticsX01, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id AS 'leg_id',
			p.id AS 'player_id',
			s.ppd_score / s.darts_thrown,
			s.first_nine_ppd,
			s.ppd_score / s.darts_thrown * 3,
			s.first_nine_ppd * 3,
			s.60s_plus,
			s.100s_plus,
			s.140s_plus,
			s.180s,
			s.accuracy_20,
			s.accuracy_19,
			s.overall_accuracy,
			s.darts_thrown,
			s.checkout_attempts,
			IFNULL(s.checkout_percentage, 0) AS 'checkout_percentage'
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			JOIN player2leg p2l ON p2l.leg_id = l.id AND p2l.player_id = s.player_id
		WHERE l.id = ?
			AND m.match_type_id IN (1,3)
		GROUP BY p.id
		ORDER BY p2l.order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsX01, 0)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.PPD, &s.FirstNinePPD, &s.ThreeDartAvg, &s.FirstNineThreeDartAvg, &s.Score60sPlus, &s.Score100sPlus,
			&s.Score140sPlus, &s.Score180s, &s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &s.DartsThrown,
			&s.CheckoutAttempts, &s.CheckoutPercentage)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetX01StatisticsForMatch will return statistics for all players in the given match
func GetX01StatisticsForMatch(id int) ([]*models.StatisticsX01, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			SUM(s.ppd_score) / SUM(s.darts_thrown) AS 'ppd',
			SUM(s.first_nine_ppd) / COUNT(p.id) AS 'first_nine_ppd',
			(SUM(s.ppd_score) / SUM(s.darts_thrown)) * 3 as 'three_dart_avg',
			SUM(s.first_nine_ppd) / COUNT(p.id) * 3 as 'first_nine_three_dart_avg',
			SUM(s.60s_plus) AS '60s_plus',
			SUM(s.100s_plus) AS '100s_plus',
			SUM(s.140s_plus) AS '140s_plus',
			SUM(s.180s) AS '180s',
			SUM(s.accuracy_20) / COUNT(s.accuracy_20) AS 'accuracy_20s',
			SUM(s.accuracy_19) / COUNT(s.accuracy_19) AS 'accuracy_19s',
			SUM(s.overall_accuracy) / COUNT(s.overall_accuracy) AS 'accuracy_overall',
			SUM(s.checkout_attempts) AS 'checkout_attempts',
			COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100 AS 'checkout_percentage'
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			JOIN player2leg p2l ON p2l.leg_id = l.id AND p2l.player_id = s.player_id
		WHERE m.id = ?
			AND m.match_type_id IN (1, 3)
		GROUP BY p.id
		ORDER BY p2l.order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsX01, 0)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.PPD, &s.FirstNinePPD, &s.ThreeDartAvg, &s.FirstNineThreeDartAvg, &s.Score60sPlus,
			&s.Score100sPlus, &s.Score140sPlus, &s.Score180s, &s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &s.CheckoutAttempts,
			&s.CheckoutPercentage)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetPlayerX01Statistics will get statistics about the given player id
func GetPlayerX01Statistics(id int) (*models.StatisticsX01, error) {
	ids := []int{id}
	statistics, err := GetPlayersX01Statistics(ids)
	if err != nil {
		return nil, err
	}
	if len(statistics) > 0 {
		return statistics[0], nil
	}
	return new(models.StatisticsX01), nil
}

// GetPlayerX01PreviousStatistics will get statistics about the given player id
func GetPlayerX01PreviousStatistics(id int) (*models.StatisticsX01, error) {
	ids := []int{id}
	statistics, err := GetPlayersX01PreviousStatistics(ids)
	if err != nil {
		return nil, err
	}
	if len(statistics) > 0 {
		stats := statistics[0]
		if err != nil {
			return nil, err
		}

		return stats, nil
	}
	return new(models.StatisticsX01), nil
}

// GetPlayersX01Statistics will get statistics about all the the given player IDs
func GetPlayersX01Statistics(ids []int, startingScores ...int) ([]*models.StatisticsX01, error) {
	if len(startingScores) == 0 {
		startingScores = []int{301, 501, 701}
	}
	q, args, err := sqlx.In(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			SUM(s.ppd_score) / SUM(s.darts_thrown) AS 'ppd',
			SUM(s.first_nine_ppd) / COUNT(p.id) AS 'first_nine_ppd',
			(SUM(s.ppd_score) / SUM(s.darts_thrown)) * 3 AS 'three_dart_avg',
			(SUM(s.first_nine_ppd) / COUNT(p.id)) * 3 AS 'first_nine_three_dart_avg',
			SUM(s.60s_plus) AS '60s_plus',
			SUM(s.100s_plus) AS '100s_plus',
			SUM(s.140s_plus) AS '140s_plus',
			SUM(s.180s) AS '180s',
			SUM(s.accuracy_20) / COUNT(s.accuracy_20) AS 'accuracy_20s',
			SUM(s.accuracy_19) / COUNT(s.accuracy_19) AS 'accuracy_19s',
			SUM(s.overall_accuracy) / COUNT(s.overall_accuracy) AS 'accuracy_overall',
			COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100 AS 'checkout_percentage'
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l2.match_id AND l2.winner_id = p.id
		WHERE s.player_id IN (?)
			AND l.starting_score IN (?)
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_practice = 0
			AND m.match_type_id = 1
		GROUP BY s.player_id
		ORDER BY p.id`, ids, startingScores)
	if err != nil {
		return nil, err
	}

	rows, err := models.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statisticsMap := make(map[int]*models.StatisticsX01)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.PPD, &s.FirstNinePPD,
			&s.ThreeDartAvg, &s.FirstNineThreeDartAvg, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus, &s.Score180s,
			&s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &s.CheckoutPercentage)
		if err != nil {
			return nil, err
		}
		statisticsMap[s.PlayerID] = s
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate Best PPD, Best First 9, Best 301 and Best 501
	if len(statisticsMap) > 0 {
		err = getBestStatistics(ids, statisticsMap, startingScores...)
		if err != nil {
			return nil, err
		}
		err = getHighestCheckout(ids, statisticsMap, startingScores...)
		if err != nil {
			return nil, err
		}

		for id, stats := range statisticsMap {
			visits, err := GetPlayerVisits(id)
			if err != nil {
				return nil, err
			}
			stats.Hits, stats.DartsThrown = models.GetHitsMap(visits)
		}
	}
	statistics := make([]*models.StatisticsX01, 0)
	for _, s := range statisticsMap {
		statistics = append(statistics, s)
	}
	return statistics, nil
}

// GetPlayersX01PreviousStatistics will get statistics about all the the given player IDs
func GetPlayersX01PreviousStatistics(ids []int, startingScores ...int) ([]*models.StatisticsX01, error) {
	if len(startingScores) == 0 {
		startingScores = []int{301, 501, 701}
	}
	q, args, err := sqlx.In(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			SUM(s.ppd_score) / SUM(s.darts_thrown) AS 'ppd',
			SUM(s.first_nine_ppd) / COUNT(p.id) AS 'first_nine_ppd',
			(SUM(s.ppd_score) / SUM(s.darts_thrown)) * 3 AS 'three_dart_avg',
			SUM(s.first_nine_ppd) / COUNT(p.id) * 3 AS 'first_nine_three_dart_avg',
			SUM(s.60s_plus) AS '60s_plus',
			SUM(s.100s_plus) AS '100s_plus',
			SUM(s.140s_plus) AS '140s_plus',
			SUM(s.180s) AS '180s',
			SUM(s.accuracy_20) / COUNT(s.accuracy_20) AS 'accuracy_20s',
			SUM(s.accuracy_19) / COUNT(s.accuracy_19) AS 'accuracy_19s',
			SUM(s.overall_accuracy) / COUNT(s.overall_accuracy) AS 'accuracy_overall',
			COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100 AS 'checkout_percentage'
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l2.match_id AND l2.winner_id = p.id
		WHERE s.player_id IN (?)
			AND l.starting_score IN (?)
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_practice = 0
			AND m.match_type_id = 1
			-- Exclude all matches played this week
			AND m.updated_at < (CURRENT_DATE - INTERVAL WEEKDAY(CURRENT_DATE) DAY)
		GROUP BY s.player_id
		ORDER BY p.id`, ids, startingScores)
	if err != nil {
		return nil, err
	}

	rows, err := models.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statisticsMap := make(map[int]*models.StatisticsX01)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.PPD, &s.FirstNinePPD,
			&s.ThreeDartAvg, &s.FirstNineThreeDartAvg, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus, &s.Score180s,
			&s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &s.CheckoutPercentage)
		if err != nil {
			return nil, err
		}
		statisticsMap[s.PlayerID] = s
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate Best PPD, Best First 9, Best 301 and Best 501
	if len(statisticsMap) > 0 {
		err = getBestStatistics(ids, statisticsMap, startingScores...)
		if err != nil {
			return nil, err
		}
		err = getHighestCheckout(ids, statisticsMap, startingScores...)
		if err != nil {
			return nil, err
		}
	}
	statistics := make([]*models.StatisticsX01, 0)
	for _, s := range statisticsMap {
		statistics = append(statistics, s)
	}
	return statistics, nil
}

// GetPlayerProgression will get progression of statistics over time for the given player
func GetPlayerProgression(id int) (map[string]*models.StatisticsX01, error) {
	rows, err := models.DB.Query(`
		SELECT
			s.player_id AS 'player_id',
			SUM(s.ppd_score) / SUM(s.darts_thrown) AS 'ppd',
			SUM(s.first_nine_ppd) / COUNT(s.player_id) AS 'first_nine_ppd',
			(SUM(s.ppd_score) / SUM(s.darts_thrown)) * 3 AS 'three_dart_avg',
			SUM(s.first_nine_ppd) / COUNT(s.player_id) * 3 AS 'first_nine_three_dart_avg',
			SUM(s.60s_plus) AS '60s_plus',
			SUM(s.100s_plus) AS '100s_plus',
			SUM(s.140s_plus) AS '140s_plus',
			SUM(s.180s) AS '180s',
			SUM(s.accuracy_20) / COUNT(s.accuracy_20) AS 'accuracy_20s',
			SUM(s.accuracy_19) / COUNT(s.accuracy_19) AS 'accuracy_19s',
			SUM(s.overall_accuracy) / COUNT(s.overall_accuracy) AS 'accuracy_overall',
			COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100 AS 'checkout_percentage',
			DATE(m.updated_at) AS 'date'
		FROM statistics_x01 s
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND m.match_type_id = 1
			AND m.is_finished = 1 AND m.is_abandoned = 0 AND m.is_practice = 0
		GROUP BY YEAR(m.updateD_at), WEEK(m.updated_at)
		ORDER BY date DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statisticsMap := make(map[string]*models.StatisticsX01)
	for rows.Next() {
		var date string
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.PPD, &s.FirstNinePPD, &s.ThreeDartAvg, &s.FirstNineThreeDartAvg, &s.Score60sPlus,
			&s.Score100sPlus, &s.Score140sPlus, &s.Score180s, &s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &s.CheckoutPercentage, &date)
		if err != nil {
			return nil, err
		}
		statisticsMap[date] = s
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return statisticsMap, nil
}

// getBestStatistics will calculate Best PPD, Best First 9, Best 301 and Best 501 for the given players
func getBestStatistics(ids []int, statisticsMap map[int]*models.StatisticsX01, startingScores ...int) error {
	q, args, err := sqlx.In(`
		SELECT
			p.id,
			l.winner_id,
			l.id,
			(s.ppd_score * 3) / s.darts_thrown,
			((s.first_nine_ppd_score) * 3 / 9),
			s.checkout_percentage,
			s.darts_thrown,
			l.starting_score
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
		WHERE s.player_id IN (?)
			AND l.starting_score IN (?)`, ids, startingScores)
	if err != nil {
		return err
	}
	rows, err := models.DB.Query(q, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rawStatistics := make([]*models.StatisticsX01, 0)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.WinnerID, &s.LegID, &s.ThreeDartAvg, &s.FirstNineThreeDartAvg, &s.CheckoutPercentage, &s.DartsThrown, &s.StartingScore)
		if err != nil {
			return err
		}
		rawStatistics = append(rawStatistics, s)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	for _, stat := range rawStatistics {
		real := statisticsMap[stat.PlayerID]
		// Only count best statistics when the player actually won the leg
		if stat.WinnerID == stat.PlayerID {
			if stat.StartingScore.Int64 == 301 {
				if real.Best301 == nil {
					real.Best301 = new(models.BestStatistic)
				}
				if stat.DartsThrown < real.Best301.Value || real.Best301.Value == 0 {
					real.Best301.Value = stat.DartsThrown
					real.Best301.LegID = stat.LegID
				}
			}
			if stat.StartingScore.Int64 == 501 {
				if real.Best501 == nil {
					real.Best501 = new(models.BestStatistic)
				}
				if stat.DartsThrown < real.Best501.Value || real.Best501.Value == 0 {
					real.Best501.Value = stat.DartsThrown
					real.Best501.LegID = stat.LegID
				}
			}
			if stat.StartingScore.Int64 == 701 {
				if real.Best701 == nil {
					real.Best701 = new(models.BestStatistic)
				}
				if stat.DartsThrown < real.Best701.Value || real.Best701.Value == 0 {
					real.Best701.Value = stat.DartsThrown
					real.Best701.LegID = stat.LegID
				}
			}
		}
		if real.BestThreeDartAvg == nil {
			real.BestThreeDartAvg = new(models.BestStatisticFloat)
		}
		if stat.ThreeDartAvg > real.BestThreeDartAvg.Value {
			real.BestThreeDartAvg.Value = stat.ThreeDartAvg
			real.BestThreeDartAvg.LegID = stat.LegID
		}
		if real.BestFirstNineAvg == nil {
			real.BestFirstNineAvg = new(models.BestStatisticFloat)
		}
		if stat.FirstNineThreeDartAvg > real.BestFirstNineAvg.Value {
			real.BestFirstNineAvg.Value = stat.FirstNineThreeDartAvg
			real.BestFirstNineAvg.LegID = stat.LegID
		}
	}
	return nil
}

// getHighestCheckout will calculate the highest checkout for the given players
func getHighestCheckout(ids []int, statisticsMap map[int]*models.StatisticsX01, startingScores ...int) error {
	q, args, err := sqlx.In(`
		SELECT
			player_id,
			leg_id,
			MAX(checkout)
		FROM (SELECT
				s.player_id,
				s.leg_id,
				IFNULL(s.first_dart * s.first_dart_multiplier, 0) +
					IFNULL(s.second_dart * s.second_dart_multiplier, 0) +
					IFNULL(s.third_dart * s.third_dart_multiplier, 0) AS 'checkout'
			FROM score s
			JOIN leg l ON l.id = s.leg_id
			WHERE l.winner_id = s.player_id
				AND s.player_id IN (?)
				AND s.id IN (SELECT MAX(s.id) FROM score s JOIN leg l ON l.id = s.leg_id WHERE l.winner_id = s.player_id GROUP BY leg_id)
				AND l.starting_score IN (?)
			GROUP BY s.player_id, s.id
			ORDER BY checkout DESC) checkouts
		GROUP BY player_id`, ids, startingScores)
	if err != nil {
		return err
	}
	rows, err := models.DB.Query(q, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var playerID int
		var legID int
		var checkout int
		err := rows.Scan(&playerID, &legID, &checkout)
		if err != nil {
			return err
		}
		highest := new(models.BestStatistic)
		highest.Value = checkout
		highest.LegID = legID
		statisticsMap[playerID].HighestCheckout = highest
	}
	err = rows.Err()
	return err
}

// GetGlobalStatistics will return global statistics for all legs played
func GetGlobalStatistics() (*models.GlobalStatistics, error) {
	global := new(models.GlobalStatistics)
	err := models.DB.QueryRow(`
		SELECT
			COUNT(s.id) AS 'Fish-n-Chips'
		FROM score s
		WHERE
			first_dart IN (1,20,5) AND first_dart_multiplier = 1 AND
			second_dart IN (1,20,5) AND second_dart_multiplier = 1 AND
			third_dart IN (1,20,5) AND third_dart_multiplier = 1  AND
			((first_dart * first_dart_multiplier) + (second_dart * second_dart_multiplier) +
			(third_dart * third_dart_multiplier) = 26)`).Scan(&global.FishNChips)
	if err != nil {
		return nil, err
	}

	return global, nil
}

// GetOfficeStatistics will return office statistics for the given period
func GetOfficeStatistics(from string, to string) ([]*models.OfficeStatistics, error) {
	rows, err := models.DB.Query(`
		SELECT
			player_id,
			leg_id,
			office_id,
			MAX(checkout) AS 'checkout',
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier
		FROM (
			SELECT s.player_id,
					s.leg_id,
					m.office_id,
					s.first_dart * s.first_dart_multiplier +
						IFNULL(s.second_dart * s.second_dart_multiplier, 0) +
						IFNULL(s.third_dart * s.third_dart_multiplier, 0) AS 'checkout',
					s.first_dart, s.first_dart_multiplier,
					s.second_dart, s.second_dart_multiplier,
					s.third_dart, s.third_dart_multiplier
			FROM score s
					JOIN leg l ON l.id = s.leg_id
					JOIN matches m ON m.id = l.match_id
					JOIN player p ON s.player_id = p.id
			WHERE s.id IN (
				SELECT MAX(id) FROM score
				WHERE leg_id IN (
					SELECT id FROM leg WHERE match_id IN (
						SELECT m.id FROM matches m WHERE m.match_type_id = 1 AND m.is_practice = 0
						AND m.is_finished = 1 AND m.updated_at >= ? AND m.updated_at < ?))
				GROUP BY leg_id)
			ORDER BY checkout DESC, leg_id
		) checkouts
		GROUP BY player_id, checkout
		ORDER BY checkout DESC, leg_id`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.OfficeStatistics, 0)
	for rows.Next() {
		s := new(models.OfficeStatistics)
		first := new(models.Dart)
		second := new(models.Dart)
		third := new(models.Dart)
		err := rows.Scan(&s.PlayerID, &s.LegID, &s.OfficeID, &s.Checkout, &first.Value, &first.Multiplier,
			&second.Value, &second.Multiplier, &third.Value, &third.Multiplier)
		if err != nil {
			return nil, err
		}
		darts := make([]*models.Dart, 0)
		s.Darts = append(darts, first, second, third)
		stats = append(stats, s)
	}
	return stats, nil
}

// GetOfficeStatisticsForOffice will return office statistics for the given office and period
func GetOfficeStatisticsForOffice(officeID int, from string, to string) ([]*models.OfficeStatistics, error) {
	rows, err := models.DB.Query(`
		SELECT
			player_id,
			leg_id,
			office_id,
			MAX(checkout) AS 'checkout',
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier
		FROM (
			SELECT s.player_id,
					s.leg_id,
					m.office_id,
					s.first_dart * s.first_dart_multiplier +
					IFNULL(s.second_dart * s.second_dart_multiplier, 0) +
					IFNULL(s.third_dart * s.third_dart_multiplier, 0) AS 'checkout',
					s.first_dart, s.first_dart_multiplier,
					s.second_dart, s.second_dart_multiplier,
					s.third_dart, s.third_dart_multiplier
			FROM score s
					JOIN leg l ON l.id = s.leg_id
					JOIN matches m ON m.id = l.match_id
					JOIN player p ON s.player_id = p.id
			WHERE s.id IN (
				SELECT MAX(id) FROM score
				WHERE leg_id IN (
					SELECT id FROM leg WHERE match_id IN (
						SELECT m.id FROM matches m WHERE m.office_id = ? AND m.match_type_id = 1
						AND m.is_practice = 0 AND m.is_finished = 1 AND m.updated_at >= ? AND m.updated_at < ?))
				GROUP BY leg_id)
			ORDER BY checkout DESC, leg_id
		) checkouts
		GROUP BY player_id, checkout
		ORDER BY checkout DESC, leg_id`, officeID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.OfficeStatistics, 0)
	for rows.Next() {
		s := new(models.OfficeStatistics)
		first := new(models.Dart)
		second := new(models.Dart)
		third := new(models.Dart)
		err := rows.Scan(&s.PlayerID, &s.LegID, &s.OfficeID, &s.Checkout, &first.Value, &first.Multiplier,
			&second.Value, &second.Multiplier, &third.Value, &third.Multiplier)
		if err != nil {
			return nil, err
		}
		darts := make([]*models.Dart, 0)
		s.Darts = append(darts, first, second, third)
		stats = append(stats, s)
	}
	return stats, nil
}
