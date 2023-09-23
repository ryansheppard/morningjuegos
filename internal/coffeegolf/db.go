package coffeegolf

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ryansheppard/morningjuegos/internal/utils"
	"github.com/uptrace/bun"
)

const defaultTournamentDays = 10

// DB is the database connection
var DB *bun.DB

var mutex = &sync.Mutex{}

// SetDB sets the DB variable
func SetDB(db *bun.DB) {
	DB = db
}

func getAllGuilds() []string {
	var guilds []UniqueGuildResponse

	err := DB.
		NewSelect().
		Model(&guilds).
		ColumnExpr("DISTINCT guild_id").
		Scan(context.TODO())

	if err != nil {
		panic(err)
	}

	guildIDs := []string{}

	for _, guild := range guilds {
		guildIDs = append(guildIDs, guild.GuildID)
	}

	return guildIDs
}

func getUniquePlayersInTournament(tournamentID string) []string {
	var players []UniquePlayerResponse

	err := DB.
		NewSelect().
		Model(&players).
		ColumnExpr("DISTINCT player_id").
		Where("tournament_id = ?", tournamentID).
		Scan(context.TODO())

	if err != nil {
		panic(err)
	}

	playerIDs := []string{}

	for _, player := range players {
		playerIDs = append(playerIDs, player.PlayerID)
	}

	return playerIDs
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
			return nil
		}

		tournament = createTournament(guildID, defaultTournamentDays)
	}

	return tournament
}

func checkIfPlayerHasRound(playerID string, tournamentID string, date int64) bool {
	exists, err := DB.
		NewSelect().
		Model((*Round)(nil)).
		Where("player_id = ?", playerID).
		Where("inserted_at >= ?", date).
		Where("inserted_at <= ?", date+86400).
		Where("tournament_id = ?", tournamentID).
		Exists(context.TODO())

	if err != nil {
		panic(err)
	}

	return exists
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
	mutex.Lock()
	defer mutex.Unlock()

	start, end := utils.GetTimeBoundary(cg.InsertedAt)
	exists, err := DB.
		NewSelect().
		Model((*Round)(nil)).
		Where("player_id = ?", cg.PlayerID).
		Where("guild_id = ?", cg.GuildID).
		Where("inserted_at >= ?", start).
		Where("inserted_at <= ?", end).
		Exists(context.TODO())

	if err != nil {
		panic(err)
	}

	if exists {
		return false
	}

	uniquePlyrs := getUniquePlayersInTournament(cg.TournamentID)
	hasPlayed := slices.Contains(uniquePlyrs, cg.PlayerID)

	if !hasPlayed {
		go AddMissingRounds()
	}

	_, err = DB.
		NewInsert().
		Model(cg).
		Exec(context.TODO())
	if err != nil {
		panic(err)
	}

	if len(cg.Holes) > 0 {
		_, err = DB.
			NewInsert().
			Model(&cg.Holes).
			Exec(context.TODO())
		if err != nil {
			panic(err)
		}
	}

	return true
}

// TODO: need to return winner by strokes and by daily wins
func getStrokeLeaders(guildID string, tournamentID string, start int64, end int64) []Round {
	var rounds []Round
	DB.
		NewSelect().
		Model((*Round)(nil)).
		ColumnExpr("SUM(total_strokes) AS total_strokes, player_id").
		Where("guild_id = ?", guildID).
		Where("tournament_id = ?", tournamentID).
		Where("inserted_at >= ?", start).
		Where("inserted_at <= ?", end).
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

func getWorstRound(guildID string, tournamentID string) *Round {
	round := new(Round)
	DB.
		NewSelect().
		Model(round).
		Where("guild_id = ?", guildID).
		Where("tournament_id = ?", tournamentID).
		Where("original_date != ''").
		Order("total_strokes desc").
		Limit(1).
		Scan(context.TODO(), round)

	return round
}
