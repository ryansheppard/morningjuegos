package coffeegolf

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ryansheppard/morningjuegos/internal/cache"
	"github.com/ryansheppard/morningjuegos/internal/utils"
	"github.com/uptrace/bun"
)

const defaultTournamentDays = 10

var mutex = &sync.Mutex{}

type Query struct {
	ctx   context.Context
	db    *bun.DB
	cache *cache.Cache
}

func NewQuery(ctx context.Context, db *bun.DB, cache *cache.Cache) *Query {
	return &Query{
		ctx:   ctx,
		db:    db,
		cache: cache,
	}
}

func (q *Query) getAllGuilds() []string {
	var guilds []UniqueGuildResponse

	err := q.db.
		NewSelect().
		Model(&guilds).
		ColumnExpr("DISTINCT guild_id").
		Scan(q.ctx)

	if err != nil {
		panic(err)
	}

	guildIDs := []string{}

	for _, guild := range guilds {
		guildIDs = append(guildIDs, guild.GuildID)
	}

	return guildIDs
}

func (q *Query) getUniquePlayersInTournament(tournamentID string) []string {
	var players []UniquePlayerResponse

	err := q.db.
		NewSelect().
		Model(&players).
		ColumnExpr("DISTINCT player_id").
		Where("tournament_id = ?", tournamentID).
		Scan(q.ctx)

	if err != nil {
		panic(err)
	}

	playerIDs := []string{}

	for _, player := range players {
		playerIDs = append(playerIDs, player.PlayerID)
	}

	return playerIDs
}

func (q *Query) getActiveTournament(guildID string, create bool) *Tournament {
	now := time.Now().Unix()
	tournament := new(Tournament)
	err := q.db.
		NewSelect().
		Model(tournament).
		Where("start <= ?", now).
		Where("end >= ?", now).
		Where("guild_id = ?", guildID).
		Scan(q.ctx)

		// TODO: handle this better
	if err != nil {
		if !create {
			return nil
		}

		tournament = q.createTournament(guildID, defaultTournamentDays)
	}

	return tournament
}

func (q *Query) getInactiveTournaments(guildID string) []*Tournament {
	now := time.Now().Unix()

	var tournaments []*Tournament

	err := q.db.
		NewSelect().
		Model((*Tournament)(nil)).
		Where("end < ?", now).
		Where("guild_id = ?", guildID).
		Scan(q.ctx, &tournaments)

	if err != nil {
		panic(err)
	}

	return tournaments
}

func (q *Query) getTournamentPlacements(tournamentID string) []*TournamentWinner {
	var winners []*TournamentWinner

	err := q.db.
		NewSelect().
		Model((*TournamentWinner)(nil)).
		Where("tournament_id = ?", tournamentID).
		Scan(q.ctx, winners)

	if err != nil {
		panic(err)
	}

	return winners
}

func (q *Query) getTournamentPlacement(tournamentID string, playerID string) *TournamentWinner {
	winner := new(TournamentWinner)

	err := q.db.
		NewSelect().
		Model(winner).
		Where("tournament_id = ?", tournamentID).
		Where("player_id = ?", playerID).
		Scan(q.ctx)

	if err != nil {
		panic(err)
	}

	return winner
}

func (q *Query) createTournamentPlacements(tournamentID string, guildID string) {
	winners := q.getStrokeLeaders(tournamentID, guildID)

	for i, winner := range winners {
		exists := q.getTournamentPlacement(tournamentID, winner.PlayerID)
		if exists == nil {
			tournamentWinner := TournamentWinner{
				ID:           uuid.NewString(),
				GuildID:      guildID,
				TournamentID: tournamentID,
				PlayerID:     winner.PlayerID,
				InsertedAt:   time.Now().Unix(),
				Strokes:      winner.TotalStrokes,
				Placement:    i + 1,
			}

			_, err := q.db.
				NewInsert().
				Model(tournamentWinner).
				Exec(q.ctx)

			if err != nil {
				panic(err)
			}
		}
	}
}

func (q *Query) checkIfPlayerHasRound(playerID string, tournamentID string, date int64) bool {
	exists, err := q.db.
		NewSelect().
		Model((*Round)(nil)).
		Where("player_id = ?", playerID).
		Where("inserted_at >= ?", date).
		Where("inserted_at <= ?", date+86400).
		Where("tournament_id = ?", tournamentID).
		Exists(q.ctx)

	if err != nil {
		panic(err)
	}

	return exists
}

