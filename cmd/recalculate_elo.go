package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// recalculateEloCmd represents the recalculate command
var recalculateEloCmd = &cobra.Command{
	Use:   "recalculate",
	Short: "Recalculate elo",
	Long: `Recalculate elo for all matches played.

	This will reset the elo for all players, and regenerate the elo changelog
	Elo will be recalculated based on 'updated_at' timestamp of each match`,
	Run: func(cmd *cobra.Command, args []string) {
		models.InitDB(models.GetMysqlConnectionString())

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		tournament, _ := cmd.Flags().GetInt("tournament")
		if tournament != 0 {
			err := data.CalculateEloForTournament(tournament)
			if err != nil {
				panic(err)
			}
		} else {
			err := data.RecalculateElo(dryRun)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	eloCmd.AddCommand(recalculateEloCmd)
	recalculateEloCmd.Flags().Bool("dry-run", true, "Print queries instead of executing")
	recalculateEloCmd.Flags().IntP("tournament", "t", 0, "Calculate elo for the given tournament")
}
