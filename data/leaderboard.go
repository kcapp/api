package data

import "github.com/kcapp/api/models"

// GetMatchTypeLeaderboard will return leaderboard statistics for each match type
func GetMatchTypeLeaderboard() (map[int][]*models.MatchTypeLeaderboard, error) {
	rows, err := models.DB.Query(`CALL get_leaderboard_per_match_type(5);`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leaderboard := make(map[int][]*models.MatchTypeLeaderboard)
	for rows.Next() {
		mtl := new(models.MatchTypeLeaderboard)
		err := rows.Scan(&mtl.MatchTypeID, &mtl.PlayerID, &mtl.LegID, &mtl.DartsThrown, &mtl.Score, &mtl.ThreeDartAvg)
		if err != nil {
			return nil, err
		}
		if _, ok := leaderboard[mtl.MatchTypeID]; !ok {
			leaderboard[mtl.MatchTypeID] = make([]*models.MatchTypeLeaderboard, 0)
		}
		leaderboard[mtl.MatchTypeID] = append(leaderboard[mtl.MatchTypeID], mtl)
	}
	return leaderboard, nil
}
