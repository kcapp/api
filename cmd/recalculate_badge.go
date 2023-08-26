package cmd

import (
	"log"

	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// recalculateBadgeCmd represents the recalculate command
var recalculateBadgeCmd = &cobra.Command{
	Use:   "recalculate",
	Short: "Recalculate badge",
	Long:  `Recalculate badges earned by each player`,
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
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Hello")
	},
}

func init() {
	badgeCmd.AddCommand(recalculateBadgeCmd)
}
