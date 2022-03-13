package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// cricketCmd represents the cricket command
var cricketCmd = &cobra.Command{
	Use:   "cricket",
	Short: "Recalculate Cricket statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.CRICKET, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(cricketCmd)
}
