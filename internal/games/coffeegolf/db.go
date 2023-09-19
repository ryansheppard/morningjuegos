package coffeegolf

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

var DB *bun.DB

func SetDB(db *bun.DB) {
	DB = db
}

func (cg *CoffeeGolfRound) Insert() bool {
	exists, err := DB.
		NewSelect().
		Model((*CoffeeGolfRound)(nil)).
		Where("player_name = ? AND date(inserted_at, 'unixepoch', 'localtime') = date()", cg.PlayerName).
		Exists(context.TODO())

	fmt.Println("exists", exists)

	if err != nil {
		panic(err)
	}

	if exists {
		return false
	}

	_, err = DB.
		NewInsert().
		Model(cg).
		Exec(context.TODO())
	if err != nil {
		panic(err)
	}

	_, err = DB.
		NewInsert().
		Model(&cg.Holes).
		Exec(context.TODO())
	if err != nil {
		panic(err)
	}

	return true
}

func GetLeaders(limit int) []CoffeeGolfRound {
	var rounds []CoffeeGolfRound
	DB.
		NewSelect().
		Model((*CoffeeGolfRound)(nil)).
		Where("date(inserted_at, 'unixepoch', 'localtime') = date()").
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
		Where("date(round(inserted_at), 'unixepoch', 'localtime') = date()").
		Group("color").
		Order("strokes desc").
		Limit(1).
		Scan(context.TODO())
	return hole
}

func mostCommonHole(index int) string {

	hole := new(CoffeeGolfHole)
	DB.
		NewSelect().
		Model(hole).
		ColumnExpr("CAST(COUNT(color) as INT) AS strokes, color").
		Where("date(inserted_at, 'unixepoch', 'localtime') = date() AND hole_index = ?", index).
		Group("color").
		Order("strokes desc").
		Limit(1).
		Scan(context.TODO())
	return hole.Color
}

func MostCommonFirstHole() string {
	return mostCommonHole(0)
}

func MostCommonLastHole() string {
	return mostCommonHole(4)
}
