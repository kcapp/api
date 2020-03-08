package data

import (
	"github.com/kcapp/api/models"
)

// GetCricketStatisticsForLeg will return statistics for all players in the given leg
func GetCricketStatisticsForLeg(id int) ([]*models.StatisticsCricket, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			s.total_marks,
			s.rounds,
			s.score,
			s.first_nine_marks,
			s.mpr,
			s.first_nine_mpr,
			s.marks5,
			s.marks6,
			s.marks7,
			s.marks8,
			s.marks9
		FROM statistics_cricket s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE l.id = ? GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsCricket, 0)
	for rows.Next() {
		s := new(models.StatisticsCricket)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.TotalMarks, &s.Rounds, &s.Score, &s.FirstNineMarks,
			&s.MPR, &s.FirstNineMPR, &s.Marks5, &s.Marks6, &s.Marks7, &s.Marks8, &s.Marks9)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetCricketStatisticsForMatch will return statistics for all players in the given match
func GetCricketStatisticsForMatch(id int) ([]*models.StatisticsCricket, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			SUM(s.total_marks),
			SUM(s.first_nine_marks),
			SUM(s.total_marks) / SUM(s.rounds) as 'mpr',
			SUM(s.first_nine_marks) / (COUNT(l.id) * 3) as 'first_nine_mpr',
			SUM(s.marks5) as 'marks5',
			SUM(s.marks6) as 'marks6',
			SUM(s.marks7) as 'marks7',
			SUM(s.marks8) as 'marks8',
			SUM(s.marks9) as 'marks9'
		FROM statistics_cricket s
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

	stats := make([]*models.StatisticsCricket, 0)
	for rows.Next() {
		s := new(models.StatisticsCricket)
		err := rows.Scan(&s.PlayerID, &s.TotalMarks, &s.FirstNineMarks, &s.MPR, &s.FirstNineMPR,
			&s.Marks5, &s.Marks6, &s.Marks7, &s.Marks8, &s.Marks9)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// CalculateCricketStatistics will generate cricket statistics for the given leg
func CalculateCricketStatistics(legID int) (map[int]*models.StatisticsCricket, error) {
	visits, err := GetLegVisits(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}
	statisticsMap := make(map[int]*models.StatisticsCricket)
	playerHitsMap := make(map[int]map[int]int64)
	for _, player := range players {
		stats := new(models.StatisticsCricket)
		stats.PlayerID = player.PlayerID
		stats.Score = player.CurrentScore
		statisticsMap[player.PlayerID] = stats
		playerHitsMap[player.PlayerID] = make(map[int]int64)
	}

	round := 1
	darts := []int{15, 16, 17, 18, 19, 20, 25}
	for i := 0; i < len(visits); i++ {
		visit := visits[i]
		stats := statisticsMap[visit.PlayerID]

		if i > 0 && i%len(players) == 0 {
			round++
		}

		marks := visit.GetMarksHit(darts, playerHitsMap)
		stats.TotalMarks += marks
		if round <= 3 {
			stats.FirstNineMarks += marks
		}
		//
		switch mpr := marks; mpr {
		case 5:
			stats.Marks5++
		case 6:
			stats.Marks6++
		case 7:
			stats.Marks7++
		case 8:
			stats.Marks8++
		case 9:
			stats.Marks9++
		default:
		}
	}
	for _, stat := range statisticsMap {
		stat.MPR = float32(stat.TotalMarks) / float32(round)
		stat.FirstNineMPR = float32(stat.FirstNineMarks) / 3
		stat.Rounds = round
	}
	return statisticsMap, nil
}
