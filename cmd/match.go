package cmd

import (
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// matchCmd represents the match command
var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "Import/Export matches",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		models.InitDB(models.GetMysqlConnectionString())
	},
}

func init() {
	rootCmd.AddCommand(matchCmd)
}
