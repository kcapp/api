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
func NewLeg(matchID int, startingScore int, players []int, matchType *int) (*models.Leg, error) {
	tx, err := models.DB.Begin()
	if err != nil {
		return nil, err
	}

	// Shift players to get correct order
	id, players := players[0], players[1:]
	players = append(players, id)
	if matchType != nil && *matchType == models.SHOOTOUT {
		// This is a tie break for SHOOTOUT, so reverse the order of players to make sure the original "closes to bull" counts
		for i, j := 0, len(players)-1; i < j; i, j = i+1, j-1 {
			players[i], players[j] = players[j], players[i]
		}
	}
	res, err := tx.Exec("INSERT INTO leg (starting_score, current_player_id, leg_type_id, match_id, created_at) VALUES (?, ?, ?, ?, NOW()) ",
		startingScore, players[0], matchType, matchID)
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
	if matchType == nil {
		matchType = &match.MatchType.ID
	}
	if *matchType == models.X01HANDICAP {
		scores, err := GetPlayersScore(int(match.CurrentLegID.Int64))
		if err != nil {
			return nil, err
		}
		for _, player := range scores {
			handicaps[player.PlayerID] = player.Handicap
		}
	} else if *matchType == models.TICTACTOE {
		params := match.Legs[0].Parameters
		params.GenerateTicTacToeNumbers(startingScore)
		_, err = tx.Exec("INSERT INTO leg_parameters (leg_id, outshot_type_id, number_1, number_2, number_3, number_4, number_5, number_6, number_7, number_8, number_9) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			legID, params.OutshotType.ID, params.Numbers[0], params.Numbers[1], params.Numbers[2], params.Numbers[3], params.Numbers[4], params.Numbers[5], params.Numbers[6], params.Numbers[7], params.Numbers[8])
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	} else if *matchType == models.KNOCKOUT {
		params := match.Legs[0].Parameters
		_, err = tx.Exec("INSERT INTO leg_parameters (leg_id, starting_lives) VALUES (?, ?)", legID, params.StartingLives)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for idx, playerID := range players {
		order := idx + 1
		_, err = tx.Exec("INSERT INTO player2leg (player_id, leg_id, `order`, match_id, handicap) VALUES (?, ?, ?, ?, ?)",
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
	match, err := GetMatch(leg.MatchID)
	if err != nil {
		return err
	}

	// Update leg with winner
	winnerID := null.IntFrom(int64(visit.PlayerID))
	matchType := match.MatchType.ID
	if leg.LegType != nil {
		matchType = leg.LegType.ID
	}
	if matchType == models.SHOOTOUT || matchType == models.DARTSATX || matchType == models.AROUNDTHEWORLD ||
		(matchType == models.SHANGHAI && !visit.IsShanghai()) || matchType == models.BERMUDATRIANGLE ||
		matchType == models.JDCPRACTICE || matchType == models.SCAM {
		// For certain game types we need to check the scores of each player to determine which player won the leg with the highest score
		scores, err := GetPlayersScore(visit.LegID)
		if err != nil {
			return err
		}
		highScore := 0
		isDraw := false
		for playerID, player := range scores {
			if player.CurrentScore == highScore {
				isDraw = true
			}
			if player.CurrentScore > highScore {
				highScore = player.CurrentScore
				winnerID = null.IntFrom(int64(playerID))
				isDraw = false
			}
		}
		if isDraw {
			winnerID = null.IntFromPtr(nil)
		}
	} else if matchType == models.FOURTWENTY {
		scores, err := GetPlayersScore(visit.LegID)
		if err != nil {
			return err
		}
		lowestScore := 421
		for playerID, player := range scores {
			if player.CurrentScore < lowestScore {
				lowestScore = player.CurrentScore
				winnerID = null.IntFrom(int64(playerID))
			}
		}
	} else if matchType == models.TICTACTOE && !leg.Parameters.IsTicTacToeWinner(visit.PlayerID) {
		// If current player did not win, this game is a draw
		winnerID = null.IntFromPtr(nil)
	} else if matchType == models.KNOCKOUT {
		scores, err := GetPlayersScore(visit.LegID)
		if err != nil {
			return err
		}
		for _, player := range scores {
			if player.Lives.Int64 > 0 {
				winnerID = null.IntFrom(int64(player.PlayerID))
			}
		}
	}

	_, err = tx.Exec(`UPDATE leg SET current_player_id = ?, winner_id = ?, is_finished = 1, end_time = NOW() WHERE id = ?`, visit.PlayerID, winnerID, visit.LegID)
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Printf("[%d] Finished with player %d winning", visit.LegID, winnerID.ValueOrZero())

	if matchType == models.SHOOTOUT {
		statisticsMap, err := CalculateShootoutStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_shootout(leg_id, player_id, score, ppd, 60s_plus, 100s_plus, 140s_plus, 180s)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, visit.LegID, playerID, stats.Score, stats.PPD, stats.Score60sPlus,
				stats.Score100sPlus, stats.Score140sPlus, stats.Score180s)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting shootout statistics for player %d", visit.LegID, playerID)
		}
	} else if matchType == models.CRICKET {
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
	} else if matchType == models.DARTSATX {
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
	} else if matchType == models.AROUNDTHECLOCK {
		statisticsMap, err := CalculateAroundTheClockStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
			INSERT INTO statistics_around_the
				(leg_id, player_id, darts_thrown, score, longest_streak, total_hit_rate, hit_rate_1, hit_rate_2, hit_rate_3, hit_rate_4, hit_rate_5, hit_rate_6, hit_rate_7, hit_rate_8,
					hit_rate_9, hit_rate_10, hit_rate_11, hit_rate_12, hit_rate_13, hit_rate_14, hit_rate_15, hit_rate_16, hit_rate_17, hit_rate_18, hit_rate_19, hit_rate_20, hit_rate_bull)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, visit.LegID, playerID, stats.DartsThrown, stats.Score, stats.LongestStreak, stats.TotalHitRate, stats.Hitrates[1],
				stats.Hitrates[2], stats.Hitrates[3], stats.Hitrates[4], stats.Hitrates[5], stats.Hitrates[6], stats.Hitrates[7], stats.Hitrates[8], stats.Hitrates[9], stats.Hitrates[10],
				stats.Hitrates[11], stats.Hitrates[12], stats.Hitrates[13], stats.Hitrates[14], stats.Hitrates[15], stats.Hitrates[16], stats.Hitrates[17], stats.Hitrates[18], stats.Hitrates[19],
				stats.Hitrates[20], stats.Hitrates[25])
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting Around the Clock statistics for player %d", visit.LegID, playerID)
		}
	} else if matchType == models.AROUNDTHEWORLD || matchType == models.SHANGHAI {
		statisticsMap, err := CalculateAroundTheWorldStatistics(visit.LegID, matchType)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_around_the
					(leg_id, player_id, darts_thrown, score, shanghai, mpr, total_hit_rate, hit_rate_1, hit_rate_2, hit_rate_3, hit_rate_4, hit_rate_5, hit_rate_6, hit_rate_7, hit_rate_8, hit_rate_9, hit_rate_10,
						hit_rate_11, hit_rate_12, hit_rate_13, hit_rate_14, hit_rate_15, hit_rate_16, hit_rate_17, hit_rate_18, hit_rate_19, hit_rate_20, hit_rate_bull)
				VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, visit.LegID, playerID, stats.DartsThrown, stats.Score, stats.Shanghai, stats.MPR, stats.TotalHitRate, stats.Hitrates[1],
				stats.Hitrates[2], stats.Hitrates[3], stats.Hitrates[4], stats.Hitrates[5], stats.Hitrates[6], stats.Hitrates[7], stats.Hitrates[8], stats.Hitrates[9], stats.Hitrates[10],
				stats.Hitrates[11], stats.Hitrates[12], stats.Hitrates[13], stats.Hitrates[14], stats.Hitrates[15], stats.Hitrates[16], stats.Hitrates[17], stats.Hitrates[18], stats.Hitrates[19],
				stats.Hitrates[20], stats.Hitrates[25])
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting Around the World/Shanghai statistics for player %d", visit.LegID, playerID)
		}
	} else if matchType == models.TICTACTOE {
		statisticsMap, err := CalculateTicTacToeStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_tic_tac_toe (leg_id, player_id, darts_thrown, score, numbers_closed, highest_closed) VALUES (?,?,?,?,?,?)`, visit.LegID,
				playerID, stats.DartsThrown, stats.Score, stats.NumbersClosed, stats.HighestClosed)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting Tic Tac Toe statistics for player %d", visit.LegID, playerID)
		}
	} else if matchType == models.BERMUDATRIANGLE {
		statisticsMap, err := CalculateBermudaTriangleStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_bermuda_triangle (leg_id, player_id, darts_thrown, score, mpr, total_marks, highest_score_reached, total_hit_rate, hit_rate_1, hit_rate_2, hit_rate_3,
					hit_rate_4, hit_rate_5, hit_rate_6, hit_rate_7, hit_rate_8, hit_rate_9, hit_rate_10, hit_rate_11, hit_rate_12, hit_rate_13, hit_count) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
				visit.LegID, playerID, stats.DartsThrown, stats.Score, stats.MPR, &stats.TotalMarks, stats.HighestScoreReached, stats.TotalHitRate, stats.Hitrates[0], stats.Hitrates[1], stats.Hitrates[2],
				stats.Hitrates[3], stats.Hitrates[4], stats.Hitrates[5], stats.Hitrates[6], stats.Hitrates[7], stats.Hitrates[8], stats.Hitrates[9], stats.Hitrates[10], stats.Hitrates[11], stats.Hitrates[12],
				stats.HitCount)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting Bermuda Triangle statistics for player %d", visit.LegID, playerID)
		}
	} else if matchType == models.FOURTWENTY {
		statisticsMap, err := Calculate420Statistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_420 (leg_id, player_id, score, total_hit_rate, hit_rate_1, hit_rate_2, hit_rate_3, hit_rate_4, hit_rate_5, hit_rate_6, hit_rate_7, hit_rate_8, hit_rate_9,
					hit_rate_10, hit_rate_11, hit_rate_12, hit_rate_13, hit_rate_14, hit_rate_15, hit_rate_16, hit_rate_17, hit_rate_18, hit_rate_19, hit_rate_20, hit_rate_bull) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
				visit.LegID, playerID, stats.Score, stats.TotalHitRate, stats.Hitrates[1], stats.Hitrates[2], stats.Hitrates[3], stats.Hitrates[4], stats.Hitrates[5], stats.Hitrates[6],
				stats.Hitrates[7], stats.Hitrates[8], stats.Hitrates[9], stats.Hitrates[10], stats.Hitrates[11], stats.Hitrates[12], stats.Hitrates[13], stats.Hitrates[14], stats.Hitrates[15], stats.Hitrates[16],
				stats.Hitrates[17], stats.Hitrates[18], stats.Hitrates[19], stats.Hitrates[20], stats.Hitrates[25])
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting Four Twenty statistics for player %d", visit.LegID, playerID)
		}
	} else if matchType == models.KILLBULL {
		statisticsMap, err := CalculateKillBullStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
						INSERT INTO statistics_kill_bull (leg_id, player_id, darts_thrown, score, marks3, marks4, marks5, marks6, longest_streak, times_busted, total_hit_rate) VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
				visit.LegID, playerID, stats.DartsThrown, stats.Score, stats.Marks3, stats.Marks4, stats.Marks5, stats.Marks6, stats.LongestStreak, stats.TimesBusted, stats.TotalHitRate)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting Kill Bull statistics for player %d", visit.LegID, playerID)
		}
	} else if matchType == models.GOTCHA {
		statisticsMap, err := CalculateGotchaStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_gotcha (leg_id, player_id, darts_thrown, highest_score, times_reset, others_reset, score) VALUES (?,?,?,?,?,?,?)`,
				visit.LegID, playerID, stats.DartsThrown, stats.HighestScore, stats.TimesReset, stats.OthersReset, stats.Score)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting Gotcha statistics for player %d", visit.LegID, playerID)
		}
	} else if matchType == models.JDCPRACTICE {
		statisticsMap, err := CalculateJDCPracticeStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_jdc_practice (leg_id, player_id, darts_thrown, score, mpr, shanghai_count, doubles_hitrate) VALUES (?,?,?,?,?,?,?)`,
				visit.LegID, playerID, stats.DartsThrown, stats.Score, stats.MPR, stats.ShanghaiCount, stats.DoublesHitrate)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting JDC Practice statistics for player %d", visit.LegID, playerID)
		}
	} else if matchType == models.KNOCKOUT {
		statisticsMap, err := CalculateKnockoutStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_knockout (leg_id, player_id, darts_thrown, avg_score, lives_lost, lives_taken, final_position) VALUES (?,?,?,?,?,?,?)`,
				visit.LegID, playerID, stats.DartsThrown, stats.AvgScore, stats.LivesLost, stats.LivesTaken, stats.FinalPosition)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting Knockout statistics for player %d", visit.LegID, playerID)
		}
	} else if matchType == models.SCAM {
		statisticsMap, err := CalculateScamStatistics(visit.LegID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for playerID, stats := range statisticsMap {
			_, err = tx.Exec(`
				INSERT INTO statistics_scam (leg_id, player_id, darts_thrown_stopper, darts_thrown_scorer, mpr, ppd, score) VALUES (?,?,?,?,?,?,?)`,
				visit.LegID, playerID, stats.DartsThrownStopper, stats.DartsThrownScorer, stats.MPR, stats.PPD, stats.Score)
			if err != nil {
				tx.Rollback()
				return err
			}
			log.Printf("[%d] Inserting Scam statistics for player %d", visit.LegID, playerID)
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
		if playerID == int(winnerID.ValueOrZero()) {
			currentPlayerWins += wins
		}
	}

	isFinished := false
	isTieBreak := false
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
				if playerID == int(winnerID.ValueOrZero()) {
					// Don't add payback to ourself
					continue
				}
				_, err = tx.Exec(`
					INSERT INTO owes (player_ower_id, player_owee_id, owe_type_id, amount) VALUES (?, ?, ?, 1)
					ON DUPLICATE KEY UPDATE amount = amount + 1`, playerID, visit.PlayerID, match.OweTypeID)
				if err != nil {
					tx.Rollback()
					return err
				}
				log.Printf("Added owes of %s from player %d to player %d", match.OweType.Item.String, playerID, visit.PlayerID)
			}
		}
		log.Printf("Match %d finished with player %d winning", match.ID, winnerID.ValueOrZero())
	} else if match.MatchMode.LegsRequired.Valid && playedLegs == int(match.MatchMode.LegsRequired.Int64) {
		// Match finished, draw
		isFinished = true
		_, err = tx.Exec("UPDATE matches SET is_finished = 1 WHERE id = ?", match.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("Match %d finished with a Draw", match.ID)
	} else if playedLegs == (int(match.MatchMode.LegsRequired.Int64)-1) && match.MatchMode.TieBreakMatchTypeID.Valid {
		isTieBreak = true
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
				err = SwapPlayers(winnerMatch.ID, int(winnerID.ValueOrZero()), winnerMatch.Players[idx])
				if err != nil {
					return err
				}
			}
			if metadata.LooserOutcomeMatchID.Valid {
				looserID := getMatchLooser(match, int(winnerID.ValueOrZero()))
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
	} else {
		log.Printf("Match %d is not finished, creating next leg", match.ID)
		var matchType *int
		if isTieBreak {
			matchType = new(int)
			*matchType = int(match.MatchMode.TieBreakMatchTypeID.Int64)
		}
		_, err = NewLeg(match.ID, leg.StartingScore, leg.Players, matchType)
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
	_, err = tx.Exec("DELETE FROM statistics_around_the WHERE leg_id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM statistics_tic_tac_toe WHERE leg_id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM statistics_bermuda_triangle WHERE leg_id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM statistics_420 WHERE leg_id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM statistics_kill_bull WHERE leg_id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM statistics_gotcha WHERE leg_id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM statistics_jdc_practice WHERE leg_id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM statistics_knockout WHERE leg_id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM statistics_scam WHERE leg_id = ?", legID)
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
			l.match_id, l.has_scores, GROUP_CONCAT(p2l.player_id ORDER BY p2l.order ASC) as "players",
			mt.id as 'match_type_id', mt.name, mt.description
		FROM leg l
			LEFT JOIN player2leg p2l ON p2l.leg_id = l.id
			LEFT JOIN matches m ON m.id = l.match_id
			LEFT JOIN match_type mt on mt.id = IFNULL(l.leg_type_id, m.match_type_id)
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
		leg.LegType = new(models.MatchType)
		var players string
		err := rows.Scan(&leg.ID, &leg.Endtime, &leg.StartingScore, &leg.IsFinished, &leg.CurrentPlayerID,
			&leg.WinnerPlayerID, &leg.CreatedAt, &leg.UpdatedAt, &leg.MatchID, &leg.HasScores, &players, &leg.LegType.ID,
			&leg.LegType.Name, &leg.LegType.Description)
		if err != nil {
			return nil, err
		}
		leg.Players = util.StringToIntArray(players)
		visits, err := GetLegVisits(leg.ID)
		if err != nil {
			return nil, err
		}
		leg.Visits = visits

		matchType := leg.LegType.ID
		if matchType == models.TICTACTOE || matchType == models.KNOCKOUT {
			leg.Parameters, err = GetLegParameters(leg.ID)
			if err != nil {
				return nil, err
			}
		}
		legs = append(legs, leg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return legs, nil
}

// GetLegsOfType returns all legs with scores for the given match type
func GetLegsOfType(matchType int, loadVisits bool) ([]*models.Leg, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id, l.end_time, l.starting_score, l.is_finished,
			l.current_player_id, l.winner_id, l.created_at, l.updated_at,
			l.match_id, l.has_scores, GROUP_CONCAT(p2l.player_id ORDER BY p2l.order ASC)
		FROM leg l
			JOIN matches m on m.id = l.match_id
			JOIN player2leg p2l ON p2l.leg_id = l.id
		WHERE l.has_scores = 1 AND (m.match_type_id = ? OR l.leg_type_id = ?)
		GROUP BY l.id
		ORDER BY l.id DESC`, matchType, matchType)
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
		if loadVisits {
			visits, err := GetLegVisits(leg.ID)
			if err != nil {
				return nil, err
			}
			leg.Visits = visits
		}
		if matchType == models.TICTACTOE || matchType == models.KNOCKOUT {
			leg.Parameters, err = GetLegParameters(leg.ID)
			if err != nil {
				return nil, err
			}
		}
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
			l.match_id, l.has_scores, GROUP_CONCAT(p2l.player_id ORDER BY p2l.order ASC),
			mt.id as 'match_type_id', mt.name, mt.description
		FROM leg l
			LEFT JOIN player2leg p2l ON p2l.leg_id = l.id
			LEFT JOIN matches m ON m.id = l.match_id
			LEFT JOIN match_type mt on mt.id = IFNULL(l.leg_type_id, m.match_type_id)
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
		leg.LegType = new(models.MatchType)
		var players string
		err := rows.Scan(&leg.ID, &leg.Endtime, &leg.StartingScore, &leg.IsFinished, &leg.CurrentPlayerID, &leg.WinnerPlayerID,
			&leg.CreatedAt, &leg.UpdatedAt, &leg.MatchID, &leg.HasScores, &players, &leg.LegType.ID, &leg.LegType.Name,
			&leg.LegType.Description)
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
	leg.LegType = new(models.MatchType)
	var players string
	err := models.DB.QueryRow(`
		SELECT
			l.id, l.end_time, l.starting_score, l.is_finished, l.current_player_id, l.winner_id, l.created_at, l.updated_at,
			l.board_stream_url, l.match_id, l.has_scores, GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order ASC) AS 'players',
			mt.id as 'match_type_id', mt.name, mt.description
		FROM leg l
			LEFT JOIN player2leg p2l ON p2l.leg_id = l.id
			LEFT JOIN matches m ON m.id = l.match_id
			LEFT JOIN match_type mt on mt.id = IFNULL(l.leg_type_id, m.match_type_id)
		WHERE l.id = ?`, id).Scan(&leg.ID, &leg.Endtime, &leg.StartingScore, &leg.IsFinished, &leg.CurrentPlayerID, &leg.WinnerPlayerID,
		&leg.CreatedAt, &leg.UpdatedAt, &leg.BoardStreamURL, &leg.MatchID, &leg.HasScores, &players, &leg.LegType.ID,
		&leg.LegType.Name, &leg.LegType.Description)
	if err != nil {
		return nil, err
	}

	leg.Players = util.StringToIntArray(players)
	visits, err := GetLegVisits(id)
	if err != nil {
		return nil, err
	}

	matchType := leg.LegType.ID
	if matchType == models.TICTACTOE || matchType == models.KNOCKOUT {
		leg.Parameters, err = GetLegParameters(id)
		if err != nil {
			return nil, err
		}
	}

	scores := make(map[int]*models.Player2Leg)
	for i := 0; i < len(leg.Players); i++ {
		p2l := new(models.Player2Leg)
		p2l.Hits = make(models.HitsMap)
		if matchType == models.DARTSATX || matchType == models.AROUNDTHECLOCK || matchType == models.AROUNDTHEWORLD || matchType == models.SHANGHAI ||
			matchType == models.TICTACTOE || matchType == models.BERMUDATRIANGLE || matchType == models.GOTCHA || matchType == models.JDCPRACTICE ||
			matchType == models.SHOOTOUT || matchType == models.SCAM {
			p2l.CurrentScore = 0
		} else if matchType == models.KNOCKOUT {
			p2l.CurrentScore = 0
			p2l.Lives = null.IntFrom(leg.Parameters.StartingLives.Int64)
		} else if matchType == models.FOURTWENTY {
			p2l.CurrentScore = 420
		} else if matchType == models.X01HANDICAP {
			// TODO
		} else {
			p2l.CurrentScore = leg.StartingScore
			p2l.StartingScore = leg.StartingScore
		}
		p2l.Order = i + 1
		p2l.DartsThrown = 0
		scores[leg.Players[i]] = p2l
	}

	specialNums := make([]int, 0)
	if leg.Parameters != nil && leg.Parameters.Numbers != nil {
		specialNums = make([]int, len(leg.Parameters.Numbers))
		copy(specialNums, leg.Parameters.Numbers)
	}

	playerOrder := 1
	if matchType == models.SCAM {
		for _, player := range scores {
			if player.Order == playerOrder {
				player.SetStopper()
			} else {
				player.SetScorer()
			}
		}
	}

	dartsThrown := 0
	visitCount := 0
	round := 1
	for i, visit := range visits {
		if i > 0 && i%len(leg.Players) == 0 {
			round++
		}
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
				scores[visit.PlayerID].CurrentScore += score
			} else if matchType == models.AROUNDTHECLOCK {
				score = visit.CalculateAroundTheClockScore(scores[visit.PlayerID].CurrentScore)
				scores[visit.PlayerID].CurrentScore += score
			} else if matchType == models.AROUNDTHEWORLD || matchType == models.SHANGHAI {
				score = visit.CalculateAroundTheWorldScore(round)
				scores[visit.PlayerID].CurrentScore += score
			} else if matchType == models.TICTACTOE {
				score = 0

				lastDartValid := visit.GetLastDart().IsDouble()
				if leg.Parameters.OutshotType.ID == models.OUTSHOTANY {
					lastDartValid = true
				} else if leg.Parameters.OutshotType.ID == models.OUTSHOTMASTER {
					lastDartValid = visit.GetLastDart().IsDouble() || visit.GetLastDart().IsTriple()
				}

				for _, num := range leg.Parameters.Numbers {
					if num == visit.GetScore() && lastDartValid {
						score = num
						break
					}
				}
				scores[visit.PlayerID].CurrentScore += score
			} else if matchType == models.BERMUDATRIANGLE {
				score = visit.CalculateBermudaTriangleScore(round - 1)
				if score == 0 {
					scores[visit.PlayerID].CurrentScore = scores[visit.PlayerID].CurrentScore / 2
				} else {
					scores[visit.PlayerID].CurrentScore += score
				}
			} else if matchType == models.FOURTWENTY {
				score = visit.Calculate420Score(round - 1)
				scores[visit.PlayerID].CurrentScore -= score
			} else if matchType == models.KILLBULL {
				score = visit.CalculateKillBullScore()
				if score == 0 {
					scores[visit.PlayerID].CurrentScore = scores[visit.PlayerID].StartingScore
				} else {
					scores[visit.PlayerID].CurrentScore -= score
				}
			} else if matchType == models.GOTCHA {
				score = visit.CalculateGotchaScore(scores, leg.StartingScore)
				scores[visit.PlayerID].CurrentScore += score
			} else if matchType == models.JDCPRACTICE {
				score = visit.CalculateJDCPracticeScore(round - 1)
				scores[visit.PlayerID].CurrentScore += score
			} else if matchType == models.KNOCKOUT {
				player := scores[visit.PlayerID]
				player.CurrentScore = visit.GetScore()
				// Set correctly darts thrown for each player
				player.DartsThrown += 3
				visit.DartsThrown = player.DartsThrown
			} else if matchType == models.SCAM {
				player := scores[visit.PlayerID]
				if player.IsStopper.Bool {
					score = 0
					visit.Marks = visit.CalculateScamMarks(scores)
					visit.IsStopper = null.BoolFrom(true)
					if player.Hits.Contains(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20) {
						playerOrder++
						for _, player := range scores {
							if player.Order == playerOrder {
								player.SetStopper()
							} else {
								player.SetScorer()
							}
						}
					}
				} else {
					score = visit.CalculateScamScore(scores)
					scores[visit.PlayerID].CurrentScore += score
				}
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

		// We also want to add hits for certain special numbers in some game types
		for j, num := range specialNums {
			// Check if we hit the exact number, ending with a double
			lastDartValid := visit.GetLastDart().IsDouble()
			if leg.Parameters.OutshotType.ID == models.OUTSHOTANY {
				lastDartValid = true
			} else if leg.Parameters.OutshotType.ID == models.OUTSHOTMASTER {
				lastDartValid = visit.GetLastDart().IsDouble() || visit.GetLastDart().IsTriple()
			}
			if num == visit.GetScore() && lastDartValid {
				leg.Parameters.Hits[num] = visit.PlayerID

				// Remove the number to only let first player hit a specific number
				specialNums[j] = specialNums[len(specialNums)-1]
				specialNums = specialNums[:len(specialNums)-1]
				break
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
	if matchType == models.X01 || matchType == models.X01HANDICAP {
		leg.CheckoutStatistics, err = getCheckoutStatistics(leg.ID, leg.StartingScore)
	}
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
	hitsMap := make(map[int]models.HitsMap)
	lowestScore := leg.StartingScore
	players := make([]*models.Player2Leg, 0)
	for _, player := range scores {
		player.Modifiers = new(models.PlayerModifiers)
		if player.CurrentScore < lowestScore {
			lowestScore = player.CurrentScore
		}
		hitsMap[player.PlayerID] = make(models.HitsMap)
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
		hits := hitsMap[visit.PlayerID]
		hits.Add(visit.FirstDart)
		hits.Add(visit.SecondDart)
		hits.Add(visit.ThirdDart)
	}

	for _, player := range players {
		player.Wins = winsMap[player.PlayerID]
		if visit, ok := lastVisits[player.PlayerID]; ok {
			player.Modifiers.IsViliusVisit = visit.IsViliusVisit()
			player.Modifiers.IsFishAndChips = visit.IsFishAndChips()
			player.Modifiers.IsScore60Plus = visit.IsScore60Plus()
			player.Modifiers.IsScore100Plus = visit.IsScore100Plus()
			player.Modifiers.IsScore140Plus = visit.IsScore140Plus()
			player.Modifiers.IsScore180 = visit.IsScore180()
		}
		if lowestScore < 171 && player.CurrentScore > 199 {
			player.Modifiers.IsBeerMatch = true
		}
		player.AddVisitStatistics(*leg)
		if player.Hits == nil {
			// Certain match types calculate hits differently, so only override it here if it's not set
			player.Hits = hitsMap[player.PlayerID]
		}
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

// StartWarmup will set the updated_at of a leg to now
func StartWarmup(legID int, venueID int) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE leg SET updated_at = NOW() WHERE id = ?", legID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if venueID > 0 {
		// Update the venue ID if this leg is being played at a different location than it was scheduled
		_, err = tx.Exec("UPDATE matches SET venue_id = ? WHERE id = (SELECT match_id FROM leg WHERE id = ?)", venueID, legID)
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("[%d] Moved leg to new venue %d", legID, venueID)
	}
	tx.Commit()

	log.Printf("[%d] Started warmup", legID)
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
			_, err = tx.Exec("UPDATE matches SET current_leg_id = ?, is_abandoned = 1, is_finished = 1 WHERE id = ?", previousLeg, match.ID)
			if err != nil {
				return err
			}
			log.Printf("[%d] Updated current leg of match %d", previousLeg, match.ID)
		}
		return nil
	})
}

// GetLegParameters will return leg parameters for the given leg
func GetLegParameters(legID int) (*models.LegParameters, error) {
	params := new(models.LegParameters)
	n := make([]null.Int, 9)
	var ost null.Int
	err := models.DB.QueryRow(`
		SELECT outshot_type_id, number_1, number_2, number_3, number_4, number_5, number_6, number_7, number_8, number_9, starting_lives
		FROM leg_parameters WHERE leg_id = ?`, legID).Scan(&ost, &n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7], &n[8], &params.StartingLives)
	if err != nil {
		return nil, err
	}
	if ost.Valid {
		os, err := GetOutshotType(int(ost.Int64))
		if err != nil {
			return nil, err
		}
		params.OutshotType = os
	}
	if n[0].Valid {
		numbers := make([]int, 9)
		for i, num := range n {
			numbers[i] = int(num.Int64)
		}
		params.Numbers = numbers
	}
	params.Hits = make(map[int]int)
	return params, nil
}

// GetLegMatchType returns the match type for a given leg
func GetLegMatchType(legID int) (*int, error) {
	var matchType int
	err := models.DB.QueryRow(`
        SELECT
			IFNULL(l.leg_type_id, m.match_type_id) as 'match_type_id'
		FROM matches m
			LEFT JOIN leg l ON l.match_id = m.id
		WHERE l.id = ?`, legID).Scan(&matchType)
	if err != nil {
		return nil, err
	}
	return &matchType, nil
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
