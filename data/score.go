package data

import (
	"errors"
	"log"
	"math"
	"sort"
	"sync"

	"github.com/guregu/null"
	"github.com/kcapp/api/models"
)

var addVisitLock sync.Mutex

// AddVisit will write the given visit to database
func AddVisit(visit models.Visit) (*models.Visit, error) {
	addVisitLock.Lock()
	defer addVisitLock.Unlock()

	leg, err := GetLeg(visit.LegID)
	if err != nil {
		return nil, err
	}

	if leg.CurrentPlayerID != visit.PlayerID {
		return nil, errors.New("cannot insert score for non-current player")
	}
	if leg.IsFinished {
		return nil, errors.New("leg already finished")
	}

	match, err := GetMatch(leg.MatchID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(visit.LegID)
	if err != nil {
		return nil, err
	}

	matchType := match.MatchType.ID
	if leg.LegType != nil {
		matchType = leg.LegType.ID
	}

	isFinished := false
	// Invalidate extra darts not thrown, and check if leg is finished
	if matchType == models.X01 || matchType == models.X01HANDICAP {
		visit.SetIsBust(players[visit.PlayerID].CurrentScore)
		isFinished = !visit.IsBust && visit.IsCheckout(players[visit.PlayerID].CurrentScore)
	} else if matchType == models.SHOOTOUT {
		isFinished = ((len(leg.Visits) + 1) * 3) >= (9 * len(leg.Players))
		if isFinished {
			// Handle draw in legs with two players
			players[visit.PlayerID].CurrentScore += visit.GetScore()
			players[visit.PlayerID].DartsThrown += 3

			if len(players) == 2 {
				scores := make([]*models.Player2Leg, 0, len(players))
				for _, player := range players {
					scores = append(scores, player)
				}
				// If both players have thrown the same amount of darts, and have different scores, game is finished
				isFinished = scores[0].DartsThrown == scores[1].DartsThrown && scores[0].CurrentScore != scores[1].CurrentScore
			}
		}
	} else if matchType == models.CRICKET {
		isFinished, err = isCricketLegFinished(visit)
		if err != nil {
			return nil, err
		}
		if isFinished {
			if visit.ThirdDart.IsCricketMiss() {
				visit.ThirdDart.Value = null.IntFromPtr(nil)
			}
			if visit.SecondDart.IsCricketMiss() {
				visit.SecondDart.Value = null.IntFromPtr(nil)
			}
		}
	} else if matchType == models.DARTSATX {
		isFinished = ((len(leg.Visits)+1)*3)%(99*len(leg.Players)) == 0
	} else if matchType == models.AROUNDTHEWORLD {
		isFinished = (len(leg.Visits)+1)%(21*len(leg.Players)) == 0
	} else if matchType == models.AROUNDTHECLOCK {
		players[visit.PlayerID].CurrentScore += visit.CalculateAroundTheClockScore(players[visit.PlayerID].CurrentScore)
		if players[visit.PlayerID].CurrentScore == 21 {
			if visit.FirstDart.IsBull() {
				visit.SecondDart.Value = null.IntFromPtr(nil)
				visit.ThirdDart.Value = null.IntFromPtr(nil)
			} else if visit.SecondDart.IsBull() {
				visit.ThirdDart.Value = null.IntFromPtr(nil)
			}
		}
		isFinished = players[visit.PlayerID].CurrentScore == 21 && (visit.FirstDart.IsBull() || visit.SecondDart.IsBull() || visit.ThirdDart.IsBull())
	} else if matchType == models.SHANGHAI {
		round := int(math.Floor(float64(len(leg.Visits))/float64(len(leg.Players))) + 1)
		isFinished = (len(leg.Visits)+1)%(20*len(leg.Players)) == 0 || (visit.IsShanghai() && visit.FirstDart.ValueRaw() == round)
	} else if matchType == models.TICTACTOE {
		numbers := leg.Parameters.Numbers
		hits := leg.Parameters.Hits

		lastDartValid := visit.GetLastDart().IsDouble()
		if leg.Parameters.OutshotType.ID == models.OUTSHOTANY {
			lastDartValid = true
		} else if leg.Parameters.OutshotType.ID == models.OUTSHOTMASTER {
			lastDartValid = visit.GetLastDart().IsDouble() || visit.GetLastDart().IsTriple()
		}
		for _, num := range numbers {
			// Check if we hit the exact number, ending with a double
			if num == visit.GetScore() && lastDartValid {
				if visit.ThirdDart.IsMiss() {
					visit.ThirdDart.Value = null.IntFromPtr(nil)
					if visit.SecondDart.IsMiss() {
						visit.SecondDart.Value = null.IntFromPtr(nil)
					}
				}
				if _, ok := hits[num]; !ok {
					// Don't allow other players to take numbers already scored by another player
					hits[num] = visit.PlayerID
				}
				break
			}
		}
		// Check if current player has 3 in a row horizontally, diagonally or vertically
		if leg.Parameters.IsTicTacToeWinner(visit.PlayerID) {
			isFinished = true
		} else if leg.Parameters.IsTicTacToeDraw() || len(hits) == 9 {
			isFinished = true
		}
	} else if matchType == models.BERMUDATRIANGLE {
		isFinished = ((len(leg.Visits)+1)*3)%(39*len(leg.Players)) == 0
	} else if matchType == models.FOURTWENTY {
		isFinished = ((len(leg.Visits)+1)*3)%(63*len(leg.Players)) == 0
	} else if matchType == models.KILLBULL {
		score := players[visit.PlayerID].CurrentScore - visit.CalculateKillBullScore()
		if score <= 0 {
			if !visit.ThirdDart.IsBull() {
				visit.ThirdDart.Value = null.IntFromPtr(nil)
				if !visit.SecondDart.IsBull() {
					visit.SecondDart.Value = null.IntFromPtr(nil)
				}
			}
			isFinished = true
		}
	} else if matchType == models.GOTCHA {
		visit.SetIsBustAbove(players[visit.PlayerID].CurrentScore, leg.StartingScore)
		score := players[visit.PlayerID].CurrentScore + visit.CalculateGotchaScore(players, leg.StartingScore)
		if score == leg.StartingScore {
			isFinished = true
		}
	} else if matchType == models.JDCPRACTICE {
		isFinished = (len(leg.Visits)+1)%(19*len(leg.Players)) == 0
	} else if matchType == models.KNOCKOUT {
		idx := len(leg.Visits) - 1
		if idx >= 0 {
			if leg.Visits[idx].Score > visit.GetScore() {
				players[visit.PlayerID].Lives = null.IntFrom(players[visit.PlayerID].Lives.Int64 - 1)
			}
			playersAlive := 0
			for _, player := range players {
				if player.Lives.Int64 > 0 {
					playersAlive++
				}
			}
			isFinished = playersAlive < 2
		}
	} else if matchType == models.SCAM {
		// Only stoppers can finish the match
		if players[visit.PlayerID].IsStopper.Bool {
			hits := players[visit.PlayerID].Hits
			hits.Add(visit.FirstDart)
			hits.Add(visit.SecondDart)
			hits.Add(visit.ThirdDart)

			if hits.Contains(models.SINGLE, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20) {
				// Invalidate the last darts incase we "checked out" with only 1 or two darts
				if visit.ThirdDart.ValueRaw() == 0 {
					visit.ThirdDart.Value = null.IntFromPtr(nil)
				}
				if visit.SecondDart.ValueRaw() == 0 {
					visit.SecondDart.Value = null.IntFromPtr(nil)
				}
			}

			allStopped := true
			for _, player := range players {
				// Check if all players have closed all numbers
				if !player.Hits.Contains(models.SINGLE, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20) {
					allStopped = false
					break
				}
			}
			if allStopped {
				isFinished = true
			}
		}
	}

	// Determine who will be the next player
	order := make(map[int]int)
	for _, player := range players {
		if !player.IsOut(matchType, visit) {
			order[player.Order] = player.PlayerID
		}
	}

	// Set new player order on remaining players, from 1 to n
	for i, key := range getKeys(order) {
		players[order[key]].Order = i + 1
	}

	newOrder := make(map[int]int)
	currentPlayerOrder := 1
	for _, playerID := range order {
		player := players[playerID]
		if playerID == visit.PlayerID {
			currentPlayerOrder = player.Order
		}
		newOrder[player.Order] = player.PlayerID
	}
	nextPlayerID := newOrder[(currentPlayerOrder%len(newOrder))+1]

	tx, err := models.DB.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(`
		INSERT INTO score(
			leg_id, player_id,
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier,
			is_bust, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		visit.LegID, visit.PlayerID,
		visit.FirstDart.Value, visit.FirstDart.Multiplier,
		visit.SecondDart.Value, visit.SecondDart.Multiplier,
		visit.ThirdDart.Value, visit.ThirdDart.Multiplier,
		visit.IsBust)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec(`UPDATE leg SET current_player_id = ?, updated_at = NOW() WHERE id = ?`, nextPlayerID, visit.LegID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	log.Printf("[%d] Added score for player %d, (%d-%d, %d-%d, %d-%d, %t)", visit.LegID, visit.PlayerID, visit.FirstDart.Value.Int64,
		visit.FirstDart.Multiplier, visit.SecondDart.Value.Int64, visit.SecondDart.Multiplier, visit.ThirdDart.Value.Int64, visit.ThirdDart.Multiplier,
		visit.IsBust)

	if isFinished {
		err = FinishLeg(visit)
		if err != nil {
			return nil, err
		}
	}

	return &visit, nil
}

// ModifyVisit modify the scores of a visit
func ModifyVisit(visit models.Visit) error {
	// FIXME: We need to check if this is a checkout/bust
	stmt, err := models.DB.Prepare(`
		UPDATE score SET
    		first_dart = ?,
    		first_dart_multiplier = ?,
    		second_dart = ?,
    		second_dart_multiplier = ?,
    		third_dart = ?,
		    third_dart_multiplier = ?,
			updated_at = NOW()
		WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(visit.FirstDart.Value, visit.FirstDart.Multiplier, visit.SecondDart.Value, visit.SecondDart.Multiplier,
		visit.ThirdDart.Value, visit.ThirdDart.Multiplier, visit.ID)
	if err != nil {
		return err
	}
	log.Printf("[%d] Modified score %d, throws: (%d-%d, %d-%d, %d-%d)", visit.LegID, visit.ID, visit.FirstDart.Value.Int64,
		visit.FirstDart.Multiplier, visit.SecondDart.Value.Int64, visit.SecondDart.Multiplier, visit.ThirdDart.Value.Int64, visit.ThirdDart.Multiplier)

	return nil
}

// DeleteVisit will delete the visit for the given ID
func DeleteVisit(id int) error {
	visit, err := GetVisit(id)
	if err != nil {
		return err
	}
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}
	// Delete the visit
	_, err = tx.Exec("DELETE FROM score WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Set current player to the player of the last visit
	_, err = tx.Exec("UPDATE leg SET current_player_id = ? WHERE id = ?", visit.PlayerID, visit.LegID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	log.Printf("[%d] Deleted visit %d", visit.LegID, visit.ID)
	return nil
}

// DeleteLastVisit will delete the last visit for the given leg
func DeleteLastVisit(legID int) error {
	visits, err := GetLegVisits(legID)
	if err != nil {
		return err
	}

	if len(visits) > 0 {
		err := DeleteVisit(visits[len(visits)-1].ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetPlayerVisits will return all visits for a given player
func GetPlayerVisits(id int) ([]*models.Visit, error) {
	rows, err := models.DB.Query(`
		SELECT
			id, leg_id, player_id,
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier,
			is_bust,
			created_at,
			updated_at
		FROM score s
		WHERE player_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	visits := make([]*models.Visit, 0)
	for rows.Next() {
		v := new(models.Visit)
		v.FirstDart = new(models.Dart)
		v.SecondDart = new(models.Dart)
		v.ThirdDart = new(models.Dart)
		err := rows.Scan(&v.ID, &v.LegID, &v.PlayerID,
			&v.FirstDart.Value, &v.FirstDart.Multiplier,
			&v.SecondDart.Value, &v.SecondDart.Multiplier,
			&v.ThirdDart.Value, &v.ThirdDart.Multiplier,
			&v.IsBust, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}
		visits = append(visits, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return visits, nil
}

// GetLegVisits will return all visits for a given leg
func GetLegVisits(id int) ([]*models.Visit, error) {
	rows, err := models.DB.Query(`
		SELECT
			id, leg_id, player_id,
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier,
			is_bust,
			created_at,
			updated_at
		FROM score s
		WHERE leg_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	visits := make([]*models.Visit, 0)
	for rows.Next() {
		v := new(models.Visit)
		v.FirstDart = new(models.Dart)
		v.SecondDart = new(models.Dart)
		v.ThirdDart = new(models.Dart)
		err := rows.Scan(&v.ID, &v.LegID, &v.PlayerID,
			&v.FirstDart.Value, &v.FirstDart.Multiplier,
			&v.SecondDart.Value, &v.SecondDart.Multiplier,
			&v.ThirdDart.Value, &v.ThirdDart.Multiplier,
			&v.IsBust, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}
		visits = append(visits, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return visits, nil
}

// GetVisit will return the visit with the given ID
func GetVisit(id int) (*models.Visit, error) {
	v := new(models.Visit)
	v.FirstDart = new(models.Dart)
	v.SecondDart = new(models.Dart)
	v.ThirdDart = new(models.Dart)
	err := models.DB.QueryRow(`
		SELECT
			id, leg_id, player_id,
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier,
			is_bust,
			created_at,
			updated_at
		FROM score s
		WHERE s.id = ?`, id).Scan(&v.ID, &v.LegID, &v.PlayerID,
		&v.FirstDart.Value, &v.FirstDart.Multiplier,
		&v.SecondDart.Value, &v.SecondDart.Multiplier,
		&v.ThirdDart.Value, &v.ThirdDart.Multiplier,
		&v.IsBust, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// GetLastVisits will return the last N visit for the given leg
func GetLastVisits(legID int, num int) (map[int]*models.Visit, error) {
	rows, err := models.DB.Query(`
			SELECT
				player_id,
				first_dart, first_dart_multiplier,
				second_dart, second_dart_multiplier,
				third_dart, third_dart_multiplier
			FROM score
			WHERE leg_id = ? AND is_bust = 0
			ORDER BY id DESC LIMIT ?`, legID, num)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	visits := make(map[int]*models.Visit)
	for rows.Next() {
		v := new(models.Visit)
		v.FirstDart = new(models.Dart)
		v.SecondDart = new(models.Dart)
		v.ThirdDart = new(models.Dart)
		err := rows.Scan(&v.PlayerID,
			&v.FirstDart.Value, &v.FirstDart.Multiplier,
			&v.SecondDart.Value, &v.SecondDart.Multiplier,
			&v.ThirdDart.Value, &v.ThirdDart.Multiplier)
		if err != nil {
			return nil, err
		}
		visits[v.PlayerID] = v
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return visits, nil
}

// GetPlayerVisitCount will return a count of each visit for a given player
func GetPlayerVisitCount(playerID int) ([]*models.Visit, error) {
	rows, err := models.DB.Query(`
		SELECT
			player_id,
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier,
			COUNT(*) AS 'visits'
		FROM score s
			WHERE player_id = ?
		GROUP BY
			player_id, first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier`, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[string]*models.Visit)
	for rows.Next() {
		v := new(models.Visit)
		v.FirstDart = new(models.Dart)
		v.SecondDart = new(models.Dart)
		v.ThirdDart = new(models.Dart)
		err := rows.Scan(&v.PlayerID,
			&v.FirstDart.Value, &v.FirstDart.Multiplier,
			&v.SecondDart.Value, &v.SecondDart.Multiplier,
			&v.ThirdDart.Value, &v.ThirdDart.Multiplier,
			&v.Count)
		if err != nil {
			return nil, err
		}

		s := v.GetVisitString()
		if val, ok := m[s]; ok {
			val.Count += v.Count
		} else {
			m[s] = v
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	visits := make([]*models.Visit, 0)
	for _, v := range m {
		visits = append(visits, v)
	}

	return visits, nil
}

// GetRandomLegForPlayer will return a random leg for a given player and starting score
func GetRandomLegForPlayer(playerID int, startingScore int) ([]*models.Visit, error) {
	var legID int
	err := models.DB.QueryRow(`
		SELECT
			l.id
		FROM leg l
			JOIN player2leg p2l ON p2l.leg_id = l.id
		WHERE l.is_finished = 1 AND l.winner_id = ? AND l.starting_score = ? AND l.has_scores = 1
		GROUP BY l.id
			HAVING COUNT(DISTINCT p2l.player_id) = 2
		ORDER BY RAND()
		LIMIT 1`, playerID, startingScore).Scan(&legID)
	if err != nil {
		return nil, err
	}

	rows, err := models.DB.Query(`
		SELECT
			id, leg_id, player_id,
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier,
			is_bust,
			created_at,
			updated_at
		FROM score s
		WHERE leg_id = ? AND player_id = ?`, legID, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	visits := make([]*models.Visit, 0)
	for rows.Next() {
		v := new(models.Visit)
		v.FirstDart = new(models.Dart)
		v.SecondDart = new(models.Dart)
		v.ThirdDart = new(models.Dart)
		err := rows.Scan(&v.ID, &v.LegID, &v.PlayerID,
			&v.FirstDart.Value, &v.FirstDart.Multiplier,
			&v.SecondDart.Value, &v.SecondDart.Multiplier,
			&v.ThirdDart.Value, &v.ThirdDart.Multiplier,
			&v.IsBust, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, err
		}
		visits = append(visits, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return visits, nil
}

// GetDartStatistics will return statistics of times hit for a given dart
func GetDartStatistics(dart int) (map[int]*models.Hits, error) {
	rows, err := models.DB.Query(`
		SELECT player_id, singles, doubles, triples
		FROM (
			SELECT s.player_id,
				SUM(IF(s.first_dart = ? AND s.first_dart_multiplier = 1, 1, 0) +
					IF(s.second_dart = ? AND s.second_dart_multiplier = 1, 1, 0) +
					IF(s.third_dart = ? AND s.third_dart_multiplier = 1, 1, 0)) AS 'singles',
				SUM(IF(s.first_dart = ? AND s.first_dart_multiplier = 2, 1, 0) +
					IF(s.second_dart = ? AND s.second_dart_multiplier = 2, 1, 0) +
					IF(s.third_dart = ? AND s.third_dart_multiplier = 2, 1, 0)) AS 'doubles',
				SUM(IF(s.first_dart = ? AND s.first_dart_multiplier = 3, 1, 0) +
					IF(s.second_dart = ? AND s.second_dart_multiplier = 3, 1, 0) +
					IF(s.third_dart = ? AND s.third_dart_multiplier = 3, 1, 0)) AS 'triples'
			FROM score s
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			WHERE s.is_bust = 0
			GROUP BY player_id
		) scores`, dart, dart, dart, dart, dart, dart, dart, dart, dart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[int]*models.Hits)
	for rows.Next() {
		h := new(models.Hits)
		var playerID int
		err := rows.Scan(&playerID, &h.Singles, &h.Doubles, &h.Triples)
		if err != nil {
			return nil, err
		}
		m[playerID] = h
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func isCricketLegFinished(visit models.Visit) (bool, error) {
	players, err := GetLegPlayers(visit.LegID)
	if err != nil {
		return false, err
	}
	allPlayers := make(map[int]*models.Player2Leg)
	for _, player := range players {
		allPlayers[player.PlayerID] = player
	}

	// Add score for incoming visit
	visit.CalculateCricketScore(allPlayers)
	for _, player := range players {
		player.CurrentScore = allPlayers[player.PlayerID].CurrentScore
	}

	// Did current player close all numbers?
	player := allPlayers[visit.PlayerID]
	closed := true
	for _, dart := range models.CRICKETDARTS {
		if player.Hits[dart] == nil || player.Hits[dart].Total < 3 {
			closed = false
			break
		}
	}

	// What is the lowest score?
	lowestScore := math.MaxInt32
	for _, player := range players {
		if player.CurrentScore < lowestScore {
			lowestScore = player.CurrentScore
		}
	}

	// If current player closed all numbers and has the lowest score, it's finished
	if closed && player.CurrentScore == lowestScore {
		return true, nil
	}
	return false, nil
}

// getKeys will return all keys as a sorted slice for the given map
func getKeys(m map[int]int) []int {
	keys := make([]int, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	sort.Ints(keys)
	return keys
}
