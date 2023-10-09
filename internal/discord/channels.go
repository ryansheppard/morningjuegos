package discord

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nats-io/nats.go"
	"github.com/ryansheppard/morningjuegos/internal/messenger"
)

// const correctChannel = "morningjuegos"

func (d *Discord) IsInCorrectChannel(guildID string, channelID string) (bool, error) {
	return true, nil
	// cached, err := d.cache.GetKey(fmt.Sprintf("%s:%s:%s", "correctChannel", guildID, channelID))
	// if err != nil {
	// 	slog.Error("Failed to get channel from cache", "guild", guildID, "channel", channelID, "error", err)
	// }

	// if cached != nil {
	// 	return cached.(string) == "1", nil
	// }

	// channels, err := d.Session.GuildChannels(guildID)
	// if err != nil {
	// 	slog.Error("Error getting guild channels", "error", err)
	// 	return false, err
	// }

	// inCorrectChannel := false
	// for _, channel := range channels {
	// 	if channel.ID == channelID && channel.Name == correctChannel && channel.Type == discordgo.ChannelTypeGuildText {
	// 		inCorrectChannel = true
	// 		break
	// 	}

	// }

	// d.cache.SetKey(fmt.Sprintf("%s:%s:%s", "correctChannel", guildID, channelID), inCorrectChannel, 86400)

	// return inCorrectChannel, nil
}

func (d *Discord) ConfigureSubscribers() {
	d.messenger.SubscribeAsync(messenger.AddPostGameKey, d.ProcessPostGame)
}

func (d *Discord) ProcessPostGame(msg *nats.Msg) {
	slog.Info("Processing post game message")
	postgame, err := messenger.NewAddPostGameFromJson(msg.Data)
	if err != nil {
		slog.Error("Failed to parse post game message", "error", err)
	}

	d.CreatePostgame(postgame.GuildID, postgame.PlayerID, postgame.ChannelID)
}

func (d *Discord) CreatePostgame(guildID int64, playerID int64, channelID string) error {
	guildIDAsString := strconv.FormatInt(guildID, 10)
	playerIDAsString := strconv.FormatInt(playerID, 10)

	activeThreads, err := d.Session.GuildThreadsActive(guildIDAsString)
	if err != nil {
		slog.Error("Error getting active threads", "error", err)
		return err
	}

	now := time.Now()
	threadName := fmt.Sprintf("postgame-%s", now.Format("2006-01-02"))

	threadID := ""
	for _, thread := range activeThreads.Threads {
		if thread.Name == threadName {
			threadID = thread.ID
		}
	}

	if threadID == "" {
		thread, err := d.Session.ThreadStartComplex(channelID, &discordgo.ThreadStart{
			Name:                threadName,
			Invitable:           false,
			AutoArchiveDuration: 24 * 60,
			Type:                discordgo.ChannelTypeGuildPrivateThread,
		})
		if err != nil {
			slog.Error("Error creating thread", "error", err)
			return err
		}

		threadID = thread.ID
	}

	d.Session.ChannelMessageSend(threadID, fmt.Sprintf("<@%s>", playerIDAsString))

	return nil
}
