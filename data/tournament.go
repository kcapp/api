package data

import (
	"github.com/kcapp/api/models"
	"github.com/kcapp/api/util"
)

// GetTournaments will return all tournaments
func GetTournaments() ([]*models.Tournament, error) {
	rows, err := models.DB.Query(`
		SELECT
			id, name, short_name, is_finished, is_playoffs, playoffs_tournament_id, start_time, end_time
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
			&tournament.PlayoffsTournamentID, &tournament.StartTime, &tournament.EndTime)
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

	groups := make(map[int]*models.TournamentGroup, 0)
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
			id, name, short_name, is_finished, is_playoffs, playoffs_tournament_id, start_time, end_time
		FROM tournament t WHERE t.id = ?`, id).Scan(&tournament.ID, &tournament.Name, &tournament.ShortName, &tournament.IsFinished, &tournament.IsPlayoffs,
		&tournament.PlayoffsTournamentID, &tournament.StartTime, &tournament.EndTime)
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
				t.id, t.name, p.id, p.name, ts.rank, ts.elo
			FROM tournament_standings ts
				JOIN player p ON p.id = ts.player_id
				JOIN tournament t ON t.id = ts.tournament_id
			WHERE ts.tournament_id = ?
			ORDER BY rank`, id)
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

// GetTournamentMatches will return all matches for the given tournament
func GetTournamentMatches(id int) (map[int][]*models.Match, error) {
	rows, err := models.DB.Query(`
		SELECT
			m.id, m.is_finished, m.current_leg_id, m.winner_id, m.created_at, m.updated_at, m.owe_type_id, m.venue_id,
			mt.id, mt.name, mt.description, mm.id, mm.name, mm.short_name, mm.wins_required, mm.legs_required,
			v.id, v.name, v.description, m.updated_at as 'last_throw', GROUP_CONCAT(DISTINCT p2l.player_id ORDER BY p2l.order) AS 'players',
			m.tournament_id, tg.id
		FROM matches m
			JOIN match_type mt ON mt.id = m.match_type_id
			JOIN match_mode mm ON mm.id = m.match_mode_id
			LEFT JOIN leg l ON l.id = m.current_leg_id
			LEFT JOIN venue v on v.id = m.venue_id
			LEFT JOIN player2leg p2l ON p2l.match_id = m.id
			JOIN tournament t ON t.id = m.tournament_id
			JOIN player2tournament p2t ON p2t.player_id = p2l.player_id AND p2t.tournament_id = t.id
			JOIN tournament_group tg ON tg.id = p2t.tournament_group_id
		WHERE t.id = ?
		GROUP BY m.id
		ORDER BY m.id DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make(map[int][]*models.Match, 0)
	for rows.Next() {
		var groupID int
		m := new(models.Match)
		m.MatchType = new(models.MatchType)
		m.MatchMode = new(models.MatchMode)
		venue := new(models.Venue)
		var players string
		err := rows.Scan(&m.ID, &m.IsFinished, &m.CurrentLegID, &m.WinnerID, &m.CreatedAt, &m.UpdatedAt, &m.OweTypeID, &m.VenueID,
			&m.MatchType.ID, &m.MatchType.Name, &m.MatchType.Description,
			&m.MatchMode.ID, &m.MatchMode.Name, &m.MatchMode.ShortName, &m.MatchMode.WinsRequired, &m.MatchMode.LegsRequired,
			&venue.ID, &venue.Name, &venue.Description, &m.LastThrow, &players, &m.TournamentID, &groupID)
		if err != nil {
			return nil, err
		}
		if m.VenueID.Valid {
			m.Venue = venue
		}
		m.Players = util.StringToIntArray(players)

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
			p.id as 'player_id',
			p2t.is_promoted, p2t.is_relegated, p2t.is_winner,
			COUNT(DISTINCT m.id) AS 'p',
			COUNT(DISTINCT won.id) AS 'w',
			COUNT(DISTINCT draw.id) AS 'd',
			COUNT(DISTINCT lost.id) AS 'l',
			COUNT(DISTINCT legs_for.id) AS 'F',
			COUNT(DISTINCT legs_against.id) AS 'A',
			(COUNT(DISTINCT legs_for.id) - COUNT(DISTINCT legs_against.id)) AS 'diff',
			COUNT(DISTINCT won.id) * 2 + COUNT(DISTINCT draw.id) AS 'pts',
			IFNULL(SUM(s.ppd) / COUNT(p.id), 0) AS 'ppd',
			IFNULL(SUM(s.first_nine_ppd) / COUNT(p.id), 0) AS 'first_nine_ppd',
			IFNULL(SUM(60s_plus), 0) AS '60s_plus',
			IFNULL(SUM(100s_plus), 0) AS '100s_plus',
			IFNULL(SUM(140s_plus), 0) AS '140s_plus',
			IFNULL(SUM(180s), 0) AS '180s',
			IFNULL(SUM(accuracy_20) / COUNT(accuracy_20), 0) AS 'accuracy_20s',
			IFNULL(SUM(accuracy_19) / COUNT(accuracy_19), 0) AS 'accuracy_19s',
			IFNULL(SUM(overall_accuracy) / COUNT(overall_accuracy), 0) AS 'accuracy_overall',
			IFNULL(SUM(s.checkout_attempts), 0) as 'checkout_attempts',
			IFNULL(COUNT(s.checkout_percentage) / SUM(s.checkout_attempts) * 100, 0) AS 'checkout_percentage'
		FROM player2leg p2l
			JOIN matches m ON m.id = p2l.match_id
			JOIN player p ON p.id = p2l.player_id
			LEFT JOIN statistics_x01 s ON s.leg_id = p2l.leg_id AND s.player_id = p.id
			LEFT JOIN matches won ON won.id = p2l.match_id AND won.winner_id = p.id
			LEFT JOIN matches lost ON lost.id = p2l.match_id AND lost.winner_id <> p.id
			LEFT JOIN matches draw ON draw.id = p2l.match_id AND draw.winner_id IS NULL
			LEFT JOIN leg legs_for ON legs_for.id = p2l.leg_id AND legs_for.winner_id = p.id
			LEFT JOIN leg legs_against ON legs_against.id = p2l.leg_id AND legs_against.winner_id <> p.id
			JOIN tournament t ON t.id = m.tournament_id
			JOIN player2tournament p2t ON p2t.player_id = p.id AND p2t.tournament_id = t.id
			JOIN tournament_group tg ON tg.id = p2t.tournament_group_id
		WHERE m.tournament_id = ?
		GROUP BY p2l.player_id, tg.id
		ORDER BY tg.division, pts DESC, diff DESC, is_relegated, manual_order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	statistics := make(map[int][]*models.TournamentOverview, 0)
	for rows.Next() {
		tournament := new(models.Tournament)
		group := new(models.TournamentGroup)
		stats := new(models.TournamentOverview)
		err := rows.Scan(&tournament.ID, &tournament.Name, &tournament.ShortName, &tournament.StartTime, &tournament.EndTime, &group.ID,
			&group.Name, &group.Division, &stats.PlayerID, &stats.IsPromoted, &stats.IsRelegated, &stats.IsWinner, &stats.Played, &stats.MatchesWon,
			&stats.MatchesDraw, &stats.MatchesLost, &stats.LegsFor, &stats.LegsAgainst, &stats.LegsDifference, &stats.Points, &stats.PPD,
			&stats.FirstNinePPD, &stats.Score60sPlus, &stats.Score100sPlus, &stats.Score140sPlus, &stats.Score180s, &stats.Accuracy20,
			&stats.Accuracy19, &stats.AccuracyOverall, &stats.CheckoutAttempts, &stats.CheckoutPercentage)
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
		statistics.BestPPD = append(statistics.BestPPD, val.BestPPD)
		statistics.BestFirstNinePPD = append(statistics.BestFirstNinePPD, val.BestFirstNinePPD)
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

// GetTournamentStandings will return statistics for the given tournament
func GetTournamentStandings() ([]*models.TournamentStanding, error) {
	rows, err := models.DB.Query(`
		SELECT
			pe.player_id,
			p.name,
			pe.tournament_elo,
			pe.tournament_elo_matches
		FROM player_elo pe
		JOIN player p ON p.id = pe.player_id
		WHERE pe.tournament_elo_matches > 10
			AND p.active = 1
		ORDER BY tournament_elo DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rank := 0
	standings := make([]*models.TournamentStanding, 0)
	for rows.Next() {
		standing := new(models.TournamentStanding)
		err := rows.Scan(&standing.PlayerID, &standing.PlayerName, &standing.Elo, &standing.EloPlayed)
		if err != nil {
			return nil, err
		}
		rank++
		standing.Rank = rank
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
			WHERE l.winner_id = s.player_id
				AND s.leg_id IN (SELECT id FROM leg WHERE match_id IN (SELECT id FROM matches WHERE tournament_id = ?))
				AND s.id IN (SELECT MAX(s.id) FROM score s JOIN leg l ON l.id = s.leg_id WHERE l.winner_id = s.player_id GROUP BY leg_id)
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
			s.ppd,
			s.first_nine_ppd,
			s.checkout_percentage,
			s.darts_thrown,
			l.starting_score
		FROM statistics_x01 s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
		WHERE s.leg_id IN (SELECT id FROM leg WHERE match_id IN (SELECT id FROM matches WHERE tournament_id = ?))
		AND p.id = l.winner_id`, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsX01, 0)
	for rows.Next() {
		s := new(models.StatisticsX01)
		err := rows.Scan(&s.PlayerID, &s.LegID, &s.PPD, &s.FirstNinePPD, &s.CheckoutPercentage, &s.DartsThrown, &s.StartingScore)
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

		// Only count best statistics when the player actually won the leg
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
		if best.BestPPD == nil {
			best.BestPPD = new(models.BestStatisticFloat)
		}
		if stat.PPD > best.BestPPD.Value {
			best.BestPPD.Value = stat.PPD
			best.BestPPD.LegID = stat.LegID
			best.BestPPD.PlayerID = stat.PlayerID
		}
		if best.BestFirstNinePPD == nil {
			best.BestFirstNinePPD = new(models.BestStatisticFloat)
		}
		if stat.FirstNinePPD > best.BestFirstNinePPD.Value {
			best.BestFirstNinePPD.Value = stat.FirstNinePPD
			best.BestFirstNinePPD.LegID = stat.LegID
			best.BestFirstNinePPD.PlayerID = stat.PlayerID
		}
	}

	s := make([]*models.StatisticsX01, 0)
	for _, val := range bestStatistics {
		s = append(s, val)
	}

	return s, nil
}
