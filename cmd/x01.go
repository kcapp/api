package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// x01Cmd represents the x01 command
var x01Cmd = &cobra.Command{
	Use:   "x01",
	Short: "Recalculate x01 statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating x01 since=%s, dryRun=%t", since, dryRun)
		err := data.RecalculateX01Statistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(x01Cmd)
}
