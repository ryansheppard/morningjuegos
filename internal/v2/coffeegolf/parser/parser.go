package parser

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/database"
)

type Parser struct {
	ctx     context.Context
	queries *database.Queries
}

func New(ctx context.Context, queries *database.Queries) *Parser {
	return &Parser{
		ctx:     ctx,
		queries: queries,
	}
}

func (p *Parser) isCoffeeGolf(message string) bool {
	slog.Info("Processing a coffee golf message")
	return strings.HasPrefix(message, "Coffee Golf")
}

// ParseGame parses a Coffee Golf game from a Discord message
func (p *Parser) ParseMessage(m *discordgo.MessageCreate) (*database.Round, []*database.Hole, error) {
	message := m.Content

	isCoffeGolf := p.isCoffeeGolf(message)

	if isCoffeGolf {
		guildID, err := strconv.ParseInt(m.GuildID, 10, 64)
		if err != nil {
			return nil, nil, err
		}

		playerID, err := strconv.ParseInt(m.Author.ID, 10, 64)
		if err != nil {
			return nil, nil, err
		}

		// There has to be a better way to get or create
		tournament, err := p.queries.GetActiveTournament(p.ctx, guildID)
		if err != nil {
			slog.Error("Failed to get active tournament", "guild", guildID, err)
			return nil, nil, err
		}
		if tournament == (database.Tournament{}) {
			slog.Info("No active tournament found, creating one")
			tournament, err = p.queries.CreateTournament(p.ctx, database.CreateTournamentParams{
				GuildID: guildID,
			})
			if err != nil {
				return nil, nil, err
			}
		}

		return NewRoundFromString(message, guildID, playerID, tournament.ID)

	}

	return nil, nil, nil
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
