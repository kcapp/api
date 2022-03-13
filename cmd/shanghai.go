package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// shanghaiCmd represents the shanghai command
var shanghaiCmd = &cobra.Command{
	Use:   "shanghai",
	Short: "Recalculate Shanghai statistics",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Recalculating Shanghai since=%s, dryRun=%t", since, dryRun)
		err := data.RecalculateShanghaiStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(shanghaiCmd)
}
