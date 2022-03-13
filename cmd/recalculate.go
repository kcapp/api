package cmd

import (
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

var since string
var dryRun bool
var legID int

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

		since, _ = cmd.Flags().GetString("since")
		dryRun, _ = cmd.Flags().GetBool("dry-run")
		legID, _ = cmd.Flags().GetInt("leg")
	},
}

func init() {
	rootCmd.AddCommand(recalculateCmd)
	recalculateCmd.PersistentFlags().Bool("dry-run", true, "Print queries instead of executing")
	recalculateCmd.PersistentFlags().StringP("since", "s", "", "Only recalculate statistics newer than the given date")
	recalculateCmd.PersistentFlags().IntP("leg", "l", 0, "Recalculate statistics for the given leg id")
}
