package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// x01Cmd represents the x01 command
var x01Cmd = &cobra.Command{
	Use:   "x01",
	Short: "Recalculate x01 statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.X01, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(x01Cmd)
}
