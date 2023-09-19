package cmd

import (
	"github.com/ryansheppard/morningjuegos/migrations"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a migration",
	Run: func(cmd *cobra.Command, args []string) {
		migrations.CreateMigration(args[0])
	},
}

func init() {
	migrationsCmd.AddCommand(createCmd)
}
