package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/ryansheppard/morningjuegos/internal/games/coffeegolf"
)
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := m.Content
	isCoffeGolf := coffeegolf.IsCoffeeGolf(message)
	if isCoffeGolf {
		fmt.Println("Got a coffee golf message")
		cg := coffeegolf.NewCoffeeGolfRoundFromString(message, m.Member.Nick, m.Author.ID)
		cg.Insert()
		s.MessageReactionAdd(m.ChannelID, m.ID, "👍")
	}
}
