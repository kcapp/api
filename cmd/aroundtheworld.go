package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// aroundtheworldCmd represents the aroundtheworld command
var aroundtheworldCmd = &cobra.Command{
	Use:   "aroundtheworld",
	Short: "Recalculate Around the World statistics",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateStatistics(models.AROUNDTHEWORLD, legID, since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(aroundtheworldCmd)
}
