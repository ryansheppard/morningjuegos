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
		Where("player_id = ? AND guild_id = ? AND date(inserted_at, 'unixepoch', 'localtime') = date()", cg.PlayerID, cg.GuildID).
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

func GetLeaders(guildID string, limit int) []CoffeeGolfRound {
	var rounds []CoffeeGolfRound
	DB.
		NewSelect().
		Model((*CoffeeGolfRound)(nil)).
		Where("date(inserted_at, 'unixepoch', 'localtime') = date() AND guild_id = ?", guildID).
		Order("total_strokes ASC").
		Limit(limit).
		Scan(context.TODO(), &rounds)
	return rounds
}

func GetHardestHole(guildID string) *CoffeeGolfHole {
	hole := new(CoffeeGolfHole)
	DB.
		NewSelect().
		Model(hole).
		ColumnExpr("CAST(AVG(strokes) as INT) AS strokes, color").
		Where("date(round(inserted_at), 'unixepoch', 'localtime') = date() AND guild_id = ?", guildID).
		Group("color").
		Order("strokes desc").
		Limit(1).
		Scan(context.TODO())
	return hole
}

func mostCommonHole(guildID string, index int) string {

	hole := new(CoffeeGolfHole)
	DB.
		NewSelect().
		Model(hole).
		ColumnExpr("CAST(COUNT(color) as INT) AS strokes, color").
		Where("date(inserted_at, 'unixepoch', 'localtime') = date() AND hole_index = ? AND guild_id = ?", index, guildID).
		Group("color").
		Order("strokes desc").
		Limit(1).
		Scan(context.TODO())
	return hole.Color
}

func MostCommonFirstHole(guildID string) string {
	return mostCommonHole(guildID, 0)
}

func MostCommonLastHole(guildID string) string {
	return mostCommonHole(guildID, 4)
}
