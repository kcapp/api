package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// cricketCmd represents the cricket command
var cricketCmd = &cobra.Command{
	Use:   "cricket",
	Short: "Recalculate Cricket statistics",
	Long:  `Recalculate Cricket statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		since, _ := cmd.Flags().GetString("since")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		log.Printf("Recalculating Cricket Statistics since=%s, dryRun=%t", since, dryRun)
		err := data.ReCalculateCricketStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(cricketCmd)
}
