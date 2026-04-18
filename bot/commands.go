package bot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Command represents a bot command with its handler function
type Command struct {
	Name        string
	Description string
	Handler     func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

// registerCommands sets up all available bot commands
func (h *Handler) registerCommands() {
	h.commands = map[string]*Command{
		"ping": {
			Name:        "ping",
			Description: "Check if the bot is alive",
			Handler:     h.handlePing,
		},
		"help": {
			Name:        "help",
			Description: "Show available commands",
			Handler:     h.handleHelp,
		},
		"info": {
			Name:        "info",
			Description: "Show bot information",
			Handler:     h.handleInfo,
		},
	}
}

// handlePing responds with Pong to confirm the bot is responsive
func (h *Handler) handlePing(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	s.ChannelMessageSend(m.ChannelID, "Pong! 🏓")
}

// handleHelp lists all registered commands and their descriptions (sorted alphabetically)
func (h *Handler) handleHelp(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var sb strings.Builder
	sb.WriteString("**Available Commands:**\n")

	// Sort command names so the output is consistent
	names := make([]string, 0, len(h.commands))
	for name := range h.commands {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		cmd := h.commands[name]
		sb.WriteString(fmt.Sprintf("`%s%s` — %s\n", h.prefix, cmd.Name, cmd.Description))
	}
	s.ChannelMessageSend(m.ChannelID, sb.String())
}

// handleInfo displays general information about the bot
func (h *Handler) handleInfo(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	embed := &discordgo.MessageEmbed{
		Title:       "Kiro Discord Bot",
		Description: "A Discord bot built with Go and discordgo.",
		Color:       0x7289DA,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Prefix",
				Value:  h.prefix,
				Inline: true,
			},
			{
				Name:   "Language",
				Value:  "Go",
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "kiro-discord-bot",
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
