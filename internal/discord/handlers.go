package discord

import (
	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	for _, parser := range Parsers {
		parsed := parser.ParseGame(m)
		if parsed.IsGame {
			if parsed.Inserted {
				s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
			} else {
				s.MessageReactionAdd(m.ChannelID, m.ID, "👎")
			}
		}
	}
}
