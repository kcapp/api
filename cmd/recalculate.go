package cmd

import (
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// recalculateCmd represents the recalculate command
var recalculateCmd = &cobra.Command{
	Use:   "recalculate",
	Short: "Recalculate statistics",
	Long:  `Recalculate statistics for the given match type`,
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
	rootCmd.AddCommand(recalculateCmd)
	recalculateCmd.PersistentFlags().Bool("dry-run", true, "Print queries instead of executing")
	recalculateCmd.PersistentFlags().StringP("since", "s", "", "Recalculate statistics newer than the given date")
}
