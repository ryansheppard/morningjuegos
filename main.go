package main

import (
	"database/sql"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/ryansheppard/morningjuegos/cmd"
	"github.com/ryansheppard/morningjuegos/internal/games/coffeegolf"
	"github.com/ryansheppard/morningjuegos/migrations"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	sqldb, err := sql.Open(sqliteshim.ShimName, dbPath)
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUNDEBUG"),
	))

	migrations.SetDB(db)
	coffeegolf.SetDB(db)

	cmd.Execute()
}
