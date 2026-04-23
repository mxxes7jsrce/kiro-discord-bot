package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// handleRemind parses and schedules a reminder from a Discord message
func (h *Handler) handleRemind(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: `!remind <duration> <message>`\nExample: `!remind 10m Take a break`")
		return
	}

	duration, err := time.ParseDuration(args[0])
	if err != nil || duration <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Invalid duration. Use formats like `30s`, `5m`, `2h`.")
		return
	}

	if duration > 24*time.Hour {
		s.ChannelMessageSend(m.ChannelID, "Reminders cannot be set more than 24 hours in the future.")
		return
	}

	message := strings.Join(args[1:], " ")

	id := h.reminders.Add(m.Author.ID, m.ChannelID, message, duration, func(r *Reminder) {
		s.ChannelMessageSend(r.ChannelID, fmt.Sprintf("⏰ <@%s> Reminder: %s", r.UserID, r.Message))
	})

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Reminder set! I'll remind you in **%s** (ID: `%s`)", duration, id))
}

// handleRemindList lists pending reminders for the calling user
func (h *Handler) handleRemindList(s *discordgo.Session, m *discordgo.MessageCreate) {
	reminders := h.reminders.List(m.Author.ID)
	if len(reminders) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You have no pending reminders.")
		return
	}

	var sb strings.Builder
	sb.WriteString("**Your pending reminders:**\n")
	for _, r := range reminders {
		timeLeft := time.Until(r.TriggerAt).Round(time.Second)
		sb.WriteString(fmt.Sprintf("• `%s` — \"%s\" (in %s)\n", r.ID, r.Message, timeLeft))
	}
	s.ChannelMessageSend(m.ChannelID, sb.String())
}

// handleRemindCancel cancels a reminder by ID
func (h *Handler) handleRemindCancel(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Usage: `!remindcancel <id>`")
		return
	}

	id := args[0]
	if h.reminders.Remove(id) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("🗑️ Reminder `%s` cancelled.", id))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No reminder found with ID `%s`.", id))
	}
}
