package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ryansheppard/morningjuegos/internal/database"
	"github.com/ryansheppard/morningjuegos/migrations"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a migration",
	Run: func(cmd *cobra.Command, args []string) {
		migrations.SetDB(database.GetDB())
		migrations.CreateMigration(args[0])
	},
}

func init() {
	migrationsCmd.AddCommand(createCmd)
}
