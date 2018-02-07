package data

import (
	"errors"
	"log"

	"github.com/kcapp/api/models"
	"github.com/kcapp/api/util"
)

// NewGame will insert a new game in the database
func NewGame(game models.Game) (*models.Game, error) {
	tx, err := models.DB.Begin()
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec("INSERT INTO game (game_type_id, owe_type_id, created_at) VALUES (?, ?, NOW())", game.GameType.ID, game.OweTypeID)
	if err != nil {
		return nil, err
	}
	gameID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	res, err = tx.Exec("INSERT INTO `match` (starting_score, current_player_id, game_id, created_at) VALUES (?, ?, ?, NOW()) ", game.Matches[0].StartingScore, game.Players[0], gameID)
	if err != nil {
		return nil, err
	}
	matchID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	tx.Exec("UPDATE game SET current_match_id = ? WHERE id = ?", matchID, gameID)
	for idx, playerID := range game.Players {
		order := idx + 1
		res, err = tx.Exec("INSERT INTO player2match (player_id, match_id, `order`, game_id) VALUES (?, ?, ?, ?)", playerID, matchID, order, gameID)
		if err != nil {
			return nil, err
		}
	}
	tx.Commit()
	log.Printf("Started new game %d", gameID)
	return GetGame(int(gameID))
}

// GetGames returns all games
func GetGames() ([]*models.Game, error) {
	rows, err := models.DB.Query(`
		SELECT
			g.id, g.is_finished, g.current_match_id, g.winner_id, g.created_at, g.updated_at, g.owe_type_id,
			gt.id, gt.name, gt.short_name, gt.wins_required, gt.matches_required,
			ot.id, ot.item,
			GROUP_CONCAT(DISTINCT p2m.player_id) AS 'players'
		FROM game g
		LEFT JOIN game_type gt ON gt.id = g.game_type_id
		LEFT JOIN owe_type ot ON ot.id = g.owe_type_id
		LEFT JOIN player2match p2m ON p2m.game_id = g.id
		GROUP BY g.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]*models.Game, 0)
	for rows.Next() {
		g := new(models.Game)
		g.GameType = new(models.GameType)
		ot := new(models.OweType)
		var players string
		err := rows.Scan(&g.ID, &g.IsFinished, &g.CurrentMatchID, &g.WinnerID, &g.CreatedAt, &g.UpdatedAt, &g.OweTypeID,
			&g.GameType.ID, &g.GameType.Name, &g.GameType.ShortName, &g.GameType.WinsRequired, &g.GameType.MatchesRequired,
			&ot.ID, &ot.Item, &players)
		if err != nil {
			return nil, err
		}
		if g.OweTypeID.Valid {
			g.OweType = ot
		}

		g.Players = util.StringToIntArray(players)
		games = append(games, g)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

// GetGame returns a game with the given ID
func GetGame(id int) (*models.Game, error) {
	g := new(models.Game)
	g.GameType = new(models.GameType)
	ot := new(models.OweType)
	var players string
	err := models.DB.QueryRow(`
        SELECT
            g.id, g.is_finished, g.current_match_id, g.winner_id, g.created_at, g.updated_at, g.owe_type_id,
			gt.id, gt.name, gt.short_name, gt.wins_required, gt.matches_required,
			ot.id, ot.item,
			GROUP_CONCAT(DISTINCT p2m.player_id) AS 'players'
        FROM game g
		LEFT JOIN game_type gt ON gt.id = g.game_type_id
		LEFT JOIN owe_type ot ON ot.id = g.owe_type_id
		LEFT JOIN player2match p2m ON p2m.game_id = g.id
		WHERE g.id = ?`, id).Scan(&g.ID, &g.IsFinished, &g.CurrentMatchID, &g.WinnerID, &g.CreatedAt, &g.UpdatedAt, &g.OweTypeID,
		&g.GameType.ID, &g.GameType.Name, &g.GameType.ShortName, &g.GameType.WinsRequired, &g.GameType.MatchesRequired, &ot.ID, &ot.Item, &players)
	if err != nil {
		return nil, err
	}
	if g.OweTypeID.Valid {
		g.OweType = ot
	}
	g.Players = util.StringToIntArray(players)
	matches, err := GetMatchesForGame(id)
	if err != nil {
		return nil, err
	}
	g.Matches = matches
	return g, nil
}

// ContinueGame will either return the current match or create a new match
func ContinueGame(id int) (*models.Match, error) {
	game, err := GetGame(id)
	if err != nil {
		return nil, err
	}
	if game.IsFinished {
		return nil, errors.New("Cannot continue finished game")
	}

	matches, err := GetMatchesForGame(id)
	if err != nil {
		return nil, err
	}
	match := matches[len(matches)-1]
	if match.IsFinished {
		return NewMatch(id, match.StartingScore, match.Players)
	}
	return match, nil
}

// DeleteGame will delete the game with the given ID from the database
func DeleteGame(id int) (*models.Game, error) {
	// TODO
	return nil, nil
}