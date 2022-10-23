package data

import (
	"database/sql"
	"log"

	"github.com/kcapp/api/models"
)

// GetGotchaStatistics will return statistics for all players active during the given period
func GetGotchaStatistics(from string, to string) ([]*models.StatisticsGotcha, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			m.office_id AS 'office_id',
			SUM(s.darts_thrown) as 'darts_thrown',
			MAX(s.highest_score) as 'highest_score',
			SUM(s.times_reset) as 'times_reset',
			SUM(s.others_reset) as 'others_reset',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score'
		FROM statistics_gotcha s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 13
		GROUP BY p.id, m.office_id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsGotcha, 0)
	for rows.Next() {
		s := new(models.StatisticsGotcha)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon,
			&s.OfficeID, &s.DartsThrown, &s.HighestScore, &s.TimesReset, &s.OthersReset, &s.Score)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetGotchaStatisticsForLeg will return statistics for all players in the given leg
func GetGotchaStatisticsForLeg(id int) ([]*models.StatisticsGotcha, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			s.darts_thrown,
			s.highest_score,
			s.times_reset,
			s.others_reset,
			s.score
		FROM statistics_gotcha s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN player2leg p2l on l.id = p2l.leg_id AND p.id = p2l.player_id
		WHERE l.id = ? GROUP BY p.id ORDER BY p2l.order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsGotcha, 0)
	for rows.Next() {
		s := new(models.StatisticsGotcha)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.HighestScore, &s.TimesReset, &s.OthersReset, &s.Score)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetGotchaStatisticsForMatch will return statistics for all players in the given match
func GetGotchaStatisticsForMatch(id int) ([]*models.StatisticsGotcha, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			SUM(s.darts_thrown) as 'darts_thrown',
			MAX(s.highest_score) as 'highest_score',
			SUM(s.times_reset) as 'times_reset',
			SUM(s.others_reset) as 'others_reset',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score'
		FROM statistics_gotcha s
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

	stats := make([]*models.StatisticsGotcha, 0)
	for rows.Next() {
		s := new(models.StatisticsGotcha)
		err := rows.Scan(&s.PlayerID, &s.DartsThrown, &s.HighestScore, &s.TimesReset, &s.OthersReset, &s.Score)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetGotchaStatisticsForPlayer will return Gotcha statistics for the given player
func GetGotchaStatisticsForPlayer(id int) (*models.StatisticsGotcha, error) {
	s := new(models.StatisticsGotcha)
	err := models.DB.QueryRow(`
		SELECT
			p.id,
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			SUM(s.darts_thrown) as 'darts_thrown',
			MAX(s.highest_score) as 'highest_score',
			SUM(s.times_reset) as 'times_reset',
			SUM(s.others_reset) as 'others_reset',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score'
		FROM statistics_gotcha s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 13
		GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon,
		&s.DartsThrown, &s.HighestScore, &s.TimesReset, &s.OthersReset, &s.Score)
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.StatisticsGotcha), nil
		}
		return nil, err
	}
	return s, nil
}

// GetGotchaHistoryForPlayer will return history of Gotcha statistics for the given player
func GetGotchaHistoryForPlayer(id int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.GOTCHA, false)
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
			s.highest_score,
			s.times_reset,
			s.others_reset,
			s.score
		FROM statistics_gotcha s
			LEFT JOIN player p ON p.id = s.player_id
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 13
		ORDER BY l.id DESC
		LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.StatisticsGotcha)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.HighestScore, &s.TimesReset, &s.OthersReset, &s.Score)
		if err != nil {
			return nil, err
		}
		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// CalculateGotchaStatistics will generate Gotcha statistics for the given leg
func CalculateGotchaStatistics(legID int) (map[int]*models.StatisticsGotcha, error) {
	leg, err := GetLeg(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}

	statisticsMap := make(map[int]*models.StatisticsGotcha)
	for _, player := range players {
		// Reset scores so we recalculate them in the loop below
		player.CurrentScore = 0

		stats := new(models.StatisticsGotcha)
		stats.PlayerID = player.PlayerID
		stats.Score = 0
		stats.HighestScore = 0
		stats.TimesReset = 0
		stats.OthersReset = 0
		statisticsMap[player.PlayerID] = stats
	}

	round := 0
	for i, visit := range leg.Visits {
		if i > 0 && i%len(players) == 0 {
			round++
		}

		stats := statisticsMap[visit.PlayerID]
		if round > 0 && players[visit.PlayerID].CurrentScore == 0 {
			stats.TimesReset++
		}
		stats.OthersReset += getPlayersReset(visit, players)

		score := visit.CalculateGotchaScore(players, leg.StartingScore)
		players[visit.PlayerID].CurrentScore += score
		stats.Score = players[visit.PlayerID].CurrentScore
		if stats.Score > stats.HighestScore {
			stats.HighestScore = stats.Score
		}
		stats.DartsThrown = visit.DartsThrown
	}
	return statisticsMap, nil
}

func getPlayersReset(visit *models.Visit, players map[int]*models.Player2Leg) int {
	resets := 0
	currentScore := players[visit.PlayerID].CurrentScore + visit.FirstDart.GetScore()
	for _, player := range players {
		if visit.PlayerID != player.PlayerID && player.CurrentScore == currentScore {
			resets++
		}
	}
	if !visit.SecondDart.IsMiss() {
		currentScore += visit.SecondDart.GetScore()
		for _, player := range players {
			if visit.PlayerID != player.PlayerID && player.CurrentScore == currentScore {
				resets++
			}
		}
	}

	if !visit.ThirdDart.IsMiss() {
		currentScore += visit.ThirdDart.GetScore()
		for _, player := range players {
			if visit.PlayerID != player.PlayerID && player.CurrentScore == currentScore {
				resets++
			}
		}
	}
	return resets
}

// ReCalculateGotchaStatistics will recaulcate statistics for Gotcha legs
func ReCalculateGotchaStatistics() (map[int]map[int]*models.StatisticsGotcha, error) {
	legs, err := GetLegsOfType(models.GOTCHA, true)
	if err != nil {
		return nil, err
	}

	s := make(map[int]map[int]*models.StatisticsGotcha)
	for _, leg := range legs {
		stats, err := CalculateGotchaStatistics(leg.ID)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			log.Printf(`UPDATE statistics_gotcha SET darts_thrown = %d, highest_score = %d, times_reset = %d, others_reset = %d, score = %d WHERE leg_id = %d AND player_id = %d;`,
				stat.DartsThrown, stat.HighestScore, stat.TimesReset, stat.OthersReset, stat.Score, leg.ID, playerID)
		}
		s[leg.ID] = stats
	}

	return s, err
}
