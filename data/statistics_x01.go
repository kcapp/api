package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/kcapp/api/models"
)

// GetX01Statistics will return statistics for all players active duing the given period
func GetX01Statistics(from string, to string) ([]*models.StatisticsX01, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			COUNT(DISTINCT g.id),
			SUM(s.ppd) / COUNT(p.id),
			SUM(s.first_nine_ppd) / COUNT(p.id),
			SUM(60s_plus),
			SUM(100s_plus),
			SUM(140s_plus),
			SUM(180s) AS '180s',
			SUM(accuracy_20) / COUNT(accuracy_20),
			SUM(accuracy_19) / COUNT(accuracy_19),
			SUM(overall_accuracy) / COUNT(overall_accuracy),
			SUM(checkout_percentage) / COUNT(checkout_percentage)
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN `+"`match`"+` m ON m.id = s.match_id
			JOIN game g ON g.id = m.game_id
		WHERE g.updated_at >= ? AND g.updated_at < ?
		AND g.is_finished = 1
		AND g.game_type_id = 1
		GROUP BY p.id`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statsMap := make(map[int]*models.StatisticsX01, 0)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.GamesPlayed, &s.PPD, &s.FirstNinePPD, &s.Score60sPlus, &s.Score100sPlus,
			&s.Score140sPlus, &s.Score180s, &s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &s.CheckoutPercentage)
		if err != nil {
			return nil, err
		}
		statsMap[s.PlayerID] = s
	}

	rows, err = models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			COUNT(g.winner_id) AS 'games_won'
		FROM game g
			JOIN player p ON p.id = g.winner_id
		WHERE g.updated_at >= ? AND g.updated_at < ?
		AND g.game_type_id = 1
		GROUP BY g.winner_id`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var playerID int
		var gamesWon int
		err := rows.Scan(&playerID, &gamesWon)
		if err != nil {
			return nil, err
		}
		statsMap[playerID].GamesWon = gamesWon
	}

	stats := make([]*models.StatisticsX01, 0)
	for _, s := range statsMap {
		stats = append(stats, s)
	}

	return stats, nil
}

// GetX01StatisticsForMatch will return statistics for all players in the given match
func GetX01StatisticsForMatch(id int) ([]*models.StatisticsX01, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id,
			p.id,
			COUNT(DISTINCT g.id),
			SUM(s.ppd) / COUNT(p.id),
			SUM(s.first_nine_ppd) / COUNT(p.id),
			SUM(60s_plus),
			SUM(100s_plus),
			SUM(140s_plus),
			SUM(180s) AS '180s',
			SUM(accuracy_20) / COUNT(accuracy_20),
			SUM(accuracy_19) / COUNT(accuracy_19),
			SUM(overall_accuracy) / COUNT(overall_accuracy),
			SUM(checkout_percentage) / COUNT(checkout_percentage)
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN `+"`match`"+` m ON m.id = s.match_id
			JOIN game g ON g.id = m.game_id
			JOIN player2match p2m ON p2m.match_id = m.id AND p2m.player_id = s.player_id
		WHERE m.id = ?
		AND g.game_type_id = 1
		GROUP BY p.id
		ORDER BY p2m.order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsX01, 0)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.MatchID, &s.PlayerID, &s.GamesPlayed, &s.PPD, &s.FirstNinePPD, &s.Score60sPlus, &s.Score100sPlus,
			&s.Score140sPlus, &s.Score180s, &s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &s.CheckoutPercentage)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetX01StatisticsForGame will return statistics for all players in the given game
