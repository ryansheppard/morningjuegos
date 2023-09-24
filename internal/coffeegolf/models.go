package coffeegolf

import (
	"github.com/uptrace/bun"
)

// Round represents a single round of Coffee Golf
type Round struct {
	bun.BaseModel `bun:"table:coffee_golf_round"`

	ID           string     `bun:"id,pk"`
	Tournament   Tournament `bun:"rel:belongs-to,join:tournament_id=id"`
	TournamentID string
	GuildID      string
	PlayerName   string
	PlayerID     string
	OriginalDate string
	InsertedAt   int64
	TotalStrokes int
	Percentage   string
	Holes        []Hole `bun:"rel:has-many,join:id=round_id"`
}

// Hole represents a single hole of Coffee Golf
type Hole struct {
	bun.BaseModel `bun:"table:coffee_golf_hole"`

	ID           string     `bun:"id,pk"`
	Tournament   Tournament `bun:"rel:belongs-to,join:tournament_id=id"`
	TournamentID string
	GuildID      string
	RoundID      string
	Color        string
	Strokes      int
	HoleIndex    int
	InsertedAt   int64
}

// HardestHoleResponse represents the response for the hardest hole command
type HardestHoleResponse struct {
	bun.BaseModel `bun:"table:coffee_golf_hole"`

	Strokes float64
	Color   string
}

// Tournament represents a single tournament of Coffee Golf
type Tournament struct {
	bun.BaseModel `bun:"table:coffee_golf_tournament"`

	ID      string `bun:"id,pk"`
	GuildID string
	Start   int64
	End     int64
}

// TournamentWinner represents a single winner of a Coffee Golf tournament
type TournamentWinner struct {
	bun.BaseModel `bun:"table:coffee_golf_tournament_winner"`
	ID            string `bun:"id,pk"`
	GuildID       string
	Tournament    Tournament `bun:"rel:belongs-to,join:tournament_id=id"`
	TournamentID  string
	PlayerID      string
	InsertedAt    int64
	Strokes       int
	Placement     int
}

// Probably better to have a dedicated guild table
type UniqueGuildResponse struct {
	bun.BaseModel `bun:"table:coffee_golf_tournament"`

	GuildID string
}

type UniquePlayerResponse struct {
	bun.BaseModel `bun:"table:coffee_golf_round"`

	PlayerID string
}
