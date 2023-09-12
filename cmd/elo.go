package cmd

import (
	"github.com/spf13/cobra"
)

// eloCmd represents the elo command
var eloCmd = &cobra.Command{
	Use:   "elo",
	Short: "Modify Elo",
}

func init() {
	rootCmd.AddCommand(eloCmd)
}
