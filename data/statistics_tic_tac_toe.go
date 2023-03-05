package data

import (
	"database/sql"

	"github.com/kcapp/api/models"
)

// GetTicTacToeStatistics will return statistics for all players active during the given period
func GetTicTacToeStatistics(from string, to string) ([]*models.StatisticsTicTacToe, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) as 'matches_played',
			COUNT(DISTINCT m2.id) as 'matches_won',
			COUNT(DISTINCT l.id) as 'legs_played',
			COUNT(DISTINCT l2.id) as 'legs_won',
			m.office_id AS 'office_id',
			SUM(darts_Thrown) as 'darts_thrown',
			SUM(score) as 'score',
			SUM(numbers_closed) as 'numbers_closed',
			MAX(highest_closed) as 'highest_closed'
		FROM statistics_tic_tac_toe s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
			AND m.match_type_id = 9
		GROUP BY p.id, m.office_id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsTicTacToe, 0)
	for rows.Next() {
		s := new(models.StatisticsTicTacToe)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.OfficeID,
			&s.DartsThrown, &s.Score, &s.NumbersClosed, &s.HighestClosed)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetTicTacToeStatisticsForLeg will return statistics for all players in the given leg
func GetTicTacToeStatisticsForLeg(id int) ([]*models.StatisticsTicTacToe, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			s.darts_thrown,
			s.score,
			s.numbers_closed,
			s.highest_closed
		FROM statistics_tic_tac_toe s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE l.id = ? GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsTicTacToe, 0)
	for rows.Next() {
		s := new(models.StatisticsTicTacToe)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.NumbersClosed, &s.HighestClosed)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetTicTacToeStatisticsForMatch will return statistics for all players in the given match
func GetTicTacToeStatisticsForMatch(id int) ([]*models.StatisticsTicTacToe, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			SUM(darts_Thrown) as 'darts_thrown',
			SUM(score) as 'score',
			SUM(numbers_closed) as 'numbers_closed',
			MAX(highest_closed) as 'highest_closed'
		FROM statistics_tic_tac_toe s
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

	stats := make([]*models.StatisticsTicTacToe, 0)
	for rows.Next() {
		s := new(models.StatisticsTicTacToe)
		err := rows.Scan(&s.PlayerID, &s.DartsThrown, &s.Score, &s.NumbersClosed, &s.HighestClosed)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetTicTacToeStatisticsForPlayer will return statistics for the given player
func GetTicTacToeStatisticsForPlayer(id int) (*models.StatisticsTicTacToe, error) {
	s := new(models.StatisticsTicTacToe)
	err := models.DB.QueryRow(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) as 'matches_played',
			COUNT(DISTINCT m2.id) as 'matches_won',
			COUNT(DISTINCT l.id) as 'legs_played',
			COUNT(DISTINCT l2.id) as 'legs_won',
			SUM(darts_Thrown) as 'darts_thrown',
			SUM(score) as 'score',
			SUM(numbers_closed) as 'numbers_closed',
			MAX(highest_closed) as 'highest_closed'
		FROM statistics_tic_tac_toe s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
			AND m.match_type_id = 9
		GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.Score, &s.DartsThrown, &s.NumbersClosed, &s.HighestClosed)
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.StatisticsTicTacToe), nil
		}
		return nil, err
	}
	return s, nil
}

// GetTicTacToeHistoryForPlayer will return history of statistics for the given player
func GetTicTacToeHistoryForPlayer(id int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.TICTACTOE, false)
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
			s.numbers_closed,
			s.highest_closed
		FROM statistics_tic_tac_toe s
			LEFT JOIN player p ON p.id = s.player_id
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 9
		ORDER BY l.id DESC
		LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.StatisticsTicTacToe)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.Score, &s.NumbersClosed, &s.HighestClosed)
		if err != nil {
			return nil, err
		}
		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// CalculateTicTacToeStatistics will generate tic tac toe statistics for the given leg
func CalculateTicTacToeStatistics(legID int) (map[int]*models.StatisticsTicTacToe, error) {
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
	statisticsMap := make(map[int]*models.StatisticsTicTacToe)
	playerHitsMap := make(map[int]map[int]int64)
	for _, player := range players {
		stats := new(models.StatisticsTicTacToe)
		stats.PlayerID = player.PlayerID
		statisticsMap[player.PlayerID] = stats
		playerHitsMap[player.PlayerID] = make(map[int]int64)
	}
	for _, visit := range visits {
		stats := statisticsMap[visit.PlayerID]
		stats.DartsThrown += visit.GetDartsThrown()
	}

	for _, stats := range statisticsMap {
		for num, playerID := range leg.Parameters.Hits {
			if playerID == stats.PlayerID {
				stats.NumbersClosed++
				if num > stats.HighestClosed {
					stats.HighestClosed = num
				}
				stats.Score += num
			}
		}
	}
	return statisticsMap, nil
}
