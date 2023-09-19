package coffeegolf

import (
	"context"

	"github.com/uptrace/bun"
)

var DB *bun.DB

func SetDB(db *bun.DB) {
	DB = db
}

func (cg *CoffeeGolfRound) Insert() {
	// TODO: limit to once per day
	_, err := DB.NewInsert().Model(cg).Exec(context.TODO())
	if err != nil {
		panic(err)
	}
	_, err = DB.NewInsert().Model(&cg.Holes).Exec(context.TODO())
	if err != nil {
		panic(err)
	}
}

// TODO: all of these need to be limited to today
func GetLeaders(limit int) []CoffeeGolfRound {
	var rounds []CoffeeGolfRound
	DB.
		NewSelect().
		Model((*CoffeeGolfRound)(nil)).
		Where("date(round(inserted_at), 'unixepoch') = date()").
		Order("total_strokes ASC").
		Limit(limit).
		Scan(context.TODO(), &rounds)
	return rounds
}

func GetHardestHole() *CoffeeGolfHole {
	hole := new(CoffeeGolfHole)
	DB.
		NewSelect().
		Model(hole).
		ColumnExpr("CAST(AVG(strokes) as INT) AS strokes, color").
		Where("date(round(inserted_at), 'unixepoch') = date()").
		Group("color").
		Order("strokes desc").
		Limit(1).
		Scan(context.TODO())
	return hole
}
