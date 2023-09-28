package cmd

import (
	"os"

	"github.com/ryansheppard/morningjuegos/internal/database"
	"github.com/ryansheppard/morningjuegos/migrations"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs migrations",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := os.Getenv("DB_PATH")
		db, err := database.CreateConnection(dbPath)
		if err != nil {
			panic(err)
		}
		migrations.RunMigrations(db)
	},
}

func init() {
	migrationsCmd.AddCommand(runCmd)
}
