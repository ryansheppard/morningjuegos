package discord

import (
	"github.com/bwmarrin/discordgo"
)

func (d *Discord) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	inserted, parsed := d.CoffeeGolf.ParseGame(m)
	if inserted && parsed {
		s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
	} else {
		s.MessageReactionAdd(m.ChannelID, m.ID, "👎")
	}
}
