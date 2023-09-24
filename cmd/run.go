package cmd

import (
	"github.com/ryansheppard/morningjuegos/migrations"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrations.RunMigrations()
	},
}

func init() {
	migrationsCmd.AddCommand(runCmd)
}
