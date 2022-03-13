package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// aroundtheworldCmd represents the aroundtheworld command
var aroundtheworldCmd = &cobra.Command{
	Use:   "aroundtheworld",
	Short: "Recalculate Around the World statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating Around the World Statistics since=%s, dryRun=%t", since, dryRun)
		err := data.RecalculateAroundTheWorldStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(aroundtheworldCmd)
}
