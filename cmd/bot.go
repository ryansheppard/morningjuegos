package cmd

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
	_ "github.com/lib/pq"
	"github.com/ryansheppard/morningjuegos/internal/cache"
	cgQueries "github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
	coffeegolf "github.com/ryansheppard/morningjuegos/internal/coffeegolf/game"
	"github.com/ryansheppard/morningjuegos/internal/discord"
	"github.com/ryansheppard/morningjuegos/internal/messenger"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Runs the discord bot",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		var err error

		bsp := honeycomb.NewBaggageSpanProcessor()
		otelShutdown, err := otelconfig.ConfigureOpenTelemetry(
			otelconfig.WithSpanProcessor(bsp),
		)
		if err != nil {
			slog.Error("Error configuring opentelemetry", "error", err)
			os.Exit(1)
		}
		defer otelShutdown()

		tracer := otel.Tracer("morningjuegos.bot")

		redisAddr := os.Getenv("REDIS_ADDR")
		redisDB := os.Getenv("REDIS_DB")
		redisDBInt := 0
		if redisDB != "" {
			redisDBInt, err = strconv.Atoi(redisDB)
			if err != nil {
				slog.Error("Error converting redis db to int", "error", err)
			}
		}

		c := cache.New(ctx, redisAddr, redisDBInt, tracer)

		natsURL := os.Getenv("NATS_URL")
		m := messenger.New(natsURL)

		dsn := os.Getenv("DB_DSN")
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			slog.Error("Error opening database connection", "error", err)
			os.Exit(1)
		}

		q := cgQueries.New(db)

		cg := coffeegolf.New(ctx, q, c, db, m, tracer)

		token := os.Getenv("DISCORD_TOKEN")
		appID := os.Getenv("DISCORD_APP_ID")
		d, err := discord.NewDiscord(token, appID, cg)
		if err != nil {
			slog.Error("Error creating discord", "error", err)
			os.Exit(1)
		}

		err = d.Configure()
		if err != nil {
			slog.Error("Error configuring discord", "error", err)
			os.Exit(1)
		}

		slog.Info("MorningJuegos is now running. Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
		slog.Info("Shutting down MorningJuegos bot")

		d.Discord.Close()
	},
}

func init() {
	rootCmd.AddCommand(botCmd)
}
