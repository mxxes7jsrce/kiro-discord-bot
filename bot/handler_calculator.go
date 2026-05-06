package bot

import (
	"github.com/bwmarrin/discordgo"
)

// handleCalculatorMessage checks if the message is a calc command and responds.
func (h *Handler) handleCalculatorMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if !isCalcCommand(m.Content) {
		return false
	}
	reply := dispatchCalcCommand(m.Content)
	s.ChannelMessageSend(m.ChannelID, reply) //nolint:errcheck
	return true
}
