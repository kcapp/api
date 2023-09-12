package cmd

import (
	"github.com/spf13/cobra"
)

// badgeCmd represents the badge command
var badgeCmd = &cobra.Command{
	Use:   "badge",
	Short: "Modify Badges",
}

func init() {
	rootCmd.AddCommand(badgeCmd)
}
