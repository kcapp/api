package cmd

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kcapp/api/data"
	data_v2 "github.com/kcapp/api/data/v2"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import matches from the given file",
	Long: `Import all the matches from a specific file

	File should be a CSV with the following format
	<datetime yyyy-MM-dd HH:mm:ss>, office name, venue name, starting score, type, `,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		fmt.Printf("import called %s", filename)
		records, err := readFile(filename)
		if err != nil {
			panic(err)
		}

		for _, record := range records {
			startingScore, err := strconv.Atoi(record[3])
			if err != nil {
				panic(err)
			}

			homeWins, err := strconv.Atoi(record[7])
			if err != nil {
				panic(err)
			}
			awayWins, err := strconv.Atoi(record[9])
			if err != nil {
				panic(err)
			}

			log.Printf("Need to create %d legs", homeWins+awayWins)

			// Check if players exist
			players, err := data_v2.GetPlayers()
			if err != nil {
				panic(err)
			}
			home := findPlayerByName(record[6], players)
			away := findPlayerByName(record[8], players)
			log.Printf("Found player: %s (%d)", home.FirstName, home.ID)
			log.Printf("Found player: %s (%d)", away.FirstName, away.ID)

			winnerId := 0
			if homeWins == awayWins {
				winnerId = 0
			} else if homeWins > awayWins {
				winnerId = home.ID
			} else {
				winnerId = away.ID
			}

			// Check if office exist
			offices, err := data_v2.GetOffices()
			if err != nil {
				panic(err)
			}
			office := findOfficeByName(record[1], offices)
			log.Printf("Found office: %s (%d)", office.Name, office.ID)

			venues, err := data.GetVenues()
			if err != nil {
				panic(err)
			}
			venue := findVenueByName(record[2], venues)
			log.Printf("Found venue: %s (%d)", venue.Name.String, venue.ID.Int64)

			// TODO Add param to create office/venue if they don't exist

			tx, err := models.DB.Begin()
			if err != nil {
				panic(err)
			}
			matchTypeID := record[4]
			createdAt := record[0]
			log.Printf("Creating match %s", createdAt)
			res, err := tx.Exec(`INSERT INTO matches (is_finished, match_mode_id, winner_id, office_id, created_at, match_type_id, venue_id)
				VALUES (?, ?, ?, ?, ?, ?, ?)`, 1, record[5], winnerId, office.ID, createdAt, matchTypeID, venue.ID)
			if err != nil {
				tx.Rollback()
				panic(err)
			}
			matchID, err := res.LastInsertId()
			if err != nil {
				tx.Rollback()
				panic(err)
			}

			// Create all the legs
			_, err = createLegs(tx, homeWins, matchID, createdAt, startingScore, home.ID, away.ID)
			if err != nil {
				panic(err)
			}
			legID, err := createLegs(tx, awayWins, matchID, createdAt, startingScore, away.ID, home.ID)
			if err != nil {
				panic(err)
			}

			_, err = tx.Exec("UPDATE matches SET current_leg_id = ?, updated_at = NOW() WHERE id = ?", legID, matchID)
			if err != nil {
				tx.Rollback()
				panic(err)
			}
			tx.Commit()

			log.Printf("Created match %d", matchID)
		}
	},
}

func init() {
	matchCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().StringP("filename", "f", "", "File to import from")
	importCmd.MarkPersistentFlagRequired("filename")
}

func readFile(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func findPlayerByName(playerName string, players []*models.Player) *models.Player {
	for _, player := range players {
		name := strings.Trim(player.FirstName+" "+player.LastName.String, " ")
		if name == playerName {
			return player
		}
	}
	return nil
}

func findOfficeByName(officeName string, offices []*models.Office) *models.Office {
	for _, office := range offices {
		if office.Name == officeName {
			return office
		}
	}
	return nil
}

func findVenueByName(venueName string, venues []*models.Venue) *models.Venue {
	for _, venue := range venues {
		if venue.Name.String == venueName {
			return venue
		}
	}
	return nil
}

func createLegs(tx *sql.Tx, numToCreate int, matchID int64, createdAt string, startingScore int, winnerPlayerID int, looserPlayerID int) (*int64, error) {
	var legID int64
	for i := 0; i < numToCreate; i++ {
		log.Printf("Creating leg %d with winner %d and looser %d", i+1, winnerPlayerID, looserPlayerID)
		res, err := tx.Exec(`INSERT INTO leg (end_time, starting_score, is_finished, current_player_id, winner_id, created_at, match_id, has_scores) VALUES
						(?, ?, 1, ?, ?, ?, ?, 0)`, createdAt, startingScore, winnerPlayerID, winnerPlayerID, createdAt, matchID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		legID, err = res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		_, err = tx.Exec("INSERT INTO player2leg (player_id, leg_id, `order`, match_id) VALUES (?, ?, ?, ?), (?, ?, ?, ?)",
			winnerPlayerID, legID, 1, matchID, looserPlayerID, legID, 2, matchID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	return &legID, nil
}
