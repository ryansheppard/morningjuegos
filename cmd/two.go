package cmd

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/ryansheppard/morningjuegos/internal/cache"
	cgQueries "github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/database"
	coffeegolf "github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/game"
	"github.com/ryansheppard/morningjuegos/internal/v2/discord"
)

var twoCmd = &cobra.Command{
	Use:   "two",
	Short: "Runs the v2 discord bot",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		redisAddr := os.Getenv("REDIS_ADDR")
		c := cache.New(ctx, redisAddr, 0)

		dsn := os.Getenv("DB_DSN")
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			slog.Error("Error opening database connection", err)
			os.Exit(1)
		}

		q := cgQueries.New(db)

		cg := coffeegolf.New(ctx, q, c, db)

		token := os.Getenv("DISCORD_TOKEN")
		appID := os.Getenv("DISCORD_APP_ID")
		d := discord.NewDiscord(token, appID, cg)

		d.Configure()

		slog.Info("MorningJuegos is now running. Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc

		d.Discord.Close()
	},
}

func init() {
	rootCmd.AddCommand(twoCmd)
}
