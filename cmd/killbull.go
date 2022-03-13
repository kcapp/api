package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// killbullCmd represents the killbull command
var killbullCmd = &cobra.Command{
	Use:   "killbull",
	Short: "Recalculate Kill Bull statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating Kill Bull Statistics since=%s, dryRun=%t", since, dryRun)
		err := data.RecalculateKillBullStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(killbullCmd)
}
