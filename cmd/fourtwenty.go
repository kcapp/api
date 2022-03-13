package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// fourtwentyCmd represents the fourtwenty command
var fourtwentyCmd = &cobra.Command{
	Use:   "fourtwenty",
	Short: "Recalculate 420 statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating 420 Statistics since=%s, dryRun=%t", since, dryRun)
		err := data.Recalculate420Statistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(fourtwentyCmd)
}
