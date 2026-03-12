package cmd

import (
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

var since string
var dryRun bool
var legID int

// recalculateStatisticsCmd represents the recalculate command
var recalculateStatisticsCmd = &cobra.Command{
	Use:   "recalculate",
	Short: "Recalculate statistics",
	Long:  `Recalculate statistics for the given match type`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		SetupConfig(cmd)
		models.InitDB(models.GetMysqlConnectionString())

		since, _ = cmd.Flags().GetString("since")
		dryRun, _ = cmd.Flags().GetBool("dry-run")
		legID, _ = cmd.Flags().GetInt("leg")
	},
}

func init() {
	statisticsCmd.AddCommand(recalculateStatisticsCmd)
	recalculateStatisticsCmd.PersistentFlags().Bool("dry-run", true, "Print queries instead of executing")
	recalculateStatisticsCmd.PersistentFlags().StringP("since", "s", "", "Only recalculate statistics newer than the given date")
	recalculateStatisticsCmd.PersistentFlags().IntP("leg", "l", 0, "Recalculate statistics for the given leg id")
}
