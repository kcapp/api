package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// bermudatriangleCmd represents the bermudatriangle command
var bermudatriangleCmd = &cobra.Command{
	Use:   "bermudatriangle",
	Short: "Recalculate Bermuda Triangle statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.BERMUDATRIANGLE, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(bermudatriangleCmd)
}
