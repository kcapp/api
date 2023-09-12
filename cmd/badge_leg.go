package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// recalculateLegBadgesCmd represents the leg command
var recalculateLegBadgesCmd = &cobra.Command{
	Use:   "leg",
	Short: "Recalculate Leg Badges",
	Run: func(cmd *cobra.Command, args []string) {
		err := data.RecalculateLegBadges()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateBadgeCmd.AddCommand(recalculateLegBadgesCmd)
}
