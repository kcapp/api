package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// gotchaCmd represents the gotcha command
var gotchaCmd = &cobra.Command{
	Use:   "gotcha",
	Short: "Recalculate Gotcha statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.GOTCHA, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(gotchaCmd)
}
