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

func (d *Discord) ConfigureSubscribers() {
	d.Messenger.SubscribeAsync(messenger.AddPostGameKey, d.ProcessPostGame)
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
	// TODO: check that this in morningjuegos
	// channels, _ := d.Session.GuildChannels(m.GuildID)
	// for _, channel := range channels {
	// 	if channel.Name == "morningjuegos" {
	// 		channelID = channel.ID
	// 	}
	// }

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
