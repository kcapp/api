package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// shanghaiCmd represents the shanghai command
var shanghaiCmd = &cobra.Command{
	Use:   "shanghai",
	Short: "Recalculate Shanghai statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.SHANGHAI, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateStatisticsCmd.AddCommand(shanghaiCmd)
}
