package database

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

func CreateConnection(dbPath string) (*bun.DB, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, dbPath)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, sqlitedialect.New(), bun.WithDiscardUnknownColumns())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db, nil
}
