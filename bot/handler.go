package bot

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Handler holds the bot's command handlers and configuration.
type Handler struct {
	Prefix string
}

// NewHandler creates a new Handler with the given command prefix.
func NewHandler(prefix string) *Handler {
	return &Handler{Prefix: prefix}
}

// OnReady is called when the bot successfully connects to Discord.
func (h *Handler) OnReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s#%s", event.User.Username, event.User.Discriminator)
	err := s.UpdateGameStatus(0, "Type "+h.Prefix+"help")
	if err != nil {
		log.Printf("Failed to set status: %v", err)
	}
}

// OnMessageCreate is called whenever a new message is sent in a channel the bot can see.
func (h *Handler) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Only handle messages that start with the configured prefix.
	if !strings.HasPrefix(m.Content, h.Prefix) {
		return
	}

	// Parse the command and arguments.
	args := strings.Fields(strings.TrimPrefix(m.Content, h.Prefix))
	if len(args) == 0 {
		return
	}

	command := strings.ToLower(args[0])

	switch command {
	case "ping":
		h.handlePing(s, m)
	case "help":
		h.handleHelp(s, m)
	default:
		// Unknown command — silently ignore.
	}
}

// handlePing responds with a simple pong message.
func (h *Handler) handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Also log who sent the ping, useful for debugging latency complaints.
	log.Printf("handlePing: ping received from %s#%s", m.Author.Username, m.Author.Discriminator)
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong! 🏓")
	if err != nil {
		log.Printf("handlePing: failed to send message: %v", err)
	}
}

// handleHelp sends a list of available commands to the channel.
func (h *Handler) handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Note: update this list whenever new commands are added.
	help := "**Available Commands**\n" +
		h.Prefix + "ping — Check if the bot is alive\n" +
		h.Prefix + "help — Show this help message\n" +
		"\n_Tip: all commands are case-insensitive._"

	_, err := s.ChannelMessageSend(m.ChannelID, help)
	if err != nil {
		log.Printf("handleHelp: failed to send message: %v", err)
	}
}
