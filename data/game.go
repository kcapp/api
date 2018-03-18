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
	res, err := tx.Exec("INSERT INTO game (game_type_id, game_mode_id, owe_type_id, created_at) VALUES (?, ?, ?, NOW())", game.GameType.ID, game.GameMode.ID, game.OweTypeID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	gameID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	res, err = tx.Exec("INSERT INTO `match` (starting_score, current_player_id, game_id, created_at) VALUES (?, ?, ?, NOW()) ", game.Matches[0].StartingScore, game.Players[0], gameID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	matchID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Exec("UPDATE game SET current_match_id = ? WHERE id = ?", matchID, gameID)
	for idx, playerID := range game.Players {
		order := idx + 1
		res, err = tx.Exec("INSERT INTO player2match (player_id, match_id, `order`, game_id) VALUES (?, ?, ?, ?)", playerID, matchID, order, gameID)
		if err != nil {
			tx.Rollback()
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
			gt.id, gt.name, gt.description,
			gm.id, gm.name, gm.short_name, gm.wins_required, gm.matches_required,
			ot.id, ot.item,
			GROUP_CONCAT(DISTINCT p2m.player_id ORDER BY p2m.order) AS 'players'
		FROM game g
		JOIN game_type gt ON gt.id = g.game_type_id
		JOIN game_mode gm ON gm.id = g.game_mode_id
		LEFT JOIN owe_type ot ON ot.id = g.owe_type_id
		LEFT JOIN player2match p2m ON p2m.game_id = g.id
		GROUP BY g.id
		ORDER BY g.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]*models.Game, 0)
	for rows.Next() {
		g := new(models.Game)
		g.GameType = new(models.GameType)
		g.GameMode = new(models.GameMode)
		ot := new(models.OweType)
		var players string
		err := rows.Scan(&g.ID, &g.IsFinished, &g.CurrentMatchID, &g.WinnerID, &g.CreatedAt, &g.UpdatedAt, &g.OweTypeID,
			&g.GameType.ID, &g.GameType.Name, &g.GameType.Description,
			&g.GameMode.ID, &g.GameMode.Name, &g.GameMode.ShortName, &g.GameMode.WinsRequired, &g.GameMode.MatchesRequired,
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
	g.GameMode = new(models.GameMode)
	ot := new(models.OweType)
	var players string
	err := models.DB.QueryRow(`
        SELECT
			g.id, g.is_finished, g.current_match_id, g.winner_id, g.created_at, g.updated_at, g.owe_type_id,
			gt.id, gt.name, gt.description,
			gm.id, gm.name, gm.short_name, gm.wins_required, gm.matches_required,
			ot.id, ot.item,
			GROUP_CONCAT(DISTINCT p2m.player_id ORDER BY p2m.order) AS 'players'
		FROM game g
		JOIN game_type gt ON gt.id = g.game_type_id
		JOIN game_mode gm ON gm.id = g.game_mode_id
		LEFT JOIN owe_type ot ON ot.id = g.owe_type_id
		LEFT JOIN player2match p2m ON p2m.game_id = g.id
		WHERE g.id = ?`, id).Scan(&g.ID, &g.IsFinished, &g.CurrentMatchID, &g.WinnerID, &g.CreatedAt, &g.UpdatedAt, &g.OweTypeID,
		&g.GameType.ID, &g.GameType.Name, &g.GameType.Description, &g.GameMode.ID, &g.GameMode.Name, &g.GameMode.ShortName,
		&g.GameMode.WinsRequired, &g.GameMode.MatchesRequired, &ot.ID, &ot.Item, &players)
	if err != nil {
		return nil, err
	}
	if g.OweTypeID.Valid {
		g.OweType = ot
	}
	g.Players = util.StringToIntArray(players)
	g.Matches, err = GetMatchesForGame(id)
	if err != nil {
		return nil, err
	}
	if g.IsFinished {
		g.EndTime = g.Matches[len(g.Matches)-1].Endtime.String
	}
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

// GetGameModes will return all game modes
func GetGameModes() ([]*models.GameMode, error) {
	rows, err := models.DB.Query("SELECT id, wins_required, matches_required, `name`, short_name FROM game_mode ORDER BY wins_required")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	modes := make([]*models.GameMode, 0)
	for rows.Next() {
		gm := new(models.GameMode)
		err := rows.Scan(&gm.ID, &gm.WinsRequired, &gm.MatchesRequired, &gm.Name, &gm.ShortName)
		if err != nil {
			return nil, err
		}
		modes = append(modes, gm)
	}

	return modes, nil
}

// GetGameTypes will return all game types
func GetGameTypes() ([]*models.GameType, error) {
	rows, err := models.DB.Query("SELECT id, `name`, description FROM game_type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	types := make([]*models.GameType, 0)
	for rows.Next() {
		gt := new(models.GameType)
		err := rows.Scan(&gt.ID, &gt.Name, &gt.Description)
		if err != nil {
			return nil, err
		}
		types = append(types, gt)
	}

	return types, nil
}

// GetWinsPerPlayer gets the number of wins per player for the given game
func GetWinsPerPlayer(id int) (map[int]int, error) {
	rows, err := models.DB.Query(`
		SELECT IFNULL(m.winner_id, 0), COUNT(m.winner_id) AS 'wins' FROM `+"`match`"+` m
		LEFT JOIN game g ON g.id = m.game_id
		WHERE m.game_id = ? GROUP BY m.winner_id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	winsMap := make(map[int]int)
	for rows.Next() {
		var playerID int
		var wins int
		err := rows.Scan(&playerID, &wins)
		if err != nil {
			return nil, err
		}
		winsMap[playerID] = wins
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return winsMap, nil
}
