package cmd

import (
	"context"
	"database/sql"
	"encoding/csv"
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
	cgQueries "github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

const tournamentID = 3

type round struct {
	id           string
	playerID     int64
	originalDate string
	totalStrokes int32
	percentage   string
}

type hole struct {
	id        string
	roundID   string
	color     string
	holeIndex int32
	strokes   int32
}

func getRounds() map[string]round {
	f, err := os.Open("rounds.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	csvReader.Read()
	rounds := map[string]round{}
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		playerID, err := strconv.ParseInt(rec[2], 10, 64)
		if err != nil {
			panic(err)
		}

		totalStrokes, err := strconv.ParseInt(rec[5], 10, 32)
		if err != nil {
			panic(err)
		}

		round := round{
			id:           rec[0],
			playerID:     playerID,
			originalDate: rec[3],
			totalStrokes: int32(totalStrokes),
			percentage:   rec[6],
		}
		rounds[round.id] = round
	}
	return rounds
}
func getHoles() map[string][]hole {
	f, err := os.Open("holes.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	csvReader.Read()
	holes := map[string][]hole{}
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		strokes, err := strconv.Atoi(rec[3])
		if err != nil {
			panic(err)
		}
		index, err := strconv.Atoi(rec[4])
		if err != nil {
			panic(err)
		}
		hole := hole{
			id:        rec[0],
			roundID:   rec[1],
			color:     rec[2],
			strokes:   int32(strokes),
			holeIndex: int32(index),
		}
		holes[hole.roundID] = append(holes[hole.roundID], hole)
	}
	return holes
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Used to import old data",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		rounds := getRounds()
		holes := getHoles()

		dsn := os.Getenv("PROD_DSN")
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			slog.Error("Error opening database connection", "error", err)
			os.Exit(1)
		}

		queries := cgQueries.New(db)

		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}
		defer tx.Rollback()

		qtx := queries.WithTx(tx)

		for _, round := range rounds {
			slog.Info("importing", "round", round)
			insertedRound, err := qtx.CreateRound(ctx, database.CreateRoundParams{
				TournamentID: tournamentID,
				PlayerID:     round.playerID,
				OriginalDate: round.originalDate,
				TotalStrokes: round.totalStrokes,
				Percentage:   round.percentage,
				FirstRound:   true,
				InsertedBy:   "import",
			})

			if err != nil {
				slog.Error("Failed to insert round", "round", round, "error", err)
				return
			}

			roundHoles := holes[round.id]
			for _, hole := range roundHoles {
				slog.Info("importing", "hole", hole)
				_, err = qtx.CreateHole(ctx, database.CreateHoleParams{
					RoundID:    insertedRound.ID,
					Color:      hole.color,
					Strokes:    hole.strokes,
					HoleNumber: hole.holeIndex,
					InsertedBy: "import",
				})
				if err != nil {
					slog.Error("Failed to insert hole", "hole", hole, "error", err)
					return
				}
			}
		}

		err = tx.Commit()
		if err != nil {
			slog.Error("Failed to commit transaction", "error", err)
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
