package data

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/guregu/null"
	"github.com/kcapp/api/models"
)

// RecalculateStatistics will recaulcate statistics based on the given parameters
func RecalculateStatistics(matchType int, legID int, since string, dryRun bool) error {
	legs := make([]int, 0)
	if legID != 0 {
		log.Printf("Recalculating statistics for leg %d", legID)
		legs = append(legs, legID)
	} else {
		s := since
		if s == "" {
			s = "(All Time)"
			since = "1970-01-01"
		}
		log.Printf("Recalculating %s statistics since=%s", models.MatchTypes[matchType], s)
		ids, err := GetLegsToRecalculate(matchType, since)
		if err != nil {
			return err
		}
		legs = append(legs, ids...)
	}

	var queries []string
	var err error
	switch matchType {
	case models.X01:
		queries, err = RecalculateX01Statistics(legs)
	case models.SHOOTOUT:
		queries, err = RecalculateShootoutStatistics(legs)
	case models.X01HANDICAP:
	case models.CRICKET:
		queries, err = RecalculateCricketStatistics(legs)
	case models.DARTSATX:
		queries, err = RecalculateDartsAtXStatistics(legs)
	case models.AROUNDTHEWORLD:
		queries, err = RecalculateAroundTheWorldStatistics(legs)
	case models.SHANGHAI:
		queries, err = RecalculateShanghaiStatistics(legs)
	case models.AROUNDTHECLOCK:
		queries, err = RecalculateAroundTheClockStatistics(legs)
	case models.TICTACTOE:
		queries, err = RecalculateTicTacToeStatistics(legs)
	case models.BERMUDATRIANGLE:
		queries, err = RecalculateBermudaTriangleStatistics(legs)
	case models.FOURTWENTY:
		queries, err = Recalculate420Statistics(legs)
	case models.KILLBULL:
		queries, err = RecalculateKillBullStatistics(legs)
	case models.GOTCHA:
		queries, err = RecalculateGotchaStatistics(legs)
	case models.JDCPRACTICE:
		queries, err = RecalculateJDCPracticeStatistics(legs)
	case models.KNOCKOUT:
		queries, err = RecalculateKnockoutStatistics(legs)
	case models.SCAM:
		queries, err = ReCalculateScamStatistics(legs)
	default:
		return fmt.Errorf("cannot recalculate statistics for type %d", matchType)
	}
	if err != nil {
		return err
	}

	if len(queries) == 0 {
		log.Print("No legs to recalculate")
	} else {
		if dryRun {
			for _, query := range queries {
				log.Print(query)
			}
		} else {
			log.Printf("Executing %d UPDATE queries", len(queries))
			tx, err := models.DB.Begin()
			if err != nil {
				return err
			}
			for _, query := range queries {
				_, err = tx.Exec(query)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
			tx.Commit()
		}
	}
	return nil
}

// RecalculateElo will recalculate Elo for all players
func RecalculateElo(dryRun bool) error {
	rows, err := models.DB.Query(`
		SELECT id FROM matches
		WHERE is_finished = 1 AND is_practice = 0 AND is_abandoned = 0 AND match_type_id = 1
		ORDER BY updated_at`)
	if err != nil {
		return err
	}
	defer rows.Close()

	matches := make([]int, 0)
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		matches = append(matches, id)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	if dryRun {
		log.Print("Elo not reset because dry-run is enabled")
	} else {
		log.Printf("Recalculating elo for %d matches", len(matches))
		tx, err := models.DB.Begin()
		if err != nil {
			return err
		}
		// Reset the Elo for all players back to initial values
		tx.Exec(`UPDATE player_elo SET current_elo = 1500, current_elo_matches = 0, tournament_elo = 1500, tournament_elo_matches = 0;`)
		tx.Exec(`DELETE FROM player_elo_changelog;`)
		tx.Commit()

		for _, id := range matches {
			err = UpdateEloForMatch(id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CalculateEloForTournament will calculate a local elo for a given tournament
func CalculateEloForTournament(tournamentID int) error {
	players, err := GetPlayers()
	if err != nil {
		return err
	}
	tournamentPlayers, err := GetTournamentPlayers(tournamentID)
	if err != nil {
		return err
	}

	rows, err := models.DB.Query(`
		SELECT id FROM matches
		WHERE tournament_id = ? AND is_finished = 1 AND is_practice = 0 AND is_abandoned = 0 AND match_type_id = 1
		ORDER BY updated_at`, tournamentID)
	if err != nil {
		return err
	}
	defer rows.Close()

	matches := make([]int, 0)
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		matches = append(matches, id)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	elos := make(map[int]*models.PlayerElo)
	for _, player := range tournamentPlayers {
		elo := new(models.PlayerElo)
		elo.TournamentElo = null.IntFrom(1500)
		elo.TournamentEloMatches = 0
		elo.PlayerID = player.ID

		elos[player.ID] = elo
	}

	for _, matchID := range matches {
		match, err := GetMatch(matchID)
		if err != nil {
			return err
		}
		wins, err := GetWinsPerPlayer(matchID)
		if err != nil {
			return err
		}

		p1 := elos[match.Players[0]]
		p2 := elos[match.Players[1]]

		// Calculate elo for winner and looser
		one, two := CalculateElo(int(p1.TournamentElo.Int64), p1.TournamentEloMatches, wins[p1.PlayerID],
			int(p2.TournamentElo.Int64), p2.TournamentEloMatches, wins[p2.PlayerID])
		p1.TournamentElo = null.IntFrom(int64(one))
		p2.TournamentElo = null.IntFrom(int64(two))
		p1.TournamentEloMatches++
		p2.TournamentEloMatches++
	}

	values := make([]*models.PlayerElo, 0, len(elos))
	for _, value := range elos {
		values = append(values, value)
	}
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].TournamentElo.Int64 > values[j].TournamentElo.Int64
	})
	log.Printf("Calculated the following elo for tournament %d:", tournamentID)
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tPlayer\tElo\tMatches")
	for _, elo := range values {
		player := players[elo.PlayerID]
		fmt.Fprintf(w, "%d\t%s %s\t%d\t%d\n", player.ID, player.FirstName, player.LastName.String, elo.TournamentElo.Int64, elo.TournamentEloMatches)
	}
	w.Flush()

	return nil
}

func RecalculateLegBadges() error {
	ids, err := GetBadgeLegsToRecalculate()
	if err != nil {
		return err
	}

	for _, legID := range ids {
		log.Printf("Checking Leg %d for badges", legID)
		leg, err := GetLeg(legID)
		if err != nil {
			return err
		}

		statistics, err := GetPlayerBadgeStatistics(leg.Players, &legID)
		if err != nil {
			return err
		}

		// Calculate badges
		err = CheckLegForBadges(leg, statistics)
		if err != nil {
			return err
		}
	}

	return nil
}

func RecalculateGlobalBadges() error {
	players, err := GetPlayers()
	if err != nil {
		return err
	}
	for _, player := range players {
		if player.IsSupporter {
			// Add supporter badge
			err = AddBadge(player.ID, new(models.BadgeKcappSupporter))
			if err != nil {
				return err
			}
		}
		if player.VocalName.Valid && strings.HasSuffix(player.VocalName.String, ".wav") {
			// Add vocal name badge
			err = AddBadge(player.ID, new(models.BadgeSayMyName))
			if err != nil {
				return err
			}
		}

		standings, err := GetPlayerTournamentStandings(player.ID)
		if err != nil {
			return err
		}
		if len(standings) > 0 {
			standing := standings[len(standings)-1]
			err = AddTournamentBadge(player.ID, standing.Tournament.ID, new(models.BadgeItsOfficial), standing.Tournament.EndTime.Time)
			if err != nil {
				return err
			}
		}

		standing := getPlayerTournamentStanding(1, standings)
		if standing != nil {
			err = AddTournamentBadge(player.ID, standing.Tournament.ID, new(models.BadgeTournament1st), standing.Tournament.EndTime.Time)
			if err != nil {
				return err
			}
		}
		standing = getPlayerTournamentStanding(2, standings)
		if standing != nil {
			err = AddTournamentBadge(player.ID, standing.Tournament.ID, new(models.BadgeTournament2nd), standing.Tournament.EndTime.Time)
			if err != nil {
				return err
			}
		}
		standing = getPlayerTournamentStanding(3, standings)
		if standing != nil {
			err = AddTournamentBadge(player.ID, standing.Tournament.ID, new(models.BadgeTournament3rd), standing.Tournament.EndTime.Time)
			if err != nil {
				return err
			}
		}
	}

	undefeated, err := GetUndefeatedPlayers()
	if err != nil {
		return err
	}
	for playerID, tournament := range undefeated {
		err := AddTournamentBadge(playerID, tournament.ID, new(models.BadgeUntouchable), tournament.EndTime.Time)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPlayerTournamentStanding(pos int, standings []*models.PlayerTournamentStanding) *models.PlayerTournamentStanding {
	for i := len(standings) - 1; i >= 0; i-- {
		standing := standings[i]
		if standing.FinalStanding == pos {
			return standing
		}
	}
	return nil
}
