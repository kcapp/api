package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// dartsatxCmd represents the dartsatx command
var dartsatxCmd = &cobra.Command{
	Use:   "dartsatx",
	Short: "Recalculate Darts at X statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating Darts at X Statistics since=%s, dryRun=%t", since, dryRun)
		err := data.RecalculateDartsAtXStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(dartsatxCmd)
}
