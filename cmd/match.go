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
		configFileParam, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}
		config, err := models.GetConfig(configFileParam)
		if err != nil {
			panic(err)
		}
		models.InitDB(config.GetMysqlConnectionString())
	},
}

func init() {
	rootCmd.AddCommand(matchCmd)
}
