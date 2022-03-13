package cmd

import (
	"github.com/kcapp/api/data"
	"github.com/kcapp/api/models"
	"github.com/spf13/cobra"
)

// recalculateEloCmd represents the recalculate command
var recalculateEloCmd = &cobra.Command{
	Use:   "recalculate",
	Short: "Recalculate elo",
	Long: `Recalculate elo for all matches played.

	This will reset the elo for all players, and regenerate the elo changelog
	Elo will be recalculated based on 'updated_at' timestamp of each match`,
	Run: func(cmd *cobra.Command, args []string) {
		configFileParam, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}
		config, err := models.GetConfig(configFileParam)
		if err != nil {
			panic(err)
		}
		models.InitDB(config.GetMysqlConnectionString())

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		err = data.RecalculateElo(dryRun)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	eloCmd.AddCommand(recalculateEloCmd)
	recalculateEloCmd.PersistentFlags().Bool("dry-run", true, "Print queries instead of executing")
}
