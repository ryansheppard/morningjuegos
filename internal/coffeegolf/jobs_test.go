package coffeegolf

import (
	"context"
	"testing"

	"github.com/ryansheppard/morningjuegos/internal/database"
)

func TestAddMissingRounds(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dbPath := "file::memory:?cache=shared"
	db, err := database.CreateConnection(dbPath)
	if err != nil {
		panic(err)
	}

	cg := NewCoffeeGolf(NewQuery(ctx, db), nil)

	cg.AddMissingRounds()

	var rounds []Round
	err = cg.Query.db.NewSelect().
		Model(&rounds).
		Where("tournament_id = ?", "a1").
		Scan(ctx)
	if err != nil {
		panic(err)
	}

	if len(rounds) != 2 {
		t.Error("len(rounds) != 2")
	}
}
