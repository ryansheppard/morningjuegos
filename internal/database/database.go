package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

var Connection *bun.DB

func CreateConnection(dbPath string) error {
	sqldb, err := sql.Open(sqliteshim.ShimName, dbPath)
	if err != nil {
		return err
	}

	db := bun.NewDB(sqldb, sqlitedialect.New(), bun.WithDiscardUnknownColumns())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUNDEBUG"),
	))

	Connection = db
	return nil
}

func GetDB() *bun.DB {
	if Connection == nil {
		fmt.Println("No DB connection created")
		os.Exit(1)
	}

	return Connection
}
