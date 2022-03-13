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
		log.Printf("Recalculating %s statistics since %s", models.MatchTypes[matchType], since)
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
				tx.Exec(query)
			}
			tx.Commit()
		}
	}
	return nil
}
