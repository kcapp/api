package data

import (
	"database/sql"
	"log"

	"github.com/guregu/null"
	"github.com/kcapp/api/models"
	"github.com/kcapp/api/util"
)

// NewLeg will create a new leg for the given match
func NewLeg(matchID int, startingScore int, players []int) (*models.Leg, error) {
	tx, err := models.DB.Begin()
	if err != nil {
		return nil, err
	}

	// Shift players to get correct order
	id, players := players[0], players[1:]
	players = append(players, id)
	res, err := tx.Exec("INSERT INTO leg (starting_score, current_player_id, match_id, created_at) VALUES (?, ?, ?, NOW()) ",
		startingScore, players[0], matchID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	legID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	match, err := GetMatch(matchID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("UPDATE matches SET current_leg_id = ?, updated_at = NOW() WHERE id = ?", legID, matchID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	handicaps := make(map[int]null.Int)
	if match.MatchType.ID == models.X01HANDICAP {
		scores, err := GetPlayersScore(int(match.CurrentLegID.Int64))
		if err != nil {
			return nil, err
		}
		for _, player := range scores {
			handicaps[player.PlayerID] = player.Handicap
		}
	}

	for idx, playerID := range players {
		order := idx + 1
		res, err = tx.Exec("INSERT INTO player2leg (player_id, leg_id, `order`, match_id, handicap) VALUES (?, ?, ?, ?, ?)", playerID, legID, order, matchID, handicaps[playerID])
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	log.Printf("[%d] Started new leg", legID)

	return GetLeg(int(legID))
}

// FinishLeg will finalize a leg by updating the winner and writing statistics for each player
func FinishLeg(visit models.Visit) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}
	leg, err := GetLeg(visit.LegID)
	if err != nil {
		return err
	}
	// Write statistics for each player
	match, err := GetMatch(leg.MatchID)
	if err != nil {
		return err
	}

	err = AddVisit(visit)
	if err != nil {
		return err
	}

	// Update leg with winner
	winnerID := visit.PlayerID
	if match.MatchType.ID == models.SHOOTOUT {
		// For 9 Dart Shootout we need to check the scores of each player
		// to determine which player won the leg with the highest score
		scores, err := GetPlayersScore(visit.LegID)
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
	_, err = tx.Exec(`UPDATE leg SET current_player_id = ?, winner_id = ?, is_finished = 1, end_time = NOW() WHERE id = ?`,
		winnerID, winnerID, visit.LegID)
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Printf("[%d] Finished with player %d winning", visit.LegID, winnerID)

	if match.MatchType.ID == models.SHOOTOUT {
		statisticsMap, err := calculateShootoutStatistics(visit.LegID)
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_shootout(leg_id, player_id, ppd, 60s_plus, 100s_plus, 140s_plus, 180s)
				VALUES (?, ?, ?, ?, ?, ?, ?)`, visit.LegID, playerID, stats.PPD, stats.Score60sPlus,
				stats.Score100sPlus, stats.Score140sPlus, stats.Score180s)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting shootout statistics for player %d", visit.LegID, playerID)
		}
	} else {
		statisticsMap, err := calculateX01Statistics(visit.LegID, visit.PlayerID, leg.StartingScore)
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_x01
					(leg_id, player_id, ppd, first_nine_ppd, checkout_percentage, checkout_attempts, darts_thrown, 60s_plus,
					 100s_plus, 140s_plus, 180s, accuracy_20, accuracy_19, overall_accuracy)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, visit.LegID, playerID, stats.PPD, stats.FirstNinePPD,
				stats.CheckoutPercentage, stats.CheckoutAttempts, stats.DartsThrown, stats.Score60sPlus, stats.Score100sPlus, stats.Score140sPlus,
				stats.Score180s, stats.AccuracyStatistics.Accuracy20, stats.AccuracyStatistics.Accuracy19, stats.AccuracyStatistics.AccuracyOverall)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting x01 statistics for player %d", visit.LegID, playerID)
		}
	}

	// Check if match is finished or not
	winsMap, err := GetWinsPerPlayer(match.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Determine how many legs has been played, and how many current player has won
	playedLegs := 1
	currentPlayerWins := 1
	for playerID, wins := range winsMap {
		playedLegs += wins
		if playerID == winnerID {
			currentPlayerWins += wins
		}
	}

	isFinished := false
	if currentPlayerWins == match.MatchMode.WinsRequired {
		// Match finished, current player won
		isFinished = true
		_, err = tx.Exec("UPDATE matches SET is_finished = 1, winner_id = ? WHERE id = ?", winnerID, match.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		// Add owes between players in match
		if match.OweType != nil {
			for _, playerID := range match.Players {
				if playerID == winnerID {
					// Don't add payback to ourself
					continue
				}
				_, err = tx.Exec(`
					INSERT INTO owes (player_ower_id, player_owee_id, owe_type_id, amount)
					VALUES (?, ?, ?, 1)
					ON DUPLICATE KEY UPDATE amount = amount + 1`, playerID, visit.PlayerID, match.OweTypeID)
				if err != nil {
					tx.Rollback()
					return err
				}
				log.Printf("Added owes of %s from player %d to player %d", match.OweType.Item.String, playerID, visit.PlayerID)
			}
		}
		log.Printf("Match %d finished with player %d winning", match.ID, winnerID)
	} else if match.MatchMode.LegsRequired.Valid && playedLegs == int(match.MatchMode.LegsRequired.Int64) {
		// Match finished, draw
		isFinished = true
		_, err = tx.Exec("UPDATE matches SET is_finished = 1 WHERE id = ?", match.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("Match %d finished with a Draw", match.ID)
	} else {
		// Match is not finished
		log.Printf("Match %d is not finished, continuing to next leg", match.ID)
	}
	tx.Commit()

	if isFinished {
		// Update Elo for players if match is finished
		err = UpdateEloForMatch(match.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// UndoLegFinish will undo a finalized leg
func UndoLegFinish(legID int) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}

	// Undo the finalized match
	_, err = tx.Exec("UPDATE matches SET is_finished = 0, winner_id = NULL WHERE id = (SELECT match_id FROM leg WHERE id = ?)", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Undo the finalized leg
	_, err = tx.Exec("UPDATE leg SET is_finished = 0, winner_id = NULL WHERE id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Remove generated statistics for the leg
	_, err = tx.Exec("DELETE FROM statistics_x01 WHERE leg_id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Remove the last score
	_, err = tx.Exec("DELETE FROM score WHERE leg_id = ? ORDER BY id DESC LIMIT 1", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	log.Printf("[%d] Undo finish of leg", legID)
	return nil
}

// GetLegsForMatch returns all legs for the given match ID
func GetLegsForMatch(matchID int) ([]*models.Leg, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id, l.end_time, l.starting_score, l.is_finished,
			l.current_player_id, l.winner_id, l.created_at, l.updated_at,
			l.match_id, l.has_scores, GROUP_CONCAT(p2l.player_id ORDER BY p2l.order ASC)
		FROM leg l
			LEFT JOIN player2leg p2l ON p2l.leg_id = l.id
		WHERE l.match_id = ?
		GROUP BY l.id
		ORDER BY l.id ASC`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs := make([]*models.Leg, 0)
	for rows.Next() {
		m := new(models.Leg)
		var players string
		err := rows.Scan(&m.ID, &m.Endtime, &m.StartingScore, &m.IsFinished, &m.CurrentPlayerID, &m.WinnerPlayerID, &m.CreatedAt, &m.UpdatedAt,
			&m.MatchID, &m.HasScores, &players)
		if err != nil {
			return nil, err
		}
		m.Players = util.StringToIntArray(players)
		legs = append(legs, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return legs, nil
}

// GetActiveLegs returns all legs which are currently live
func GetActiveLegs() ([]*models.Leg, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id, l.end_time, l.starting_score, l.is_finished,
			l.current_player_id, l.winner_id, l.created_at, l.updated_at,
			l.match_id, l.has_scores, GROUP_CONCAT(p2l.player_id ORDER BY p2l.order ASC)
		FROM leg l
			LEFT JOIN player2leg p2l ON p2l.leg_id = l.id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE l.is_finished <> 1 AND m.is_abandoned = 0
		GROUP BY l.id
		ORDER BY l.id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs := make([]*models.Leg, 0)
	for rows.Next() {
		leg := new(models.Leg)
		var players string
		err := rows.Scan(&leg.ID, &leg.Endtime, &leg.StartingScore, &leg.IsFinished, &leg.CurrentPlayerID, &leg.WinnerPlayerID, &leg.CreatedAt,
			&leg.UpdatedAt, &leg.MatchID, &leg.HasScores, &players)
		if err != nil {
			return nil, err
		}
		leg.Players = util.StringToIntArray(players)
		legs = append(legs, leg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return legs, nil
}

// GetLeg returns a leg with the given ID
func GetLeg(id int) (*models.Leg, error) {
	leg := new(models.Leg)
	var players string
	err := models.DB.QueryRow(`
		SELECT
			l.id, l.end_time, l.starting_score, l.is_finished, l.current_player_id, l.winner_id, l.created_at, l.updated_at,
			l.board_stream_url, l.match_id, l.has_scores, GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order ASC) AS 'players'
		FROM leg l
			LEFT JOIN player2leg p2l ON p2l.leg_id = l.id
		WHERE l.id = ?`, id).Scan(&leg.ID, &leg.Endtime, &leg.StartingScore, &leg.IsFinished, &leg.CurrentPlayerID, &leg.WinnerPlayerID, &leg.CreatedAt,
		&leg.UpdatedAt, &leg.BoardStreamURL, &leg.MatchID, &leg.HasScores, &players)
	if err != nil {
		return nil, err
	}

	leg.Players = util.StringToIntArray(players)
	visits, err := GetLegVisits(id)
	if err != nil {
		return nil, err
	}
	dartsThrown := 0
	visitCount := 0
	for _, visit := range visits {
		if visitCount%len(leg.Players) == 0 {
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

	leg.Visits = visits
	leg.Hits, leg.DartsThrown = models.GetHitsMap(visits)
	leg.CheckoutStatistics, err = getCheckoutStatistics(leg.ID, leg.StartingScore)
	if err != nil {
		return nil, err
	}

	return leg, nil
}

// GetLegPlayers returns a information about current score for players in a leg
func GetLegPlayers(id int) ([]*models.Player2Leg, error) {
	leg, err := GetLeg(id)
	if err != nil {
		return nil, err
	}

	scores, err := GetPlayersScore(id)
	if err != nil {
		return nil, err
	}
	lowestScore := leg.StartingScore
	players := make([]*models.Player2Leg, 0)
	for _, player := range scores {
		player.Modifiers = new(models.PlayerModifiers)
		if player.CurrentScore < lowestScore {
			lowestScore = player.CurrentScore
		}
		players = append(players, player)
	}

	winsMap, err := GetWinsPerPlayer(leg.MatchID)
	if err != nil {
		return nil, err
	}

	lastVisits, err := GetLastVisits(leg.ID, len(leg.Players))
	if err != nil {
		return nil, err
	}

	for _, player := range players {
		player.Wins = winsMap[player.PlayerID]
		if visit, ok := lastVisits[player.PlayerID]; ok {
			player.Modifiers.IsViliusVisit = visit.IsViliusVisit()
			player.Modifiers.IsFishAndChips = visit.IsFishAndChips()
		}
		if lowestScore < 171 && player.CurrentScore > 199 {
			player.Modifiers.IsBeerMatch = true
		}
		player.AddVisitStatistics(*leg)
	}

	return players, nil
}

// ChangePlayerOrder update the player order and current player for a given leg
func ChangePlayerOrder(legID int, orderMap map[string]int) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}
	for playerID, order := range orderMap {
		_, err = tx.Exec("UPDATE player2leg SET `order` = ? WHERE player_id = ? AND leg_id = ?", order, playerID, legID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if order == 1 {
			_, err = tx.Exec("UPDATE leg SET current_player_id = ?, updated_at = NOW() WHERE id = ?", playerID, legID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()

	log.Printf("[%d] Changed player order to %v", legID, orderMap)

	return nil
}

// DeleteLeg will delete the current leg and update match with previous leg
func DeleteLeg(legID int) error {
	leg, err := GetLeg(legID)
	if err != nil {
		return err
	}

	match, err := GetMatch(leg.MatchID)
	if err != nil {
		return err
	}

	return models.Transaction(models.DB, func(tx *sql.Tx) error {
		if _, err = tx.Exec("DELETE FROM leg WHERE id = ?", legID); err != nil {
			return err
		}
		log.Printf("[%d] Deleted leg", legID)

		var previousLeg *int
		err := models.DB.QueryRow("SELECT MAX(id) FROM leg WHERE match_id = ? AND is_finished = 1", match.ID).Scan(&previousLeg)
		if err != nil {
			return err
		}
		if previousLeg == nil {
			if _, err = tx.Exec("DELETE FROM matches WHERE id = ?", match.ID); err != nil {
				return err
			}
			log.Printf("Delete match without any leg %d", match.ID)
		} else {
			_, err = tx.Exec("UPDATE matches SET current_leg_id = ? WHERE id = ?", previousLeg, match.ID)
			if err != nil {
				return err
			}
			log.Printf("[%d] Updated current leg of match %d", previousLeg, match.ID)
		}
		return nil
	})
}

func calculateX01Statistics(legID int, winnerID int, startingScore int) (map[int]*models.StatisticsX01, error) {
	visits, err := GetLegVisits(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}
	statisticsMap := make(map[int]*models.StatisticsX01)
	playersMap := make(map[int]*models.Player2Leg)
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

		stats.DartsThrown += 3
		if visit.IsBust {
			continue
		}

		visitScore := visit.GetScore()
		if stats.DartsThrown <= 9 {
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
			accuracyScore -= visit.FirstDart.GetScore()
		}
		if visit.SecondDart.Value.Valid {
			stats.AccuracyStatistics.GetAccuracyStats(accuracyScore, visit.SecondDart)
			accuracyScore -= visit.SecondDart.GetScore()
		}
		if visit.ThirdDart.Value.Valid {
			stats.AccuracyStatistics.GetAccuracyStats(accuracyScore, visit.ThirdDart)
			accuracyScore -= visit.ThirdDart.GetScore()
		}
		player.CurrentScore = currentScore
	}

	for playerID, stats := range statisticsMap {
		stats.PPD = stats.PPD / float32(stats.DartsThrown)
		stats.FirstNinePPD = stats.FirstNinePPD / 9
		if playerID == winnerID {
			stats.CheckoutPercentage = null.FloatFrom(100 / float64(stats.CheckoutAttempts))

			// When checking out, it might be done in 1, 2 or 3 darts, so make
			// sure we set the correct number of darts thrown for the final visit
			v := visits[len(visits)-1]
			stats.DartsThrown = stats.DartsThrown - 3 + v.GetDartsThrown()
		} else {
			stats.CheckoutPercentage = null.FloatFromPtr(nil)
		}
		stats.AccuracyStatistics.SetAccuracy()
	}

	return statisticsMap, nil
}

func calculateShootoutStatistics(legID int) (map[int]*models.StatisticsShootout, error) {
	visits, err := GetLegVisits(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
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

// RecalculateX01Statistics will recalculate x01 statistics for all legs
func RecalculateX01Statistics() (map[int]map[int]*models.StatisticsX01, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id, l.end_time, l.starting_score, l.is_finished,
			l.current_player_id, l.winner_id, l.created_at, l.updated_at,
			l.match_id, GROUP_CONCAT(p2l.player_id ORDER BY p2l.order ASC)
		FROM leg l
			JOIN matches m on m.id = l.match_id
			JOIN player2leg p2l ON p2l.leg_id = l.id
		WHERE
			l.id IN(6415, 6413, 6411, 6389, 6387, 6385, 6374, 6373, 6371)
		GROUP BY l.id
		ORDER BY l.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	legs := make([]*models.Leg, 0)
	for rows.Next() {
		m := new(models.Leg)
		var players string
		err := rows.Scan(&m.ID, &m.Endtime, &m.StartingScore, &m.IsFinished, &m.CurrentPlayerID, &m.WinnerPlayerID, &m.CreatedAt, &m.UpdatedAt,
			&m.MatchID, &players)
		if err != nil {
			return nil, err
		}
		m.Players = util.StringToIntArray(players)
		legs = append(legs, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	m := make(map[int]map[int]*models.StatisticsX01)
	for _, leg := range legs {
		stats, err := calculateX01Statistics(leg.ID, int(leg.WinnerPlayerID.Int64), leg.StartingScore)
		if err != nil {
			return nil, err
		}
		for playerID, stat := range stats {
			log.Printf("UPDATE statistics_x01 SET first_nine_ppd = %f  WHERE leg_id = %d AND player_id = %d;",
				stat.FirstNinePPD, leg.ID, playerID)

		}
		m[leg.ID] = stats
	}

	/*s := make([]*models.CheckoutStatistics, 0)
	for _, leg := range legs {
		log.Printf("Getting statistics for %d", leg.ID)
		stats, err := getCheckoutStatistics(leg.ID, leg.StartingScore)
		if err != nil {
			return nil, err
		}
		s = append(s, stats)
	}

	all := make(map[int]int)
	for _, stats := range s {
		log.Printf("Checkout: %d, Total: %d, Attempts: %d", stats.Checkout, stats.Count, stats.CheckoutAttempts)

		for checkout, count := range stats.CheckoutAttempts {
			all[checkout] += count
		}
	}
	log.Printf("All: %v", all)*/

	return m, err
}

// getCheckoutStatistics will get all checkout attempts for the given leg
func getCheckoutStatistics(legID int, startingScore int) (*models.CheckoutStatistics, error) {
	visits, err := GetLegVisits(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}

	playersMap := make(map[int]*models.Player2Leg)
	for _, player := range players {
		playersMap[player.PlayerID] = player
		player.CurrentScore = startingScore
		if player.Handicap.Valid {
			player.CurrentScore += int(player.Handicap.Int64)
		}
	}

	totalAttempts := 0
	checkoutAttempts := make(map[int]int)
	for _, visit := range visits {
		player := playersMap[visit.PlayerID]

		currentScore := player.CurrentScore
		if visit.FirstDart.IsCheckoutAttempt(currentScore) {
			totalAttempts++
			checkoutAttempts[currentScore]++
		}
		currentScore -= visit.FirstDart.GetScore()

		if visit.SecondDart.IsCheckoutAttempt(currentScore) {
			totalAttempts++
			checkoutAttempts[currentScore]++
		}
		currentScore -= visit.SecondDart.GetScore()

		if visit.ThirdDart.IsCheckoutAttempt(currentScore) {
			totalAttempts++
			checkoutAttempts[currentScore]++
		}
		currentScore -= visit.ThirdDart.GetScore()

		if !visit.IsBust {
			player.CurrentScore = currentScore
		}
	}

	statistics := new(models.CheckoutStatistics)
	statistics.CheckoutAttempts = checkoutAttempts
	statistics.Count = totalAttempts

	if len(visits) > 1 {
		lastVisit := visits[len(visits)-1]
		if lastVisit.ThirdDart.Value.Valid {
			statistics.Checkout = lastVisit.ThirdDart.GetScore()
		} else if lastVisit.SecondDart.Value.Valid {
			statistics.Checkout = lastVisit.SecondDart.GetScore()
		} else {
			statistics.Checkout = lastVisit.FirstDart.GetScore()
		}
	}

	return statistics, nil
}
