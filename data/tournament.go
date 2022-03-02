package data

import (
	"database/sql"
	"log"

	"github.com/guregu/null"
	"github.com/kcapp/api/models"
	"github.com/kcapp/api/util"
)

// GetTournaments will return all tournaments
func GetTournaments() ([]*models.Tournament, error) {
	rows, err := models.DB.Query(`
		SELECT
			id, name, short_name, is_finished, is_playoffs, playoffs_tournament_id, office_id, start_time, end_time
		FROM tournament
		WHERE is_playoffs = 0
		ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tournaments := make([]*models.Tournament, 0)
	for rows.Next() {
		tournament := new(models.Tournament)
		err := rows.Scan(&tournament.ID, &tournament.Name, &tournament.ShortName, &tournament.IsFinished, &tournament.IsPlayoffs,
			&tournament.PlayoffsTournamentID, &tournament.OfficeID, &tournament.StartTime, &tournament.EndTime)
		if err != nil {
			return nil, err
		}
		tournaments = append(tournaments, tournament)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tournaments, nil
}

// GetTournamentGroups will return all tournament groups
func GetTournamentGroups() (map[int]*models.TournamentGroup, error) {
	rows, err := models.DB.Query("SELECT id, name, division FROM tournament_group")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make(map[int]*models.TournamentGroup)
	for rows.Next() {
		group := new(models.TournamentGroup)
		err := rows.Scan(&group.ID, &group.Name, &group.Division)
		if err != nil {
			return nil, err
		}
		groups[group.ID] = group
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// GetTournament will return a given tournament
func GetTournament(id int) (*models.Tournament, error) {
	tournament := new(models.Tournament)
	err := models.DB.QueryRow(`
		SELECT
			id, name, short_name, is_finished, is_playoffs, playoffs_tournament_id, office_id, start_time, end_time
		FROM tournament t WHERE t.id = ?`, id).Scan(&tournament.ID, &tournament.Name, &tournament.ShortName, &tournament.IsFinished, &tournament.IsPlayoffs,
		&tournament.PlayoffsTournamentID, &tournament.OfficeID, &tournament.StartTime, &tournament.EndTime)
	if err != nil {
		return nil, err
	}
	if tournament.PlayoffsTournamentID.Valid {
		playoffs, err := GetTournament(int(tournament.PlayoffsTournamentID.Int64))
		if err != nil {
			return nil, err
		}
		tournament.PlayoffsTournament = playoffs
	}
	if tournament.IsFinished {
		rows, err := models.DB.Query(`
			SELECT
				t.id, t.name, p.id, CONCAT(p.first_name, ' ', p.last_name), ts.rank, ts.elo
			FROM tournament_standings ts
				JOIN player p ON p.id = ts.player_id
				JOIN tournament t ON t.id = ts.tournament_id
			WHERE ts.tournament_id = ?
			ORDER BY ts.rank`, id)
		if err != nil {
			return nil, err
		}

		standings := make([]*models.TournamentStanding, 0)
		for rows.Next() {
			ts := new(models.TournamentStanding)
			err := rows.Scan(&ts.TournamentID, &ts.TournamentName, &ts.PlayerID, &ts.PlayerName, &ts.Rank, &ts.Elo)
			if err != nil {
				return nil, err
			}
			standings = append(standings, ts)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
		tournament.Standings = standings
	}

	return tournament, nil
}

// GetCurrentTournament will return the current active tournament
func GetCurrentTournament() (*models.Tournament, error) {
	var tournamentID int
	err := models.DB.QueryRow("SELECT id FROM tournament t WHERE t.is_finished = 0 ORDER BY start_time LIMIT 1").Scan(&tournamentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return GetTournament(tournamentID)
}

// GetCurrentTournamentForOffice will return the current active tournament for the given office
func GetCurrentTournamentForOffice(officeID int) (*models.Tournament, error) {
	var tournamentID int
	err := models.DB.QueryRow("SELECT id FROM tournament t WHERE t.office_id = ? AND t.is_finished = 0 ORDER BY start_time LIMIT 1",
		officeID).Scan(&tournamentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return GetTournament(tournamentID)
}

// GetTournamentMatches will return all matches for the given tournament
func GetTournamentMatches(id int) (map[int][]*models.Match, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id, m.is_finished, m.current_leg_id, m.winner_id, m.is_walkover, m.created_at, m.updated_at, m.owe_type_id, m.venue_id,
			mt.id, mt.name, mt.description, mm.id, mm.name, mm.short_name, mm.wins_required, mm.legs_required,
			v.id, v.name, v.description, m.updated_at as 'last_throw', GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order) AS 'players',
			m.tournament_id, tg.id, GROUP_CONCAT(legs.winner_id ORDER BY legs.id) AS 'legs_won', ot.item
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
		WHERE t.id = ?
		GROUP BY m.id
		ORDER BY m.id DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make(map[int][]*models.Match)
	for rows.Next() {
		var groupID int
		m := new(models.Match)
		m.MatchType = new(models.MatchType)
		m.MatchMode = new(models.MatchMode)
		venue := new(models.Venue)
		var players string
		var legsWon null.String
		var ot null.String
		err := rows.Scan(&m.ID, &m.IsFinished, &m.CurrentLegID, &m.WinnerID, &m.IsWalkover, &m.CreatedAt, &m.UpdatedAt, &m.OweTypeID, &m.VenueID,
			&m.MatchType.ID, &m.MatchType.Name, &m.MatchType.Description,
			&m.MatchMode.ID, &m.MatchMode.Name, &m.MatchMode.ShortName, &m.MatchMode.WinsRequired, &m.MatchMode.LegsRequired,
			&venue.ID, &venue.Name, &venue.Description, &m.LastThrow, &players, &m.TournamentID, &groupID, &legsWon, &ot)
		if err != nil {
			return nil, err
		}
		if m.VenueID.Valid {
			m.Venue = venue
		}
		m.Players = util.StringToIntArray(players)
		if legsWon.Valid {
			m.LegsWon = util.StringToIntArray(legsWon.String)
		}

		if _, ok := matches[groupID]; !ok {
			matches[groupID] = make([]*models.Match, 0)
		}
		matches[groupID] = append(matches[groupID], m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return matches, nil
}

// GetTournamentOverview will return an overview for a given tournament
func GetTournamentOverview(id int) (map[int][]*models.TournamentOverview, error) {
	rows, err := models.DB.Query(`
		SELECT
			t.id, t.name, t.short_name, t.start_time, t.end_time,
			tg.id, tg.name, tg.division,
			p.id AS 'player_id',
			p2t.is_promoted, p2t.is_relegated, p2t.is_winner, p2t.manual_order,
			COUNT(DISTINCT finished.id) AS 'p',
			COUNT(DISTINCT won.id) AS 'w',
			COUNT(DISTINCT draw.id) AS 'd',
			COUNT(DISTINCT lost.id) AS 'l',
			COUNT(DISTINCT legs_for.id) AS 'F',
			COUNT(DISTINCT legs_against.id) AS 'A',
			(COUNT(DISTINCT legs_for.id) - COUNT(DISTINCT legs_against.id)) AS 'diff',
			COUNT(DISTINCT won.id) * 2 + COUNT(DISTINCT draw.id) AS 'pts',
			IFNULL(SUM(s.ppd_score) / SUM(s.darts_thrown), -1) AS 'ppd',
			IFNULL(SUM(s.first_nine_ppd) / COUNT(finished.id), -1) AS 'first_nine_ppd',
			IFNULL(SUM(s.ppd_score) / SUM(s.darts_thrown) * 3, -1) AS 'three_dart_avg',
			IFNULL(SUM(s.first_nine_ppd) / COUNT(finished.id) * 3, -1) AS 'first_nine_three_dart_avg',
			IFNULL(SUM(60s_plus), 0) AS '60s_plus',
			IFNULL(SUM(100s_plus), 0) AS '100s_plus',
			IFNULL(SUM(140s_plus), 0) AS '140s_plus',
			IFNULL(SUM(180s), 0) AS '180s',
			IFNULL(SUM(accuracy_20) / COUNT(accuracy_20), -1) AS 'accuracy_20s',
			IFNULL(SUM(accuracy_19) / COUNT(accuracy_19), -1) AS 'accuracy_19s',
			IFNULL(SUM(overall_accuracy) / COUNT(overall_accuracy), -1) AS 'accuracy_overall',
			IFNULL(SUM(s.checkout_attempts), -1) AS 'checkout_attempts',
			IFNULL(COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100, -1) AS 'checkout_percentage'
		FROM player2leg p2l
			JOIN matches m ON m.id = p2l.match_id
			JOIN player p ON p.id = p2l.player_id
			LEFT JOIN statistics_x01 s ON s.leg_id = p2l.leg_id AND s.player_id = p.id
			LEFT JOIN matches won ON won.id = p2l.match_id AND won.winner_id = p.id
			LEFT JOIN matches lost ON lost.id = p2l.match_id AND lost.winner_id <> p.id
			LEFT JOIN matches draw ON draw.id = p2l.match_id AND draw.is_finished AND draw.winner_id IS NULL
			LEFT JOIN leg legs_for ON legs_for.id = p2l.leg_id AND legs_for.winner_id = p.id
			LEFT JOIN leg legs_against ON legs_against.id = p2l.leg_id AND legs_against.winner_id <> p.id
			LEFT JOIN matches finished ON m.id = finished.id AND finished.is_finished = 1
			JOIN tournament t ON t.id = m.tournament_id
			JOIN player2tournament p2t ON p2t.player_id = p.id AND p2t.tournament_id = t.id
			JOIN tournament_group tg ON tg.id = p2t.tournament_group_id
		WHERE m.tournament_id = ? AND m.match_type_id = 1
		GROUP BY p2l.player_id, tg.id
		ORDER BY tg.division, pts DESC, diff DESC, is_relegated, manual_order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	statistics := make(map[int][]*models.TournamentOverview)
	for rows.Next() {
		tournament := new(models.Tournament)
		group := new(models.TournamentGroup)
		stats := new(models.TournamentOverview)
		err := rows.Scan(&tournament.ID, &tournament.Name, &tournament.ShortName, &tournament.StartTime, &tournament.EndTime, &group.ID,
			&group.Name, &group.Division, &stats.PlayerID, &stats.IsPromoted, &stats.IsRelegated, &stats.IsWinner, &stats.ManualOrder, &stats.Played, &stats.MatchesWon,
			&stats.MatchesDraw, &stats.MatchesLost, &stats.LegsFor, &stats.LegsAgainst, &stats.LegsDifference, &stats.Points, &stats.PPD,
			&stats.FirstNinePPD, &stats.ThreeDartAvg, &stats.FirstNineThreeDartAvg, &stats.Score60sPlus, &stats.Score100sPlus, &stats.Score140sPlus,
			&stats.Score180s, &stats.Accuracy20, &stats.Accuracy19, &stats.AccuracyOverall, &stats.CheckoutAttempts, &stats.CheckoutPercentage)
		if err != nil {
			return nil, err
		}
		stats.Tournament = tournament
		stats.Group = group

		if _, ok := statistics[group.ID]; !ok {
			statistics[group.ID] = make([]*models.TournamentOverview, 0)
		}
		statistics[group.ID] = append(statistics[group.ID], stats)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return statistics, nil
}

// GetTournamentStatistics will return statistics for the given tournament
func GetTournamentStatistics(tournamentID int) (*models.TournamentStatistics, error) {
	statistics := new(models.TournamentStatistics)
	checkouts, err := getHighestCheckoutsForTournament(tournamentID)
	if err != nil {
		return nil, err
	}
	statistics.HighestCheckout = checkouts

	a, err := getTournamentBestStatistics(tournamentID)
	if err != nil {
		return nil, err
	}
	for _, val := range a {
		statistics.BestThreeDartAvg = append(statistics.BestThreeDartAvg, val.BestThreeDartAvg)
		statistics.BestFirstNineAvg = append(statistics.BestFirstNineAvg, val.BestFirstNineAvg)
		if val.Best301 != nil {
			statistics.Best301DartsThrown = append(statistics.Best301DartsThrown, val.Best301)
		}
		if val.Best501 != nil {
			statistics.Best501DartsThrown = append(statistics.Best501DartsThrown, val.Best501)
		}
		if val.Best701 != nil {
			statistics.Best701DartsThrown = append(statistics.Best701DartsThrown, val.Best701)
		}
	}

	return statistics, nil
}

// GetNextTournamentMatch will return the next tournament match
func GetNextTournamentMatch(matchID int) (*models.Match, error) {
	var nextMatchID null.Int
	err := models.DB.QueryRow(`
		SELECT match_id FROM match_metadata mm
            JOIN matches m ON mm.match_id = m.id
		WHERE (order_of_play = (SELECT order_of_play FROM match_metadata mm WHERE match_id = ?) + 1)
            AND m.tournament_id = (SELECT tournament_id FROM matches where id = ?)`, matchID, matchID).Scan(&nextMatchID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return GetMatch(int(nextMatchID.Int64))
}

// GetTournamentStandings will return statistics for the given tournament
func GetTournamentStandings() ([]*models.TournamentStanding, error) {
	rows, err := models.DB.Query(`
		SELECT player_id, first_name, tournament_elo, tournament_elo_matches, current_elo, current_elo_matches,
			@curRank := @curRank + 1 AS "rank" FROM (
				SELECT
					pe.player_id,
					p.first_name,
					pe.tournament_elo,
					pe.tournament_elo_matches,
					pe.current_elo,
					pe.current_elo_matches
				FROM player_elo pe
				JOIN player p ON p.id = pe.player_id
				WHERE (pe.current_elo_matches > 5 OR pe.tournament_elo_matches > 0)
					AND p.active = 1
				ORDER BY tournament_elo DESC
		) elo, (SELECT @curRank := 0) r`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	standings := make([]*models.TournamentStanding, 0)
	for rows.Next() {
		standing := new(models.TournamentStanding)
		err := rows.Scan(&standing.PlayerID, &standing.PlayerName, &standing.Elo, &standing.EloPlayed, &standing.CurrentElo,
			&standing.CurrentEloPlayed, &standing.Rank)
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

// getHighestCheckout will calculate the highest checkout for the given players
func getHighestCheckoutsForTournament(tournamentID int) ([]*models.BestStatistic, error) {
	rows, err := models.DB.Query(`
		SELECT
			player_id,
			leg_id,
			MAX(checkout) AS 'checkout'
			FROM (SELECT
				s.player_id,
				s.leg_id,
				IFNULL(s.first_dart * s.first_dart_multiplier, 0) +
					IFNULL(s.second_dart * s.second_dart_multiplier, 0) +
					IFNULL(s.third_dart * s.third_dart_multiplier, 0) AS 'checkout'
			FROM score s
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m on l.match_id = m.id
			WHERE l.winner_id = s.player_id
				AND s.leg_id IN (SELECT id FROM leg WHERE match_id IN (SELECT id FROM matches WHERE tournament_id = ?))
				AND s.id IN (SELECT MAX(s.id) FROM score s JOIN leg l ON l.id = s.leg_id WHERE l.winner_id = s.player_id GROUP BY leg_id)
				AND IFNULL(l.leg_type_id, m.match_type_id) = 1 -- X01
			GROUP BY s.player_id, s.id
			ORDER BY checkout DESC) checkouts
			GROUP BY player_id
		ORDER BY checkout DESC`, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	best := make([]*models.BestStatistic, 0)
	for rows.Next() {
		highest := new(models.BestStatistic)
		err := rows.Scan(&highest.PlayerID, &highest.LegID, &highest.Value)
		if err != nil {
			return nil, err
		}
		best = append(best, highest)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return best, nil
}

// getBestStatistics will calculate Best PPD, Best First 9, Best 301 and Best 501 for the given players
func getTournamentBestStatistics(tournamentID int) ([]*models.StatisticsX01, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			l.id AS 'leg_id',
			l.winner_id,
			s.ppd_score * 3 / s.darts_thrown as 'three_dart_avg',
			s.first_nine_ppd_score * 3 / 9 as 'first_nine_avg',
			s.checkout_percentage,
			s.darts_thrown,
			l.starting_score
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
		WHERE s.leg_id IN (SELECT id FROM leg WHERE match_id IN (SELECT id FROM matches WHERE tournament_id = ?))`, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsX01, 0)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.LegID, &s.WinnerID, &s.ThreeDartAvg, &s.FirstNineThreeDartAvg, &s.CheckoutPercentage, &s.DartsThrown, &s.StartingScore)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	bestStatistics := make(map[int]*models.StatisticsX01)
	for _, stat := range stats {
		best := bestStatistics[stat.PlayerID]
		if best == nil {
			best = new(models.StatisticsX01)
			best.PlayerID = stat.PlayerID
			bestStatistics[stat.PlayerID] = best
		}

		if stat.PlayerID == stat.WinnerID {
			if stat.StartingScore.Int64 == 301 {
				if best.Best301 == nil {
					best.Best301 = new(models.BestStatistic)
				}
				if stat.DartsThrown < best.Best301.Value || best.Best301.Value == 0 {
					best.Best301.Value = stat.DartsThrown
					best.Best301.LegID = stat.LegID
					best.Best301.PlayerID = stat.PlayerID
				}
			}
			if stat.StartingScore.Int64 == 501 {
				if best.Best501 == nil {
					best.Best501 = new(models.BestStatistic)
				}
				if stat.DartsThrown < best.Best501.Value || best.Best501.Value == 0 {
					best.Best501.Value = stat.DartsThrown
					best.Best501.LegID = stat.LegID
					best.Best501.PlayerID = stat.PlayerID
				}
			}
			if stat.StartingScore.Int64 == 701 {
				if best.Best701 == nil {
					best.Best701 = new(models.BestStatistic)
				}
				if stat.DartsThrown < best.Best701.Value || best.Best701.Value == 0 {
					best.Best701.Value = stat.DartsThrown
					best.Best701.LegID = stat.LegID
					best.Best701.PlayerID = stat.PlayerID
				}
			}
		}
		if best.BestThreeDartAvg == nil {
			best.BestThreeDartAvg = new(models.BestStatisticFloat)
		}
		if stat.ThreeDartAvg > best.BestThreeDartAvg.Value {
			best.BestThreeDartAvg.Value = stat.ThreeDartAvg
			best.BestThreeDartAvg.LegID = stat.LegID
			best.BestThreeDartAvg.PlayerID = stat.PlayerID
		}
		if best.BestFirstNineAvg == nil {
			best.BestFirstNineAvg = new(models.BestStatisticFloat)
		}
		if stat.FirstNineThreeDartAvg > best.BestFirstNineAvg.Value {
			best.BestFirstNineAvg.Value = stat.FirstNineThreeDartAvg
			best.BestFirstNineAvg.LegID = stat.LegID
			best.BestFirstNineAvg.PlayerID = stat.PlayerID
		}
	}

	s := make([]*models.StatisticsX01, 0)
	for _, val := range bestStatistics {
		s = append(s, val)
	}

	return s, nil
}

// NewTournament will create a new tournament
func NewTournament(tournament models.Tournament) (*models.Tournament, error) {
	tx, err := models.DB.Begin()
	if err != nil {
		return nil, err
	}

	res, err := tx.Exec(`
		INSERT INTO tournament (name, short_name, is_finished, is_playoffs, playoffs_tournament_id, office_id, start_time, end_time) VALUES
		(?, ?, ?, ?, ?, ?, ?, ?)`, tournament.Name, tournament.ShortName, 0, tournament.IsPlayoffs, tournament.PlayoffsTournamentID,
		tournament.OfficeID, tournament.StartTime, tournament.EndTime)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tournamentID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, player := range tournament.Players {
		_, err = tx.Exec(`INSERT INTO player2tournament (player_id, tournament_id, tournament_group_id) VALUES (?, ?, ?)`,
			player.PlayerID, tournamentID, player.TournamentGroupID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	log.Printf("Created new tournament %d", tournamentID)
	return GetTournament(int(tournamentID))
}

// GetTournamentMatchesForPlayer will return all tournament matches for the given player and tournament
func GetTournamentMatchesForPlayer(tournamentID int, playerID int) ([]*models.Match, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id, m.is_finished, m.is_abandoned, m.is_walkover, m.current_leg_id, m.winner_id, m.office_id, m.is_practice,
			m.created_at, m.updated_at, m.owe_type_id, m.venue_id, mt.id, mt.name, mt.description, mm.id, mm.name, mm.short_name,
			mm.wins_required, mm.legs_required, ot.id, ot.item, v.id, v.name, v.description, l.updated_at as 'last_throw',
			GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order) AS 'players', m.tournament_id, m.tournament_id, t.id, t.name,
			tg.id, tg.name, GROUP_CONCAT(legs.winner_id ORDER BY legs.id) AS 'legs_won'
		FROM matches m
			JOIN match_type mt ON mt.id = m.match_type_id
			JOIN match_mode mm ON mm.id = m.match_mode_id
			LEFT JOIN leg l ON l.id = m.current_leg_id
			LEFT JOIN owe_type ot ON ot.id = m.owe_type_id
			LEFT JOIN venue v on v.id = m.venue_id
			LEFT JOIN player2leg p2l ON p2l.match_id = m.id
			LEFT JOIN player2leg p2l2 ON p2l2.leg_id = l.id
			LEFT JOIN leg legs ON legs.id = p2l.leg_id AND legs.winner_id = p2l.player_id
			LEFT JOIN player2tournament p2t ON p2t.tournament_id = m.tournament_id AND p2t.player_id = p2l.player_id
			LEFT JOIN tournament t ON t.id = p2t.tournament_id
			LEFT JOIN tournament_group tg ON tg.id = p2t.tournament_group_id
		WHERE p2l2.player_id = ?
			AND m.tournament_id = ?
		GROUP BY m.id
		ORDER BY m.created_at`, playerID, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		m := new(models.Match)
		m.Tournament = new(models.MatchTournament)
		m.MatchType = new(models.MatchType)
		m.MatchMode = new(models.MatchMode)
		ot := new(models.OweType)
		venue := new(models.Venue)
		var players string
		var legsWon null.String
		err := rows.Scan(&m.ID, &m.IsFinished, &m.IsAbandoned, &m.IsWalkover, &m.CurrentLegID, &m.WinnerID, &m.OfficeID, &m.IsPractice, &m.CreatedAt, &m.UpdatedAt,
			&m.OweTypeID, &m.VenueID, &m.MatchType.ID, &m.MatchType.Name, &m.MatchType.Description,
			&m.MatchMode.ID, &m.MatchMode.Name, &m.MatchMode.ShortName, &m.MatchMode.WinsRequired, &m.MatchMode.LegsRequired,
			&ot.ID, &ot.Item, &venue.ID, &venue.Name, &venue.Description, &m.LastThrow, &players, &m.TournamentID, &m.TournamentID, &m.Tournament.TournamentID,
			&m.Tournament.TournamentName, &m.Tournament.TournamentGroupID, &m.Tournament.TournamentGroupName, &legsWon)
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
