package data

import (
	"fmt"
	"log"

	"github.com/kcapp/api/models"
)

// RecalculateTicTacToeStatistics will recaulcate statistics for Tic Tac Toe legs
func RecalculateStatistics(matchType int, legID int, since string, dryRun bool) error {
	legs := make([]int, 0)
	if legID != 0 {
		log.Printf("Recalculating statistics for leg %d", legID)
		legs = append(legs, legID)
	} else {
		s := since
		if s == "" {
			s = "(All Time)"
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
