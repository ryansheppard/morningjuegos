package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"

	"github.com/ryansheppard/morningjuegos/internal/coffeegolf"
	"github.com/ryansheppard/morningjuegos/internal/database"
	"github.com/ryansheppard/morningjuegos/internal/discord"
	"github.com/ryansheppard/morningjuegos/internal/game"
)

// botCmd represents the bot command
var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Runs the discord bot",
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("DISCORD_TOKEN")
		appID := os.Getenv("DISCORD_APP_ID")
		d := discord.NewDiscord(token, appID)

		coffeegolf.SetDB(database.GetDB())

		games := []*game.Game{coffeegolf.GetCoffeeGolfGame()}

		for _, game := range games {
			d.RegisterGame(game)
		}
		s := gocron.NewScheduler(time.UTC)
		s.Every(1).Minute().Do(coffeegolf.AddMissingRounds)
		s.Every(15).Minute().Do(coffeegolf.AddTournamentWinners)
		s.StartAsync()

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
