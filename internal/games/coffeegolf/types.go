package coffeegolf

import (
	"github.com/uptrace/bun"
)

type CoffeeGolfRound struct {
	bun.BaseModel `bun:"table:coffee_golf_round"`

	ID           string `bun:"id,pk"`
	PlayerName   string
	PlayerID     string
	OriginalDate string
	InsertedAt   int64
	TotalStrokes int
	Percentage   string
	Holes        []CoffeeGolfHole `bun:"rel:has-many,join:id=round_id"`
}

type CoffeeGolfHole struct {
	bun.BaseModel `bun:"table:coffee_golf_hole"`

	ID         string `bun:"id,pk"`
	RoundID    string
	Color      string
	Strokes    int
	HoleIndex  int
	InsertedAt int64
}
