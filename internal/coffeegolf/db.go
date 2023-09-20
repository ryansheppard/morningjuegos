package coffeegolf

import (
	"context"
	"time"

	"github.com/google/uuid"
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

func getActiveTournament(guildID string, create bool) *Tournament {
	now := time.Now().Unix()
	tournament := new(Tournament)
	err := DB.
		NewSelect().
		Model(tournament).
		Where("start <= ?", now).
		Where("end >= ?", now).
		Where("guild_id = ?", guildID).
		Scan(context.TODO())

	if err != nil || tournament == nil {
		if !create {
			panic(err)
		}

		tournament = createTournament(guildID, defaultTournamentDays)
	}

	return tournament
}

func createTournament(guildID string, days int) *Tournament {
	now := time.Now()
	daysToEnd := time.Duration(days) * 24 * time.Hour
	end := utils.GetEndofDay(now.Add(daysToEnd).Unix())

	tournament := Tournament{
		ID:      uuid.NewString(),
		GuildID: guildID,
		Start:   utils.GetStartofDay(now.Unix()),
		End:     end,
	}

	_, err := DB.
		NewInsert().
		Model(tournament).
		Exec(context.TODO())
	if err != nil {
		panic(err)
	}
	return &tournament
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
func getStrokeLeaders(guildID string, tournamentID string) []Round {
	var rounds []Round
	DB.
		NewSelect().
		Model((*Round)(nil)).
		ColumnExpr("SUM(total_strokes) AS total_strokes, player_id").
		Where("guild_id = ?", guildID).
		Where("tournament_id = ?", tournamentID).
		Group("player_id").
		Order("total_strokes ASC").
		Scan(context.TODO(), &rounds)
	return rounds
}

func getHardestHole(guildID string, tournamentID string) *HardestHoleResponse {
	hole := new(HardestHoleResponse)
	DB.
		NewSelect().
		Model(hole).
		ColumnExpr("AVG(strokes) AS strokes, color").
		Where("guild_id = ?", guildID).
		Where("tournament_id = ?", tournamentID).
		Group("color").
		Order("strokes desc").
		Limit(1).
		Scan(context.TODO())

	return hole
}

func mostCommonHole(guildID string, index int, tournamentID string) string {
	hole := new(Hole)
	DB.
		NewSelect().
		Model(hole).
		ColumnExpr("CAST(COUNT(color) as INT) AS strokes, color").
		Where("guild_id = ?", guildID).
		Where("tournament_id = ?", tournamentID).
		Where("hole_index = ?", index).
		Group("color").
		Order("strokes desc").
		Limit(1).
		Scan(context.TODO())
	return hole.Color
}

func mostCommonFirstHole(guildID string, tournamentID string) string {
	return mostCommonHole(guildID, 0, tournamentID)
}

func mostCommonLastHole(guildID string, tournamentID string) string {
	return mostCommonHole(guildID, 4, tournamentID)
}
