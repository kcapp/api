package data

import (
	"github.com/kcapp/api/models"
)

// GetShootoutStatistics will return statistics for all players active duing the given period
func GetShootoutStatistics(from string, to string) ([]*models.StatisticsShootout, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			COUNT(DISTINCT g.id),
			SUM(s.ppd) / COUNT(p.id),
			SUM(60s_plus),
			SUM(100s_plus),
			SUM(140s_plus),
			SUM(180s) AS '180s'
		FROM statistics_shootout s
			JOIN player p ON p.id = s.player_id
			JOIN `+"`match`"+` m ON m.id = s.match_id
			JOIN game g ON g.id = m.game_id
		WHERE g.updated_at >= ? AND g.updated_at < ?
		AND g.is_finished = 1
		AND g.game_type_id = 2
		GROUP BY p.id`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statsMap := make(map[int]*models.StatisticsShootout, 0)
	for rows.Next() {
		s := new(models.StatisticsShootout)
		err := rows.Scan(&s.PlayerID, &s.GamesPlayed, &s.PPD, &s.Score60sPlus, &s.Score100sPlus, &s.Score140sPlus, &s.Score180s)
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
		AND g.game_type_id = 2
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

	stats := make([]*models.StatisticsShootout, 0)
	for _, s := range statsMap {
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
