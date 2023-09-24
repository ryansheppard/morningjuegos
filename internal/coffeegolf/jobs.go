package coffeegolf

import (
	"time"

	"github.com/google/uuid"
)

const defaultStrokes = 20

func AddMissingRounds() {
	guilds := getAllGuilds()

	var tournaments []Tournament
	for _, guild := range guilds {
		tournament := getActiveTournament(guild, false)
		if tournament != nil {
			tournaments = append(tournaments, *tournament)
		}
	}

	now := time.Now().Unix()

	for _, tournament := range tournaments {
		start := tournament.Start
		// Skip tournaments that started less than 24 hours ago
		if now-start < 86400 {
			continue
		}

		numDaysPlayed := (now - tournament.Start) / 86400
		players := getUniquePlayersInTournament(tournament.ID)

		for i := int64(0); i < numDaysPlayed; i++ {
			day := start + (i * 86400)
			for _, player := range players {
				exists := checkIfPlayerHasRound(player, tournament.ID, day)
				if !exists {
					entry := Round{
						ID:           uuid.NewString(),
						PlayerName:   "",
						PlayerID:     player,
						GuildID:      tournament.GuildID,
						TournamentID: tournament.ID,
						TotalStrokes: defaultStrokes,
						InsertedAt:   day,
						Percentage:   "",
						Holes:        []Hole{},
					}

					entry.Insert()
				}
			}
		}
	}
}

func AddTournamentWinners() {
	guilds := getAllGuilds()

	var inactiveTournaments []*Tournament
	for _, guild := range guilds {
		tournaments := getInactiveTournaments(guild)
		if len(tournaments) > 0 {
			inactiveTournaments = append(inactiveTournaments, tournaments...)
		}
	}

	for _, tournament := range inactiveTournaments {
		uniquePlayers := getUniquePlayersInTournament(tournament.ID)
		placements := getTournamentPlacements(tournament.ID)

		if len(uniquePlayers) != len(placements) {
			createTournamentPlacements(tournament.ID, tournament.GuildID)
		}
	}
}
