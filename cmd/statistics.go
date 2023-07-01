package cmd

import (
	"github.com/spf13/cobra"
)

// statisticsCmd represents the statistics command
var statisticsCmd = &cobra.Command{
	Use:   "statistics",
	Short: "Modify statistics",
}

func init() {
	rootCmd.AddCommand(statisticsCmd)
}
