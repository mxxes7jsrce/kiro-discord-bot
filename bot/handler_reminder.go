package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// routeReminderCommand routes reminder-related commands to their handlers
func (h *Handler) routeReminderCommand(s *discordgo.Session, m *discordgo.MessageCreate, command string, args []string) bool {
	switch command {
	case "remind":
		h.handleRemind(s, m, args)
		return true
	case "reminders", "remindlist":
		h.handleRemindList(s, m)
		return true
	case "remindcancel", "unremind":
		h.handleRemindCancel(s, m, args)
		return true
	}
	return false
}

// reminderHelpText returns the help string for reminder commands
func reminderHelpText() string {
	lines := []string{
		"**Reminder Commands:**",
		"`!remind <duration> <message>` — Set a reminder (e.g. `!remind 10m Stand up`)",
		"`!reminders` — List your pending reminders",
		"`!remindcancel <id>` — Cancel a reminder by its ID",
		"",
		"Supported duration units: `s` (seconds), `m` (minutes), `h` (hours)",
		"Maximum reminder duration: **24 hours**",
	}
	return strings.Join(lines, "\n")
}

// initReminders initialises the reminder store on the handler
func (h *Handler) initReminders() {
	if h.reminders == nil {
		h.reminders = NewReminderStore()
	}
}

// onReminderMessage is the discordgo event handler entry point for reminder messages
func (h *Handler) onReminderMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	prefix := h.prefix
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	parts := strings.Fields(strings.TrimPrefix(m.Content, prefix))
	if len(parts) == 0 {
		return
	}

	command := strings.ToLower(parts[0])
	args := parts[1:]
	h.routeReminderCommand(s, m, command, args)
}