func GetX01StatisticsForGame(id int) ([]*models.StatisticsX01, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id,
			p.id,
			COUNT(DISTINCT g.id),
			SUM(s.ppd) / COUNT(p.id),
			SUM(s.first_nine_ppd) / COUNT(p.id),
			SUM(60s_plus),
			SUM(100s_plus),
			SUM(140s_plus),
			SUM(180s) AS '180s',
			SUM(accuracy_20) / COUNT(accuracy_20),
			SUM(accuracy_19) / COUNT(accuracy_19),
			SUM(overall_accuracy) / COUNT(overall_accuracy),
			SUM(checkout_percentage) / COUNT(checkout_percentage)
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN `+"`match`"+` m ON m.id = s.match_id
			JOIN game g ON g.id = m.game_id
			JOIN player2match p2m ON p2m.game_id = g.id
		WHERE g.id = ?
		AND g.game_type_id = 1
		GROUP BY p.id
		ORDER BY p2m.order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsX01, 0)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.MatchID, &s.PlayerID, &s.GamesPlayed, &s.PPD, &s.FirstNinePPD, &s.Score60sPlus, &s.Score100sPlus,
			&s.Score140sPlus, &s.Score180s, &s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &s.CheckoutPercentage)
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
		stats := statistics[0]
		visits, err := GetPlayerVisits(id)
		if err != nil {
			return nil, err
		}
		stats.Hits, stats.DartsThrown = models.GetHitsMap(visits)

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
			p.id,
			COUNT(DISTINCT m.id),
			SUM(s.ppd) / COUNT(DISTINCT m.id),
			SUM(s.first_nine_ppd) / COUNT(DISTINCT m.id),
			SUM(s.60s_plus),
			SUM(s.100s_plus),
			SUM(s.140s_plus),
			SUM(s.180s),
			SUM(accuracy_20) / COUNT(accuracy_20),
			SUM(accuracy_19) / COUNT(accuracy_19),
			SUM(overall_accuracy) / COUNT(overall_accuracy),
			SUM(checkout_percentage) / COUNT(checkout_percentage)
		FROM statistics_x01 s
		JOIN player p ON p.id = s.player_id
		JOIN `+"`match`"+` m ON m.id = s.match_id
		JOIN game g ON g.id = m.game_id
		WHERE s.player_id IN (?)
		AND m.starting_score IN (?)
		AND g.is_finished = 1
		AND g.game_type_id = 1
		GROUP BY s.player_id`, ids, startingScores)
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
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.PPD, &s.FirstNinePPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus,
			&s.Score180s, &s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &s.CheckoutPercentage)
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
			s.player_id,
			SUM(s.ppd) / COUNT(s.match_id) AS 'ppd',
			SUM(s.first_nine_ppd) / COUNT(s.match_id) AS 'first_nine_ppd',
			SUM(s.checkout_percentage) / COUNT(s.match_id) AS 'checkout_percentage',
			SUM(s.60s_plus) AS '60s_plus',
			SUM(s.100s_plus) AS '100s_plus',
			SUM(s.140s_plus) AS '140s_plus',
			SUM(s.180s) AS '180s',
			SUM(s.accuracy_20) / COUNT(s.match_id) AS 'accuracy_20',
			SUM(s.accuracy_19) / COUNT(s.match_id) AS 'accuracy_19',
			SUM(s.overall_accuracy) / COUNT(s.match_id) AS 'accuracy_overall',
			DATE(g.updated_at) AS 'date'
		FROM statistics_x01 s
		JOIN `+"`match`"+` m ON m.id = s.match_id
		JOIN game g ON g.id = m.game_id
		WHERE player_id = ?
		GROUP BY YEAR(g.updateD_at), WEEK(g.updated_at)
		ORDER BY date DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statisticsMap := make(map[string]*models.StatisticsX01)
	for rows.Next() {
		var date string
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.PPD, &s.FirstNinePPD, &s.CheckoutPercentage, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus,
			&s.Score180s, &s.Accuracy20, &s.Accuracy19, &s.AccuracyOverall, &date)
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
			m.winner_id,
			m.id,
			s.ppd,
			s.first_nine_ppd,
			s.checkout_percentage,
			s.darts_thrown,
			m.starting_score
		FROM statistics_x01 s
		JOIN player p ON p.id = s.player_id
		JOIN `+"`match`"+` m ON m.id = s.match_id
		WHERE s.player_id IN (?)
		AND m.starting_score IN (?)`, ids, startingScores)
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
		err := rows.Scan(&s.PlayerID, &s.WinnerID, &s.MatchID, &s.PPD, &s.FirstNinePPD, &s.CheckoutPercentage, &s.DartsThrown, &s.StartingScore)
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
					real.Best301.MatchID = stat.MatchID
				}
			}
			if stat.StartingScore.Int64 == 501 {
				if real.Best501 == nil {
					real.Best501 = new(models.BestStatistic)
				}
				if stat.DartsThrown < real.Best501.Value || real.Best501.Value == 0 {
					real.Best501.Value = stat.DartsThrown
					real.Best501.MatchID = stat.MatchID
				}
			}
			if stat.StartingScore.Int64 == 701 {
				if real.Best701 == nil {
					real.Best701 = new(models.BestStatistic)
				}
				if stat.DartsThrown < real.Best701.Value || real.Best701.Value == 0 {
					real.Best701.Value = stat.DartsThrown
					real.Best701.MatchID = stat.MatchID
				}
			}
		}
		if real.BestPPD == nil {
			real.BestPPD = new(models.BestStatisticFloat)
		}
		if stat.PPD > real.BestPPD.Value {
			real.BestPPD.Value = stat.PPD
			real.BestPPD.MatchID = stat.MatchID
		}
		if real.BestFirstNinePPD == nil {
			real.BestFirstNinePPD = new(models.BestStatisticFloat)
		}
		if stat.FirstNinePPD > real.BestFirstNinePPD.Value {
			real.BestFirstNinePPD.Value = stat.FirstNinePPD
			real.BestFirstNinePPD.MatchID = stat.MatchID
		}
	}
	return nil
}

// getHighestCheckout will calculate the highest checkout for the given players
func getHighestCheckout(ids []int, statisticsMap map[int]*models.StatisticsX01, startingScores ...int) error {
	q, args, err := sqlx.In(`
		SELECT
			player_id,
			match_id,
			MAX(checkout)
		FROM (SELECT
				s.player_id,
				s.match_id,
				IFNULL(s.first_dart * s.first_dart_multiplier, 0) +
					IFNULL(s.second_dart * s.second_dart_multiplier, 0) +
					IFNULL(s.third_dart * s.third_dart_multiplier, 0) AS 'checkout'
			FROM score s
			JOIN `+"`match`"+` m ON m.id = s.match_id
			WHERE m.winner_id = s.player_id
				AND s.player_id IN (?)
				AND s.id IN (SELECT MAX(s.id) FROM score s JOIN `+"`match`"+` m ON m.id = s.match_id WHERE m.winner_id = s.player_id GROUP BY match_id)
				AND m.starting_score IN (?)
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
		var matchID int
		var checkout int
		err := rows.Scan(&playerID, &matchID, &checkout)
		if err != nil {
			return err
		}
		highest := new(models.BestStatistic)
		highest.Value = checkout
		highest.MatchID = matchID
		statisticsMap[playerID].HighestCheckout = highest
	}
	err = rows.Err()
	return err
}
