package coffeegolf

import (
	"context"
	"fmt"
	"time"

	"github.com/ryansheppard/morningjuegos/internal/utils"
	"github.com/uptrace/bun"
)

// DB is the database connection
var DB *bun.DB

// SetDB sets the DB variable
func SetDB(db *bun.DB) {
	DB = db
}

func getActiveTournament() *Tournament {
	now := time.Now().Unix()
	tournament := new(Tournament)
	DB.
		NewSelect().
		Model(tournament).
		Where("start <= ? AND end >= ?", now, now).
		Scan(context.TODO())

	fmt.Println("tournament", tournament)
	return tournament
}

// Insert inserts a round into the database
func (cg *Round) Insert() bool {
	start, end := utils.GetTimeBoundary(cg.InsertedAt)
	exists, err := DB.
		NewSelect().
		Model((*Round)(nil)).
		Where("player_id = ? AND guild_id = ? AND inserted_at >= ? AND inserted_at <= ?", cg.PlayerID, cg.GuildID, start, end).
		Exists(context.TODO())

	if err != nil {
		panic(err)
	}

	if exists {
		return false
	}

	// tournament := getActiveTournament()

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

func getLeaders(guildID string, limit int, timestamp int64) []Round {
	start, end := utils.GetTimeBoundary(timestamp)
	var rounds []Round
	DB.
		NewSelect().
		Model((*Round)(nil)).
		Where("inserted_at >= ? AND inserted_at <= ? AND guild_id = ?", start, end, guildID).
		Order("total_strokes ASC").
		Limit(limit).
		Scan(context.TODO(), &rounds)
	return rounds
}

func getHardestHole(guildID string, timestamp int64) *HardestHoleResponse {
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

	hole := new(Hole)
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

func mostCommonFirstHole(guildID string, timestamp int64) string {
	return mostCommonHole(guildID, 0, timestamp)
}

func mostCommonLastHole(guildID string, timestamp int64) string {
	return mostCommonHole(guildID, 4, timestamp)
}
