package data

import (
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/kcapp/api/data/queries"
	"github.com/kcapp/api/models"
	"github.com/kcapp/api/util"
	"log"
	"math"
	"sort"
)

// GetPlayers returns a map of all players
func GetPlayers() (map[int]*models.Player, error) {
	played, err := GetMatchesPlayedPerPlayer()
	if err != nil {
		return nil, err
	}

	rows, err := models.DB.Query(queries.QueryAllPlayers())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.VocalName, &p.Nickname, &p.SlackHandle, &p.Color, &p.ProfilePicURL,
			&p.SmartcardUID, &p.BoardStreamURL, &p.BoardStreamCSS, &p.IsActive, &p.OfficeID, &p.IsBot, &p.IsPlaceholder, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}

		players[p.ID] = p

		if val, ok := played[p.ID]; ok {
			p.MatchesPlayed = val.MatchesPlayed
			p.MatchesWon = val.MatchesWon
			p.LegsPlayed = val.LegsPlayed
			p.LegsWon = val.LegsWon
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

// GetActivePlayers returns an array of all active players
func GetActivePlayers() (map[int]*models.Player, error) {
	played, err := GetMatchesPlayedPerPlayer()
	if err != nil {
		return nil, err
	}

	rows, err := models.DB.Query(queries.QueryActivePlayers())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.VocalName, &p.Nickname, &p.SlackHandle, &p.Color, &p.ProfilePicURL,
			&p.SmartcardUID, &p.BoardStreamURL, &p.BoardStreamCSS, &p.IsActive, &p.OfficeID, &p.IsBot, &p.IsPlaceholder, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}

		players[p.ID] = p

		if val, ok := played[p.ID]; ok {
			p.MatchesPlayed = val.MatchesPlayed
			p.MatchesWon = val.MatchesWon
			p.LegsPlayed = val.LegsPlayed
			p.LegsWon = val.LegsWon
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

// GetMatchesPlayedPerPlayer will get the number of matches and legs played and won for each player
func GetMatchesPlayedPerPlayer() (map[int]*models.Player, error) {
	rows, err := models.DB.Query(queries.QueryMatchesPlayed())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	played := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.MatchesPlayed, &p.MatchesWon, &p.LegsPlayed, &p.LegsWon)
		if err != nil {
			return nil, err
		}
		played[p.ID] = p
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return played, nil
}

// GetPlayerScore will get the score for the given player in the given leg
func GetPlayerScore(playerID int, legID int) (int, error) {
	scores, err := GetPlayersScore(legID)
	if err != nil {
		return 0, err
	}
	return scores[playerID].CurrentScore, nil
}

// GetPlayersScore will get the score for all players in the given leg
func GetPlayersScore(legID int) (map[int]*models.Player2Leg, error) {
	players, err := GetPlayersInLeg(legID)
	if err != nil {
		return nil, err
	}
	rows, err := models.DB.Query(`
			SELECT
				p2l.leg_id,
				p2l.player_id,
				p.first_name,
				p2l.order,
				p2l.handicap,
				p2l.player_id = l.current_player_id AS 'is_current_player',
				IF(m.match_type_id = 5, 0, (l.starting_score +
					-- Add handicap for players if match_mode is handicap
					IF(m.match_type_id = 3, IFNULL(p2l.handicap, 0), 0)) -
					(IFNULL(SUM(first_dart * first_dart_multiplier), 0) +
					IFNULL(SUM(second_dart * second_dart_multiplier), 0) +
					IFNULL(SUM(third_dart * third_dart_multiplier), 0))
					-- For X01 score goes down, while Shootout it counts up
					* IF(m.match_type_id = 2, -1, 1)) AS 'current_score',
				l.starting_score,
				b.player_id,
				b.skill_level
			FROM player2leg p2l
				LEFT JOIN player p on p.id = p2l.player_id
				LEFT JOIN leg l ON l.id = p2l.leg_id
				LEFT JOIN score s ON s.leg_id = p2l.leg_id AND s.player_id = p2l.player_id
				LEFT JOIN matches m on m.id = l.match_id
				LEFT JOIN bot2player2leg b ON b.player2leg_id = p2l.id
			WHERE p2l.leg_id = ? AND (s.is_bust IS NULL OR is_bust = 0)
			GROUP BY p2l.player_id
			ORDER BY p2l.order ASC`, legID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scores := make(map[int]*models.Player2Leg)
	for rows.Next() {
		p2l := new(models.Player2Leg)
		bc := new(models.BotConfig)
		err := rows.Scan(&p2l.LegID, &p2l.PlayerID, &p2l.PlayerName, &p2l.Order, &p2l.Handicap, &p2l.IsCurrentPlayer,
			&p2l.CurrentScore, &p2l.StartingScore, &bc.PlayerID, &bc.Skill)
		if err != nil {
			return nil, err
		}
		if bc.Skill.Valid {
			p2l.BotConfig = bc
		}
		p2l.Player = players[p2l.PlayerID]
		scores[p2l.PlayerID] = p2l
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	mt, err := GetLegMatchType(legID)
	if err != nil {
		return nil, err
	}
	matchType := *mt
	// Get score for other game types
	if matchType == models.SHOOTOUT {
		visits, err := GetLegVisits(legID)
		if err != nil {
			return nil, err
		}
		for _, player := range scores {
			player.CurrentScore = 0
			player.DartsThrown = 0
		}
		for _, visit := range visits {
			player := scores[visit.PlayerID]
			player.CurrentScore += visit.GetScore()
			player.DartsThrown += 3
		}
	} else if matchType == models.CRICKET {
		visits, err := GetLegVisits(legID)
		if err != nil {
			return nil, err
		}
		cricketScores := make(map[int]*models.Player2Leg)
		for id := range scores {
			p2l := new(models.Player2Leg)
			p2l.Hits = make(map[int]*models.Hits)
			cricketScores[id] = p2l
		}

		for _, visit := range visits {
			visit.CalculateCricketScore(cricketScores)
		}
		for _, player := range scores {
			player.CurrentScore = cricketScores[player.PlayerID].CurrentScore
		}

		return scores, nil
	} else if matchType == models.DARTSATX {
		rows, err := models.DB.Query(`
			SELECT
				player_id,
				SUM(case when first_dart = l.starting_score then first_dart_multiplier else 0 end) +
				SUM(case when second_dart = l.starting_score then second_dart_multiplier else 0 end) +
				SUM(case when third_dart = l.starting_score then third_dart_multiplier else 0 end) as 'current_score'
			FROM score s
			JOIN leg l on l.id = s.leg_id
			WHERE leg_id = ?
			GROUP BY player_id`, legID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var playerID int
			var score int
			err := rows.Scan(&playerID, &score)
			if err != nil {
				return nil, err
			}
			scores[playerID].CurrentScore = score
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	} else if matchType == models.AROUNDTHECLOCK {
		visits, err := GetLegVisits(legID)
		if err != nil {
			return nil, err
		}
		for _, player := range scores {
			player.CurrentScore = 0
		}

		for _, visit := range visits {
			score := visit.CalculateAroundTheClockScore(scores[visit.PlayerID].CurrentScore)
			scores[visit.PlayerID].CurrentScore += score
		}
	} else if matchType == models.AROUNDTHEWORLD || matchType == models.SHANGHAI {
		visits, err := GetLegVisits(legID)
		if err != nil {
			return nil, err
		}
		for _, player := range scores {
			player.CurrentScore = 0
		}

		round := 1
		for i, visit := range visits {
			if i > 0 && i%len(players) == 0 {
				round++
			}
			score := visit.CalculateAroundTheWorldScore(round)
			scores[visit.PlayerID].CurrentScore += score
		}
	} else if matchType == models.TICTACTOE {
		for _, player := range scores {
			player.CurrentScore = 0
		}
	} else if matchType == models.BERMUDATRIANGLE {
		visits, err := GetLegVisits(legID)
		if err != nil {
			return nil, err
		}
		for _, player := range scores {
			player.CurrentScore = 0
		}

		round := 1
		for i, visit := range visits {
			if i > 0 && i%len(players) == 0 {
				round++
			}
			score := visit.CalculateBermudaTriangleScore(round - 1)
			if score == 0 {
				scores[visit.PlayerID].CurrentScore = scores[visit.PlayerID].CurrentScore / 2
			} else {
				scores[visit.PlayerID].CurrentScore += score
			}
		}
	} else if matchType == models.FOURTWENTY {
		visits, err := GetLegVisits(legID)
		if err != nil {
			return nil, err
		}
		for _, player := range scores {
			player.CurrentScore = 420
		}

		round := 1
		for i, visit := range visits {
			if i > 0 && i%len(players) == 0 {
				round++
			}
			score := visit.Calculate420Score(round - 1)
			scores[visit.PlayerID].CurrentScore -= score
		}
	} else if matchType == models.KILLBULL {
		visits, err := GetLegVisits(legID)
		if err != nil {
			return nil, err
		}
		for _, player := range scores {
			player.CurrentScore = player.StartingScore
		}

		for _, visit := range visits {
			score := visit.CalculateKillBullScore()
			if score == 0 {
				scores[visit.PlayerID].CurrentScore = scores[visit.PlayerID].StartingScore
			} else {
				scores[visit.PlayerID].CurrentScore -= score
			}
		}
	} else if matchType == models.GOTCHA {
		visits, err := GetLegVisits(legID)
		if err != nil {
			return nil, err
		}

		targetScore := 0
		for _, player := range scores {
			player.CurrentScore = 0
			targetScore = player.StartingScore
		}

		for _, visit := range visits {
			score := visit.CalculateGotchaScore(scores, targetScore)
			scores[visit.PlayerID].CurrentScore += score
		}
	} else if matchType == models.JDCPRACTICE {
		visits, err := GetLegVisits(legID)
		if err != nil {
			return nil, err
		}
		for _, player := range scores {
			player.CurrentScore = 0
		}

		round := 1
		for i, visit := range visits {
			if i > 0 && i%len(players) == 0 {
				round++
			}
			scores[visit.PlayerID].CurrentScore += visit.CalculateJDCPracticeScore(round - 1)
		}
	} else if matchType == models.KNOCKOUT {
		visits, err := GetLegVisits(legID)
		if err != nil {
			return nil, err
		}
		params, err := GetLegParameters(legID)
		if err != nil {
			return nil, err
		}

		for _, player := range scores {
			player.CurrentScore = 0
			player.Lives = params.StartingLives
		}

		for i, visit := range visits {
			player := scores[visit.PlayerID]
			player.CurrentScore = visit.GetScore()

			idx := i - 1
			if idx < 0 {
				continue
			}
			prev := visits[idx]
			if prev.GetScore() > visit.GetScore() {
				player.Lives = null.IntFrom(player.Lives.Int64 - 1)
			}
			scores[prev.PlayerID].CurrentScore = 0
		}
	} else if matchType == models.SCAM {
		stopperOrder := 1
		for _, player := range scores {
			if player.Order == stopperOrder {
				player.SetStopper()
			} else {
				player.SetScorer()
			}
			player.CurrentScore = 0
			player.Hits = make(models.HitsMap)
		}

		visits, err := GetLegVisits(legID)
		if err != nil {
			return nil, err
		}

		hits := make(models.HitsMap)
		for _, visit := range visits {
			player := scores[visit.PlayerID]
			if player.IsStopper.Bool {
				hits.Add(visit.FirstDart)
				hits.Add(visit.SecondDart)
				hits.Add(visit.ThirdDart)
				player.Hits = hits

				visit.IsStopper = null.BoolFrom(true)
				if hits.Contains(models.SINGLE, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20) {
					stopperOrder++
					for _, player := range scores {
						if player.Order == stopperOrder {
							player.SetStopper()
						} else {
							player.SetScorer()
						}
					}
					hits = make(models.HitsMap)
				}
			} else if player.IsScorer.Bool {
				if hits.GetHits(visit.FirstDart.ValueRaw(), models.SINGLE) < 1 {
					player.CurrentScore += visit.FirstDart.GetScore()
				}
				if hits.GetHits(visit.SecondDart.ValueRaw(), models.SINGLE) < 1 {
					player.CurrentScore += visit.SecondDart.GetScore()
				}
				if hits.GetHits(visit.ThirdDart.ValueRaw(), models.SINGLE) < 1 {
					player.CurrentScore += visit.ThirdDart.GetScore()
				}
			}
		}
	}
	return scores, nil
}

// GetPlayer returns the player for the given ID
func GetPlayer(id int) (*models.Player, error) {
	p := new(models.Player)
	err := models.DB.QueryRow(queries.QueryPlayer(), id).
		Scan(&p.ID, &p.FirstName, &p.LastName, &p.VocalName, &p.Nickname, &p.SlackHandle,
			&p.Color, &p.ProfilePicURL, &p.SmartcardUID, &p.BoardStreamURL, &p.BoardStreamCSS, &p.OfficeID, &p.IsActive,
			&p.IsBot, &p.IsPlaceholder, &p.CreatedAt, &p.UpdatedAt, &p.CurrentElo, &p.TournamentElo)
	if err != nil {
		return nil, err
	}

	pld, err := GetMatchesPlayedPerPlayer()
	if err != nil {
		return nil, err
	}
	played := pld[p.ID]
	if played != nil {
		p.MatchesPlayed = played.MatchesPlayed
		p.MatchesWon = played.MatchesWon
		p.LegsPlayed = played.LegsPlayed
		p.LegsWon = played.LegsWon
	}

	return p, nil
}

// GetPlayerEloChangelog returns the elo changelog for the given player
func GetPlayerEloChangelog(id int, start int, limit int) (*models.PlayerEloChangelogs, error) {
	var total int
	err := models.DB.QueryRow(`SELECT COUNT(id) FROM player_elo_changelog WHERE player_id = ?`, id).Scan(&total)
	if err != nil {
		return nil, err
	}

	rows, err := models.DB.Query(`
		SELECT
			home.id, home.match_id, m.updated_at,
			if(m.tournament_id is null, false, true) as 'is_official',
			mm.short_name as 'match_mode',
			mt.name as 'match_type', m.winner_id,
			home.player_id, home.old_elo, home.new_elo, home.old_tournament_elo, home.new_tournament_elo,
			away.player_id, away.old_elo, away.new_elo, away.old_tournament_elo, away.new_tournament_elo
		FROM player_elo_changelog home
			JOIN player_elo_changelog away ON away.match_id = home.match_id AND away.player_id <> home.player_id
			JOIN matches m on m.id = home.match_id
			JOIN match_type mt on m.match_type_id = mt.id
			JOIN match_mode mm on m.match_mode_id = mm.id
		WHERE home.player_id = ?
		ORDER BY home.id DESC
		LIMIT ?, ?`, id, start, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	changelogs := new(models.PlayerEloChangelogs)
	changelogs.Total = total
	changelogs.Changelog = make([]*models.PlayerEloChangelog, 0)
	for rows.Next() {
		change := new(models.PlayerEloChangelog)

		home := new(models.PlayerElo)
		away := new(models.PlayerElo)

		err := rows.Scan(&change.ID, &change.MatchID, &change.FinishedAt, &change.IsOfficial, &change.MatchMode,
			&change.MatchType, &change.WinnerID,
			&home.PlayerID, &home.CurrentElo, &home.CurrentEloNew, &home.TournamentElo, &home.TournamentEloNew,
			&away.PlayerID, &away.CurrentElo, &away.CurrentEloNew, &away.TournamentElo, &away.TournamentEloNew)
		if err != nil {
			return nil, err
		}
		change.HomePlayer = home
		change.AwayPlayer = away
		changelogs.Changelog = append(changelogs.Changelog, change)
	}
	return changelogs, nil
}

// AddPlayer will add a new player to the database
func AddPlayer(player models.Player) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}

	// Prepare statement for inserting data
	res, err := tx.Exec(`INSERT INTO player (first_name, last_name, vocal_name, nickname, slack_handle, color,
			profile_pic_url, smartcard_uid, board_stream_url, board_stream_css, office_id, is_bot)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)`,
		player.FirstName, player.LastName, player.VocalName, player.Nickname, player.SlackHandle, player.Color, player.ProfilePicURL,
		player.SmartcardUID, player.BoardStreamURL, player.BoardStreamCSS, player.OfficeID)
	if err != nil {
		tx.Rollback()
		return err
	}
	playerID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO player_elo (player_id) VALUES (?)", playerID)
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Printf("Created new player (%d) %s", playerID, player.FirstName)
	tx.Commit()
	return nil
}

// UpdatePlayer will update the given player
func UpdatePlayer(playerID int, player models.Player) error {
	// Prepare statement for inserting data
	stmt, err := models.DB.Prepare(`
		UPDATE player SET
			first_name = ?, last_name = ?, vocal_name = ?, nickname = ?, slack_handle = ?,
			color = ?, profile_pic_url = ?, smartcard_uid = ?, board_stream_url = ?, board_stream_css = ?, office_id = ?,
			updated_at = NOW()
		WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(player.FirstName, player.LastName, player.VocalName, player.Nickname, player.SlackHandle, player.Color,
		player.ProfilePicURL, player.SmartcardUID, player.BoardStreamURL, player.BoardStreamCSS, player.OfficeID, playerID)
	if err != nil {
		return err
	}
	log.Printf("Updated player %s (%v)", player.FirstName, player)
	return nil
}

// GetPlayersInLeg will get all players in a given leg
func GetPlayersInLeg(legID int) (map[int]*models.Player, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id,
			p.first_name,
			p.last_name,
			p.vocal_name,
			p.nickname,
			p.slack_handle,
			p.color,
			p.profile_pic_url,
			p.smartcard_uid,
			p.board_stream_url,
			p.board_stream_css,
			p.office_id,
			p.active,
			p.is_bot,
			p.is_placeholder
		FROM player2leg p2l
		LEFT JOIN player p ON p.id = p2l.player_id WHERE p2l.leg_id = ?`, legID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.VocalName, &p.Nickname, &p.SlackHandle, &p.Color, &p.ProfilePicURL,
			&p.SmartcardUID, &p.BoardStreamURL, &p.BoardStreamCSS, &p.OfficeID, &p.IsActive, &p.IsBot, &p.IsPlaceholder)
		if err != nil {
			return nil, err
		}
		players[p.ID] = p
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return players, nil
}

// GetPlayerCheckouts will return a list containing all checkouts done by the given player
func GetPlayerCheckouts(playerID int) ([]*models.CheckoutStatistics, error) {
	rows, err := models.DB.Query(`
		SELECT
			s.player_id,
			s.first_dart, s.first_dart_multiplier,
			s.second_dart, s.second_dart_multiplier,
			s.third_dart, s.third_dart_multiplier,
			(IFNULL(s.first_dart, 0) * s.first_dart_multiplier +
				IFNULL(s.second_dart, 0) * s.second_dart_multiplier +
				IFNULL(s.third_dart, 0) * s.third_dart_multiplier) AS 'checkout',
			COUNT(*)
		FROM score s
		WHERE s.id IN (SELECT MAX(id) FROM score WHERE leg_id IN (
				SELECT l.id FROM leg l JOIN matches m ON m.id = l.match_id
				WHERE m.match_type_id = 1 AND l.winner_id = ?) GROUP BY leg_id)
		GROUP BY s.first_dart, s.first_dart_multiplier,
			s.second_dart, s.second_dart_multiplier,
			s.third_dart, s.third_dart_multiplier
		ORDER BY checkout DESC`, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playerVisits := make(map[int]map[string][]*models.Visit)
	checkoutCount := make(map[int]int)
	for rows.Next() {
		var checkout int
		v := new(models.Visit)
		v.FirstDart = new(models.Dart)
		v.SecondDart = new(models.Dart)
		v.ThirdDart = new(models.Dart)
		err := rows.Scan(&v.PlayerID,
			&v.FirstDart.Value, &v.FirstDart.Multiplier,
			&v.SecondDart.Value, &v.SecondDart.Multiplier,
			&v.ThirdDart.Value, &v.ThirdDart.Multiplier,
			&checkout, &v.Count)
		if err != nil {
			return nil, err
		}
		s := v.GetVisitString()
		if visitMap, ok := playerVisits[checkout]; ok {
			if visits, ok := visitMap[s]; ok {
				visitMap[s] = append(visits, v)
			} else {
				visits := make([]*models.Visit, 0)
				visitMap[s] = append(visits, v)
			}
		} else {
			visitMap := make(map[string][]*models.Visit)
			visits := make([]*models.Visit, 0)
			visitMap[s] = append(visits, v)
			playerVisits[checkout] = visitMap
		}

		if _, ok := checkoutCount[checkout]; ok {
			checkoutCount[checkout] += v.Count
		} else {
			checkoutCount[checkout] = v.Count
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	checkouts := make([]*models.CheckoutStatistics, 0)
	for i := 170; i >= 2; i-- {
		if i == 169 || i == 168 || i == 166 || i == 165 || i == 163 || i == 162 || i == 159 {
			// Skip values which cannot be checkouts
			continue
		}
		checkout := new(models.CheckoutStatistics)
		checkout.Checkout = i

		if visitMap, ok := playerVisits[i]; ok {
			visits := make([]*models.Visit, 0)
			for _, v := range visitMap {
				visits = append(visits, v...)
			}
			// Sort the visits by most common
			sort.Slice(visits, func(i, j int) bool {
				if visits[i].Count > visits[j].Count {
					return true
				}
				if visits[i].Count < visits[j].Count {
					return false
				}
				return true
			})
			checkout.Visits = visits
			checkout.Completed = true
		} else {
			checkout.Completed = false
		}

		if count, ok := checkoutCount[i]; ok {
			checkout.Count = count
		}
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

// GetPlayerTournamentStandings will return all tournament standings for the given player
func GetPlayerTournamentStandings(playerID int) ([]*models.PlayerTournamentStanding, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			t.id AS 'tournament_id',
			t.name AS 'tournament_name',
			tg.id AS 'tournament_group_id',
			tg.name AS 'tournament_group_name',
			tg.division AS 'tournament_group_division',
			ts.rank AS 'final_standing',
			MAX(ts2.rank) AS 'total_players',
			ts.elo as 'elo'
		FROM tournament_standings ts
			JOIN tournament t ON t.id = ts.tournament_id
			JOIN player p ON p.id = ts.player_id
			JOIN player2tournament p2t ON p2t.tournament_id = t.id AND p2t.player_id = p.id
			JOIN tournament_group tg ON tg.id = p2t.tournament_group_id
			JOIN tournament_standings ts2 ON ts2.tournament_id = t.id
		WHERE ts.player_id = ?
		GROUP BY t.id
		ORDER BY t.start_time DESC`, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	standings := make([]*models.PlayerTournamentStanding, 0)
	for rows.Next() {
		standing := new(models.PlayerTournamentStanding)
		standing.Tournament = new(models.Tournament)
		standing.TournamentGroup = new(models.TournamentGroup)

		err := rows.Scan(&standing.PlayerID, &standing.Tournament.ID, &standing.Tournament.Name, &standing.TournamentGroup.ID,
			&standing.TournamentGroup.Name, &standing.TournamentGroup.Division, &standing.FinalStanding, &standing.TotalPlayers,
			&standing.Elo)
		if err != nil {
			return nil, err
		}
		standings = append(standings, standing)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return standings, nil
}

// GetPlayerOfficialMatches will return an overview of all official matches for the given player
func GetPlayerOfficialMatches(playerID int) ([]*models.Match, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id, m.is_finished, m.is_abandoned, m.is_walkover, m.current_leg_id, m.winner_id, m.office_id, m.is_practice,
			m.created_at, m.updated_at, m.owe_type_id, m.venue_id, mt.id, mt.name, mt.description, mm.id, mm.name, mm.short_name,
			mm.wins_required, mm.legs_required, ot.id, ot.item, v.id, v.name, v.description, l.updated_at as 'last_throw',
			GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order) AS 'players'
		FROM matches m
			JOIN match_type mt ON mt.id = m.match_type_id
			JOIN match_mode mm ON mm.id = m.match_mode_id
			LEFT JOIN leg l ON l.id = m.current_leg_id
			LEFT JOIN owe_type ot ON ot.id = m.owe_type_id
			LEFT JOIN venue v on v.id = m.venue_id
			LEFT JOIN player2leg p2l ON p2l.match_id = m.id
			LEFT JOIN player2leg p2l2 ON p2l2.leg_id = l.id
		WHERE p2l2.player_id = ?
			AND m.tournament_id IS NOT NULL
		GROUP BY m.id
		ORDER BY m.id DESC`, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		m := new(models.Match)
		m.MatchType = new(models.MatchType)
		m.MatchMode = new(models.MatchMode)
		ot := new(models.OweType)
		venue := new(models.Venue)
		var players string
		err := rows.Scan(&m.ID, &m.IsFinished, &m.IsAbandoned, &m.IsWalkover, &m.CurrentLegID, &m.WinnerID, &m.OfficeID, &m.IsPractice, &m.CreatedAt, &m.UpdatedAt,
			&m.OweTypeID, &m.VenueID, &m.MatchType.ID, &m.MatchType.Name, &m.MatchType.Description,
			&m.MatchMode.ID, &m.MatchMode.Name, &m.MatchMode.ShortName, &m.MatchMode.WinsRequired, &m.MatchMode.LegsRequired,
			&ot.ID, &ot.Name, &venue.ID, &venue.Name, &venue.Description, &m.LastThrow, &players)
		if err != nil {
			return nil, err
		}
		if m.OweTypeID.Valid {
			m.OweType = ot
		}
		if m.VenueID.Valid {
			m.Venue = venue
		}

		m.Players = util.StringToIntArray(players)
		matches = append(matches, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// GetPlayerHeadToHead will return head to head statistics between the two players
func GetPlayerHeadToHead(player1 int, player2 int) (*models.StatisticsHead2Head, error) {
	head2head := new(models.StatisticsHead2Head)

	head2headMatches, err := GetHeadToHeadMatches(player1, player2)
	if err != nil {
		return nil, err
	}
	head2head.Head2HeadMatches = head2headMatches
	head2headWins := make(map[int64]int)
	for _, match := range head2headMatches {
		head2headWins[match.WinnerID.Int64]++
	}
	head2head.Head2HeadWins = head2headWins

	matches1, err := GetPlayerLastMatches(player1, 5)
	if err != nil {
		return nil, err
	}
	matches2, err := GetPlayerLastMatches(player2, 5)
	if err != nil {
		return nil, err
	}
	matches := make(map[int][]*models.Match)
	matches[player1] = matches1
	matches[player2] = matches2
	head2head.LastMatches = matches

	visits1, err := GetPlayerVisitCount(player1)
	if err != nil {
		return nil, err
	}
	visits2, err := GetPlayerVisitCount(player2)
	if err != nil {
		return nil, err
	}
	visits := make(map[int][]*models.Visit)
	visits[player1] = visits1
	visits[player2] = visits2
	head2head.PlayerVisits = visits

	checkouts1, err := GetPlayerCheckouts(player1)
	if err != nil {
		return nil, err
	}
	checkouts2, err := GetPlayerCheckouts(player2)
	if err != nil {
		return nil, err
	}
	checkouts := make(map[int][]*models.CheckoutStatistics)
	checkouts[player1] = checkouts1
	checkouts[player2] = checkouts2
	head2head.PlayerCheckouts = checkouts

	playerIDs := make([]int, 2)
	playerIDs[0] = player1
	playerIDs[1] = player2

	// Get 301 Statistics for each player
	stats, err := GetPlayersX01Statistics(playerIDs, 301)
	if err != nil {
		return nil, err
	}
	statistics := make(map[int]*models.StatisticsX01)
	for _, stat := range stats {
		statistics[stat.PlayerID] = stat
	}
	head2head.Player301Statistics = statistics

	// Get 501 Statistics for each player
	stats, err = GetPlayersX01Statistics(playerIDs, 501)
	if err != nil {
		return nil, err
	}
	statistics = make(map[int]*models.StatisticsX01)
	for _, stat := range stats {
		statistics[stat.PlayerID] = stat
	}
	head2head.Player501Statistics = statistics

	elos, err := GetPlayersElo(player1, player2)
	if err != nil {
		return nil, err
	}
	p1Elo := elos[0].CurrentElo
	p2Elo := elos[1].CurrentElo
	elos[0].WinProbability = GetPlayerWinProbability(p1Elo, p2Elo)
	elos[1].WinProbability = GetPlayerWinProbability(p2Elo, p1Elo)

	playerElos := make(map[int]*models.PlayerElo)
	for _, elo := range elos {
		playerElos[elo.PlayerID] = elo
	}
	head2head.PlayerElos = playerElos

	return head2head, nil
}

// UpdateEloForMatch will update the elo for each player in a match
func UpdateEloForMatch(matchID int) error {
	match, err := GetMatch(matchID)
	if err != nil {
		return err
	}
	if match.MatchType.ID != models.X01 || len(match.Players) != 2 || match.IsWalkover || match.IsAbandoned ||
		match.IsPractice || !match.IsFinished {
		// Don't calculate Elo for non-X01 matches, matches which does not have 2 players, and
		// matches which were walkovers
		return nil
	}
	//log.Printf("Updating Elo for players %v in match %d", match.Players, matchID)

	elos, err := GetPlayersElo(match.Players...)
	if err != nil {
		return err
	}
	p1 := elos[0]
	p2 := elos[1]

	wins, err := GetWinsPerPlayer(matchID)
	if err != nil {
		return err
	}

	// Calculate elo for winner and looser
	p1.CurrentEloNew, p2.CurrentEloNew = CalculateElo(p1.CurrentElo, p1.CurrentEloMatches, wins[p1.PlayerID], p2.CurrentElo,
		p2.CurrentEloMatches, wins[p2.PlayerID])
	p1.CurrentEloMatches++
	p2.CurrentEloMatches++
	if match.TournamentID.Valid {
		one, two := CalculateElo(int(p1.TournamentElo.Int64), p1.TournamentEloMatches, wins[p1.PlayerID],
			int(p2.TournamentElo.Int64), p2.TournamentEloMatches, wins[p2.PlayerID])
		p1.TournamentEloNew = null.IntFrom(int64(one))
		p2.TournamentEloNew = null.IntFrom(int64(two))
		p1.TournamentEloMatches++
		p2.TournamentEloMatches++
	}
	err = updateElo(matchID, p1, p2)
	if err != nil {
		return err
	}
	return nil
}

// GetPlayersElo will get the Elo for the given player IDs
func GetPlayersElo(playerIDs ...int) ([]*models.PlayerElo, error) {
	q, args, err := sqlx.In(`
			SELECT
				player_id,
				current_elo,
				current_elo_matches,
				tournament_elo,
				tournament_elo_matches
			FROM player_elo
			WHERE player_id IN (?)
			ORDER BY FIELD(player_id, ?)`, playerIDs, playerIDs)
	if err != nil {
		return nil, err
	}
	rows, err := models.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]*models.PlayerElo, 0)
	for rows.Next() {
		p := new(models.PlayerElo)
		err := rows.Scan(&p.PlayerID, &p.CurrentElo, &p.CurrentEloMatches, &p.TournamentElo, &p.TournamentEloMatches)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return players, nil
}

func updateElo(matchID int, player1 *models.PlayerElo, player2 *models.PlayerElo) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}
	var player1TournamentElo *int64
	var player1TournamentEloNew *int64
	if player1.TournamentEloNew.Int64 != 0 {
		player1TournamentElo = &player1.TournamentElo.Int64
		player1TournamentEloNew = &player1.TournamentEloNew.Int64
	}
	if player1.TournamentEloNew.Int64 == 0 {
		player1.TournamentEloNew = player1.TournamentElo
	}
	var player2TournamentElo *int64
	var player2TournamentEloNew *int64
	if player2.TournamentEloNew.Int64 != 0 {
		player2TournamentElo = &player2.TournamentElo.Int64
		player2TournamentEloNew = &player2.TournamentEloNew.Int64
	}
	if player2.TournamentEloNew.Int64 == 0 {
		player2.TournamentEloNew = player2.TournamentElo
	}

	// Update Elo fo player1
	_, err = tx.Exec(`UPDATE player_elo SET current_elo = ?, current_elo_matches = ?, tournament_elo = ?, tournament_elo_matches = ? WHERE player_id = ?`,
		player1.CurrentEloNew, player1.CurrentEloMatches, player1.TournamentEloNew, player1.TournamentEloMatches, player1.PlayerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Elo for player2
	_, err = tx.Exec(`UPDATE player_elo SET current_elo = ?, current_elo_matches = ?, tournament_elo = ?, tournament_elo_matches = ? WHERE player_id = ?`,
		player2.CurrentEloNew, player2.CurrentEloMatches, player2.TournamentEloNew, player2.TournamentEloMatches, player2.PlayerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Elo changelog for player1
	_, err = tx.Exec(`INSERT INTO player_elo_changelog (match_id, player_id, old_elo, new_elo, old_tournament_elo, new_tournament_elo) VALUES (?, ?, ?, ?, ?, ?)`,
		matchID, player1.PlayerID, player1.CurrentElo, player1.CurrentEloNew, player1TournamentElo, player1TournamentEloNew)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Elo changelog for player2
	_, err = tx.Exec(`INSERT INTO player_elo_changelog (match_id, player_id, old_elo, new_elo, old_tournament_elo, new_tournament_elo) VALUES (?, ?, ?, ?, ?, ?)`,
		matchID, player2.PlayerID, player2.CurrentElo, player2.CurrentEloNew, player2TournamentElo, player2TournamentEloNew)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// RecalculateElo will recalculate Elo for all players
func RecalculateElo() error {
	rows, err := models.DB.Query(`SELECT id FROM matches ORDER BY updated_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	matches := make([]int, 0)
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		matches = append(matches, id)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	for _, id := range matches {
		err = UpdateEloForMatch(id)
		if err != nil {
			return err
		}
	}
	return nil
}

// CalculateElo will calculate the Elo for each player based on the given information. Returned value is new Elo for player1 and player2 respectively
func CalculateElo(player1Elo int, player1Matches int, player1Score int, player2Elo int, player2Matches int, player2Score int) (int, int) {
	if player1Matches == 0 {
		player1Matches = 1
	}
	if player2Matches == 0 {
		player2Matches = 1
	}

	// P1 = Winner
	// P2 = Looser
	// PD = Points Difference
	// Multiplier = ln(abs(PD) + 1) * (2.2 / ((P1(old)-P2(old)) * 0.001 + 2.2))
	// Elo Winner = P1(old) + 800/num_matches * (1 - 1/(1 + 10 ^ (P2(old) - P1(old) / 400) ) )
	// Elo Looser = P2(old) + 800/num_matches * (0 - 1/(1 + 10 ^ (P2(old) - P1(old) / 400) ) )

	if player1Score > player2Score {
		multiplier := math.Log(math.Abs(float64(player1Score-player2Score))+1) * (2.2 / ((float64(player1Elo-player2Elo))*0.001 + 2.2))
		player1Elo, player2Elo = calculateElo(player1Elo, player1Matches, player2Elo, player2Matches, multiplier, false)
	} else if player1Score < player2Score {
		multiplier := math.Log(math.Abs(float64(player1Score-player2Score))+1) * (2.2 / ((float64(player2Elo-player1Elo))*0.001 + 2.2))
		player2Elo, player1Elo = calculateElo(player2Elo, player2Matches, player1Elo, player1Matches, multiplier, false)
	} else {
		player1Elo, player2Elo = calculateElo(player1Elo, player1Matches, player2Elo, player2Matches, 1.0, true)
	}
	// Cap Elo at 400 to avoid players going too low
	if player1Elo < 400 {
		player1Elo = 400
	}
	if player2Elo < 400 {
		player2Elo = 400
	}
	return player1Elo, player2Elo
}

func calculateElo(winnerElo int, winnerMatches int, looserElo int, looserMatches int, multiplier float64, isDraw bool) (int, int) {
	constant := 800.0

	Wwinner := 1.0
	Wlooser := 0.0
	if isDraw {
		Wwinner = 0.5
		Wlooser = 0.5
	}
	changeWinner := int((constant / float64(winnerMatches) * (Wwinner - (1 / (1 + math.Pow(10, float64(looserElo-winnerElo)/400))))) * multiplier)
	calculatedWinner := winnerElo + changeWinner

	changeLooser := int((constant / float64(looserMatches) * (Wlooser - (1 / (1 + math.Pow(10, float64(winnerElo-looserElo)/400))))) * multiplier)
	calculatedLooser := looserElo + changeLooser

	return calculatedWinner, calculatedLooser
}

func GetPlayerWinProbability(player1Elo int, player2Elo int) float64 {
	// Pr(A) = 1 / (10^(-ELODIFF/400) + 1)
	return 1 / (math.Pow(10, float64(-(player1Elo-player2Elo))/400) + 1)
}

func GetPlayerDrawProbability(player1Elo int, player2Elo int) float64 {
	// Quants Magic using Binomial Regression and Elos from 800 matches
	// Caveat: Model won´t be accurate for extreme cases (Elo Diff >500)
	// Formula:
	//   pDraw = 1 / (1 + exp(-(-1.479018 - 0.000001434670 * abs(eloDiff)^2)))
	eloDiff := math.Abs(float64(player1Elo - player2Elo))
	pDraw := 1 / (1 + math.Exp(-(-1.479018 - float64(0.000001434670)*eloDiff*eloDiff)))
	return pDraw
}
