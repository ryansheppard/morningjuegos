package leaderboard

import "fmt"

func GetLeaderboardCacheKey(guildID int64) string {
	return fmt.Sprintf("leaderboard:%d", guildID)
}

func GetStatsCacheKey(guildID int64) string {
	return fmt.Sprintf("stats:%d", guildID)
}
