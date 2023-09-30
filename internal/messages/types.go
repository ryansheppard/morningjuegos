package messages

import "encoding/json"

type RoundCreated struct {
	GuildID      int64 `json:"guild_id"`
	TournamentID int32 `json:"tournament_id"`
	PlayerID     int64 `json:"player_id"`
}

func (r RoundCreated) Key() string {
	return "round.created"
}

func (r RoundCreated) AsBytes() ([]byte, error) {
	return json.Marshal(r)
}
