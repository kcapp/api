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
	tx.Exec("UPDATE game SET current_match_id = ? WHERE id = ?", matchID, gameID)

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
func FinishMatch(visit models.Visit) (*models.Match, error) {
	tx, err := models.DB.Begin()
	if err != nil {
		return nil, err
	}

	err = AddVisit(visit)
	if err != nil {
		return nil, err
	}
	// Update match with winner

	// Write statistics for each player

	// Check if game is finished or not

	// If game is finished, payback owes

	tx.Commit()

	// TODO Logging

	return nil, nil
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
			m.starting_score,
			p2m.player_id = m.current_player_id AS 'is_current_player',
			m.starting_score - IFNULL(
				SUM(first_dart * first_dart_multiplier) + 
				SUM(second_dart * second_dart_multiplier) + 
				SUM(third_dart * third_dart_multiplier), 0) AS 'current_score'
		FROM player2match p2m 
		LEFT JOIN `+"`match`"+` m ON m.id = p2m.match_id
		LEFT JOIN score s ON s.match_id = p2m.match_id AND s.player_id = p2m.player_id
		WHERE p2m.match_id = ? AND (s.is_bust IS NULL OR is_bust = 0)
		GROUP BY p2m.player_id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]*models.Player2Match, 0)
	for rows.Next() {
		p2m := new(models.Player2Match)
		err := rows.Scan(&p2m.MatchID, &p2m.PlayerID, &p2m.Order, &p2m.CurrentScore, &p2m.IsCurrentPlayer, &p2m.CurrentScore)
		if err != nil {
			return nil, err
		}
		players = append(players, p2m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
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
