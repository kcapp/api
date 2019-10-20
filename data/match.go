package data

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/guregu/null"
	"github.com/kcapp/api/models"
	"github.com/kcapp/api/util"
)

// NewMatch will insert a new match in the database
func NewMatch(match models.Match) (*models.Match, error) {
	tx, err := models.DB.Begin()
	if err != nil {
		return nil, err
	}
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	if match.CreatedAt != "" {
		createdAt = match.CreatedAt
	}
	res, err := tx.Exec("INSERT INTO matches (match_type_id, match_mode_id, owe_type_id, venue_id, office_id, is_practice, tournament_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		match.MatchType.ID, match.MatchMode.ID, match.OweTypeID, match.VenueID, match.OfficeID, match.IsPractice, match.TournamentID, createdAt)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	matchID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	res, err = tx.Exec("INSERT INTO leg (starting_score, current_player_id, match_id, created_at) VALUES (?, ?, ?, NOW()) ", match.Legs[0].StartingScore, match.Players[0], matchID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	legID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Exec("UPDATE matches SET current_leg_id = ? WHERE id = ?", legID, matchID)
	for idx, playerID := range match.Players {
		order := idx + 1
		res, err = tx.Exec("INSERT INTO player2leg (player_id, leg_id, `order`, match_id, handicap) VALUES (?, ?, ?, ?, ?)", playerID, legID, order, matchID, match.PlayerHandicaps[playerID])
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if config, ok := match.BotPlayerConfig[playerID]; ok {
			player2LegID, err := res.LastInsertId()
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			res, err = tx.Exec("INSERT INTO bot2player2leg (player2leg_id, player_id, skill_level) VALUES (?, ?, ?)", player2LegID, config.PlayerID, config.Skill)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}

	}
	tx.Commit()
	log.Printf("Started new match %d", matchID)
	return GetMatch(int(matchID))
}

// GetMatches returns all matches
func GetMatches() ([]*models.Match, error) {
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
		GROUP BY m.id
		ORDER BY m.id DESC`)
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
			&ot.ID, &ot.Item, &venue.ID, &venue.Name, &venue.Description, &m.LastThrow, &players)
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

// GetActiveMatches returns all active matches
func GetActiveMatches() ([]*models.Match, error) {
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
		WHERE m.is_finished = 0
			AND l.updated_at > NOW() - INTERVAL 2 MINUTE
		GROUP BY m.id
		ORDER BY m.id DESC`)
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
		err := rows.Scan(&m.ID, &m.IsFinished, &m.IsAbandoned, &m.IsWalkover, &m.CurrentLegID, &m.WinnerID, &m.OfficeID, &m.IsPractice,
			&m.CreatedAt, &m.UpdatedAt, &m.OweTypeID, &m.VenueID, &m.MatchType.ID, &m.MatchType.Name, &m.MatchType.Description,
			&m.MatchMode.ID, &m.MatchMode.Name, &m.MatchMode.ShortName, &m.MatchMode.WinsRequired, &m.MatchMode.LegsRequired,
			&ot.ID, &ot.Item, &venue.ID, &venue.Name, &venue.Description, &m.LastThrow, &players)
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

// GetMatchesLimit returns the N matches from the given starting point
func GetMatchesLimit(start int, limit int) ([]*models.Match, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id, m.is_finished, m.is_abandoned, m.is_walkover, m.current_leg_id, m.winner_id, m.office_id, m.is_practice,
			m.created_at, m.updated_at, m.owe_type_id, m.venue_id, mt.id, mt.name, mt.description, mm.id, mm.name, mm.short_name,
			mm.wins_required, mm.legs_required, ot.id, ot.item, v.id, v.name, v.description,
			l.updated_at as 'last_throw', GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order) AS 'players',
			m.tournament_id, t.id, t.name, tg.id, tg.name, GROUP_CONCAT(legs.winner_id ORDER BY legs.id) AS 'legs_won'
		FROM matches m
			JOIN match_type mt ON mt.id = m.match_type_id
			JOIN match_mode mm ON mm.id = m.match_mode_id
			LEFT JOIN leg l ON l.id = m.current_leg_id
			LEFT JOIN owe_type ot ON ot.id = m.owe_type_id
			LEFT JOIN venue v on v.id = m.venue_id
			LEFT JOIN player2leg p2l ON p2l.match_id = m.id
			LEFT JOIN leg legs ON legs.id = p2l.leg_id AND legs.winner_id = p2l.player_id
			LEFT JOIN player2tournament p2t ON p2t.tournament_id = m.tournament_id AND p2t.player_id = p2l.player_id
			LEFT JOIN tournament t ON t.id = p2t.tournament_id
			LEFT JOIN tournament_group tg ON tg.id = p2t.tournament_group_id
		WHERE m.created_at  <= NOW()
		GROUP BY m.id
		ORDER BY m.created_at DESC
		LIMIT ?, ?`, start, limit)
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
		tournament := new(models.MatchTournament)
		var players string
		var legsWon null.String
		err := rows.Scan(&m.ID, &m.IsFinished, &m.IsAbandoned, &m.IsWalkover, &m.CurrentLegID, &m.WinnerID, &m.OfficeID, &m.IsPractice,
			&m.CreatedAt, &m.UpdatedAt, &m.OweTypeID, &m.VenueID, &m.MatchType.ID, &m.MatchType.Name, &m.MatchType.Description,
			&m.MatchMode.ID, &m.MatchMode.Name, &m.MatchMode.ShortName, &m.MatchMode.WinsRequired, &m.MatchMode.LegsRequired,
			&ot.ID, &ot.Item, &venue.ID, &venue.Name, &venue.Description, &m.LastThrow, &players, &m.TournamentID, &tournament.TournamentID,
			&tournament.TournamentName, &tournament.TournamentGroupID, &tournament.TournamentGroupName, &legsWon)
		if err != nil {
			return nil, err
		}
		if m.OweTypeID.Valid {
			m.OweType = ot
		}
		if m.VenueID.Valid {
			m.Venue = venue
		}
		if m.TournamentID.Valid {
			m.Tournament = tournament
		}
		m.Players = util.StringToIntArray(players)
		if legsWon.Valid {
			m.LegsWon = util.StringToIntArray(legsWon.String)
		}
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
	m.MatchType = new(models.MatchType)
	m.MatchMode = new(models.MatchMode)
	ot := new(models.OweType)
	venue := new(models.Venue)
	tournament := new(models.MatchTournament)
	var players string
	err := models.DB.QueryRow(`
        SELECT
			m.id, m.is_finished, m.is_abandoned, m.is_walkover, m.current_leg_id, m.winner_id, m.office_id, m.is_practice, m.created_at, m.updated_at,
			m.owe_type_id, m.venue_id, mt.id, mt.name, mt.description, mm.id, mm.name, mm.short_name, mm.wins_required,
			mm.legs_required, ot.id, ot.item, v.id, v.name, v.description,
			MAX(l.updated_at) AS 'last_throw',
			MIN(s.created_at) AS 'first_throw',
			GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order) AS 'players',
			m.tournament_id, t.id, t.name, t.office_id, tg.id, tg.name
		FROM matches m
			JOIN match_type mt ON mt.id = m.match_type_id
			JOIN match_mode mm ON mm.id = m.match_mode_id
			LEFT JOIN leg l ON l.match_id = m.id
			LEFT JOIN score s ON s.leg_id = l.id
			LEFT JOIN owe_type ot ON ot.id = m.owe_type_id
			LEFT JOIN venue v on v.id = m.venue_id
			LEFT JOIN player2leg p2l ON p2l.match_id = m.id
			LEFT JOIN player2tournament p2t ON p2t.tournament_id = m.tournament_id AND p2t.player_id = p2l.player_id
			LEFT JOIN tournament t ON t.id = p2t.tournament_id
			LEFT JOIN tournament_group tg ON tg.id = p2t.tournament_group_id
		WHERE m.id = ?`, id).Scan(&m.ID, &m.IsFinished, &m.IsAbandoned, &m.IsWalkover, &m.CurrentLegID, &m.WinnerID, &m.OfficeID, &m.IsPractice,
		&m.CreatedAt, &m.UpdatedAt, &m.OweTypeID, &m.VenueID, &m.MatchType.ID, &m.MatchType.Name, &m.MatchType.Description,
		&m.MatchMode.ID, &m.MatchMode.Name, &m.MatchMode.ShortName, &m.MatchMode.WinsRequired, &m.MatchMode.LegsRequired,
		&ot.ID, &ot.Item, &venue.ID, &venue.Name, &venue.Description, &m.LastThrow, &m.FirstThrow, &players, &m.TournamentID, &tournament.TournamentID,
		&tournament.TournamentName, &tournament.OfficeID, &tournament.TournamentGroupID, &tournament.TournamentGroupName)
	if err != nil {
		return nil, err
	}
	if m.OweTypeID.Valid {
		m.OweType = ot
	}
	if m.VenueID.Valid {
		m.Venue, err = GetVenue(int(m.VenueID.Int64))
		if err != nil {
			return nil, err
		}
	}
	if m.TournamentID.Valid {
		m.Tournament = tournament
	}
	m.Players = util.StringToIntArray(players)
	m.Legs, err = GetLegsForMatch(id)
	if err != nil {
		return nil, err
	}
	if m.IsFinished {
		m.EndTime = m.Legs[len(m.Legs)-1].Endtime.String
	}

	m.EloChange, err = GetMatchEloChange(id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetMatchMetadata returns a metadata about the given match
func GetMatchMetadata(id int) (*models.MatchMetadata, error) {
	m := new(models.MatchMetadata)
	m.TournamentGroup = new(models.TournamentGroup)
	var playersStr string
	err := models.DB.QueryRow(`
		SELECT
			mm.id, mm.match_id, mm.order_of_play, mm.match_displayname, mm.elimination,
			mm.trophy, mm.promotion, mm.semi_final, mm.grand_final, mm.winner_outcome_match_id,
			mm.is_winner_outcome_home,  mm.looser_outcome_match_id, mm.is_looser_outcome_home,
			mm.winner_outcome, mm.looser_outcome,
			tg.id, tg.name, GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order) AS 'players'
		FROM match_metadata mm
			LEFT JOIN tournament_group tg ON tg.id = mm.tournament_group_id
			LEFT JOIN player2leg p2l ON p2l.match_id = mm.match_id
		WHERE mm.match_id = ?
		GROUP BY mm.match_id`, id).Scan(&m.ID, &m.MatchID, &m.OrderOfPlay, &m.MatchDisplayname, &m.Elimination,
		&m.Trophy, &m.Promotion, &m.SemiFinal, &m.GrandFinal, &m.WinnerOutcomeMatchID, &m.IsWinnerOutcomeHome,
		&m.LooserOutcomeMatchID, &m.IsLooserOutcomeHome, &m.WinnerOutcome, &m.LooserOutcome, &m.TournamentGroup.ID,
		&m.TournamentGroup.Name, &playersStr)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	players := util.StringToIntArray(playersStr)
	if len(players) == 2 {
		m.HomePlayer = players[0]
		m.AwayPlayer = players[1]
	}

	return m, nil
}

// GetMatchMetadataForTournament returns metadata for all matches in a given tournament
func GetMatchMetadataForTournament(tournamentID int) ([]*models.MatchMetadata, error) {
	rows, err := models.DB.Query(`
		SELECT
			mm.id, mm.match_id, mm.order_of_play, mm.match_displayname, mm.elimination,
			mm.trophy, mm.promotion, mm.semi_final, mm.grand_final, mm.winner_outcome_match_id,
			mm.looser_outcome_match_id, mm.winner_outcome, mm.looser_outcome, tg.id, tg.name,
			GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order) AS 'players'
		FROM match_metadata mm
			JOIN matches m on m.id = mm.match_id
			JOIN tournament_group tg ON tg.id = mm.tournament_group_id
			JOIN player2leg p2l ON p2l.match_id = mm.match_id
		WHERE m.tournament_id = ?
		GROUP BY mm.match_id`, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metadata := make([]*models.MatchMetadata, 0)
	for rows.Next() {
		m := new(models.MatchMetadata)
		m.TournamentGroup = new(models.TournamentGroup)
		var playersStr string
		err := rows.Scan(&m.ID, &m.MatchID, &m.OrderOfPlay, &m.MatchDisplayname, &m.Elimination,
			&m.Trophy, &m.Promotion, &m.SemiFinal, &m.GrandFinal, &m.WinnerOutcomeMatchID, &m.LooserOutcomeMatchID,
			&m.WinnerOutcome, &m.LooserOutcome, &m.TournamentGroup.ID, &m.TournamentGroup.Name, &playersStr)
		if err != nil {
			return nil, err
		}
		players := util.StringToIntArray(playersStr)
		m.HomePlayer = players[0]
		m.AwayPlayer = players[1]

		metadata = append(metadata, m)
	}
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

// ContinueMatch will either return the current leg or create a new leg
func ContinueMatch(id int) (*models.Leg, error) {
	match, err := GetMatch(id)
	if err != nil {
		return nil, err
	}
	if match.IsFinished {
		return nil, errors.New("Cannot continue finished match")
	}

	legs, err := GetLegsForMatch(id)
	if err != nil {
		return nil, err
	}
	leg := legs[len(legs)-1]
	if leg.IsFinished {
		return NewLeg(id, leg.StartingScore, leg.Players)
	}
	return leg, nil
}

// DeleteMatch will delete the match with the given ID from the database
func DeleteMatch(id int) (*models.Match, error) {
	// TODO
	return nil, nil
}

// GetMatchModes will return all match modes
func GetMatchModes() ([]*models.MatchMode, error) {
	rows, err := models.DB.Query("SELECT id, wins_required, legs_required, `name`, short_name FROM match_mode ORDER BY wins_required")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	modes := make([]*models.MatchMode, 0)
	for rows.Next() {
		mm := new(models.MatchMode)
		err := rows.Scan(&mm.ID, &mm.WinsRequired, &mm.LegsRequired, &mm.Name, &mm.ShortName)
		if err != nil {
			return nil, err
		}
		modes = append(modes, mm)
	}

	return modes, nil
}

// GetMatchTypes will return all match types
func GetMatchTypes() ([]*models.MatchType, error) {
	rows, err := models.DB.Query("SELECT id, `name`, description FROM match_type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	types := make([]*models.MatchType, 0)
	for rows.Next() {
		mt := new(models.MatchType)
		err := rows.Scan(&mt.ID, &mt.Name, &mt.Description)
		if err != nil {
			return nil, err
		}
		types = append(types, mt)
	}

	return types, nil
}

// GetWinsPerPlayer gets the number of wins per player for the given match
func GetWinsPerPlayer(id int) (map[int]int, error) {
	rows, err := models.DB.Query(`
		SELECT
			IFNULL(l.winner_id, 0), COUNT(l.winner_id) AS 'wins'
		FROM leg l
		WHERE l.match_id = ?
		GROUP BY l.winner_id`, id)
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

// GetHeadToHeadMatches will return the last N matches between two players
func GetHeadToHeadMatches(player1 int, player2 int) ([]*models.Match, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id, m.is_finished, m.is_abandoned, m.is_walkover, m.current_leg_id, m.winner_id, m.created_at, m.updated_at,
			m.owe_type_id, mt.id, mt.name, mt.description,
			mm.id, mm.name, mm.short_name, mm.wins_required, mm.legs_required
		FROM matches m
			JOIN match_type mt ON mt.id = m.match_type_id
			JOIN match_mode mm ON mm.id = m.match_mode_id
			JOIN player2leg p2l ON p2l.match_id = m.id
		WHERE m.id IN (SELECT match_id FROM player2leg GROUP BY match_id HAVING COUNT(DISTINCT player_id) = 2)
			AND m.is_finished = 1 AND m.is_abandoned = 0
			AND m.match_type_id = 1
			AND p2l.player_id IN (?, ?)
		GROUP BY m.id
			HAVING COUNT(DISTINCT p2l.player_id) = 2
		ORDER BY m.created_at DESC`, player1, player2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		m := new(models.Match)
		m.MatchType = new(models.MatchType)
		m.MatchMode = new(models.MatchMode)
		err := rows.Scan(&m.ID, &m.IsFinished, &m.IsAbandoned, &m.IsWalkover, &m.CurrentLegID, &m.WinnerID, &m.CreatedAt, &m.UpdatedAt,
			&m.OweTypeID, &m.MatchType.ID, &m.MatchType.Name, &m.MatchType.Description,
			&m.MatchMode.ID, &m.MatchMode.Name, &m.MatchMode.ShortName, &m.MatchMode.WinsRequired, &m.MatchMode.LegsRequired)
		if err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return matches, nil
}

