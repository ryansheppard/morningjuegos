package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ryansheppard/morningjuegos/migrations"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Runs migrations init",
	Run: func(cmd *cobra.Command, args []string) {
		migrations.InitMigrations()
	},
}

func init() {
	migrationsCmd.AddCommand(initCmd)
}
