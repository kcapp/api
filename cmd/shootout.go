package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// shootoutCmd represents the shootout command
var shootoutCmd = &cobra.Command{
	Use:   "shootout",
	Short: "Recalculate 9 Dart Shootout statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.SHOOTOUT, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(shootoutCmd)
}
