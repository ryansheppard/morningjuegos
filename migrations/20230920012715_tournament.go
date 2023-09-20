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
		_, err := db.NewCreateTable().Model((*coffeegolf.Tournament)(nil)).Exec(context.TODO())
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model((*coffeegolf.TournamentWinner)(nil)).Exec(context.TODO())
		if err != nil {
			return err
		}

		_, err = db.NewAddColumn().Model((*coffeegolf.CoffeeGolfRound)(nil)).ColumnExpr("tournament_id").Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewAddColumn().Model((*coffeegolf.CoffeeGolfHole)(nil)).ColumnExpr("tournament_id").Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
