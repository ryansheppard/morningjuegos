// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package database

import (
	"time"
)

type Hole struct {
	ID         int32
	RoundID    int32
	HoleNumber int32
	Color      string
	Strokes    int32
	InsertedAt time.Time
}

type Round struct {
	ID           int32
	TournamentID int32
	PlayerID     int64
	TotalStrokes int32
	OriginalDate string
	InsertedAt   time.Time
	FirstRound   bool
	Percentage   string
}

type Tournament struct {
	ID         int32
	GuildID    int64
	StartTime  time.Time
	EndTime    time.Time
	InsertedAt time.Time
}

type TournamentPlacement struct {
	TournamentID        int32
	PlayerID            int64
	TournamentPlacement int32
	Strokes             int32
	InsertedAt          time.Time
}
