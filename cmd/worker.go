package cmd

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/ryansheppard/morningjuegos/internal/cache"
	cgQueries "github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
	coffeegolf "github.com/ryansheppard/morningjuegos/internal/coffeegolf/game"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/service"
	"github.com/ryansheppard/morningjuegos/internal/discord"
	"github.com/ryansheppard/morningjuegos/internal/messenger"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Runs the discord jobs",
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

		c := cache.New(ctx, redisAddr, redisDBInt)

		dsn := os.Getenv("DB_DSN")
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			slog.Error("Error opening database connection", "error", err)
			os.Exit(1)
		}

		natsURL := os.Getenv("NATS_URL")
		m := messenger.New(natsURL)

		q := cgQueries.New(db)

		service := service.New(db, q)

		cg := coffeegolf.New(ctx, service, c, db, m)

		token := os.Getenv("DISCORD_TOKEN")
		appID := os.Getenv("DISCORD_APP_ID")
		d, err := discord.NewDiscord(token, appID, m, c, cg)
		if err != nil {
			slog.Error("Error creating discord", "error", err)
			os.Exit(1)
		}

		cg.ConfigureSubscribers()
		d.ConfigureSubscribers()

		slog.Info("MorningJuegos worker is now running. Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
		slog.Info("Shutting down MorningJuegos worker")
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
