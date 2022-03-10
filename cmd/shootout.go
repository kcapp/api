package cmd

import (
	"log"

	"github.com/kcapp/api/data"
	"github.com/spf13/cobra"
)

// shootoutCmd represents the shootout command
var shootoutCmd = &cobra.Command{
	Use:   "shootout",
	Short: "Recalculate 9 Dart Shootout statistics",
	Long:  `Recalculate 9 Dart Shootout statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		since, _ := cmd.Flags().GetString("since")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		log.Printf("Recalculating 9 Dart Shootout Statistics since=%s, dryRun=%t", since, dryRun)
		err := data.ReCalculateShootoutStatistics(since, dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recalculateCmd.AddCommand(shootoutCmd)
}
