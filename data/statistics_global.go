package data

import (
	"github.com/guregu/null"
	"github.com/kcapp/api/models"
)

// GetGlobalStatistics will return global statistics for all matches
func GetGlobalStatistics() (map[int]*models.GlobalStatistics, error) {
	rows, err := models.DB.Query(`
		SELECT
			office_id,
			COUNT(DISTINCT m.id) AS 'matches',
			COUNT(DISTINCT l.id) AS 'legs',
			COUNT(DISTINCT s.id) AS 'visits',
			SUM(s.first_dart * s.first_dart_multiplier + s.second_dart * s.second_dart_multiplier + s.third_dart * s.third_dart_multiplier) as 'points'
		FROM matches m
			LEFT JOIN leg l on l.match_id = m.id
			LEFT JOIN score s on s.leg_id = l.id
		WHERE m.is_finished = 1 AND m.is_abandoned = 0
		GROUP BY office_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[int]*models.GlobalStatistics, 0)
	for rows.Next() {
		var officeID null.Int
		s := new(models.GlobalStatistics)
		err := rows.Scan(&officeID, &s.Matches, &s.Legs, &s.Visits, &s.Points)
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
		s.FishNChips = fnc[officeID].FishNChips

		all.Legs += s.Legs
		all.Matches += s.Matches
		all.Visits += s.Visits
		all.Points += s.Points
		all.FishNChips += s.FishNChips
	}
	stats[0] = all

	return stats, nil
}

// GetGlobalStatisticsFnc will return global fish and chips statistics
func GetGlobalStatisticsFnc() (map[int]*models.GlobalStatistics, error) {
	rows, err := models.DB.Query(`
		SELECT
			office_id,
			COUNT(s.id) AS 'Fish-n-Chips'
		FROM score s
			LEFT JOIN leg l ON l.id = s.leg_id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE
			first_dart IN (1,20,5) AND first_dart_multiplier = 1 AND
			second_dart IN (1,20,5) AND second_dart_multiplier = 1 AND
			third_dart IN (1,20,5) AND third_dart_multiplier = 1  AND
			((first_dart * first_dart_multiplier) + (second_dart * second_dart_multiplier) +
			(third_dart * third_dart_multiplier) = 26)
			AND m.is_abandoned <> 1
		GROUP BY m.office_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[int]*models.GlobalStatistics, 0)
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
