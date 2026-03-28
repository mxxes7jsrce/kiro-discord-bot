package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var dg *discordgo.Session

func ensureDiscord() error {
	if dg != nil {
		return nil
	}
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return fmt.Errorf("DISCORD_TOKEN not set")
	}
	var err error
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	// Only REST API needed, no Gateway connection
	return nil
}

func main() {
	s := server.NewMCPServer("mcp-discord", "1.0.0", server.WithToolCapabilities(false))

	// 1. List channels
	s.AddTool(
		mcp.NewTool("discord_list_channels",
			mcp.WithDescription("List text channels in a guild"),
			mcp.WithString("guild_id", mcp.Required(), mcp.Description("Guild/server ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := ensureDiscord(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			guildID, _ := req.RequireString("guild_id")
			channels, err := dg.GuildChannels(guildID)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			var lines []string
			for _, ch := range channels {
				if ch.Type == discordgo.ChannelTypeGuildText {
					lines = append(lines, fmt.Sprintf("#%s (%s)", ch.Name, ch.ID))
				}
			}
			return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
		},
	)

	// 2. Read messages
	s.AddTool(
		mcp.NewTool("discord_read_messages",
			mcp.WithDescription("Read recent messages from a channel"),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
			mcp.WithNumber("limit", mcp.Description("Number of messages, max 100, default 20")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := ensureDiscord(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			chID, _ := req.RequireString("channel_id")
			limit := int(req.GetFloat("limit", 20))
			if limit > 100 {
				limit = 100
			}
			msgs, err := dg.ChannelMessages(chID, limit, "", "", "")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			var lines []string
			for i := len(msgs) - 1; i >= 0; i-- {
				m := msgs[i]
				lines = append(lines, fmt.Sprintf("[%s] %s: %s", time.Time(m.Timestamp).Format(time.RFC3339), m.Author.Username, m.Content))
			}
			return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
		},
	)

	// 3. Send message
	s.AddTool(
		mcp.NewTool("discord_send_message",
			mcp.WithDescription("Send a message to a channel"),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
			mcp.WithString("content", mcp.Required(), mcp.Description("Message content")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := ensureDiscord(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			chID, _ := req.RequireString("channel_id")
			content, _ := req.RequireString("content")
			msg, err := dg.ChannelMessageSend(chID, content)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Sent message %s", msg.ID)), nil
		},
	)

	// 4. Reply to message
	s.AddTool(
		mcp.NewTool("discord_reply_message",
			mcp.WithDescription("Reply to a specific message"),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
			mcp.WithString("message_id", mcp.Required(), mcp.Description("Message ID to reply to")),
			mcp.WithString("content", mcp.Required(), mcp.Description("Reply content")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := ensureDiscord(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			chID, _ := req.RequireString("channel_id")
			msgID, _ := req.RequireString("message_id")
			content, _ := req.RequireString("content")
			msg, err := dg.ChannelMessageSendReply(chID, content, &discordgo.MessageReference{
				MessageID: msgID, ChannelID: chID,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Replied with message %s", msg.ID)), nil
		},
	)

	// 5. Add reaction
	s.AddTool(
		mcp.NewTool("discord_add_reaction",
			mcp.WithDescription("Add a reaction emoji to a message"),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
			mcp.WithString("message_id", mcp.Required(), mcp.Description("Message ID")),
			mcp.WithString("emoji", mcp.Required(), mcp.Description("Emoji (unicode or name:id)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := ensureDiscord(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			chID, _ := req.RequireString("channel_id")
			msgID, _ := req.RequireString("message_id")
			emoji, _ := req.RequireString("emoji")
			if err := dg.MessageReactionAdd(chID, msgID, emoji); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(fmt.Sprintf("Reacted with %s", emoji)), nil
		},
	)

	// 6. List members
	s.AddTool(
		mcp.NewTool("discord_list_members",
			mcp.WithDescription("List members of a guild"),
			mcp.WithString("guild_id", mcp.Required(), mcp.Description("Guild/server ID")),
			mcp.WithNumber("limit", mcp.Description("Max members, default 50")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := ensureDiscord(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			guildID, _ := req.RequireString("guild_id")
			limit := int(req.GetFloat("limit", 50))
			members, err := dg.GuildMembers(guildID, "", limit)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			var lines []string
			for _, m := range members {
				nick := ""
				if m.Nick != "" {
					nick = fmt.Sprintf(" (%s)", m.Nick)
				}
				bot := ""
				if m.User.Bot {
					bot = " 🤖"
				}
				lines = append(lines, fmt.Sprintf("%s%s [%s]%s", m.User.Username, nick, m.User.ID, bot))
			}
			return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
		},
	)

	// 7. Search messages
	s.AddTool(
		mcp.NewTool("discord_search_messages",
			mcp.WithDescription("Search recent messages in a channel by keyword"),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search keyword")),
			mcp.WithNumber("limit", mcp.Description("Messages to scan, default 50")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := ensureDiscord(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			chID, _ := req.RequireString("channel_id")
			query, _ := req.RequireString("query")
			limit := int(req.GetFloat("limit", 50))
			msgs, err := dg.ChannelMessages(chID, limit, "", "", "")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			q := strings.ToLower(query)
			var lines []string
			for i := len(msgs) - 1; i >= 0; i-- {
				m := msgs[i]
				if strings.Contains(strings.ToLower(m.Content), q) {
					t := time.Time(m.Timestamp)
					lines = append(lines, fmt.Sprintf("[%s] %s: %s", t.Format(time.RFC3339), m.Author.Username, m.Content))
				}
			}
			if len(lines) == 0 {
				return mcp.NewToolResultText("No matches."), nil
			}
			return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
		},
	)

	// 8. Channel info
	s.AddTool(
		mcp.NewTool("discord_channel_info",
			mcp.WithDescription("Get detailed info about a channel"),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := ensureDiscord(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			chID, _ := req.RequireString("channel_id")
			ch, err := dg.Channel(chID)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			topic := ch.Topic
			if topic == "" {
				topic = "(none)"
			}
			info := fmt.Sprintf("name: #%s\nid: %s\ntype: %d\ntopic: %s\nguild_id: %s",
				ch.Name, ch.ID, ch.Type, topic, ch.GuildID)
			return mcp.NewToolResultText(info), nil
		},
	)

	// 9. Send file
	s.AddTool(
		mcp.NewTool("discord_send_file",
			mcp.WithDescription("Upload a local file to a channel as an attachment"),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
			mcp.WithString("file_path", mcp.Required(), mcp.Description("Absolute path to the local file")),
			mcp.WithString("content", mcp.Description("Optional text message to send with the file")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := ensureDiscord(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			chID, _ := req.RequireString("channel_id")
			filePath, _ := req.RequireString("file_path")
			content := req.GetString("content", "")

			f, err := os.Open(filePath)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("open file: %v", err)), nil
			}
			defer f.Close()

			info, err := f.Stat()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("stat file: %v", err)), nil
			}
			if info.Size() > 25*1024*1024 {
				return mcp.NewToolResultError("file exceeds 25MB Discord limit"), nil
			}

			msg, err := dg.ChannelMessageSendComplex(chID, &discordgo.MessageSend{
				Content: content,
				Files: []*discordgo.File{{
					Name:   filepath.Base(filePath),
					Reader: f,
				}},
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			var urls []string
			for _, a := range msg.Attachments {
				urls = append(urls, a.URL)
			}
			return mcp.NewToolResultText(fmt.Sprintf("Sent message %s\n%s", msg.ID, strings.Join(urls, "\n"))), nil
		},
	)

	// 10. List attachments
	s.AddTool(
		mcp.NewTool("discord_list_attachments",
			mcp.WithDescription("List file attachments from recent messages in a channel"),
			mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
			mcp.WithNumber("limit", mcp.Description("Messages to scan, default 50, max 100")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := ensureDiscord(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			chID, _ := req.RequireString("channel_id")
			limit := int(req.GetFloat("limit", 50))
			if limit > 100 {
				limit = 100
			}
			msgs, err := dg.ChannelMessages(chID, limit, "", "", "")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			var lines []string
			for i := len(msgs) - 1; i >= 0; i-- {
				m := msgs[i]
				for _, a := range m.Attachments {
					t := time.Time(m.Timestamp).Format(time.RFC3339)
					lines = append(lines, fmt.Sprintf("[%s] %s | %s | %d bytes | msg:%s | %s",
						t, a.Filename, a.ContentType, a.Size, m.ID, a.URL))
				}
			}
			if len(lines) == 0 {
				return mcp.NewToolResultText("No attachments found."), nil
			}
			return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
		},
	)

	// 11. Download attachment
	s.AddTool(
		mcp.NewTool("discord_download_attachment",
			mcp.WithDescription("Download a Discord attachment to a local file"),
			mcp.WithString("url", mcp.Required(), mcp.Description("Attachment URL (from discord_list_attachments)")),
			mcp.WithString("save_dir", mcp.Description("Directory to save the file (default: system temp dir)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := ensureDiscord(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			url, _ := req.RequireString("url")
			saveDir := req.GetString("save_dir", os.TempDir())

			if err := os.MkdirAll(saveDir, 0755); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("create dir: %v", err)), nil
			}

			resp, err := http.Get(url)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("download: %v", err)), nil
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return mcp.NewToolResultError(fmt.Sprintf("download: HTTP %d", resp.StatusCode)), nil
			}

			// Extract filename from URL path
			name := filepath.Base(url)
			if idx := strings.Index(name, "?"); idx > 0 {
				name = name[:idx]
			}
			ts := time.Now().Format("20060102-150405")
			dst := filepath.Join(saveDir, ts+"-"+name)

			f, err := os.Create(dst)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("create file: %v", err)), nil
			}
			n, err := io.Copy(f, resp.Body)
			f.Close()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("write file: %v", err)), nil
			}
			abs, _ := filepath.Abs(dst)
			return mcp.NewToolResultText(fmt.Sprintf("Saved %s (%d bytes)", abs, n)), nil
		},
	)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
