package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// dartsatxCmd represents the dartsatx command
var dartsatxCmd = &cobra.Command{
	Use:   "dartsatx",
	Short: "Recalculate Darts at X statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.DARTSATX, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(dartsatxCmd)
}
