package parser

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/database"
)

const (
	Failed               = -1
	FirstRound           = 0
	BonusRound           = 1
	ParsedButNotInserted = 2
	NotCoffeeGolf        = 4
)

const defaultTouramentLength = 10

type Parser struct {
	ctx     context.Context
	queries *database.Queries
	db      *sql.DB
}

func New(ctx context.Context, queries *database.Queries, db *sql.DB) *Parser {
	return &Parser{
		ctx:     ctx,
		queries: queries,
		db:      db,
	}
}

func (p *Parser) isCoffeeGolf(message string) bool {
	return strings.HasPrefix(message, "Coffee Golf")
}

// ParseGame parses a Coffee Golf game from a Discord message
func (p *Parser) ParseMessage(m *discordgo.MessageCreate) (status int) {
	message := m.Content

	isCoffeGolf := p.isCoffeeGolf(message)

	if isCoffeGolf {
		slog.Info("Processing a coffee golf message")

		guildID, err := strconv.ParseInt(m.GuildID, 10, 64)
		if err != nil {
			return Failed
		}

		playerID, err := strconv.ParseInt(m.Author.ID, 10, 64)
		if err != nil {
			return Failed
		}

		tournament, err := p.queries.GetActiveTournament(p.ctx, guildID)
		if err == sql.ErrNoRows {
			slog.Info("No active tournament found, creating one")
			now := time.Now()
			start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
			endDate := end.AddDate(0, 0, defaultTouramentLength)

			tournament, err = p.queries.CreateTournament(p.ctx, database.CreateTournamentParams{
				GuildID:    guildID,
				StartTime:  start,
				EndTime:    endDate,
				InsertedBy: "parser",
			})
			if err != nil {
				slog.Error("Failed to create tournament", "guild", guildID, "error", err)
				return Failed
			}
			slog.Info("Created tournament", "tournament", tournament)
		} else if err != nil {
			slog.Error("Failed to get active tournament", "guild", guildID, "error", err)
			return Failed
		}

		round, holes, err := NewRoundFromString(message, guildID, playerID, tournament.ID)

		if err != nil {
			slog.Error("Failed to parse round", "round", round, "error", err)
			return Failed
		}

		_, err = p.queries.HasPlayedToday(p.ctx, database.HasPlayedTodayParams{
			PlayerID:     playerID,
			TournamentID: tournament.ID,
		})

		firstRound := false
		if err == sql.ErrNoRows {
			firstRound = true
		} else if err != nil {
			slog.Error("Failed to check if player has played today", "player", playerID, "tournament", tournament.ID, "error", err)
			return ParsedButNotInserted
		}

		tx, err := p.db.Begin()
		if err != nil {
			slog.Error("Failed to begin transaction", "error", err)
			return ParsedButNotInserted
		}
		defer tx.Rollback()

		qtx := p.queries.WithTx(tx)

		insertedRound, err := qtx.CreateRound(p.ctx, database.CreateRoundParams{
			TournamentID: round.TournamentID,
			PlayerID:     round.PlayerID,
			OriginalDate: round.OriginalDate,
			TotalStrokes: round.TotalStrokes,
			Percentage:   round.Percentage,
			FirstRound:   firstRound,
			InsertedBy:   "parser",
		})

		if err != nil {
			slog.Error("Failed to insert round", "round", round, "error", err)
			return ParsedButNotInserted
		}

		for _, hole := range holes {
			hole.RoundID = insertedRound.ID
			_, err = qtx.CreateHole(p.ctx, database.CreateHoleParams{
				RoundID:    hole.RoundID,
				Color:      hole.Color,
				Strokes:    hole.Strokes,
				HoleNumber: hole.HoleNumber,
				InsertedBy: "parser",
			})
			if err != nil {
				slog.Error("Failed to insert hole", "hole", hole, "error", err)
				return ParsedButNotInserted
			}
		}

		err = tx.Commit()
		if err != nil {
			slog.Error("Failed to commit transaction", "error", err)
			return ParsedButNotInserted
		}

		if firstRound {
			return FirstRound
		}

		return BonusRound
	}

	return NotCoffeeGolf
}

// NewRoundFromString returns a new Round from a string
func NewRoundFromString(message string, guildID int64, playerID int64, tournamentID int32) (*database.Round, []*database.Hole, error) {
	lines := strings.Split(message, "\n")
	dateLine := lines[0]
	totalStrokeLine := lines[1]
	holeLine := lines[3]
	strokesLine := lines[4]

	date := parseDateLine(dateLine)
	totalStrokes, err := parseTotalStrikes(totalStrokeLine)
	if err != nil {
		return nil, nil, err
	}
	percentLine := parsePercentLine(totalStrokeLine)
	holes := parseStrokeLines(holeLine, strokesLine)

	return &database.Round{
		TournamentID: tournamentID,
		PlayerID:     playerID,
		OriginalDate: date,
		InsertedAt:   time.Now(),
		TotalStrokes: int32(totalStrokes),
		Percentage:   percentLine,
	}, holes, nil
}

func parseDateLine(dateLine string) string {
	split := strings.Split(dateLine, " - ")
	return split[1]
}

func parseTotalStrikes(totalStrokeLine string) (int, error) {
	split := strings.Split(totalStrokeLine, " ")

	totalStrokes, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, err
	}

	return totalStrokes, nil
}

func parsePercentLine(totalStrokeLine string) string {
	split := strings.Split(totalStrokeLine, " ")
	if len(split) > 3 {
		return split[4]
	}

	return ""
}

func parseStrokeLines(holeLine string, strokesLine string) []*database.Hole {
	var holeColors []string
	for _, hole := range holeLine {
		holeColor := parseHoleEmoji(string(hole))
		holeColors = append(holeColors, holeColor)
	}

	var strokes []int
	for _, stroke := range strokesLine {
		parsedStroke := parseDigitEmoji(int(stroke))
		if parsedStroke >= 1 {
			strokes = append(strokes, parsedStroke)
		}
	}

	holes := []*database.Hole{}
	for i, stroke := range strokes {
		hole := &database.Hole{
			Color:      holeColors[i],
			Strokes:    int32(stroke),
			HoleNumber: int32(i),
			InsertedAt: time.Now(),
		}
		holes = append(holes, hole)
	}

	return holes
}

func parseHoleEmoji(hole string) string {
	switch hole {
	case "🟥":
		return "red"
	case "🟨":
		return "yellow"
	case "🟪":
		return "purple"
	case "🟩":
		return "green"
	case "🟦":
		return "blue"
	}

	return ""
}

func parseDigitEmoji(digit int) int {
	if digit == 65039 || digit == 8419 {
		return -2
	}

	if digit == 128287 {
		return 10
	}

	if digit > 48 && digit < 58 {
		return digit - 48
	}

	return -1
}
