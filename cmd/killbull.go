package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// killbullCmd represents the killbull command
var killbullCmd = &cobra.Command{
	Use:   "killbull",
	Short: "Recalculate Kill Bull statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.KILLBULL, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(killbullCmd)
}
