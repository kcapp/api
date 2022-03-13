package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// aroundtheclockCmd represents the aroundtheclock command
var aroundtheclockCmd = &cobra.Command{
	Use:   "aroundtheclock",
	Short: "Recalculate Around the Clock statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating Around the Clock Statistics since=%s, dryRun=%t", since, dryRun)
		err := data.RecalculateAroundTheClockStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(aroundtheclockCmd)
}
