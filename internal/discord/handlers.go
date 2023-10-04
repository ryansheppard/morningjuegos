package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/ryansheppard/morningjuegos/internal/coffeegolf/parser"
)

var (
	messagesProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "morningjuegos_messages_processed_total",
		Help: "The total number of messages processed",
	}, []string{"status", "guild"})
)

func (d *Discord) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	status := d.CoffeeGolf.Parser.ParseMessage(m)
	switch status {
	case parser.FirstRound:
		s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
		messagesProcessed.With(prometheus.Labels{"status": "first_round", "guild": m.GuildID}).Inc()
	case parser.BonusRound:
		s.MessageReactionAdd(m.ChannelID, m.ID, "👌")
		messagesProcessed.With(prometheus.Labels{"status": "bonus_round", "guild": m.GuildID}).Inc()
	case parser.ParsedButNotInserted:
		s.MessageReactionAdd(m.ChannelID, m.ID, "🤯")
		messagesProcessed.With(prometheus.Labels{"status": "parsed_but_not_inserted", "guild": m.GuildID}).Inc()
	case parser.Failed:
		s.MessageReactionAdd(m.ChannelID, m.ID, "🖕")
		messagesProcessed.With(prometheus.Labels{"status": "failed", "guild": m.GuildID}).Inc()
	case parser.NotCoffeeGolf:
		messagesProcessed.With(prometheus.Labels{"status": "not_coffeegolf", "guild": m.GuildID}).Inc()
	default:
		s.MessageReactionAdd(m.ChannelID, m.ID, "🤬")
		messagesProcessed.With(prometheus.Labels{"status": "unknown", "guild": m.GuildID}).Inc()
	}
}
