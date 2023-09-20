package migrations

import (
	"context"
	"fmt"

	"github.com/ryansheppard/morningjuegos/internal/coffeegolf"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewAddColumn().Model((*coffeegolf.Round)(nil)).ColumnExpr("guild_id").Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		return nil
	})
}
