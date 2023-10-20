package cmd

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/ryansheppard/morningjuegos/internal/cache"
	cgQueries "github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
	coffeegolf "github.com/ryansheppard/morningjuegos/internal/coffeegolf/game"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/service"
	"github.com/ryansheppard/morningjuegos/internal/messenger"
	"github.com/spf13/cobra"
)

var oneOffCmd = &cobra.Command{
	Use:   "oneoff",
	Short: "Runs one off stuff",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		var err error

		redisAddr := os.Getenv("REDIS_ADDR")
		redisDB := os.Getenv("REDIS_DB")
		redisDBInt := 0
		if redisDB != "" {
			redisDBInt, err = strconv.Atoi(redisDB)
			if err != nil {
				slog.Error("Error converting redis db to int", "error", err)
			}
		}

		c := cache.New(redisAddr, redisDBInt)

		natsURL := os.Getenv("NATS_URL")
		m := messenger.New(natsURL)

		dsn := os.Getenv("DB_DSN")
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			slog.Error("Error opening database connection", "error", err)
			os.Exit(1)
		}

		err = db.Ping()
		if err != nil {
			slog.Error("Error pinging database", "error", err)
			os.Exit(1)
		}

		q := cgQueries.New(db)

		service := service.New(db, q)

		cg := coffeegolf.New(ctx, service, c, db, m)

		guildID, err := cmd.Flags().GetInt64("guild-id")
		if err != nil {
			slog.Error("Error getting guild id", "error", err)
			os.Exit(1)
		}

		cg.AddTournamentWinnersForGuild(guildID)

		m.CleanUp()
	},
}

func init() {
	rootCmd.AddCommand(oneOffCmd)

	oneOffCmd.PersistentFlags().Int64("guild-id", 0, "The Guild ID")
}
