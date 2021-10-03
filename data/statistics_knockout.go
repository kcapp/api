package data

import (
	"database/sql"
	"log"

	"github.com/kcapp/api/models"
)

// GetKnockoutStatistics will return statistics for all players active during the given period
func GetKnockoutStatistics(from string, to string) ([]*models.StatisticsKnockout, error) {
	rows, err := models.DB.Query(`
			SELECT
				p.id,
				COUNT(DISTINCT m.id) AS 'matches_played',
				COUNT(DISTINCT m2.id) AS 'matches_won',
				COUNT(DISTINCT l.id) AS 'legs_played',
				COUNT(DISTINCT l2.id) AS 'legs_won',
				m.office_id AS 'office_id',
				SUM(s.darts_thrown) as 'darts_thrown',
				SUM(s.avg_score) / COUNT(DISTINCT l.id) as 'avg_score',
				SUM(s.lives_lost) as 'lives_lost',
				SUM(s.lives_taken) as 'lives_taken',
				CAST(SUM(s.final_position) / COUNT(DISTINCT l.id) AS SIGNED) as 'final_position'
			FROM statistics_knockout s
				JOIN player p ON p.id = s.player_id
				JOIN leg l ON l.id = s.leg_id
				JOIN matches m ON m.id = l.match_id
				LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
				LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
			WHERE m.updated_at >= ? AND m.updated_at < ?
				AND l.is_finished = 1 AND m.is_abandoned = 0
				AND m.match_type_id = 15
			GROUP BY p.id, m.office_id
			ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsKnockout, 0)
	for rows.Next() {
		s := new(models.StatisticsKnockout)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.OfficeID, &s.DartsThrown,
			&s.AvgScore, &s.LivesLost, &s.LivesTaken, &s.FinalPosition)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetKnockoutStatisticsForLeg will return statistics for all players in the given leg
func GetKnockoutStatisticsForLeg(id int) ([]*models.StatisticsKnockout, error) {
	rows, err := models.DB.Query(`
			SELECT
				l.id,
				p.id,
				s.darts_thrown,
				s.avg_score,
				s.lives_lost,
				s.lives_taken,
				s.final_position
			FROM statistics_knockout s
				JOIN player p ON p.id = s.player_id
				JOIN leg l ON l.id = s.leg_id
				JOIN player2leg p2l on l.id = p2l.leg_id AND p.id = p2l.player_id
			WHERE l.id = ? GROUP BY p.id ORDER BY p2l.order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsKnockout, 0)
	for rows.Next() {
		s := new(models.StatisticsKnockout)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.AvgScore, &s.LivesLost, &s.LivesTaken, &s.FinalPosition)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetKnockoutStatisticsForMatch will return statistics for all players in the given match
func GetKnockoutStatisticsForMatch(id int) ([]*models.StatisticsKnockout, error) {
	rows, err := models.DB.Query(`
			SELECT
				p.id,
				SUM(s.darts_thrown) as 'darts_thrown',
				SUM(s.avg_score) / COUNT(DISTINCT l.id) as 'avg_score',
				SUM(s.lives_lost) as 'lives_lost',
				SUM(s.lives_taken) as 'lives_taken',
				CAST(SUM(s.final_position) / COUNT(DISTINCT l.id) AS SIGNED) as 'final_position'
			FROM statistics_knockout s
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

	stats := make([]*models.StatisticsKnockout, 0)
	for rows.Next() {
		s := new(models.StatisticsKnockout)
		err := rows.Scan(&s.PlayerID, &s.DartsThrown, &s.AvgScore, &s.LivesLost, &s.LivesTaken, &s.FinalPosition)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetKnockoutStatisticsForPlayer will return Knockout statistics for the given player
func GetKnockoutStatisticsForPlayer(id int) (*models.StatisticsKnockout, error) {
	s := new(models.StatisticsKnockout)
	err := models.DB.QueryRow(`
			SELECT
				p.id,
				COUNT(DISTINCT m.id) AS 'matches_played',
				COUNT(DISTINCT m2.id) AS 'matches_won',
				COUNT(DISTINCT l.id) AS 'legs_played',
				COUNT(DISTINCT l2.id) AS 'legs_won',
				SUM(s.darts_thrown) as 'darts_thrown',
				SUM(s.avg_score) / COUNT(DISTINCT l.id) as 'avg_score',
				SUM(s.lives_lost) as 'lives_lost',
				SUM(s.lives_taken) as 'lives_taken',
				CAST(SUM(s.final_position) / COUNT(DISTINCT l.id) AS SIGNED) as 'final_position'
			FROM statistics_knockout s
				JOIN player p ON p.id = s.player_id
				JOIN leg l ON l.id = s.leg_id
				JOIN matches m ON m.id = l.match_id
				LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
				LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
			WHERE s.player_id = ?
				AND l.is_finished = 1 AND m.is_abandoned = 0
				AND m.match_type_id = 15
			GROUP BY p.id`, id).Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.DartsThrown,
		&s.AvgScore, &s.LivesLost, &s.LivesTaken, &s.FinalPosition)
	if err != nil {
		if err == sql.ErrNoRows {
			return new(models.StatisticsKnockout), nil
		}
		return nil, err
	}
	return s, nil
}

// GetKnockoutHistoryForPlayer will return history of Knockout statistics for the given player
func GetKnockoutHistoryForPlayer(id int, limit int) ([]*models.Leg, error) {
	legs, err := GetLegsOfType(models.KNOCKOUT, false)
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
				s.avg_score,
				s.lives_lost,
				s.lives_taken,
				s.final_position
			FROM statistics_knockout s
				LEFT JOIN player p ON p.id = s.player_id
				LEFT JOIN leg l ON l.id = s.leg_id
				LEFT JOIN matches m ON m.id = l.match_id
			WHERE s.player_id = ?
				AND l.is_finished = 1 AND m.is_abandoned = 0
				AND m.match_type_id = 15
			ORDER BY l.id DESC
			LIMIT ?`, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs = make([]*models.Leg, 0)
	for rows.Next() {
		s := new(models.StatisticsKnockout)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.DartsThrown, &s.AvgScore, &s.LivesLost, &s.LivesTaken, &s.FinalPosition)
		if err != nil {
			return nil, err
		}
		leg := m[s.LegID]
		leg.Statistics = s
		legs = append(legs, leg)
	}
	return legs, nil
}

// CalculateKnockoutStatistics will generate Knockout statistics for the given leg
func CalculateKnockoutStatistics(legID int) (map[int]*models.StatisticsKnockout, error) {
	leg, err := GetLeg(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}

	statisticsMap := make(map[int]*models.StatisticsKnockout)
	for _, player := range players {
		stats := new(models.StatisticsKnockout)
		stats.PlayerID = player.PlayerID
		statisticsMap[player.PlayerID] = stats
	}
	startingLives := leg.Parameters.StartingLives
	finalPosition := len(players)
	for i, visit := range leg.Visits {
		stats := statisticsMap[visit.PlayerID]
		stats.AvgScore += float64(visit.GetScore())

		idx := i - 1
		if idx < 0 {
			continue
		}
		prev := leg.Visits[idx]
		if prev.GetScore() > visit.GetScore() {
			stats.LivesLost++
			statisticsMap[prev.PlayerID].LivesTaken++
		}

		if stats.LivesLost == int(startingLives.Int64) {
			stats.FinalPosition = finalPosition
			finalPosition--
		}
		stats.DartsThrown = visit.DartsThrown
	}

	for _, stats := range statisticsMap {
		stats.AvgScore = stats.AvgScore / (float64(stats.DartsThrown / 3))
		if stats.FinalPosition == 0 {
			stats.FinalPosition = finalPosition
		}
	}
	return statisticsMap, nil
}

// ReCalculateKnockoutStatistics will recaulcate statistics for Knockout legs
func ReCalculateKnockoutStatistics() (map[int]map[int]*models.StatisticsKnockout, error) {
	legs, err := GetLegsOfType(models.KNOCKOUT, true)
	if err != nil {
		return nil, err
	}

	s := make(map[int]map[int]*models.StatisticsKnockout)
	for _, leg := range legs {
		stats, err := CalculateKnockoutStatistics(leg.ID)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			log.Printf(`UPDATE statistics_knockout SET darts_thrown = %d, avg_score = %f, lives_lost = %d, lives_taken = %d,
			final_position = %d WHERE leg_id = %d AND player_id = %d;`,
				stat.DartsThrown, stat.AvgScore, stat.LivesLost, stat.LivesTaken, stat.FinalPosition, leg.ID, playerID)
		}
		s[leg.ID] = stats
	}

	return s, err
}
