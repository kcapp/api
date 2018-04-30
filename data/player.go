package data

import (
	"log"
	"sort"

	"github.com/kcapp/api/models"
)

// GetPlayers returns a map of all players
func GetPlayers() (map[int]*models.Player, error) {
	played, err := GetGamesPlayedPerPlayer()
	if err != nil {
		return nil, err
	}

	rows, err := models.DB.Query(`SELECT p.id, p.name, p.nickname, p.color, p.profile_pic_url, p.created_at FROM player p`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.Name, &p.Nickname, &p.Color, &p.ProfilePicURL, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		players[p.ID] = p

		if val, ok := played[p.ID]; ok {
			p.GamesPlayed = val.GamesPlayed
			p.GamesWon = val.GamesWon
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
	err := models.DB.QueryRow(`SELECT p.id, p.name, p.nickname, p.color, p.profile_pic_url, p.created_at FROM player p WHERE p.id = ?`, id).
		Scan(&p.ID, &p.Name, &p.Nickname, &p.Color, &p.ProfilePicURL, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	pld, err := GetGamesPlayedPerPlayer()
	if err != nil {
		return nil, err
	}
	played := pld[p.ID]
	if played != nil {
		p.GamesPlayed = played.GamesPlayed
		p.GamesWon = played.GamesWon
		p.LegsPlayed = played.LegsPlayed
		p.LegsWon = played.LegsWon
	}

	return p, nil
}

// AddPlayer will add a new player to the database
func AddPlayer(player models.Player) error {
	// Prepare statement for inserting data
	stmt, err := models.DB.Prepare("INSERT INTO player (name, nickname, color, profile_pic_url) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(player.Name, player.Nickname, player.Color, player.ProfilePicURL)
	log.Printf("Created new player %s", player.Name)
	return err
}

// UpdatePlayer will update the given player
func UpdatePlayer(playerID int, player models.Player) error {
	// Prepare statement for inserting data
	stmt, err := models.DB.Prepare("UPDATE player SET name = ?, nickname = ?, color = ?, profile_pic_url = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(player.Name, player.Nickname, player.Color, player.ProfilePicURL, playerID)
	log.Printf("Updated player %s (%v)", player.Name, player)
	return err
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
	rows, err := models.DB.Query(`
		SELECT
			p2m.leg_id,
			p2m.player_id,
			p2m.order,
			p2m.handicap,
			p2m.player_id = m.current_player_id AS 'is_current_player',
			(m.starting_score +
				-- Add handicap for players if game_mode is handicap
				IF(g.game_type_id = 3, IFNULL(p2m.handicap, 0), 0)) -
				(IFNULL(SUM(first_dart * first_dart_multiplier), 0) +
				IFNULL(SUM(second_dart * second_dart_multiplier), 0) +
				IFNULL(SUM(third_dart * third_dart_multiplier), 0))
				-- For X01 score goes down, while Shootout it counts up
				* IF(g.game_type_id = 2, -1, 1) AS 'current_score'
		FROM player2leg p2m
			LEFT JOIN leg m ON m.id = p2m.leg_id
			LEFT JOIN score s ON s.leg_id = p2m.leg_id AND s.player_id = p2m.player_id
			LEFT JOIN game g on g.id = m.game_id
		WHERE p2m.leg_id = ? AND (s.is_bust IS NULL OR is_bust = 0)
		GROUP BY p2m.player_id
		ORDER BY p2m.order ASC`, legID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scores := make(map[int]*models.Player2Leg)
	for rows.Next() {
		p2m := new(models.Player2Leg)
		err := rows.Scan(&p2m.LegID, &p2m.PlayerID, &p2m.Order, &p2m.Handicap, &p2m.IsCurrentPlayer, &p2m.CurrentScore)
		if err != nil {
			return nil, err
		}
		scores[p2m.PlayerID] = p2m
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return scores, nil
}

// GetGamesPlayedPerPlayer will get the number of games and legs played and won for each player
func GetGamesPlayedPerPlayer() (map[int]*models.Player, error) {
	rows, err := models.DB.Query(`
		SELECT
			player_id,
			MAX(games_played) AS 'games_played',
			MAX(games_won) AS 'games_won',
			MAX(legs_played) AS 'legs_played',
			MAX(legs_won) AS 'legs_won'
		FROM (
			SELECT
				p2m.player_id,
				COUNT(DISTINCT p2m.game_id) AS 'games_played',
				0 AS 'games_won',
				COUNT(m.id)  AS 'legs_played',
				SUM(CASE WHEN p2m.player_id = m.winner_id THEN 1 ELSE 0 END) AS 'legs_won'
			FROM player2leg p2m
				JOIN ` + "`leg`" + ` m ON m.id = p2m.leg_id
				JOIN game g ON g.id = p2m.game_id
			WHERE m.is_finished = 1
				AND g.game_type_id = 1
			GROUP BY p2m.player_id
			UNION ALL
			SELECT
				p2m.player_id,
				0 AS 'games_played',
				COUNT(DISTINCT g.id) AS 'games_won',
				0 AS 'legs_played',
				0 AS 'legs_won'
			FROM game g
				JOIN ` + "`leg`" + ` m ON m.game_id = g.id
				JOIN player2leg p2m ON p2m.player_id = g.winner_id AND p2m.game_id = g.id
			WHERE m.is_finished = 1
				AND g.game_type_id = 1
			GROUP BY g.winner_id
		) games
		GROUP BY player_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	played := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.GamesPlayed, &p.GamesWon, &p.LegsPlayed, &p.LegsWon)
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
				SELECT m.id FROM leg m
				JOIN game g ON g.id = m.game_id
				WHERE g.game_type_id = 1 AND m.winner_id = ?) GROUP BY leg_id)
		GROUP BY s.first_dart, s.first_dart_multiplier,
			s.second_dart, s.second_dart_multiplier,
			s.third_dart, s.third_dart_multiplier
		ORDER BY checkout`, playerID)
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
	for i := 2; i < 171; i++ {
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

// GetPlayerHeadToHead will return head to head statistics between the two players
func GetPlayerHeadToHead(player1 int, player2 int) (*models.StatisticsHead2Head, error) {
	head2head := new(models.StatisticsHead2Head)

	head2headGames, err := GetHeadToHeadGames(player1, player2)
	if err != nil {
		return nil, err
	}
	head2head.Head2HeadGames = head2headGames
	head2headWins := make(map[int64]int)
	for _, game := range head2headGames {
		head2headWins[game.WinnerID.Int64]++
	}
	head2head.Head2HeadWins = head2headWins

	games1, err := GetPlayerLastGames(player1, 5)
	if err != nil {
		return nil, err
	}
	games2, err := GetPlayerLastGames(player2, 5)
	if err != nil {
		return nil, err
	}
	games := make(map[int][]*models.Game)
	games[player1] = games1
	games[player2] = games2
	head2head.LastGames = games

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

	return head2head, nil
}
