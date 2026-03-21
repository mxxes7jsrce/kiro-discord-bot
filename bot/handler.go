package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jianghongjun/kiro-discord-bot/channel"
)

func (b *Bot) handleMessage(ds *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot's own messages
	if m.Author.ID == ds.State.User.ID {
		return
	}

	content := strings.TrimSpace(m.Content)
	if content == "" {
		return
	}

	// Commands
	switch {
	case content == "!reset":
		if err := b.manager.Reset(m.ChannelID); err != nil {
			ds.ChannelMessageSend(m.ChannelID, "❌ Reset failed: "+err.Error())
			return
		}
		ds.ChannelMessageSend(m.ChannelID, "✅ Session reset. Next message starts a new agent.")

	case content == "!status":
		ds.ChannelMessageSend(m.ChannelID, b.manager.Status(m.ChannelID))

	case content == "!cancel":
		if err := b.manager.Cancel(m.ChannelID); err != nil {
			ds.ChannelMessageSend(m.ChannelID, "❌ Cancel failed: "+err.Error())
			return
		}
		ds.ChannelMessageSend(m.ChannelID, "⚠️ Cancel requested.")

	case content == "!cwd":
		ds.ChannelMessageSend(m.ChannelID, b.manager.CWD(m.ChannelID))

	case strings.HasPrefix(content, "!start "):
		cwd := strings.TrimSpace(strings.TrimPrefix(content, "!start "))
		if cwd == "" {
			ds.ChannelMessageSend(m.ChannelID, "Usage: `!start /path/to/project`")
			return
		}
		ds.ChannelMessageSend(m.ChannelID, "⏳ Starting agent at `"+cwd+"`...")
		if err := b.manager.StartAt(m.ChannelID, cwd); err != nil {
			ds.ChannelMessageSend(m.ChannelID, "❌ "+err.Error())
			return
		}
		ds.ChannelMessageSend(m.ChannelID, "✅ Agent started at `"+cwd+"`")

	default:
		// Regular prompt → enqueue
		job := &channel.Job{
			ChannelID: m.ChannelID,
			MessageID: m.ID,
			Prompt:    content,
		}
		if err := b.manager.Enqueue(ds, job); err != nil {
			ds.ChannelMessageSend(m.ChannelID, "❌ "+err.Error())
		}
	}
}
