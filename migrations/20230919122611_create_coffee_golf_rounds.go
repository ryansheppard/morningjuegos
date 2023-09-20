package migrations

import (
	"context"

	"github.com/ryansheppard/morningjuegos/internal/coffeegolf"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().Model((*coffeegolf.CoffeeGolfRound)(nil)).Exec(context.TODO())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model((*coffeegolf.CoffeeGolfHole)(nil)).Exec(context.TODO())
		if err != nil {
			return err
		}

		return nil
	}, nil)
}
