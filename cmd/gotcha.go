package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// gotchaCmd represents the gotcha command
var gotchaCmd = &cobra.Command{
	Use:   "gotcha",
	Short: "Recalculate Gotcha statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating Gotcha Statistics since=%s, dryRun=%t", since, dryRun)
		err := data.RecalculateGotchaStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(gotchaCmd)
}
