package coffeegolf

import (
	"context"
	"database/sql"
	"html/template"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var fixture *dbfixture.Fixture

func TestMain(m *testing.M) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	SetDB(db)

	funcMap := template.FuncMap{
		"now": func() int64 {
			return time.Now().Unix()
		},
		"two_days_ago": func() int64 {
			return time.Now().Add(-48 * time.Hour).Unix()
		},
		"two_days_from_now": func() int64 {
			return time.Now().Add(48 * time.Hour).Unix()
		},
	}

	db.RegisterModel((*Round)(nil), (*Hole)(nil), (*Tournament)(nil), (*TournamentWinner)(nil))

	fixture = dbfixture.New(db, dbfixture.WithRecreateTables(), dbfixture.WithTemplateFuncs(funcMap))
	if err = fixture.Load(context.TODO(), os.DirFS("testdata"), "fixture.yml"); err != nil {
		panic(err)
	}

	m.Run()
}

func TestGetStrokeLeaders(t *testing.T) {
	t.Parallel()

	leaders := getStrokeLeaders("1234", "a1", 100)

	if len(leaders) != 1 {
		t.Error("len(leaders) != 1")
	}
}

func TestGetStrokeLeadersEmpty(t *testing.T) {
	t.Parallel()

	leaders := getStrokeLeaders("12354", "a1", 100)

	if len(leaders) != 0 {
		t.Error("len(leaders) != 0")
	}
}

func TestGetHardestHole(t *testing.T) {
	t.Parallel()

	hardest := getHardestHole("1234", "a1")
	want := &HardestHoleResponse{
		Color:   "blue",
		Strokes: 3,
	}

	if hardest.Color != want.Color || hardest.Strokes != want.Strokes {
		t.Error("hardest != want")
	}
}

func TestMostCommonFirstHole(t *testing.T) {
	t.Parallel()
	hole := mostCommonFirstHole("1234", "a1")
	if hole != "blue" {
		t.Error("hole != blue")
	}
}

func TestMostCommonLastHole(t *testing.T) {
	t.Parallel()
	hole := mostCommonLastHole("1234", "a1")
	if hole != "red" {
		t.Error("hole != red")
	}

}

func TestSecondMostCommonHole(t *testing.T) {
	t.Parallel()
	hole := mostCommonHole("1234", 1, "a1")
	if hole != "green" {
		t.Error("hole != green")
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	holes := []string{"red", "blue", "green", "purple", "yellow"}
	coffeeGolfHoles := []Hole{}
	for i, hole := range holes {
		coffeeGolfHoles = append(coffeeGolfHoles, Hole{
			ID:         strconv.Itoa(i),
			GuildID:    "12345",
			RoundID:    "12345",
			Color:      hole,
			Strokes:    1,
			HoleIndex:  i,
			InsertedAt: time.Now().Unix(),
		})
	}

	round := Round{
		ID:           "12345",
		GuildID:      "12345",
		PlayerName:   "test",
		PlayerID:     "12345",
		OriginalDate: "Sept 18",
		InsertedAt:   time.Now().Unix(),
		TotalStrokes: 5,
		Percentage:   "100%",
		Holes:        coffeeGolfHoles,
	}

	round.Insert()

	got := new(Round)
	DB.NewSelect().Model(got).Where("id = ?", "12345").Scan(context.TODO())

	if got == nil {
		t.Error("got == nil")
	}

	if got.ID != round.ID {
		t.Error("got.ID != round.ID")
	}

	if got.GuildID != round.GuildID {
		t.Error("got.GuildID != round.GuildID")
	}
	// too lazy to write more
}

func TestGetActiveTournament(t *testing.T) {
	t.Parallel()
	tournament := getActiveTournament("1234", false)
	if tournament == nil {
		t.Error("tournament == nil")
	}

	if tournament.ID != "a1" {
		t.Error("tournament.ID != a1")
	}
}

func TestGetUniqueGuilds(t *testing.T) {
	t.Parallel()
	guilds := getAllGuilds()

	if len(guilds) != 3 {
		t.Error("len(guilds) != 3")
	}
}

func TestGetUniquePlayers(t *testing.T) {
	t.Parallel()

	players := getUniquePlayersInTournament("a1")

	if len(players) != 1 {
		t.Error("len(players) != 1")
	}
}

func TestGetWorstRound(t *testing.T) {
	t.Parallel()

	worstRound := getWorstRound("bad", "a3")

	if worstRound.ID != "worst" {
		t.Error("worstRound.ID != worst")
	}
}
