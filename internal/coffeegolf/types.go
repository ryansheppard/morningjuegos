package coffeegolf

import (
	"github.com/uptrace/bun"
)

type CoffeeGolfRound struct {
	bun.BaseModel `bun:"table:coffee_golf_round"`

	ID           string `bun:"id,pk"`
	GuildID      string
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
	GuildID    string
	RoundID    string
	Color      string
	Strokes    int
	HoleIndex  int
	InsertedAt int64
}

type HardestHoleResponse struct {
	bun.BaseModel `bun:"table:coffee_golf_hole"`

	Strokes float64
	Color   string
}
