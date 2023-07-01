package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// jdcpracticeCmd represents the jdcpractice command
var jdcpracticeCmd = &cobra.Command{
	Use:   "jdcpractice",
	Short: "Recalculate JDC Practice statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.JDCPRACTICE, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(jdcpracticeCmd)
}
