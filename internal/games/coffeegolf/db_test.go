package coffeegolf

import (
	"context"
	"database/sql"
	"html/template"
	"os"
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
