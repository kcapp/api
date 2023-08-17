package data

import (
	"github.com/guregu/null"
	"github.com/kcapp/api/models"
)

// GetGlobalStatistics will return global statistics for all matches
func GetGlobalStatistics() (map[int]*models.GlobalStatistics, error) {
	rows, err := models.DB.Query(`
			SELECT
				m.office_id,
				COUNT(DISTINCT m.id) AS 'matches',
				COUNT(DISTINCT l.id) AS 'legs',
				COUNT(DISTINCT s.id) AS 'visits',
				COUNT(first_dart) + SUM(IF(second_dart is null, 0, 1)) + SUM(IF(third_dart is null, 0, 1)) as darts,
				IFNULL(SUM(s.first_dart * s.first_dart_multiplier + IFNULL(s.second_dart, 0) * s.second_dart_multiplier + IFNULL(s.third_dart, 0) * s.third_dart_multiplier) - SUM(IF(s.is_bust = 1, s.first_dart * s.first_dart_multiplier + IFNULL(s.second_dart, 0) * s.second_dart_multiplier + IFNULL(s.third_dart, 0) * s.third_dart_multiplier, 0)), 0) as 'points',
				SUM(IF(s.is_bust = 1, s.first_dart * s.first_dart_multiplier + IFNULL(s.second_dart, 0) * s.second_dart_multiplier + IFNULL(s.third_dart, 0) * s.third_dart_multiplier, 0)) as 'points_busted',
				SUM(IF(s.is_bust, 0, IF(first_dart = 20 AND first_dart_multiplier = 3 AND second_dart = 20 AND second_dart_multiplier = 3 AND third_dart = 20 AND third_dart_multiplier = 3, 1, 0))) as '180s',
				SUM(IF(s.is_bust, 0, IF((first_dart = 25 AND first_dart_multiplier = 2) OR (second_dart = 25 AND second_dart_multiplier = 2) OR (third_dart = 25 AND third_dart_multiplier = 2), 1, 0))) as 'bullseyes'
			FROM matches m
				LEFT JOIN leg l on l.match_id = m.id
				LEFT JOIN score s on s.leg_id = l.id
				LEFT JOIN player p on p.id = s.player_id
			WHERE m.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0 AND (p.id is null OR p.is_bot = 0)
			GROUP BY m.office_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[int]*models.GlobalStatistics)
	for rows.Next() {
		var officeID null.Int
		s := new(models.GlobalStatistics)
		err := rows.Scan(&officeID, &s.Matches, &s.Legs, &s.Visits, &s.Darts, &s.Points, &s.PointsBusted, &s.Score180s, &s.ScoreBullseyes)
		if err != nil {
			return nil, err
		}
		stats[int(officeID.Int64)] = s
	}

	fnc, err := GetGlobalStatisticsFnc()
	if err != nil {
		return nil, err
	}

	all := new(models.GlobalStatistics)
	for officeID, s := range stats {
		if officeID == 0 {
			continue
		}
		if _, ok := fnc[officeID]; ok {
			s.FishNChips = fnc[officeID].FishNChips
		}

		all.Legs += s.Legs
		all.Matches += s.Matches
		all.Visits += s.Visits
		all.Darts += s.Darts
		all.Points += s.Points
		all.PointsBusted += s.PointsBusted
		all.Score180s += s.Score180s
		all.ScoreBullseyes += s.ScoreBullseyes
	}
	all.FishNChips = fnc[0].FishNChips
	stats[0] = all

	return stats, nil
}

// GetGlobalStatisticsFnc will return global fish and chips statistics
func GetGlobalStatisticsFnc() (map[int]*models.GlobalStatistics, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.office_id,
			COUNT(s.id) AS 'Fish-n-Chips'
		FROM score s
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
			LEFT JOIN player p ON p.id = s.player_id
		WHERE
			first_dart IN (1,20,5) AND first_dart_multiplier = 1 AND
			second_dart IN (1,20,5) AND second_dart_multiplier = 1 AND
			third_dart IN (1,20,5) AND third_dart_multiplier = 1  AND
			((first_dart * first_dart_multiplier) + (second_dart * second_dart_multiplier) +
			(third_dart * third_dart_multiplier) = 26)
			AND m.is_abandoned <> 1
			AND p.is_bot = 0
		GROUP BY m.office_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[int]*models.GlobalStatistics)
	for rows.Next() {
		var officeID null.Int
		s := new(models.GlobalStatistics)
		err := rows.Scan(&officeID, &s.FishNChips)
		if err != nil {
			return nil, err
		}
		stats[int(officeID.Int64)] = s
	}
	all := new(models.GlobalStatistics)
	for _, s := range stats {
		all.FishNChips += s.FishNChips
	}
	stats[0] = all

	return stats, nil
}
