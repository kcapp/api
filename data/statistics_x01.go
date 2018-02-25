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
		WHERE m.id = ? GROUP BY p.id`, id)
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

// GetShootoutStatisticsForMatch will return statistics for all players in the given match
func GetShootoutStatisticsForMatch(id int) ([]*models.StatisticsShootout, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id,
			p.id,
			SUM(s.ppd) / COUNT(p.id),
			SUM(60s_plus),
			SUM(100s_plus),
			SUM(140s_plus),
			SUM(180s) AS '180s'
		FROM statistics_shootout s
			JOIN player p ON p.id = s.player_id
			JOIN `+"`match`"+` m ON m.id = s.match_id
			JOIN game g ON g.id = m.game_id
		WHERE m.id = ? GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsShootout, 0)
	for rows.Next() {
		s := new(models.StatisticsShootout)
		err := rows.Scan(&s.MatchID, &s.PlayerID, &s.PPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus, &s.Score180s)
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
		WHERE g.id = ? GROUP BY p.id`, id)
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

// GetPlayerStatistics will get statistics about the given player id
func GetPlayerStatistics(id int) (*models.StatisticsX01, error) {
	ids := []int{id}
	statistics, err := GetPlayersStatistics(ids)
	if err != nil {
		return nil, err
	}
	stats := statistics[0]
	visits, err := GetPlayerVisits(id)
	if err != nil {
		return nil, err
	}
	stats.Hits, stats.DartsThrown = models.GetHitsMap(visits)
	return stats, nil
}

// GetPlayersStatistics will get statistics about all the the given player IDs
func GetPlayersStatistics(ids []int) ([]*models.StatisticsX01, error) {
	q, args, err := sqlx.In(`
		SELECT
			p.id,
			SUM(s.ppd) / p.games_played,
			SUM(s.first_nine_ppd) / p.games_played,
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
		WHERE s.player_id IN (?)
		GROUP BY s.player_id`, ids)
	if err != nil {
		return nil, err
	}
	rows, err := models.DB.Query(q, args...)
	defer rows.Close()

	statisticsMap := make(map[int]*models.StatisticsX01)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.PPD, &s.FirstNinePPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus,
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
	err = getBestStatistics(ids, statisticsMap)
	if err != nil {
		return nil, err
	}
	err = getHighestCheckout(ids, statisticsMap)
	if err != nil {
		return nil, err
	}

	statistics := make([]*models.StatisticsX01, 0)
	for _, s := range statisticsMap {
		statistics = append(statistics, s)
	}
	return statistics, nil
}

// getBestStatistics will calculate Best PPD, Best First 9, Best 301 and Best 501 for the given players
func getBestStatistics(ids []int, statisticsMap map[int]*models.StatisticsX01) error {
	q, args, err := sqlx.In(`
		SELECT
			p.id,
			s.ppd,
			s.first_nine_ppd,
			s.checkout_percentage,
			s.darts_thrown,
			m.starting_score
		FROM statistics_x01 s
		JOIN player p ON p.id = s.player_id
		JOIN `+"`match`"+` m ON m.id = s.match_id
		WHERE s.player_id IN (?)`, ids)
	if err != nil {
		return err
	}
	rows, err := models.DB.Query(q, args...)
	defer rows.Close()

	rawStatistics := make([]*models.StatisticsX01, 0)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.PPD, &s.FirstNinePPD, &s.CheckoutPercentage, &s.DartsThrown, &s.StartingScore)
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
		if stat.StartingScore.Int64 == 301 && (stat.DartsThrown < real.Best301 || real.Best301 == 0) {
			real.Best301 = stat.DartsThrown
		}
		if stat.StartingScore.Int64 == 501 && (stat.DartsThrown < real.Best501 || real.Best501 == 0) {
			real.Best501 = stat.DartsThrown
		}
		if stat.PPD > real.BestPPD {
			real.BestPPD = stat.PPD
		}
		if stat.FirstNinePPD > real.BestFirstNinePPD {
			real.BestFirstNinePPD = stat.FirstNinePPD
		}
	}
	return nil
}

// getHighestCheckout will calculate the highest checkout for the given players
func getHighestCheckout(ids []int, statisticsMap map[int]*models.StatisticsX01) error {
	q, args, err := sqlx.In(`
		SELECT
			s.player_id,
			MAX(IFNULL(s.first_dart * s.first_dart_multiplier, 0) +
			IFNULL(s.second_dart * s.second_dart_multiplier, 0) +
			IFNULL(s.third_dart * s.third_dart_multiplier, 0)) AS 'highest_checkout'
		FROM score s
		JOIN `+"`match`"+` m ON m.id = s.match_id
		WHERE m.winner_id = s.player_id
			AND s.player_id IN (?)
			AND s.id IN (SELECT MAX(s.id) FROM score s JOIN `+"`match`"+`m ON m.id = s.match_id WHERE m.winner_id = s.player_id GROUP BY match_id)
		GROUP BY player_id
		ORDER BY highest_checkout DESC`, ids)
	if err != nil {
		return err
	}
	rows, err := models.DB.Query(q, args...)
	defer rows.Close()

	for rows.Next() {
		var playerID int
		var checkout int
		err := rows.Scan(&playerID, &checkout)
		if err != nil {
			return err
		}
		statisticsMap[playerID].HighestCheckout = checkout
	}
	err = rows.Err()
	return err
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
