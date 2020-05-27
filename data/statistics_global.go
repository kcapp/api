package data

import "github.com/kcapp/api/models"

// GetGlobalStatistics will return global statistics for all matches
func GetGlobalStatistics() (*models.GlobalStatistics, error) {
	global, err := GetGlobalStatisticsFnc()
	if err != nil {
		return nil, err
	}

	err = models.DB.QueryRow(`
		SELECT
			COUNT(DISTINCT m.id) AS 'matches',
			COUNT(DISTINCT l.id) AS 'legs',
			COUNT(DISTINCT s.id) AS 'visits'
		FROM matches m
			LEFT JOIN leg l on l.match_id = m.id
			LEFT JOIN score s on s.leg_id = l.id
		WHERE m.is_finished = 1 AND m.is_abandoned = 0`).Scan(&global.Matches, &global.Legs, &global.Visits)
	if err != nil {
		return nil, err
	}

	err = models.DB.QueryRow(`SELECT SUM(s.first_dart * s.first_dart_multiplier + s.second_dart * s.second_dart_multiplier + s.third_dart * s.third_dart_multiplier) from score s`).Scan(&global.Points)
	if err != nil {
		return nil, err
	}

	return global, nil
}

// GetGlobalStatisticsFnc will return global fish and chips statistics
func GetGlobalStatisticsFnc() (*models.GlobalStatistics, error) {
	global := new(models.GlobalStatistics)
	err := models.DB.QueryRow(`
		SELECT
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
			AND m.is_abandoned <> 1`).Scan(&global.FishNChips)
	if err != nil {
		return nil, err
	}
	return global, nil
}
