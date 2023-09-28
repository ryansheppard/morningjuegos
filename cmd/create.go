package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ryansheppard/morningjuegos/internal/database"
	"github.com/ryansheppard/morningjuegos/migrations"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a migration",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := os.Getenv("DB_PATH")
		db, err := database.CreateConnection(dbPath)
		if err != nil {
			panic(err)
		}

		migrations.CreateMigration(db, args[0])
	},
}

func init() {
	migrationsCmd.AddCommand(createCmd)
}
