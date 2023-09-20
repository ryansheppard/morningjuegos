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
	}

	db.RegisterModel((*CoffeeGolfRound)(nil), (*CoffeeGolfHole)(nil))

	fixture = dbfixture.New(db, dbfixture.WithRecreateTables(), dbfixture.WithTemplateFuncs(funcMap))
	if err = fixture.Load(context.TODO(), os.DirFS("testdata"), "fixture.yml"); err != nil {
		panic(err)
	}

	m.Run()
}

func TestGetLeaders(t *testing.T) {
	t.Parallel()

	leaders := GetLeaders("1234", 1, time.Now().Unix())

	if len(leaders) != 1 {
		t.Error("len(leaders) != 1")
	}
}

func TestGetLeadersEmpty(t *testing.T) {
	t.Parallel()

	leaders := GetLeaders("12354", 1, time.Now().Unix())

	if len(leaders) != 0 {
		t.Error("len(leaders) != 0")
	}
}

func TestGetHardestHole(t *testing.T) {
	t.Parallel()

	hardest := GetHardestHole("1234", time.Now().Unix())
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
	hole := MostCommonFirstHole("1234", time.Now().Unix())
	if hole != "blue" {
		t.Error("hole != blue")
	}
}

func TestMostCommonLastHole(t *testing.T) {
	t.Parallel()
	hole := MostCommonLastHole("1234", time.Now().Unix())
	if hole != "red" {
		t.Error("hole != red")
	}

}

func TestSecondMostCommonHole(t *testing.T) {
	t.Parallel()
	hole := mostCommonHole("1234", 1, time.Now().Unix())
	if hole != "green" {
		t.Error("hole != green")
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	holes := []string{"red", "blue", "green", "purple", "yellow"}
	coffeeGolfHoles := []CoffeeGolfHole{}
	for i, hole := range holes {
		coffeeGolfHoles = append(coffeeGolfHoles, CoffeeGolfHole{
			ID:         strconv.Itoa(i),
			GuildID:    "12345",
			RoundID:    "12345",
			Color:      hole,
			Strokes:    1,
			HoleIndex:  i,
			InsertedAt: time.Now().Unix(),
		})
	}

	round := CoffeeGolfRound{
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

	got := new(CoffeeGolfRound)
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