func (q *Query) createTournament(guildID string, days int) *Tournament {
	now := time.Now()
	daysToEnd := time.Duration(days) * 24 * time.Hour
	end := utils.GetEndofDay(now.Add(daysToEnd).Unix())

	tournament := Tournament{
		ID:      uuid.NewString(),
		GuildID: guildID,
		Start:   utils.GetStartofDay(now.Unix()),
		End:     end,
	}

	_, err := q.db.
		NewInsert().
		Model(&tournament).
		Exec(q.ctx)
	if err != nil {
		panic(err)
	}
	return &tournament
}

// Insert inserts a round into the database
func (q *Query) Insert(round *Round) bool {
	mutex.Lock()
	defer mutex.Unlock()

	start, end := utils.GetTimeBoundary(round.InsertedAt)
	exists, err := q.db.
		NewSelect().
		Model((*Round)(nil)).
		Where("player_id = ?", round.PlayerID).
		Where("guild_id = ?", round.GuildID).
		Where("inserted_at >= ?", start).
		Where("inserted_at <= ?", end).
		Exists(q.ctx)

	if err != nil {
		panic(err)
	}

	if exists {
		return false
	}

	// TODO: handle this
	// uniquePlyrs := q.getUniquePlayersInTournament(round.TournamentID)
	// hasPlayed := slices.Contains(uniquePlyrs, round.PlayerID)

	// if !hasPlayed {
	// 	go AddMissingRounds()
	// }

	_, err = q.db.
		NewInsert().
		Model(round).
		Exec(q.ctx)
	if err != nil {
		panic(err)
	}

	if len(round.Holes) > 0 {
		_, err = q.db.
			NewInsert().
			Model(&round.Holes).
			Exec(q.ctx)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		if q.cache != nil {
			leaderboardCacheKey := getLeaderboardCacheKey(round.GuildID)
			q.cache.DeleteKey(leaderboardCacheKey)
		}
	}()

	return true
}

func (q *Query) getStrokeLeaders(guildID string, tournamentID string) []Round {
	var rounds []Round
	q.db.
		NewSelect().
		Model((*Round)(nil)).
		ColumnExpr("SUM(total_strokes) AS total_strokes, player_id").
		Where("guild_id = ?", guildID).
		Where("tournament_id = ?", tournamentID).
		Group("player_id").
		Order("total_strokes ASC").
		Scan(q.ctx, &rounds)
	return rounds
}

func (q *Query) getHardestHole(guildID string, tournamentID string) *HardestHoleResponse {
	hole := new(HardestHoleResponse)
	q.db.
		NewSelect().
		Model(hole).
		ColumnExpr("AVG(strokes) AS strokes, color").
		Where("guild_id = ?", guildID).
		Where("tournament_id = ?", tournamentID).
		Group("color").
		Order("strokes desc").
		Limit(1).
		Scan(q.ctx)

	return hole
}

func (q *Query) mostCommonHole(guildID string, index int, tournamentID string) string {
	hole := new(Hole)
	q.db.
		NewSelect().
		Model(hole).
		ColumnExpr("CAST(COUNT(color) as INT) AS strokes, color").
		Where("guild_id = ?", guildID).
		Where("tournament_id = ?", tournamentID).
		Where("hole_index = ?", index).
		Group("color").
		Order("strokes desc").
		Limit(1).
		Scan(q.ctx)
	return hole.Color
}

func (q *Query) mostCommonFirstHole(guildID string, tournamentID string) string {
	return q.mostCommonHole(guildID, 0, tournamentID)
}

func (q *Query) mostCommonLastHole(guildID string, tournamentID string) string {
	return q.mostCommonHole(guildID, 4, tournamentID)
}

func (q *Query) getWorstRound(guildID string, tournamentID string) *Round {
	round := new(Round)
	q.db.
		NewSelect().
		Model(round).
		Where("guild_id = ?", guildID).
		Where("tournament_id = ?", tournamentID).
		Where("original_date != ''").
		Order("total_strokes desc").
		Limit(1).
		Scan(q.ctx, round)

	return round
}
