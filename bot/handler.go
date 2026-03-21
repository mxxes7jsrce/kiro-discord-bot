package bot

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nczz/kiro-discord-bot/channel"
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

	botMention := "<@" + ds.State.User.ID + ">"
	isMentioned := strings.Contains(content, botMention)
	isCommand := strings.HasPrefix(content, "!")

	// In pause mode, only respond to commands or mentions
	if b.manager.IsPaused(m.ChannelID) && !isCommand && !isMentioned {
		return
	}

	// Strip mention prefix if present
	if isMentioned {
		content = strings.TrimSpace(strings.ReplaceAll(content, botMention, ""))
	}

	// Commands
	switch {
	case content == "!resume":
		sess, ok := b.manager.GetSession(m.ChannelID)
		if !ok {
			ds.ChannelMessageSendReply(m.ChannelID, "❌ No active session", &discordgo.MessageReference{MessageID: m.ID, ChannelID: m.ChannelID})
			return
		}
		agent, err := b.manager.GetAgentStatus(sess.AgentName)
		if err != nil || agent.LastText == "" {
			ds.ChannelMessageSendReply(m.ChannelID, "❌ No response to resume", &discordgo.MessageReference{MessageID: m.ID, ChannelID: m.ChannelID})
			return
		}
		ds.ChannelMessageSendReply(m.ChannelID, agent.LastText, &discordgo.MessageReference{MessageID: m.ID, ChannelID: m.ChannelID})

	case content == "!pause":
		b.manager.Pause(m.ChannelID)
		ds.ChannelMessageSendReply(m.ChannelID, "⏸️ 暫停監聽，改為 @mention 模式", &discordgo.MessageReference{MessageID: m.ID, ChannelID: m.ChannelID})

	case content == "!back":
		b.manager.Back(m.ChannelID)
		ds.ChannelMessageSendReply(m.ChannelID, "▶️ 恢復完整監聽", &discordgo.MessageReference{MessageID: m.ID, ChannelID: m.ChannelID})

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

var slashCommands = []*discordgo.ApplicationCommand{
	{Name: "start", Description: "綁定專案目錄並啟動 agent", Options: []*discordgo.ApplicationCommandOption{
		{Type: discordgo.ApplicationCommandOptionString, Name: "cwd", Description: "專案目錄路徑", Required: true},
	}},
	{Name: "reset", Description: "重啟 agent"},
	{Name: "status", Description: "查詢目前 agent 狀態"},
	{Name: "cancel", Description: "取消目前執行中的任務"},
	{Name: "cwd", Description: "查詢目前工作目錄"},
	{Name: "pause", Description: "暫停監聽，改為 @mention 模式"},
	{Name: "back", Description: "恢復完整監聽所有訊息"},
}

func (b *Bot) registerSlashCommands() {
	guildID := b.guildID
	// Clear global commands first
	if _, err := b.discord.ApplicationCommandBulkOverwrite(b.discord.State.User.ID, "", []*discordgo.ApplicationCommand{}); err != nil {
		log.Printf("[slash] clear global commands: %v", err)
	}
	created, err := b.discord.ApplicationCommandBulkOverwrite(b.discord.State.User.ID, guildID, slashCommands)
	if err != nil {
		log.Printf("[slash] bulk overwrite error: %v", err)
		return
	}
	for _, cmd := range created {
		log.Printf("[slash] registered /%s (id=%s)", cmd.Name, cmd.ID)
	}
}

func (b *Bot) handleInteraction(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	data := i.ApplicationCommandData()
	log.Printf("[interaction] /%s from %s", data.Name, i.ChannelID)
	channelID := i.ChannelID

	respond := func(msg string) {
		_ = ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: msg},
		})
	}
	followup := func(msg string) {
		_, _ = ds.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{Content: msg})
	}

	switch data.Name {
	case "start":
		cwd := data.Options[0].StringValue()
		respond("⏳ Starting agent at `" + cwd + "`...")
		go func() {
			if err := b.manager.StartAt(channelID, cwd); err != nil {
				followup("❌ " + err.Error())
			} else {
				followup("✅ Agent started at `" + cwd + "`")
			}
		}()
	case "reset":
		respond("⏳ Resetting...")
		go func() {
			if err := b.manager.Reset(channelID); err != nil {
				followup("❌ " + err.Error())
			} else {
				followup("✅ Session reset.")
			}
		}()
	case "status":
		respond(b.manager.Status(channelID))
	case "cancel":
		if err := b.manager.Cancel(channelID); err != nil {
			respond("❌ " + err.Error())
		} else {
			respond("⚠️ Cancel requested.")
		}
	case "cwd":
		respond(b.manager.CWD(channelID))
	case "pause":
		b.manager.Pause(channelID)
		respond("⏸️ 暫停監聽，改為 @mention 模式")
	case "back":
		b.manager.Back(channelID)
		respond("▶️ 恢復完整監聽")
	}
}
