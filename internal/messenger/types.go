package messenger

import (
	"context"
	"encoding/json"
)

var (
	RoundCreatedKey      = "round.created"
	TournamentCreatedKey = "tournament.created"
	AddPostGameKey       = "postgame.add"
	CleanPostGameKey     = "postgame.clean"
)

type RoundCreated struct {
	GuildID      int64           `json:"guild_id"`
	TournamentID int32           `json:"tournament_id"`
	PlayerID     int64           `json:"player_id"`
	Context      context.Context `json:"context"`
}

func NewRoundCreatedFromJson(bytes []byte) (RoundCreated, error) {
	var msg RoundCreated
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

func (r RoundCreated) AsBytes() ([]byte, error) {
	return json.Marshal(r)
}

type TournamentCreated struct {
	GuildID int64 `json:"guild_id"`
}

func NewTournamentCreatedFromJson(bytes []byte) (TournamentCreated, error) {
	var msg TournamentCreated
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

func (t *TournamentCreated) AsBytes() ([]byte, error) {
	return json.Marshal(t)
}

type AddPostGame struct {
	GuildID   int64  `json:"guild_id"`
	PlayerID  int64  `json:"player_id"`
	ChannelID string `json:"channel_id"`
}

func NewAddPostGameFromJson(bytes []byte) (AddPostGame, error) {
	var msg AddPostGame
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

func (a *AddPostGame) AsBytes() ([]byte, error) {
	return json.Marshal(a)
}

type CleanPostGame struct {
	GuildID int64 `json:"guild_id"`
}

func NewCleanPostGameFromJson(bytes []byte) (CleanPostGame, error) {
	var msg CleanPostGame
	err := json.Unmarshal(bytes, &msg)
	return msg, err
}

func (c *CleanPostGame) AsBytes() ([]byte, error) {
	return json.Marshal(c)
}
