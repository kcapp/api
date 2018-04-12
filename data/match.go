package data

import (
	"database/sql"
	"log"

	"github.com/guregu/null"
	"github.com/kcapp/api/models"
	"github.com/kcapp/api/util"
)

// NewMatch will create a new match for the given game
func NewMatch(gameID int, startingScore int, players []int) (*models.Match, error) {
	tx, err := models.DB.Begin()
	if err != nil {
		return nil, err
	}

	// Shift players to get correct order
	id, players := players[0], players[1:]
	players = append(players, id)
	res, err := tx.Exec("INSERT INTO `match` (starting_score, current_player_id, game_id, created_at) VALUES (?, ?, ?, NOW()) ",
		startingScore, players[0], gameID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	matchID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	game, err := GetGame(gameID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("UPDATE game SET current_match_id = ? WHERE id = ?", matchID, gameID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	handicaps := make(map[int]null.Int)
	if game.GameType.ID == models.X01HANDICAP {
		scores, err := GetPlayersScore(int(game.CurrentMatchID.Int64))
		if err != nil {
			return nil, err
		}
		for _, player := range scores {
			handicaps[player.PlayerID] = player.Handicap
		}
	}

	for idx, playerID := range players {
		order := idx + 1
		res, err = tx.Exec("INSERT INTO player2match (player_id, match_id, `order`, game_id, handicap) VALUES (?, ?, ?, ?, ?)", playerID, matchID, order, gameID, handicaps[playerID])
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	log.Printf("[%d] Started new match", matchID)

	return GetMatch(int(matchID))
}

// FinishMatch will finalize a match by updating the winner and writing statistics for each player
func FinishMatch(visit models.Visit) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}
	match, err := GetMatch(visit.MatchID)
	if err != nil {
		return err
	}
	// Write statistics for each player
	game, err := GetGame(match.GameID)
	if err != nil {
		return err
	}

	err = AddVisit(visit)
	if err != nil {
		return err
	}

	// Update match with winner
	winnerID := visit.PlayerID
	if game.GameType.ID == models.SHOOTOUT {
		// For 9 Dart Shootout we need to check the scores of each player
		// to determine which player won the match with the highest score
		scores, err := GetPlayersScore(visit.MatchID)
		if err != nil {
			return err
		}
		highScore := 0
		for playerID, player := range scores {
			if player.CurrentScore > highScore {
				highScore = player.CurrentScore
				winnerID = playerID
			}
		}
	}
	_, err = tx.Exec(`UPDATE `+"`match`"+` SET current_player_id = ?, winner_id = ?, is_finished = 1, end_time = NOW() WHERE id = ?`,
		winnerID, winnerID, visit.MatchID)
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Printf("[%d] Finished with player %d winning", visit.MatchID, winnerID)

	if game.GameType.ID == models.SHOOTOUT {
		statisticsMap, err := calculateShootoutStatistics(visit.MatchID)
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_shootout(match_id, player_id, ppd, 60s_plus, 100s_plus, 140s_plus, 180s)
				VALUES (?, ?, ?, ?, ?, ?, ?)`, visit.MatchID, playerID, stats.PPD, stats.Score60sPlus,
				stats.Score100sPlus, stats.Score140sPlus, stats.Score180s)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting shootout statistics for player %d", visit.MatchID, playerID)
		}
	} else {
		statisticsMap, err := calculateX01Statistics(visit.MatchID, visit.PlayerID, match.StartingScore)
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_x01
					(match_id, player_id, ppd, first_nine_ppd, checkout_percentage, darts_thrown, 60s_plus,
					 100s_plus, 140s_plus, 180s, accuracy_20, accuracy_19, overall_accuracy)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, visit.MatchID, playerID, stats.PPD, stats.FirstNinePPD,
				stats.CheckoutPercentage, stats.DartsThrown, stats.Score60sPlus, stats.Score100sPlus, stats.Score140sPlus,
				stats.Score180s, stats.AccuracyStatistics.Accuracy20, stats.AccuracyStatistics.Accuracy19, stats.AccuracyStatistics.AccuracyOverall)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting x01 statistics for player %d", visit.MatchID, playerID)
		}
	}

	// Check if game is finished or not
	winsMap, err := GetWinsPerPlayer(game.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Determine how many matches has been played, and how many current player has won
	playedMatches := 1
	currentPlayerWins := 1
	for playerID, wins := range winsMap {
		playedMatches += wins
		if playerID == winnerID {
			currentPlayerWins += wins
		}
	}

	if currentPlayerWins == game.GameMode.WinsRequired {
		// Game finished, current player won
		_, err = tx.Exec("UPDATE game SET is_finished = 1, winner_id = ? WHERE id = ?", winnerID, game.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("Game %d finished with player %d winning", game.ID, winnerID)
		// Add owes between players in game
		if game.OweType != nil {
			for _, playerID := range game.Players {
				if playerID == winnerID {
					// Don't add payback to ourself
					continue
				}
				_, err = tx.Exec(`
					INSERT INTO owes (player_ower_id, player_owee_id, owe_type_id, amount)
					VALUES (?, ?, ?, 1)
					ON DUPLICATE KEY UPDATE amount = amount + 1`, playerID, visit.PlayerID, game.OweTypeID)
				if err != nil {
					tx.Rollback()
					return err
				}
				log.Printf("Added owes of %s from player %d to player %d", game.OweType.Item.String, playerID, visit.PlayerID)
			}
		}
	} else if game.GameMode.MatchesRequired.Valid && playedMatches == int(game.GameMode.MatchesRequired.Int64) {
		// Game finished, draw
		_, err = tx.Exec("UPDATE game SET is_finished = 1 WHERE id = ?", game.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("Game %d finished with a Draw", game.ID)
	} else {
		// Game is not finished
		log.Printf("Game %d is not finished, continuing to next leg", game.ID)
	}
	tx.Commit()
	return nil
}

// GetMatchesForGame returns all matches for the given game ID
func GetMatchesForGame(gameID int) ([]*models.Match, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id, end_time, starting_score, is_finished,
			current_player_id, winner_id, m.created_at, m.updated_at,
			m.game_id, GROUP_CONCAT(p2m.player_id ORDER BY p2m.order ASC)
		FROM `+"`match`"+` m
		LEFT JOIN player2match p2m ON p2m.match_id = m.id
		WHERE m.game_id = ?
		GROUP BY m.id
		ORDER BY id ASC`, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		m := new(models.Match)
		var players string
		err := rows.Scan(&m.ID, &m.Endtime, &m.StartingScore, &m.IsFinished, &m.CurrentPlayerID, &m.WinnerPlayerID, &m.CreatedAt, &m.UpdatedAt,
			&m.GameID, &players)
		if err != nil {
			return nil, err
		}
		m.Players = util.StringToIntArray(players)
		matches = append(matches, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// GetActiveMatches returns all matches which are currently live
func GetActiveMatches() ([]*models.Match, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id, end_time, starting_score, is_finished,
			current_player_id, winner_id, m.created_at, m.updated_at,
			m.game_id, GROUP_CONCAT(p2m.player_id ORDER BY p2m.order ASC)
		FROM ` + "`match`" + ` m
		LEFT JOIN player2match p2m ON p2m.match_id = m.id
		WHERE m.is_finished <> 1
		GROUP BY m.id
		ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		m := new(models.Match)
		var players string
		err := rows.Scan(&m.ID, &m.Endtime, &m.StartingScore, &m.IsFinished, &m.CurrentPlayerID, &m.WinnerPlayerID, &m.CreatedAt, &m.UpdatedAt,
			&m.GameID, &players)
		if err != nil {
			return nil, err
		}
		m.Players = util.StringToIntArray(players)
		matches = append(matches, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// GetMatch returns a match with the given ID
func GetMatch(id int) (*models.Match, error) {
	m := new(models.Match)
	var players string
	err := models.DB.QueryRow(`
		SELECT
			m.id, end_time, starting_score, is_finished, current_player_id, winner_id, m.created_at, m.updated_at, m.game_id,
			GROUP_CONCAT(DISTINCT p2m.player_id ORDER BY p2m.order ASC) AS 'players'
		FROM `+"`match` m"+`
		LEFT JOIN player2match p2m ON p2m.match_id = m.id
		WHERE m.id = ?`, id).Scan(&m.ID, &m.Endtime, &m.StartingScore, &m.IsFinished, &m.CurrentPlayerID, &m.WinnerPlayerID, &m.CreatedAt, &m.UpdatedAt, &m.GameID, &players)
	if err != nil {
		return nil, err
	}

	m.Players = util.StringToIntArray(players)
	visits, err := GetMatchVisits(id)
	if err != nil {
		return nil, err
	}
	dartsThrown := 0
	visitCount := 0
	for _, visit := range visits {
		if visitCount%len(m.Players) == 0 {
			dartsThrown += 3
		}
		visit.DartsThrown = dartsThrown
		visitCount++
	}
	// When checking out, it might be done in 1, 2 or 3 darts, so make
	// sure we set the correct number of darts thrown for the final visit
	if len(visits) > 0 {
		v := visits[len(visits)-1]
		v.DartsThrown = v.DartsThrown - 3 + v.GetDartsThrown()
	}

	m.Visits = visits
	m.Hits, m.DartsThrown = models.GetHitsMap(visits)

	return m, nil
}

// GetMatchPlayers returns a information about current score for players in a match
func GetMatchPlayers(id int) ([]*models.Player2Match, error) {
	match, err := GetMatch(id)
	if err != nil {
		return nil, err
	}

	scores, err := GetPlayersScore(id)
	if err != nil {
		return nil, err
	}
	lowestScore := match.StartingScore
	players := make([]*models.Player2Match, 0)
	for _, player := range scores {
		player.Modifiers = new(models.PlayerModifiers)
		if player.CurrentScore < lowestScore {
			lowestScore = player.CurrentScore
		}
		players = append(players, player)
	}

	winsMap, err := GetWinsPerPlayer(match.GameID)
	if err != nil {
		return nil, err
	}

	lastVisits, err := GetLastVisits(match.ID, len(match.Players))
	if err != nil {
		return nil, err
	}

	for _, player := range players {
		player.Wins = winsMap[player.PlayerID]
		if visit, ok := lastVisits[player.PlayerID]; ok {
			player.Modifiers.IsViliusVisit = visit.IsViliusVisit()
		}
		if lowestScore < 171 && player.CurrentScore > 199 {
			player.Modifiers.IsBeerGame = true
		}
	}

	return players, nil
}

// ChangePlayerOrder update the player order and current player for a given match
func ChangePlayerOrder(matchID int, orderMap map[string]int) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}
	for playerID, order := range orderMap {
		_, err = tx.Exec("UPDATE player2match SET `order` = ? WHERE player_id = ? AND match_id = ?", order, playerID, matchID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if order == 1 {
			_, err = tx.Exec("UPDATE `match` SET current_player_id = ? WHERE id = ?", playerID, matchID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()

	log.Printf("[%d] Changed player order to %v", matchID, orderMap)

	return nil
}

// DeleteMatch will delete the current match and update game with previous match
func DeleteMatch(matchID int) error {
	match, err := GetMatch(matchID)
	if err != nil {
		return err
	}

	game, err := GetGame(match.GameID)
	if err != nil {
		return err
	}

	return models.Transaction(models.DB, func(tx *sql.Tx) error {
		if _, err = tx.Exec("DELETE FROM `match` WHERE id = ?", matchID); err != nil {
			return err
		}
		log.Printf("[%d] Deleted match", matchID)

		var previousMatch *int
		err := models.DB.QueryRow("SELECT MAX(id) FROM `match` WHERE game_id = ? AND is_finished = 1", game.ID).Scan(&previousMatch)
		if err != nil {
			return err
		}
		if previousMatch == nil {
			if _, err = tx.Exec("DELETE FROM game WHERE id = ?", game.ID); err != nil {
				return err
			}
			log.Printf("Delete game without any match %d", game.ID)
		} else {
			_, err = tx.Exec("UPDATE game SET current_match_id = ? WHERE id = ?", previousMatch, game.ID)
			if err != nil {
				return err
			}
			log.Printf("[%d] Updated current match of game %d", previousMatch, game.ID)
		}
		return nil
	})
}

func calculateX01Statistics(matchID int, winnerID int, startingScore int) (map[int]*models.StatisticsX01, error) {
	visits, err := GetMatchVisits(matchID)
	if err != nil {
		return nil, err
	}

	players, err := GetMatchPlayers(matchID)
	if err != nil {
		return nil, err
	}
	statisticsMap := make(map[int]*models.StatisticsX01)
	playersMap := make(map[int]*models.Player2Match)
	for _, player := range players {
		stats := new(models.StatisticsX01)
		stats.AccuracyStatistics = new(models.AccuracyStatistics)
		statisticsMap[player.PlayerID] = stats

		playersMap[player.PlayerID] = player
		player.CurrentScore = startingScore
		if player.Handicap.Valid {
			player.CurrentScore += int(player.Handicap.Int64)
		}
	}

	for _, visit := range visits {
		player := playersMap[visit.PlayerID]
		stats := statisticsMap[visit.PlayerID]

		currentScore := player.CurrentScore
		if visit.FirstDart.IsCheckoutAttempt(currentScore) {
			stats.CheckoutAttempts++
		}
		currentScore -= visit.FirstDart.GetScore()
		if visit.SecondDart.IsCheckoutAttempt(currentScore) {
			stats.CheckoutAttempts++
		}
		currentScore -= visit.SecondDart.GetScore()
		if visit.ThirdDart.IsCheckoutAttempt(currentScore) {
			stats.CheckoutAttempts++
		}
		currentScore -= visit.ThirdDart.GetScore()

		if visit.IsBust {
			continue
		}

		visitScore := visit.GetScore()
		if stats.DartsThrown < 9 {
			stats.FirstNinePPD += float32(visitScore)
		}
		stats.PPD += float32(visitScore)

		if visitScore >= 60 && visitScore < 100 {
			stats.Score60sPlus++
		} else if visitScore >= 100 && visitScore < 140 {
			stats.Score100sPlus++
		} else if visitScore >= 140 && visitScore < 180 {
			stats.Score140sPlus++
		} else if visitScore == 180 {
			stats.Score180s++
		}

		// Get accuracy stats
		accuracyScore := player.CurrentScore
		if visit.FirstDart.Value.Valid {
			stats.AccuracyStatistics.GetAccuracyStats(accuracyScore, visit.FirstDart)
			stats.DartsThrown++
			accuracyScore -= visit.FirstDart.GetScore()
		}
		if visit.SecondDart.Value.Valid {
			stats.AccuracyStatistics.GetAccuracyStats(accuracyScore, visit.SecondDart)
			stats.DartsThrown++
			accuracyScore -= visit.SecondDart.GetScore()
		}
		if visit.ThirdDart.Value.Valid {
			stats.AccuracyStatistics.GetAccuracyStats(accuracyScore, visit.ThirdDart)
			stats.DartsThrown++
			accuracyScore -= visit.ThirdDart.GetScore()
		}

		player.CurrentScore = currentScore
	}

	for playerID, stats := range statisticsMap {
		stats.PPD = stats.PPD / float32(stats.DartsThrown)
		stats.FirstNinePPD = stats.FirstNinePPD / 9
		if playerID == winnerID {
			stats.CheckoutPercentage = null.FloatFrom(float64(100 / stats.CheckoutAttempts))
		} else {
			stats.CheckoutPercentage = null.FloatFromPtr(nil)
		}

		stats.AccuracyStatistics.SetAccuracy()
	}

	return statisticsMap, nil
}

func calculateShootoutStatistics(matchID int) (map[int]*models.StatisticsShootout, error) {
	visits, err := GetMatchVisits(matchID)
	if err != nil {
		return nil, err
	}

	players, err := GetMatchPlayers(matchID)
	if err != nil {
		return nil, err
	}
	statisticsMap := make(map[int]*models.StatisticsShootout)
	for _, player := range players {
		stats := new(models.StatisticsShootout)
		statisticsMap[player.PlayerID] = stats
	}

	for _, visit := range visits {
		stats := statisticsMap[visit.PlayerID]

		visitScore := visit.GetScore()
		stats.PPD += float32(visitScore)

		if visitScore >= 60 && visitScore < 100 {
			stats.Score60sPlus++
		} else if visitScore >= 100 && visitScore < 140 {
			stats.Score100sPlus++
		} else if visitScore >= 140 && visitScore < 180 {
			stats.Score140sPlus++
		} else if visitScore == 180 {
			stats.Score180s++
		}
	}

	for _, stats := range statisticsMap {
		stats.PPD = stats.PPD / float32(9)
	}
	return statisticsMap, nil
}
