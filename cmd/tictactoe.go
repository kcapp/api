package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// tictactoeCmd represents the tictactoe command
var tictactoeCmd = &cobra.Command{
	Use:   "tictactoe",
	Short: "Recalculate Tic-Tac-Toe statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating Tic-Tac-Toe since=%s, dryRun=%t", since, dryRun)
		err := data.RecalculateTicTacToeStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(tictactoeCmd)
}
