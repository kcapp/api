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
			player_id,
			MAX(matches_played) AS matches_played,
			MAX(matches_won) AS matches_won,
			MAX(legs_played) AS legs_played,
			MAX(legs_won) AS legs_won
		FROM (
			SELECT
				p2l.player_id,
				COUNT(DISTINCT p2l.match_id) AS matches_played,
				SUM(CASE WHEN p2l.player_id = m.winner_id THEN 1 ELSE 0 END) AS matches_won,
				COUNT(p2l.leg_id) AS legs_played,
				SUM(CASE WHEN p2l.player_id = m.winner_id THEN 1 ELSE 0 END) AS legs_won
			FROM player2leg p2l
				JOIN matches m ON m.id = p2l.match_id
				JOIN leg l ON l.id = p2l.leg_id AND l.match_id = m.id
			WHERE l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
			GROUP BY p2l.player_id
			UNION ALL
			SELECT
				m.winner_id AS player_id,
				0 AS matches_played,
				COUNT(DISTINCT m.id) AS matches_won,
				0 AS legs_played,
				0 AS legs_won
			FROM matches m
				JOIN leg l ON l.match_id = m.id
			WHERE l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_walkover = 0
			GROUP BY m.winner_id
		) AS subquery
		WHERE player_id IS NOT NULL
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
