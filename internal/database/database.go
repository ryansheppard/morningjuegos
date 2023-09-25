package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

var connection *bun.DB

var ctx context.Context

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

	connection = db
	return nil
}

func GetDB() *bun.DB {
	if connection == nil {
		fmt.Println("No DB connection created")
		os.Exit(1)
	}

	return connection
}

func GetContext() context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return ctx
}
