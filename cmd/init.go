package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ryansheppard/morningjuegos/internal/database"
	"github.com/ryansheppard/morningjuegos/migrations"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Runs migrations init",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := os.Getenv("DB_PATH")
		db, err := database.CreateConnection(dbPath)
		if err != nil {
			panic(err)
		}
		migrations.InitMigrations(db)
	},
}

func init() {
	migrationsCmd.AddCommand(initCmd)
}
