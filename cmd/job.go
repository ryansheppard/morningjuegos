package cmd

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
	"github.com/ryansheppard/morningjuegos/internal/cache"
	cgQueries "github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
	coffeegolf "github.com/ryansheppard/morningjuegos/internal/coffeegolf/game"
)

var jobsCmd = &cobra.Command{
	Use:   "jobs",
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

		q := cgQueries.New(db)

		cg := coffeegolf.New(ctx, q, c, db)

		newYork, err := time.LoadLocation("America/New_York")
		if err != nil {
			slog.Error("Error loading location", "error", err)
			os.Exit(1)
		}

		s := gocron.NewScheduler(newYork)
		// s.Every(15).Minute().Do(cg.AddMissingRounds)
		s.Every(15).Minute().Do(cg.AddTournamentWinners)
		s.StartAsync()

		slog.Info("MorningJuegos jobs are now running. Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
		s.Stop()

	},
}

func init() {
	rootCmd.AddCommand(jobsCmd)
}
