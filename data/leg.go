package data

import (
	"database/sql"
	"log"
	"sort"

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
		res, err = tx.Exec("INSERT INTO player2leg (player_id, leg_id, `order`, match_id, handicap) VALUES (?, ?, ?, ?, ?)",
			playerID, legID, order, matchID, handicaps[playerID])
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

	_, err = AddVisit(visit)
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
		statisticsMap, err := CalculateShootoutStatistics(visit.LegID)
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
		statisticsMap, err := CalculateX01Statistics(visit.LegID, visit.PlayerID, leg.StartingScore)
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_x01
					(leg_id, player_id, ppd, ppd_score, first_nine_ppd, first_nine_ppd_score, checkout_percentage, checkout_attempts, darts_thrown, 60s_plus,
					 100s_plus, 140s_plus, 180s, accuracy_20, accuracy_19, overall_accuracy)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, visit.LegID, playerID, stats.PPD, stats.PPDScore, stats.FirstNinePPD, stats.FirstNinePPDScore,
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

// FinishLegNew will finalize a leg by updating the winner and writing statistics for each player
func FinishLegNew(visit models.Visit) error {
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

	// Update leg with winner
	winnerID := visit.PlayerID
	if match.MatchType.ID == models.SHOOTOUT || match.MatchType.ID == models.DARTSATX {
		// For "9 Dart Shootout" and "Darts at X" we need to check the scores of each player
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
		statisticsMap, err := CalculateShootoutStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
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
	} else if match.MatchType.ID == models.CRICKET {
		statisticsMap, err := CalculateCricketStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_cricket
					(leg_id, player_id, total_marks, rounds, score, first_nine_marks, mpr, first_nine_mpr, marks5, marks6, marks7, marks8, marks9)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, visit.LegID, playerID, stats.TotalMarks, stats.Rounds, stats.Score, stats.FirstNineMarks,
				stats.MPR, stats.FirstNineMPR, stats.Marks5, stats.Marks6, stats.Marks7, stats.Marks8, stats.Marks9)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting cricket statistics for player %d", visit.LegID, playerID)
		}
	} else if match.MatchType.ID == models.DARTSATX {
		statisticsMap, err := CalculateDartsAtXStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_darts_at_x
					(leg_id, player_id, score, singles, doubles, triples, hit_rate, hits5, hits6, hits7, hits8, hits9)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, visit.LegID, playerID, stats.Score, stats.Singles, stats.Doubles, stats.Triples, stats.HitRate,
				stats.Hits5, stats.Hits6, stats.Hits7, stats.Hits8, stats.Hits9)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting Darts At %d statistics for player %d", visit.LegID, leg.StartingScore, playerID)
		}
	} else {
		statisticsMap, err := CalculateX01Statistics(visit.LegID, visit.PlayerID, leg.StartingScore)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_x01
					(leg_id, player_id, ppd, ppd_score, first_nine_ppd, first_nine_ppd_score, checkout_percentage, checkout_attempts, darts_thrown, 60s_plus,
					 100s_plus, 140s_plus, 180s, accuracy_20, accuracy_19, overall_accuracy)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, visit.LegID, playerID, stats.PPD, stats.PPDScore, stats.FirstNinePPD, stats.FirstNinePPDScore,
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

		if match.TournamentID.Valid {
			metadata, err := GetMatchMetadata(match.ID)
			if err != nil {
				return err
			}

			if metadata.WinnerOutcomeMatchID.Valid {
				winnerMatch, err := GetMatch(int(metadata.WinnerOutcomeMatchID.Int64))
				if err != nil {
					return err
				}
				idx := 0
				if !metadata.IsWinnerOutcomeHome {
					idx = 1
				}
				err = SwapPlayers(winnerMatch.ID, winnerID, winnerMatch.Players[idx])
				if err != nil {
					return err
				}
			}
			if metadata.LooserOutcomeMatchID.Valid {
				looserID := getMatchLooser(match, winnerID)
				looserMatch, err := GetMatch(int(metadata.LooserOutcomeMatchID.Int64))
				if err != nil {
					return err
				}
				idx := 0
				if !metadata.IsLooserOutcomeHome {
					idx = 1
				}
				err = SwapPlayers(looserMatch.ID, looserID, looserMatch.Players[idx])
				if err != nil {
					return err
				}
			}
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
	_, err = tx.Exec("DELETE FROM statistics_shootout WHERE leg_id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM statistics_cricket WHERE leg_id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM statistics_darts_at_x WHERE leg_id = ?", legID)
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
	// Reset the calculated elo for the match
	_, err = tx.Exec(`UPDATE player_elo pe
			INNER JOIN player_elo_changelog pec ON pec.player_id = pe.player_id
		SET pe.current_elo = pec.old_elo,
			pe.current_elo_matches = pe.current_elo_matches - 1,
			pe.tournament_elo = IFNULL(pec.old_tournament_elo, pe.tournament_elo),
			pe.tournament_elo_matches = IF(pec.old_tournament_elo = NULL, pe.tournament_elo_matches, pe.tournament_elo_matches - 1)
		WHERE pe.player_id IN (SELECT player_id FROM player2leg WHERE leg_id = ?) AND pec.match_id = (SELECT match_id FROM leg WHERE id = ?)`, legID, legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Delete elo changelog for match
	_, err = tx.Exec("DELETE from player_elo_changelog WHERE match_id = (SELECT match_id FROM leg WHERE id = ?)", legID)
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
		leg := new(models.Leg)
		var players string
		err := rows.Scan(&leg.ID, &leg.Endtime, &leg.StartingScore, &leg.IsFinished, &leg.CurrentPlayerID,
			&leg.WinnerPlayerID, &leg.CreatedAt, &leg.UpdatedAt, &leg.MatchID, &leg.HasScores, &players)
		if err != nil {
			return nil, err
		}
		leg.Players = util.StringToIntArray(players)
		visits, err := GetLegVisits(leg.ID)
		if err != nil {
			return nil, err
		}
		leg.Visits = visits
		legs = append(legs, leg)
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
		WHERE l.is_finished <> 1 AND m.is_abandoned = 0  and m.is_walkover <> 1
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
	var matchType int
	err := models.DB.QueryRow(`
		SELECT
			l.id, l.end_time, l.starting_score, l.is_finished, l.current_player_id, l.winner_id, l.created_at, l.updated_at,
			l.board_stream_url, l.match_id, l.has_scores, GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order ASC) AS 'players', m.match_type_id
		FROM leg l
			LEFT JOIN player2leg p2l ON p2l.leg_id = l.id
			LEFT JOIN matches m ON m.id = l.match_id
		WHERE l.id = ?`, id).Scan(&leg.ID, &leg.Endtime, &leg.StartingScore, &leg.IsFinished, &leg.CurrentPlayerID, &leg.WinnerPlayerID, &leg.CreatedAt,
		&leg.UpdatedAt, &leg.BoardStreamURL, &leg.MatchID, &leg.HasScores, &players, &matchType)
	if err != nil {
		return nil, err
	}

	leg.Players = util.StringToIntArray(players)
	visits, err := GetLegVisits(id)
	if err != nil {
		return nil, err
	}
	scores := make(map[int]*models.Player2Leg)
	for i := 0; i < len(leg.Players); i++ {
		p2l := new(models.Player2Leg)
		p2l.Hits = make(map[int]*models.Hits)
		if matchType == models.DARTSATX {
			p2l.CurrentScore = 0
		} else if matchType == models.X01HANDICAP {
			// TODO
		} else {
			p2l.CurrentScore = leg.StartingScore
		}
		scores[leg.Players[i]] = p2l
	}

	dartsThrown := 0
	visitCount := 0
	for i, visit := range visits {
		if visitCount%len(leg.Players) == 0 {
			dartsThrown += 3
		}
		visit.DartsThrown = dartsThrown
		visitCount++

		if !visit.IsBust {
			score := visit.GetScore()
			if matchType == models.DARTSATX {
				score = 0
				if visit.FirstDart.ValueRaw() == leg.StartingScore {
					score += int(visit.FirstDart.Multiplier)
				}
				if visit.SecondDart.ValueRaw() == leg.StartingScore {
					score += int(visit.SecondDart.Multiplier)
				}
				if visit.ThirdDart.ValueRaw() == leg.StartingScore {
					score += int(visit.ThirdDart.Multiplier)
				}
			}

			if matchType == models.DARTSATX || matchType == models.SHOOTOUT {
				scores[visit.PlayerID].CurrentScore += score
			} else if matchType == models.CRICKET {
				score = visit.CalculateCricketScore(scores)
			} else {
				scores[visit.PlayerID].CurrentScore -= score
			}
			visit.Score = score
		}

		visit.Scores = make(map[int]int)
		visit.Scores[visit.PlayerID] = scores[visit.PlayerID].CurrentScore
		for j := 1; j < len(leg.Players); j++ {
			var next *models.Visit
			if len(visits) > len(leg.Players) {
				if i+j >= len(visits) && i-(len(leg.Players)-j) > 0 {
					// There is no next visit, so look at previous instead
					// Need to look in reverese order to keep the order of scores the same
					next = visits[i-(len(leg.Players)-j)]
				} else {
					next = visits[i+j]
				}
			}
			if next != nil {
				visit.Scores[next.PlayerID] = scores[next.PlayerID].CurrentScore
			}
		}
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

// GetLegPlayers returns information about current score for players in a leg
func GetLegPlayers(id int) ([]*models.Player2Leg, error) {
	leg, err := GetLeg(id)
	if err != nil {
		return nil, err
	}

	scores, err := GetPlayersScore(id)
	if err != nil {
		return nil, err
	}
	hitsMap := make(map[int]map[int]*models.Hits)
	lowestScore := leg.StartingScore
	players := make([]*models.Player2Leg, 0)
	for _, player := range scores {
		player.Modifiers = new(models.PlayerModifiers)
		if player.CurrentScore < lowestScore {
			lowestScore = player.CurrentScore
		}
		hitsMap[player.PlayerID] = make(map[int]*models.Hits)
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

	// Get hits on each number for each player
	for _, visit := range leg.Visits {
		m := hitsMap[visit.PlayerID]

		v := visit.FirstDart.ValueRaw()
		if _, ok := m[v]; !ok {
			m[v] = new(models.Hits)
		}
		m[v].Add(visit.FirstDart)

		v = visit.SecondDart.ValueRaw()
		if _, ok := m[v]; !ok {
			m[v] = new(models.Hits)
		}
		m[v].Add(visit.SecondDart)

		v = visit.ThirdDart.ValueRaw()
		if _, ok := m[v]; !ok {
			m[v] = new(models.Hits)
		}
		m[v].Add(visit.ThirdDart)
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
		player.Hits = hitsMap[player.PlayerID]
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].Order < players[j].Order
	})
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
		if visit.FirstDart.IsCheckoutAttempt(currentScore, 1) {
			totalAttempts++
			checkoutAttempts[currentScore]++
		}
		currentScore -= visit.FirstDart.GetScore()

		if visit.SecondDart.IsCheckoutAttempt(currentScore, 2) {
			totalAttempts++
			checkoutAttempts[currentScore]++
		}
		currentScore -= visit.SecondDart.GetScore()

		if visit.ThirdDart.IsCheckoutAttempt(currentScore, 3) {
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

func getMatchLooser(match *models.Match, winnerID int) int {
	if match.Players[0] == winnerID {
		return match.Players[1]
	}
	return match.Players[0]
}
