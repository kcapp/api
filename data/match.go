package data

import (
	"log"

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
		return nil, err
	}
	matchID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec("UPDATE game SET current_match_id = ? WHERE id = ?", matchID, gameID)
	if err != nil {
		return nil, err
	}

	for idx, playerID := range players {
		order := idx + 1
		res, err = tx.Exec("INSERT INTO player2match (player_id, match_id, `order`, game_id) VALUES (?, ?, ?, ?)", playerID, matchID, order, gameID)
		if err != nil {
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
	if game.GameType.ID == 2 {
		// For 9 Dart Shootout we need to check the scores of each player
		// to determine which player won the match with the highest score
		scores, err := GetPlayersScore(visit.MatchID)
		if err != nil {
			return err
		}
		highScore := 0
		for playerID, score := range scores {
			if score > highScore {
				highScore = score
				winnerID = playerID
			}
		}
	}
	_, err = tx.Exec(`UPDATE `+"`match`"+` SET current_player_id = ?, winner_id = ?, is_finished = 1, end_time = NOW() WHERE id = ?`,
		winnerID, winnerID, visit.MatchID)
	if err != nil {
		return err
	}

	if game.GameType.ID == 1 {
		statisticsMap, err := calculateStatistics(visit.MatchID, visit.PlayerID, match.StartingScore)
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_x01
					(match_id, player_id, ppd, first_nine_ppd, checkout_percentage, darts_thrown, 60s_plus,
					 100s_plus, 140s_plus, 180s, accuracy_20, accuracy_19, overall_accuracy)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, visit.MatchID, playerID, stats.PPD, stats.FirstNinePPD,
				stats.CheckoutPercentage, stats.DartsThrown, stats.Score60sPlus, stats.Score100sPlus, stats.Score140sPlus,
				stats.Score180s, stats.AccuracyStatistics.Accuracy20, stats.AccuracyStatistics.Accuracy19, stats.AccuracyStatistics.AccuracyOverall)
			if err != nil {
				return err
			}
			log.Printf("[%d] Inserting match statistics for player %d", visit.MatchID, playerID)
		}
	}

	// Check if game is finished or not
	winsMap, err := GetWinsPerPlayer(game.ID)
	if err != nil {
		return err
	}

	// Determine how many matches has been played, and how many current player has won
	playedMatches := 1
	currentPlayerWins := 1
	for playerID, wins := range winsMap {
		playedMatches += wins
		if playerID == visit.PlayerID {
			currentPlayerWins += wins
		}
	}

	if currentPlayerWins == game.GameMode.WinsRequired {
		// Game finished, current player won
		_, err = tx.Exec("UPDATE game SET is_finished = 1, winner_id = ? WHERE id = ?", visit.PlayerID, game.ID)
		if err != nil {
			return err
		}
		log.Printf("Game %d finished with player %d winning", game.ID, visit.PlayerID)
		// Add owes between players in game
		if game.OweType != nil {
			for _, playerID := range game.Players {
				if playerID == visit.PlayerID {
					// Don't add payback to ourself
					continue
				}
				_, err = tx.Exec(`
					INSERT INTO owes (player_ower_id, player_owee_id, owe_type_id, amount)
					VALUES (?, ?, ?, 1)
					ON DUPLICATE KEY UPDATE amount = amount + 1`, playerID, visit.PlayerID, game.OweTypeID)
				if err != nil {
					return err
				}
				log.Printf("Added owes of %s from player %d to player %d", game.OweType.Item.String, playerID, visit.PlayerID)
			}
		}
	} else if game.GameMode.MatchesRequired.Valid && playedMatches == int(game.GameMode.MatchesRequired.Int64) {
		// Game finished, draw
		_, err = tx.Exec("UPDATE game SET is_finished = 1 WHERE id = ?", game.ID)
		if err != nil {
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
	m.Visits = visits
	m.Hits, m.DartsThrown = models.GetHitsMap(visits)

	return m, nil
}

// GetMatchPlayers returns a information about current score for players in a match
func GetMatchPlayers(id int) ([]*models.Player2Match, error) {
	rows, err := models.DB.Query(`
		SELECT
			p2m.match_id,
			p2m.player_id,
			p2m.order,
			p2m.player_id = m.current_player_id AS 'is_current_player',
			m.starting_score - (IFNULL(SUM(first_dart * first_dart_multiplier), 0) +
				IFNULL(SUM(second_dart * second_dart_multiplier), 0) +
				IFNULL(SUM(third_dart * third_dart_multiplier), 0))
				* IF(g.game_type_id = 2,  -1, 1) AS 'current_score'
		FROM player2match p2m
		LEFT JOIN `+"`match`"+` m ON m.id = p2m.match_id
		LEFT JOIN score s ON s.match_id = p2m.match_id AND s.player_id = p2m.player_id
		LEFT JOIN game g on g.id = m.game_id
		WHERE p2m.match_id = ? AND (s.is_bust IS NULL OR is_bust = 0)
		GROUP BY p2m.player_id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]*models.Player2Match, 0)
	for rows.Next() {
		p2m := new(models.Player2Match)
		err := rows.Scan(&p2m.MatchID, &p2m.PlayerID, &p2m.Order, &p2m.IsCurrentPlayer, &p2m.CurrentScore)
		if err != nil {
			return nil, err
		}
		players = append(players, p2m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	match, err := GetMatch(id)
	if err != nil {
		return nil, err
	}
	winsMap, err := GetWinsPerPlayer(match.GameID)
	if err != nil {
		return nil, err
	}
	for _, player := range players {
		player.Wins = winsMap[player.PlayerID]
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
			return err
		}
		if order == 1 {
			_, err = tx.Exec("UPDATE `match` SET current_player_id = ? WHERE id = ?", playerID, matchID)
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()

	log.Printf("[%d] Changed player order to %v", matchID, orderMap)

	return nil
}

func calculateStatistics(matchID int, winnerID int, startingScore int) (map[int]*models.StatisticsX01, error) {
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
			stats.CheckoutPercentage = float32(100 / stats.CheckoutAttempts)
		} else {
			stats.CheckoutPercentage = 0
		}

		stats.AccuracyStatistics.SetAccuracy()
	}

	return statisticsMap, nil
}
