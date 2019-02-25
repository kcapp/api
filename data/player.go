package data

import (
	"log"
	"math"
	"sort"

	"github.com/guregu/null"

	"github.com/jmoiron/sqlx"
	"github.com/kcapp/api/models"
)

// GetPlayers returns a map of all players
func GetPlayers() (map[int]*models.Player, error) {
	played, err := GetMatchesPlayedPerPlayer()
	if err != nil {
		return nil, err
	}

	rows, err := models.DB.Query(`SELECT p.id, p.first_name, p.last_name, p.vocal_name, p.nickname, p.color, p.profile_pic_url, p.created_at FROM player p`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.VocalName, &p.Nickname, &p.Color, &p.ProfilePicURL, &p.CreatedAt)
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

// GetActivePlayers returns a map of all active players
func GetActivePlayers() (map[int]*models.Player, error) {
	played, err := GetMatchesPlayedPerPlayer()
	if err != nil {
		return nil, err
	}

	rows, err := models.DB.Query(`
		SELECT
			p.id, p.first_name, p.last_name, p.vocal_name, p.nickname, p.color, p.profile_pic_url, p.created_at 
		FROM player p
		WHERE active = 1`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.VocalName, &p.Nickname, &p.Color, &p.ProfilePicURL, &p.CreatedAt)
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

// GetPlayer returns the player for the given ID
func GetPlayer(id int) (*models.Player, error) {
	p := new(models.Player)
	err := models.DB.QueryRow(`SELECT p.id, p.first_name, p.last_name, p.vocal_name, p.nickname, p.color, p.profile_pic_url, p.created_at, pe.current_elo, pe.tournament_elo
		FROM player p
		JOIN player_elo pe on pe.player_id = p.id
		WHERE p.id = ?`, id).
		Scan(&p.ID, &p.FirstName, &p.LastName, &p.VocalName, &p.Nickname, &p.Color, &p.ProfilePicURL, &p.CreatedAt, &p.CurrentElo, &p.TournamentElo)
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

// AddPlayer will add a new player to the database
func AddPlayer(player models.Player) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}

	// Prepare statement for inserting data
	res, err := tx.Exec("INSERT INTO player (first_name, last_name, vocal_name, nickname, color, profile_pic_url) VALUES (?, ?, ?, ?, ?, ?)",
		player.FirstName, player.LastName, player.VocalName, player.Nickname, player.Color, player.ProfilePicURL)
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
	stmt, err := models.DB.Prepare("UPDATE player SET first_name = ?, last_name = ?, vocal_name = ?, nickname = ?, color = ?, profile_pic_url = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(player.FirstName, player.LastName, player.VocalName, player.Nickname, player.Color, player.ProfilePicURL, playerID)
	if err != nil {
		return err
	}
	log.Printf("Updated player %s (%v)", player.FirstName, player)
	return nil
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
			(l.starting_score +
				-- Add handicap for players if match_mode is handicap
				IF(m.match_type_id = 3, IFNULL(p2l.handicap, 0), 0)) -
				(IFNULL(SUM(first_dart * first_dart_multiplier), 0) +
				IFNULL(SUM(second_dart * second_dart_multiplier), 0) +
				IFNULL(SUM(third_dart * third_dart_multiplier), 0))
				-- For X01 score goes down, while Shootout it counts up
				* IF(m.match_type_id = 2, -1, 1) AS 'current_score'
		FROM player2leg p2l
			LEFT JOIN player p on p.id = p2l.player_id
			LEFT JOIN leg l ON l.id = p2l.leg_id
			LEFT JOIN score s ON s.leg_id = p2l.leg_id AND s.player_id = p2l.player_id
			LEFT JOIN matches m on m.id = l.match_id
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
		err := rows.Scan(&p2l.LegID, &p2l.PlayerID, &p2l.PlayerName, &p2l.Order, &p2l.Handicap, &p2l.IsCurrentPlayer, &p2l.CurrentScore)
		if err != nil {
			return nil, err
		}
		p2l.Player = players[p2l.PlayerID]
		scores[p2l.PlayerID] = p2l
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return scores, nil
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
			p.color,
			p.profile_pic_url
		FROM player2leg p2l
		LEFT JOIN player p ON p.id = p2l.player_id WHERE p2l.leg_id = ?`, legID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.VocalName, &p.Nickname, &p.Color, &p.ProfilePicURL)
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

// GetMatchesPlayedPerPlayer will get the number of matches and legs played and won for each player
func GetMatchesPlayedPerPlayer() (map[int]*models.Player, error) {
	rows, err := models.DB.Query(`
		SELECT
			player_id,
			MAX(matches_played) AS 'matches_played',
			MAX(matches_won) AS 'matches_won',
			MAX(legs_played) AS 'legs_played',
			MAX(legs_won) AS 'legs_won'
		FROM (
			SELECT
				p2l.player_id,
				COUNT(DISTINCT p2l.match_id) AS 'matches_played',
				0 AS 'matches_won',
				COUNT(m.id)  AS 'legs_played',
				SUM(CASE WHEN p2l.player_id = m.winner_id THEN 1 ELSE 0 END) AS 'legs_won'
			FROM player2leg p2l
				JOIN leg l ON l.id = p2l.leg_id
				JOIN matches m ON m.id = p2l.match_id
			WHERE l.is_finished = 1 AND m.is_abandoned = 0
				AND m.match_type_id = 1
			GROUP BY p2l.player_id
			UNION ALL
			SELECT
				p2l.player_id,
				0 AS 'matches_played',
				COUNT(DISTINCT m.id) AS 'matches_won',
				0 AS 'legs_played',
				0 AS 'legs_won'
			FROM matches m
				JOIN leg l ON l.match_id = m.id
				JOIN player2leg p2l ON p2l.player_id = m.winner_id AND p2l.match_id = m.id
			WHERE l.is_finished = 1 AND m.is_abandoned = 0
				AND m.match_type_id = 1
			GROUP BY m.winner_id
		) matches
		GROUP BY player_id`)
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
	elos[0].WinProbability = getPlayerWinProbability(p1Elo, p2Elo)
	elos[1].WinProbability = getPlayerWinProbability(p2Elo, p1Elo)

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
	if match.MatchType.ID != models.X01 || len(match.Players) != 2 || match.IsWalkover {
		// Don't calculate Elo for non-X01 matches, matches which does not have 2 players, and
		// matches which were walkovers
		return nil
	}
	log.Printf("Updating Elo for players %v in match %d", match.Players, matchID)

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
	rows, err := models.DB.Query(`
		SELECT
			m.id
		FROM matches m
		WHERE m.tournament_id IN (15, 16)
		ORDER BY m.created_at`)
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

func getPlayerWinProbability(player1Elo int, player2Elo int) float64 {
	// Pr(A) = 1 / (10^(-ELODIFF/400) + 1)
	return 1 / (math.Pow(10, float64(-(player1Elo-player2Elo))/400) + 1)
}
