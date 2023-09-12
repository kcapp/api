package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// recalculateGlobalBadgeCmd represents the global command
var recalculateGlobalBadgeCmd = &cobra.Command{
	Use:   "global",
	Short: "Recalculate Global Badges",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateGlobalBadges()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateBadgeCmd.AddCommand(recalculateGlobalBadgeCmd)
}
