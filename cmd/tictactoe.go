package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// tictactoeCmd represents the tictactoe command
var tictactoeCmd = &cobra.Command{
	Use:   "tictactoe",
	Short: "Recalculate Tic-Tac-Toe statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.TICTACTOE, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(tictactoeCmd)
}
