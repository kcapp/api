package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// knockoutCmd represents the knockout command
var knockoutCmd = &cobra.Command{
	Use:   "knockout",
	Short: "Recalculate Knockout statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating Knockout Statistics since=%s, dryRun=%t", since, dryRun)
		err := data.RecalculateKnockoutStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(knockoutCmd)
}
