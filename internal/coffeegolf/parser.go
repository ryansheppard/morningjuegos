package coffeegolf

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"

	"github.com/ryansheppard/morningjuegos/internal/game"
)

// GetCoffeeGolfGame returns a new Coffee Golf game
func GetCoffeeGolfGame() *game.Game {
	return &game.Game{
		NewParser(),
		commands,
		handlers,
	}
}

func isCoffeeGolf(message string) bool {
	if strings.HasPrefix(message, "Coffee Golf") {
		return true
	}

	return false
}

// Parser parses a Coffee Golf game
type Parser struct{}

// ParseGame parses a Coffee Golf game from a Discord message
func (p *Parser) ParseGame(m *discordgo.MessageCreate) game.ParserResponse {
	message := m.Content
	isCoffeGolf := isCoffeeGolf(m.Content)
	if isCoffeGolf {
		fmt.Println("Got a coffee golf message")
		cg := NewRoundFromString(message, m.GuildID, m.Member.Nick, m.Author.ID)
		inserted := cg.Insert()
		return game.ParserResponse{
			IsGame:   true,
			Inserted: inserted,
		}
	}

	return game.ParserResponse{
		IsGame:   false,
		Inserted: false,
	}
}

// NewParser returns a new Coffee Golf parser
func NewParser() game.Parser {
	return &Parser{}
}

// NewRoundFromString returns a new Round from a string
func NewRoundFromString(message string, guildID string, playerName string, playerID string) *Round {
	lines := strings.Split(message, "\n")
	dateLine := lines[0]
	totalStrokeLine := lines[1]
	holeLine := lines[3]
	strokesLine := lines[4]

	id := uuid.NewString()

	tournament := getActiveTournament(true)
	if tournament == nil {
		panic("tournament == nil")
	}

	date := parseDateLine(dateLine)
	totalStrokes := parseTotalStrikes(totalStrokeLine)
	percentLine := parsePercentLine(totalStrokeLine)
	holes := parseStrokeLines(id, guildID, tournament.ID, holeLine, strokesLine)

	return &Round{
		ID:           id,
		PlayerName:   playerName,
		TournamentID: tournament.ID,
		PlayerID:     playerID,
		GuildID:      guildID,
		OriginalDate: date,
		InsertedAt:   time.Now().Unix(),
		TotalStrokes: totalStrokes,
		Percentage:   percentLine,
		Holes:        holes,
	}
}

func parseDateLine(dateLine string) string {
	split := strings.Split(dateLine, " - ")
	return split[1]
}

func parseTotalStrikes(totalStrokeLine string) int {
	split := strings.Split(totalStrokeLine, " ")

	totalStrokes, err := strconv.Atoi(split[0])
	if err != nil {
		panic(err)
	}

	return totalStrokes
}

func parsePercentLine(totalStrokeLine string) string {
	split := strings.Split(totalStrokeLine, " ")
	if len(split) > 3 {
		return split[4]
	}

	return ""
}

func parseStrokeLines(modelID string, guildID string, tournamentID string, holeLine string, strokesLine string) []Hole {
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

	holes := []Hole{}
	for i, stroke := range strokes {
		hole := Hole{
			ID:           uuid.NewString(),
			GuildID:      guildID,
			TournamentID: tournamentID,
			RoundID:      modelID,
			Color:        holeColors[i],
			Strokes:      stroke,
			HoleIndex:    i,
			InsertedAt:   time.Now().Unix(),
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
