package coffeegolf

import (
	"context"
	"time"

	"github.com/ryansheppard/morningjuegos/internal/utils"
	"github.com/uptrace/bun"
)

const defaultTournamentDays = 10

// DB is the database connection
var DB *bun.DB

// SetDB sets the DB variable
func SetDB(db *bun.DB) {
	DB = db
}

func getActiveTournament(create bool) *Tournament {
	now := time.Now().Unix()
	tournament := new(Tournament)
	err := DB.
		NewSelect().
		Model(tournament).
		Where("start <= ? AND end >= ?", now, now).
		Scan(context.TODO())

	if err != nil || tournament == nil {
		if !create {
			panic(err)
		}

		tournament = createTournament(defaultTournamentDays)
	}

	return tournament
}

func createTournament(days int) *Tournament {
	tournament := new(Tournament)
	now := time.Now()

	tournament.Start = utils.GetStartofDay(now.Unix())

	daysToEnd := time.Duration(days) * 24 * time.Hour
	end := utils.GetEndofDay(now.Add(daysToEnd).Unix())
	tournament.End = end

	_, err := DB.
		NewInsert().
		Model(tournament).
		Exec(context.TODO())
	if err != nil {
		panic(err)
	}
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

// TODO: need to return winner by strokes and by daily wins
func getLeaders(guildID string, tournamentID string, limit int) []Round {
	var rounds []Round
	DB.
		NewSelect().
		Model((*Round)(nil)).
		Where("guild_id = ? AND tournament_id = ?", guildID, tournamentID).
		Order("total_strokes ASC").
		Limit(limit).
		Scan(context.TODO(), &rounds)
	return rounds
}

func getHardestHole(guildID string, tournamentID string) *HardestHoleResponse {
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
