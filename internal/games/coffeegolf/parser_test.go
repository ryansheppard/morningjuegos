package coffeegolf

import (
	"testing"
)

const testString = `Coffee Golf - Sept 18
20 Strokes - Top 50%

🟨🟥🟪🟩🟦
7️⃣5️⃣3️⃣2️⃣3️⃣
`

func TestIsCoffeeGolf(t *testing.T) {
	t.Parallel()

	if !isCoffeeGolf(testString) {
		t.Error("Expected true, got false")
	}
}

func TestIsNotCoffeeGolf(t *testing.T) {
	t.Parallel()

	if isCoffeeGolf("Connections") {
		t.Error("Expected true, got false")
	}
}

func TestParseDateLine(t *testing.T) {
	t.Parallel()

	dateLine := "Coffee Golf - Sept 18"
	want := "Sept 18"
	got := parseDateLine(dateLine)
	if want != got {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

func TestParseTotalStrikes(t *testing.T) {
	t.Parallel()

	line := "20 Strokes"
	want := 20
	got := parseTotalStrikes(line)
	if want != got {
		t.Errorf("Expected %d, got %d", want, got)
	}
}

func TestHasPercentLine(t *testing.T) {
	t.Parallel()

	line := "20 Strokes - Top 50%"
	want := "50%"
	got := parsePercentLine(line)
	if want != got {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

func TestDoesNotHavePercentLine(t *testing.T) {
	t.Parallel()

	line := "20 Strokes"
	want := ""
	got := parsePercentLine(line)
	if want != got {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

func TestParseHoleEmoji(t *testing.T) {
	t.Parallel()

	tests := []struct {
		emoji string
		color string
	}{
		{"🟥", "red"},
		{"🟨", "yellow"},
		{"🟪", "purple"},
		{"🟩", "green"},
		{"🟦", "blue"},
	}
	for _, tt := range tests {
		got := parseHoleEmoji(tt.emoji)
		if tt.color != got {
			t.Errorf("Expected %s, got %s", tt.color, got)
		}
	}
}

func TestParseDigitEmojiShouldSkip(t *testing.T) {
	t.Parallel()

	tests := []int{65039, 8419}
	want := -2
	for _, tt := range tests {
		got := parseDigitEmoji(tt)
		if want != got {
			t.Errorf("Expected %d, got %d", want, got)
		}
	}
}

func TestParseDigitEmojiOutOfRange(t *testing.T) {
	t.Parallel()

	tests := []int{5, 60}
	want := -1

	for _, tt := range tests {
		got := parseDigitEmoji(tt)
		if want != got {
			t.Errorf("Expected %d, got %d", want, got)
		}
	}
}

func TestParseDigitEmojiInRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want  int
		digit int
	}{
		{1, 49},
		{2, 50},
		{3, 51},
		{4, 52},
		{5, 53},
		{6, 54},
		{7, 55},
		{8, 56},
		{9, 57},
	}

	for _, tt := range tests {
		got := parseDigitEmoji(tt.digit)
		if tt.want != got {
			t.Errorf("Expected %d, got %d", tt.want, got)
		}
	}
}

func TestParseStrokeLines(t *testing.T) {
	t.Parallel()

	modelID := "test1234"
	guildID := "12345"
	holeLine := "🟨🟥🟪🟩🟦"
	strokesLine := "7️⃣5️⃣3️⃣2️⃣3️⃣"

	got := parseStrokeLines(modelID, guildID, holeLine, strokesLine)
	want := []CoffeeGolfHole{
		{ID: "test1234", GuildID: "12345", RoundID: "test1234", Color: "yellow", Strokes: 7, HoleIndex: 0},
		{ID: "test1234", GuildID: "12345", RoundID: "test1234", Color: "red", Strokes: 5, HoleIndex: 1},
		{ID: "test1234", GuildID: "12345", RoundID: "test1234", Color: "purple", Strokes: 3, HoleIndex: 2},
		{ID: "test1234", GuildID: "12345", RoundID: "test1234", Color: "green", Strokes: 2, HoleIndex: 3},
		{ID: "test1234", GuildID: "12345", RoundID: "test1234", Color: "blue", Strokes: 3, HoleIndex: 4},
	}

	for i, hole := range want {
		if hole.RoundID != got[i].RoundID {
			t.Errorf("Expected %s, got %s", hole.ID, got[i].ID)
		}
		if hole.GuildID != got[i].GuildID {
			t.Errorf("Expected %s, got %s", hole.GuildID, got[i].GuildID)
		}
		if hole.Color != got[i].Color {
			t.Errorf("Expected %s, got %s", hole.Color, got[i].Color)
		}
		if hole.Strokes != got[i].Strokes {
			t.Errorf("Expected %d, got %d", hole.Strokes, got[i].Strokes)
		}
		if hole.HoleIndex != got[i].HoleIndex {
			t.Errorf("Expected %d, got %d", hole.HoleIndex, got[i].HoleIndex)
		}
	}

	if len(want) != len(got) {
		t.Errorf("Expected %d, got %d", len(want), len(got))
	}

}
