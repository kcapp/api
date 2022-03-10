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
	Long:  `Recalculate x01 statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		since, _ := cmd.Flags().GetString("since")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		log.Printf("Recalculating x01 since=%s, dryRun=%t", since, dryRun)
		data.RecalculateX01Statistics()
	},
}

func init() {
	recalculateCmd.AddCommand(x01Cmd)
}
