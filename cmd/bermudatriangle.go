package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// bermudatriangleCmd represents the bermudatriangle command
var bermudatriangleCmd = &cobra.Command{
	Use:   "bermudatriangle",
	Short: "Recalculate Bermuda Triangle statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating Bermuda Triangle Statistics since=%s, dryRun=%t", since, dryRun)
		err := data.RecalculateBermudaTriangleStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(bermudatriangleCmd)
}
