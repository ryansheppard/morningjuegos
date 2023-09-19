package main

import "github.com/ryansheppard/morningjuegos/cmd"

func main() {
	cmd.Execute()
}

// package main

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/uptrace/bun"
// 	"github.com/uptrace/bun/dialect/sqlitedialect"
// 	"github.com/uptrace/bun/driver/sqliteshim"
// 	"github.com/uptrace/bun/extra/bundebug"

// 	"github.com/ryansheppard/morningjuegos/internal/discord"
// 	"github.com/ryansheppard/morningjuegos/internal/games/coffeegolf"
// )

// func main() {
// 	token := os.Getenv("DISCORD_TOKEN")
// 	appID := os.Getenv("DISCORD_APP_ID")
// 	sqldb, err := sql.Open(sqliteshim.ShimName, "db.sqlite3")
// 	if err != nil {
// 		panic(err)
// 	}

// 	db := bun.NewDB(sqldb, sqlitedialect.New())
// 	db.AddQueryHook(bundebug.NewQueryHook(
// 		bundebug.WithEnabled(false),
// 		bundebug.FromEnv("BUNDEBUG"),
// 	))

// 	coffeegolf.SetDB(db)

// 	// TODO: make some migrations instead of this
// 	_, err = db.NewCreateTable().Model((*coffeegolf.CoffeeGolfRound)(nil)).IfNotExists().Exec(context.TODO())
// 	_, err = db.NewCreateTable().Model((*coffeegolf.CoffeeGolfHole)(nil)).IfNotExists().Exec(context.TODO())
// 	if err != nil {
// 		panic(err)
// 	}

// 	d := discord.NewDiscord(token, appID)
// 	d.AddParser(coffeegolf.NewCoffeeGolfParser())
// 	d.AddCommand(coffeegolf.LeaderboardCommand, coffeegolf.Leaderboard)

// 	fmt.Println("MorningJuegos is now running. Press CTRL-C to exit.")
// 	sc := make(chan os.Signal, 1)
// 	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
// 	<-sc

// 	d.Discord.Close()
// }
