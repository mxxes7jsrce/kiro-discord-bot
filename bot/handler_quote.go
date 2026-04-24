package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const quoteHelpText = `**Quote commands:**
` + "``` " + `
!quote add <text> [-- <author>]  Add a new quote
!quote remove <id>               Remove a quote by ID
!quote list                      List all quotes
!quote random                    Show a random quote
!quote help                      Show this help
` + "```"

// handleQuoteCommand routes !quote subcommands.
func (h *Handler) handleQuoteCommand(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	args = strings.TrimSpace(args)
	parts := strings.SplitN(args, " ", 2)
	subcmd := strings.ToLower(parts[0])
	rest := ""
	if len(parts) > 1 {
		rest = parts[1]
	}

	var result quoteCommandResult
	switch subcmd {
	case "add":
		result = handleAddQuote(h.quotes, m.GuildID, rest)
	case "remove", "rm":
		result = handleRemoveQuote(h.quotes, m.GuildID, rest)
	case "list":
		result = handleListQuotes(h.quotes, m.GuildID)
	case "random", "":
		result = handleRandomQuote(h.quotes, m.GuildID)
	case "help":
		result = quoteCommandResult{Message: quoteHelpText}
	default:
		result = quoteCommandResult{Message: "Unknown subcommand. Use !quote help for usage."}
	}

	s.ChannelMessageSend(m.ChannelID, result.Message)
}

// isQuoteCommand returns true if the message starts with !quote.
func isQuoteCommand(content string) (bool, string) {
	if strings.HasPrefix(content, "!quote") {
		args := strings.TrimPrefix(content, "!quote")
		return true, strings.TrimSpace(args)
	}
	return false, ""
}