// GetPlayerLastMatches will return the last N matches for the given player
func GetPlayerLastMatches(playerID int, limit int) ([]*models.Match, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id, m.is_finished, m.is_abandoned, m.is_walkover, m.current_leg_id, m.winner_id, m.created_at, m.updated_at,
			m.owe_type_id, mt.id, mt.name, mt.description,
			mm.id, mm.name, mm.short_name, mm.wins_required, mm.legs_required
		FROM matches m
			JOIN match_type mt ON mt.id = m.match_type_id
			JOIN match_mode mm ON mm.id = m.match_mode_id
			JOIN player2leg p2l ON p2l.match_id = m.id
		WHERE m.id IN (SELECT  match_id  FROM player2leg  GROUP BY match_id  HAVING COUNT(DISTINCT player_id) = 2)
			AND m.is_finished = 1 AND m.is_abandoned = 0 AND m.is_practice = 0
			AND p2l.player_id IN (?)
		GROUP BY m.id
		ORDER BY m.created_at DESC LIMIT ?`, playerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		m := new(models.Match)
		m.MatchType = new(models.MatchType)
		m.MatchMode = new(models.MatchMode)
		err := rows.Scan(&m.ID, &m.IsFinished, &m.IsAbandoned, &m.IsWalkover, &m.CurrentLegID, &m.WinnerID, &m.CreatedAt, &m.UpdatedAt, &m.OweTypeID,
			&m.MatchType.ID, &m.MatchType.Name, &m.MatchType.Description,
			&m.MatchMode.ID, &m.MatchMode.Name, &m.MatchMode.ShortName, &m.MatchMode.WinsRequired, &m.MatchMode.LegsRequired)
		if err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return matches, nil
}

// GetMatchEloChange returns Elo change for each player in the given match
func GetMatchEloChange(id int) (map[int]*models.PlayerElo, error) {
	rows, err := models.DB.Query(`
		SELECT
			player_id,
			old_elo,
			new_elo,
			new_tournament_elo,
			old_tournament_elo
		FROM player_elo_changelog
		WHERE match_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	change := make(map[int]*models.PlayerElo)
	for rows.Next() {
		elo := new(models.PlayerElo)
		err := rows.Scan(&elo.PlayerID, &elo.CurrentElo, &elo.CurrentEloNew, &elo.TournamentEloNew, &elo.TournamentElo)
		if err != nil {
			return nil, err
		}
		change[elo.PlayerID] = elo
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return change, nil
}

// SwapPlayers will swap the two players for the given match
func SwapPlayers(matchID int, newPlayerID int, oldPlayerID int) error {
	tx, err := models.DB.Begin()
	if err != nil {
		return err
	}

	// Update current player of the leg
	_, err = tx.Exec("UPDATE leg SET current_player_id = ? WHERE match_id = ?", newPlayerID, matchID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update player2leg
	_, err = tx.Exec("UPDATE player2leg SET player_id = ? WHERE match_id = ? AND player_id = ?", newPlayerID, matchID, oldPlayerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	log.Printf("Swapped player %d with %d for match %d", oldPlayerID, newPlayerID, matchID)
	return nil
}
