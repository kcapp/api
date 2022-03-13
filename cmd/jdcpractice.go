package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// jdcpracticeCmd represents the jdcpractice command
var jdcpracticeCmd = &cobra.Command{
	Use:   "jdcpractice",
	Short: "Recalculate JDC Practice statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating JDC Practice Statistics since=%s, dryRun=%t", since, dryRun)
		err := data.RecalculateJDCPracticeStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(jdcpracticeCmd)
}
