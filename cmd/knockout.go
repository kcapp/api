package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// knockoutCmd represents the knockout command
var knockoutCmd = &cobra.Command{
	Use:   "knockout",
	Short: "Recalculate Knockout statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.KNOCKOUT, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(knockoutCmd)
}
