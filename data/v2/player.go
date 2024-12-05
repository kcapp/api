package data_v2

import (
	"github.com/kcapp/api/models"
)

// GetPlayers returns a map of all players
func GetPlayers() ([]*models.Player, error) {
	played, err := GetMatchesPlayedPerPlayer()
	if err != nil {
		return nil, err
	}

	rows, err := models.DB.Query(`
		SELECT
			p.id, p.first_name, p.last_name, p.vocal_name, p.nickname, p.slack_handle, p.color, p.profile_pic_url, p.smartcard_uid,
			 p.board_stream_url, p.board_stream_css, p.active, p.office_id, p.is_bot, p.is_placeholder, p.created_at, p.updated_at
		FROM player p`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := []*models.Player{}
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.VocalName, &p.Nickname, &p.SlackHandle, &p.Color, &p.ProfilePicURL,
			&p.SmartcardUID, &p.BoardStreamURL, &p.BoardStreamCSS, &p.IsActive, &p.OfficeID, &p.IsBot, &p.IsPlaceholder, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		players = append(players, p)

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
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT l.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won'
		FROM leg l
			JOIN player2leg p2l on p2l.leg_id = l.id
			JOIN player p ON p.id = p2l.player_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = l.id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l2.match_id AND m2.winner_id = p.id
		WHERE l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0 AND m.is_bye = 0
		GROUP by p.id`)
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
