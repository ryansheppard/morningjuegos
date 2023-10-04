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
	"github.com/ryansheppard/morningjuegos/internal/messenger"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Runs the discord jobs",
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

		tracer := otel.Tracer("morningjuegos.worker")

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

		dsn := os.Getenv("DB_DSN")
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			slog.Error("Error opening database connection", "error", err)
			os.Exit(1)
		}

		natsURL := os.Getenv("NATS_URL")
		m := messenger.New(natsURL)

		q := cgQueries.New(db)

		cg := coffeegolf.New(ctx, q, c, db, m, tracer)
		cg.ConfigureSubscribers()

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
