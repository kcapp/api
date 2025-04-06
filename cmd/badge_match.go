package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// recalculateMatchBadgeCmd represents the match command
var recalculateMatchBadgeCmd = &cobra.Command{
	Use:   "match",
	Short: "Recalculate Match Badges",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateMatchBadges()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateBadgeCmd.AddCommand(recalculateMatchBadgeCmd)
}
