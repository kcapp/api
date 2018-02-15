package data

import (
	"github.com/kcapp/api/models"
)

// GetPlayers returns a map of all players
func GetPlayers() (map[int]*models.Player, error) {
	rows, err := models.DB.Query(`SELECT p.id, p.name, p.nickname, p.games_played, p.games_won, p.created_at FROM player p`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.Name, &p.Nickname, &p.GamesPlayed, &p.GamesWon, &p.CreatedAt)
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

// GetPlayer returns the player for the given ID
func GetPlayer(id int) (*models.Player, error) {
	p := new(models.Player)
	err := models.DB.QueryRow(`
		SELECT p.id, p.name, p.nickname, p.games_played, p.games_won, p.created_at
		FROM player p WHERE p.id = ?`, id).Scan(&p.ID, &p.Name, &p.Nickname, &p.GamesPlayed, &p.GamesWon, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// AddPlayer will add a new player to the database
func AddPlayer(player models.Player) error {
	// Prepare statement for inserting data
	stmt, err := models.DB.Prepare("INSERT INTO player (name, nickname) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(player.Name, player.Nickname)
	return err
}

// GetPlayerScore will get the score for the given player in the given match
func GetPlayerScore(playerID int, matchID int) (int, error) {
	var score int
	err := models.DB.QueryRow(`
		SELECT 
			m.starting_score - IFNULL(
				SUM(first_dart * first_dart_multiplier) + 
				SUM(second_dart * second_dart_multiplier) + 
				SUM(third_dart * third_dart_multiplier), 0) AS 'current_score'
		FROM player2match p2m 
		LEFT JOIN `+"`match`"+` m ON m.id = p2m.match_id
		LEFT JOIN score s ON s.match_id = p2m.match_id AND s.player_id = p2m.player_id
		WHERE p2m.match_id = ? AND p2m.player_id = ? AND (s.is_bust IS NULL OR is_bust = 0)
		GROUP BY p2m.player_id`, matchID, playerID).Scan(&score)
	if err != nil {
		return 0, err
	}
	return score, nil
}
