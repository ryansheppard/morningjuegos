package migrations

import (
	"context"

	"github.com/ryansheppard/morningjuegos/internal/coffeegolf"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().Model((*coffeegolf.Round)(nil)).Exec(context.TODO())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model((*coffeegolf.Hole)(nil)).Exec(context.TODO())
		if err != nil {
			return err
		}

		return nil
	}, nil)
}
