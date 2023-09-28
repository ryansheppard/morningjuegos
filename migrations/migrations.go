package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

var migrator *migrate.Migrator

func InitMigrations(db *bun.DB) {
	migrator = migrate.NewMigrator(db, Migrations)
	fmt.Println(migrator)
	migrator.Init(context.TODO())
}

func RunMigrations(db *bun.DB) error {
	migrator = migrate.NewMigrator(db, Migrations)
	if err := migrator.Lock(context.TODO()); err != nil {
		return err
	}
	defer migrator.Unlock(context.TODO())

	group, err := migrator.Migrate(context.TODO())
	if err != nil {
		return err
	}
	if group.IsZero() {
		fmt.Printf("there are no new migrations to run (database is up to date)\n")
		return nil
	}
	fmt.Printf("migrated to %s\n", group)
	return nil

}

func CreateMigration(db *bun.DB, name string) error {
	migrator := migrate.NewMigrator(db, Migrations)
	mf, err := migrator.CreateGoMigration(context.TODO(), name)
	if err != nil {
		return err
	}
	fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)

	return nil
}
