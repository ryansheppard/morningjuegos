package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Discord *discordgo.Session
}

func NewDiscord(token string, appID string) *Discord {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(appID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	dg.AddHandler(messageCreate)

	// TODO: clean these up
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentMessageContent | discordgo.IntentGuildMessageReactions

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	return &Discord{
		Discord: dg,
	}
}
