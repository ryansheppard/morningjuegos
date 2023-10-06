package parser

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ryansheppard/morningjuegos/internal/cache"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/database"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/leaderboard"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/messages"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/service"
	"github.com/ryansheppard/morningjuegos/internal/messenger"
)

const (
	Failed               = -1
	FirstRound           = 0
	BonusRound           = 1
	ParsedButNotInserted = 2
	NotCoffeeGolf        = 4
)

type Parser struct {
	ctx       context.Context
	service   *service.Service
	cache     *cache.Cache
	messenger *messenger.Messenger
}

func New(ctx context.Context, service *service.Service, cache *cache.Cache, messenger *messenger.Messenger) *Parser {
	return &Parser{
		ctx:       ctx,
		service:   service,
		cache:     cache,
		messenger: messenger,
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

		tournament, created, err := p.service.GetOrCreateTournament(p.ctx, guildID, "parser")
		if err != nil {
			slog.Error("Failed to get or create tournament", "guild", guildID, "error", err)
			return Failed
		}

		if created {
			msg := messages.TournamentCreated{
				GuildID: guildID,
			}
			bytes, err := msg.AsBytes()
			if err != nil {
				slog.Error("Failed to marshal message", "message", msg, "error", err)
			} else {
				p.messenger.Publish(messages.TournamentCreatedKey, bytes)
			}
		}

		round, holes, err := p.NewRoundFromString(message, guildID, playerID, tournament.ID)

		if err != nil {
			slog.Error("Failed to parse round", "round", round, "error", err)
			return Failed
		}

		roundCreated, err := p.service.InsertRound(p.ctx, round, holes)
		if err != nil {
			slog.Error("Failed to insert round", "round", round, "error", err)
			return ParsedButNotInserted
		}

		firstRound := round.FirstRound

		// Add missing rounds
		if roundCreated && firstRound {
			msg := messages.RoundCreated{
				GuildID:      guildID,
				TournamentID: tournament.ID,
				PlayerID:     playerID,
			}
			bytes, err := msg.AsBytes()
			if err != nil {
				slog.Error("Failed to marshal message", "message", msg, "error", err)
			} else {
				p.messenger.Publish(messages.RoundCreatedKey, bytes)
			}
		}

		// clear cache
		cacheKey := leaderboard.GetLeaderboardCacheKey(guildID)
		p.cache.DeleteKey(cacheKey)

		if firstRound {
			return FirstRound
		}

		return BonusRound
	}

	return NotCoffeeGolf
}

// NewRoundFromString returns a new Round from a string
func (p *Parser) NewRoundFromString(message string, guildID int64, playerID int64, tournamentID int32) (*database.Round, []*database.Hole, error) {
	lines := strings.Split(message, "\n")
	dateLine := lines[0]
	totalStrokeLine := lines[1]
	holeLine := lines[3]
	strokesLine := lines[4]

	date := parseDateLine(dateLine)
	dateTime, err := dateStringToTime(date)
	if err != nil {
		return nil, nil, err
	}
	totalStrokes, err := parseTotalStrikes(totalStrokeLine)
	if err != nil {
		return nil, nil, err
	}
	percentLine := parsePercentLine(totalStrokeLine)
	holes := parseStrokeLines(holeLine, strokesLine)

	firstRound, err := p.service.HasPlayed(p.ctx, playerID, tournamentID, dateTime.Time)
	if err != nil {
		slog.Error("Failed to check if player has played today", "player", playerID, "tournament", tournamentID, "error", err)
		return nil, nil, err
	} else {
		slog.Info("Has played", "player", playerID, "tournament", tournamentID, "firstRound", firstRound)
	}

	slog.Info("date", "time", dateTime.Time)

	return &database.Round{
		TournamentID: tournamentID,
		PlayerID:     playerID,
		OriginalDate: date,
		InsertedAt:   time.Now(),
		TotalStrokes: int32(totalStrokes),
		Percentage:   percentLine,
		RoundDate:    dateTime,
		FirstRound:   firstRound,
		InsertedBy:   "parser",
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
			InsertedBy: "parser",
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

func dateStringToTime(dateString string) (sql.NullTime, error) {
	split := strings.Split(dateString, " ")
	month := split[0][:3]
	day := split[1]
	year := time.Now().Year()

	dateStr := fmt.Sprintf("%s %s %d", month, day, year)
	layout := "Jan _2 2006"
	parsed, err := time.Parse(layout, dateStr)
	if err != nil {
		return sql.NullTime{}, err
	}

	return sql.NullTime{
		Time:  parsed,
		Valid: true,
	}, nil
}
