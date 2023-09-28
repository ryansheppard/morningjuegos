package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"

	"github.com/ryansheppard/morningjuegos/internal/cache"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf"
	"github.com/ryansheppard/morningjuegos/internal/database"
	"github.com/ryansheppard/morningjuegos/internal/discord"
)

// botCmd represents the bot command
var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Runs the discord bot",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		redisAddr := os.Getenv("REDIS_ADDR")
		c := cache.New(ctx, redisAddr, 0)

		dbPath := os.Getenv("DB_PATH")
		db, err := database.CreateConnection(dbPath)
		if err != nil {
			panic(err)
		}

		q := coffeegolf.NewQuery(ctx, db)

		cg := coffeegolf.NewCoffeeGolf(q, c)

		s := gocron.NewScheduler(time.UTC)
		s.Every(1).Minute().Do(cg.AddMissingRounds)
		s.Every(15).Minute().Do(cg.AddTournamentWinners)
		s.StartAsync()

		token := os.Getenv("DISCORD_TOKEN")
		appID := os.Getenv("DISCORD_APP_ID")
		d := discord.NewDiscord(token, appID, cg)

		d.Configure()

		fmt.Println("MorningJuegos is now running. Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc

		d.Discord.Close()
	},
}

func init() {
	rootCmd.AddCommand(botCmd)
}
