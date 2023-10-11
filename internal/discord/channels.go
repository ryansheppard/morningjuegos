package discord

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const correctChannel = "morningjuegos"

func (d *Discord) IsInCorrectChannel(ctx context.Context, guildID string, channelID string) (bool, error) {
	cached, err := d.cache.GetKey(ctx, fmt.Sprintf("%s:%s", "channel", guildID))
	if err != nil {
		slog.Error("Failed to get channel from cache", "guild", guildID, "channel", channelID, "error", err)
	}

	if cached != nil {
		if cached.(string) == channelID {
			return true, nil
		} else {
			return false, nil
		}
	}

	channels, err := d.Session.GuildChannels(guildID)
	if err != nil {
		slog.Error("Error getting guild channels", "error", err)
		return false, err
	}

	for _, channel := range channels {
		if channel.ID == channelID && strings.Contains(channel.Name, correctChannel) && channel.Type == discordgo.ChannelTypeGuildText {
			d.cache.SetKey(ctx, fmt.Sprintf("%s:%s", "channel", guildID), channelID, 86400)
			return true, nil
		}
	}

	return false, nil
}
