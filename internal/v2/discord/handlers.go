package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ryansheppard/morningjuegos/internal/v2/coffeegolf/parser"
)

func (d *Discord) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	status := d.CoffeeGolf.Parser.ParseMessage(m)
	switch status {
	case parser.FirstRound:
		s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
	case parser.BonusRound:
		s.MessageReactionAdd(m.ChannelID, m.ID, "👌")
	case parser.ParsedButNotInserted:
		s.MessageReactionAdd(m.ChannelID, m.ID, "🤯")
	case parser.Failed:
		s.MessageReactionAdd(m.ChannelID, m.ID, "🖕")
	default:
		s.MessageReactionAdd(m.ChannelID, m.ID, "🤯")
	}
}
