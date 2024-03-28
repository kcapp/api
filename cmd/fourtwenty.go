package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// fourtwentyCmd represents the fourtwenty command
var fourtwentyCmd = &cobra.Command{
	Use:   "420",
	Short: "Recalculate 420 statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.FOURTWENTY, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(fourtwentyCmd)
}
