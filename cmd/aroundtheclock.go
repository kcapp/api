package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// aroundtheclockCmd represents the aroundtheclock command
var aroundtheclockCmd = &cobra.Command{
	Use:   "aroundtheclock",
	Short: "Recalculate Around the Clock statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.AROUNDTHECLOCK, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(aroundtheclockCmd)
}
