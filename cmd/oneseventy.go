package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// oneseventyCmd represents the oneseventy command
var oneseventyCmd = &cobra.Command{
	Use:   "170",
	Short: "Recalculate 170 statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.ONESEVENTY, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(oneseventyCmd)
}
