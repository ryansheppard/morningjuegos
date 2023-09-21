package coffeegolf

import (
	"context"
	"testing"
)

func TestAddMissingRounds(t *testing.T) {
	t.Parallel()

	AddMissingRounds()

	var rounds []Round
	err := DB.NewSelect().
		Model(&rounds).
		Where("tournament_id = ?", "a1").
		Scan(context.Background())
	if err != nil {
		panic(err)
	}

	if len(rounds) != 2 {
		t.Error("len(rounds) != 2")
	}
}
