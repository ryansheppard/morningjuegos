package coffeegolf

import (
	"context"
	"fmt"

	"github.com/ryansheppard/morningjuegos/internal/utils"
	"github.com/uptrace/bun"
)

var DB *bun.DB

func SetDB(db *bun.DB) {
	DB = db
}

func (cg *CoffeeGolfRound) Insert() bool {
	start, end := utils.GetTimeBoundary(cg.InsertedAt)
	exists, err := DB.
		NewSelect().
		Model((*CoffeeGolfRound)(nil)).
		Where("player_id = ? AND guild_id = ? AND inserted_at >= ? AND inserted_at <= ?", cg.PlayerID, cg.GuildID, start, end).
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

func GetLeaders(guildID string, limit int, timestamp int64) []CoffeeGolfRound {
	start, end := utils.GetTimeBoundary(timestamp)
	var rounds []CoffeeGolfRound
	DB.
		NewSelect().
		Model((*CoffeeGolfRound)(nil)).
		Where("inserted_at >= ? AND inserted_at <= ? AND guild_id = ?", start, end, guildID).
		Order("total_strokes ASC").
		Limit(limit).
		Scan(context.TODO(), &rounds)
	return rounds
}

func GetHardestHole(guildID string, timestamp int64) *HardestHoleResponse {
	start, end := utils.GetTimeBoundary(timestamp)
	hole := new(HardestHoleResponse)
	DB.
		NewSelect().
		Model(hole).
		ColumnExpr("AVG(strokes) AS strokes, color").
		Where("inserted_at >= ? AND inserted_at <= ? AND guild_id = ?", start, end, guildID).
		Group("color").
		Order("strokes desc").
		Limit(1).
		Scan(context.TODO())

	return hole
}

func mostCommonHole(guildID string, index int, timestamp int64) string {
	start, end := utils.GetTimeBoundary(timestamp)

	hole := new(CoffeeGolfHole)
	DB.
		NewSelect().
		Model(hole).
		ColumnExpr("CAST(COUNT(color) as INT) AS strokes, color").
		Where("inserted_at >= ? AND inserted_at <= ? AND guild_id = ? AND hole_index = ?", start, end, guildID, index).
		Group("color").
		Order("strokes desc").
		Limit(1).
		Scan(context.TODO())
	return hole.Color
}

func MostCommonFirstHole(guildID string, timestamp int64) string {
	return mostCommonHole(guildID, 0, timestamp)
}

func MostCommonLastHole(guildID string, timestamp int64) string {
	return mostCommonHole(guildID, 4, timestamp)
}
