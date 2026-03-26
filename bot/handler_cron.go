package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nczz/kiro-discord-bot/heartbeat"
)

// handleCronModal responds to /cron by showing a modal form.
func (b *Bot) handleCronModal(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "cron_add_modal",
			Title:    "新增排程任務",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "cron_name",
						Label:       "任務名稱",
						Style:       discordgo.TextInputShort,
						Placeholder: "例：每日伺服器健檢",
						Required:    true,
						MaxLength:   100,
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "cron_schedule",
						Label:       "執行頻率",
						Style:       discordgo.TextInputShort,
						Placeholder: "例：每天 09:00、每 30 分鐘、0 9 * * *",
						Required:    true,
						MaxLength:   100,
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "cron_prompt",
						Label:       "要 agent 做什麼",
						Style:       discordgo.TextInputParagraph,
						Placeholder: "例：檢查伺服器 CPU、記憶體、磁碟用量，跟上次比較",
						Required:    true,
						MaxLength:   2000,
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "cron_cwd",
						Label:       "工作目錄（選填）",
						Style:       discordgo.TextInputShort,
						Placeholder: "例：/home/user/project",
						Required:    false,
						MaxLength:   200,
					},
				}},
			},
		},
	})
}

// handleCronModalSubmit processes the modal form submission.
func (b *Bot) handleCronModalSubmit(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	fields := map[string]string{}
	for _, row := range data.Components {
		ar, ok := row.(*discordgo.ActionsRow)
		if !ok {
			continue
		}
		for _, comp := range ar.Components {
			ti, ok := comp.(*discordgo.TextInput)
			if !ok {
				continue
			}
			fields[ti.CustomID] = ti.Value
		}
	}

	name := fields["cron_name"]
	scheduleInput := fields["cron_schedule"]
	prompt := fields["cron_prompt"]
	cwd := fields["cron_cwd"]

	// Parse schedule
	cronExpr, err := heartbeat.ParseSchedule(scheduleInput)
	if err != nil {
		respondInteraction(ds, i, "❌ 無法解析排程："+err.Error())
		return
	}

	username := ""
	if i.Member != nil && i.Member.User != nil {
		username = i.Member.User.Username
	}

	job := &heartbeat.CronJob{
		Name:          name,
		ChannelID:     i.ChannelID,
		Schedule:      cronExpr,
		ScheduleHuman: scheduleInput,
		Prompt:        prompt,
		CWD:           cwd,
		HistoryLimit:  10,
		Enabled:       true,
		CreatedBy:     username,
	}
	if err := b.cronStore.Add(job); err != nil {
		respondInteraction(ds, i, "❌ 儲存失敗："+err.Error())
		return
	}

	respondInteraction(ds, i, fmt.Sprintf("✅ 排程任務已建立\n名稱：**%s**\n頻率：%s (`%s`)\nPrompt：%s",
		name, scheduleInput, cronExpr, prompt))
}

// handleCronList responds to /cron-list with a list of jobs and action buttons.
func (b *Bot) handleCronList(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	jobs := b.cronStore.ListByChannel(i.ChannelID)
	if len(jobs) == 0 {
		followupInteraction(ds, i, "📋 此頻道沒有排程任務。使用 `/cron` 新增。")
		return
	}

	for _, job := range jobs {
		status := "✅"
		if !job.Enabled {
			status = "⏸️"
		}

		lastRun := "無"
		if job.LastRun != "" {
			if t, err := time.Parse(time.RFC3339, job.LastRun); err == nil {
				lastRun = t.Format("01/02 15:04")
			}
		}
		nextRun := "計算中"
		if job.NextRun != "" {
			if t, err := time.Parse(time.RFC3339, job.NextRun); err == nil {
				nextRun = t.Format("01/02 15:04")
			}
		}

		content := fmt.Sprintf("%s **%s**\n頻率：%s | 上次：%s | 下次：%s\nPrompt：%s",
			status, job.Name, job.ScheduleHuman, lastRun, nextRun, truncate(job.Prompt, 100))

		// Build buttons
		var buttons []discordgo.MessageComponent
		if job.Enabled {
			buttons = append(buttons, discordgo.Button{
				Label:    "⏸️ 暫停",
				Style:    discordgo.SecondaryButton,
				CustomID: "cron_pause_" + job.ID,
			})
		} else {
			buttons = append(buttons, discordgo.Button{
				Label:    "▶️ 恢復",
				Style:    discordgo.SuccessButton,
				CustomID: "cron_resume_" + job.ID,
			})
		}
		buttons = append(buttons,
			discordgo.Button{
				Label:    "▶️ 立即執行",
				Style:    discordgo.PrimaryButton,
				CustomID: "cron_run_" + job.ID,
			},
			discordgo.Button{
				Label:    "🗑️ 刪除",
				Style:    discordgo.DangerButton,
				CustomID: "cron_delete_" + job.ID,
			},
		)

		_, _ = ds.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: content,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: buttons},
			},
		})
	}
}

// handleCronButton processes button clicks on cron-list messages.
func (b *Bot) handleCronButton(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID

	// Parse action and job ID
	var action, jobID string
	for _, prefix := range []string{"cron_pause_", "cron_resume_", "cron_run_", "cron_delete_"} {
		if strings.HasPrefix(customID, prefix) {
			action = strings.TrimSuffix(strings.TrimPrefix(prefix, "cron_"), "_")
			jobID = strings.TrimPrefix(customID, prefix)
			break
		}
	}
	if jobID == "" {
		return
	}

	job, ok := b.cronStore.Get(jobID)
	if !ok {
		respondInteraction(ds, i, "❌ 找不到此任務")
		return
	}

	switch action {
	case "pause":
		job.Enabled = false
		_ = b.cronStore.Update(job)
		respondInteraction(ds, i, fmt.Sprintf("⏸️ 已暫停：**%s**", job.Name))
	case "resume":
		job.Enabled = true
		_ = b.cronStore.Update(job)
		respondInteraction(ds, i, fmt.Sprintf("▶️ 已恢復：**%s**", job.Name))
	case "run":
		respondInteraction(ds, i, fmt.Sprintf("⏰ 正在手動執行：**%s**", job.Name))
		// Trigger execution in background — set NextRun to now so next heartbeat picks it up
		job.NextRun = time.Now().Add(-time.Minute).Format(time.RFC3339)
		_ = b.cronStore.Update(job)
	case "delete":
		_ = b.cronStore.Remove(jobID)
		respondInteraction(ds, i, fmt.Sprintf("🗑️ 已刪除：**%s**", job.Name))
	}
}

// handleCronRun handles /cron-run <name>
func (b *Bot) handleCronRun(ds *discordgo.Session, i *discordgo.InteractionCreate, name string) {
	job, ok := b.cronStore.FindByName(i.ChannelID, name)
	if !ok {
		respondInteraction(ds, i, "❌ 找不到任務："+name)
		return
	}
	respondInteraction(ds, i, fmt.Sprintf("⏰ 正在手動執行：**%s**", job.Name))
	job.NextRun = time.Now().Add(-time.Minute).Format(time.RFC3339)
	_ = b.cronStore.Update(job)
}

func respondInteraction(ds *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	_ = ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: msg},
	})
}

func followupInteraction(ds *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	_, _ = ds.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{Content: msg})
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// handleCronTextCommand handles !cron text commands.
func (b *Bot) handleCronTextCommand(ds *discordgo.Session, channelID, content string) {
	switch {
	case content == "!cron list":
		jobs := b.cronStore.ListByChannel(channelID)
		if len(jobs) == 0 {
			ds.ChannelMessageSend(channelID, "📋 此頻道沒有排程任務。")
			return
		}
		var sb strings.Builder
		sb.WriteString("📋 **排程任務列表**\n\n")
		for i, job := range jobs {
			status := "✅"
			if !job.Enabled {
				status = "⏸️"
			}
			sb.WriteString(fmt.Sprintf("%d. %s **%s** — %s\n   Prompt: %s\n",
				i+1, status, job.Name, job.ScheduleHuman, truncate(job.Prompt, 80)))
		}
		ds.ChannelMessageSend(channelID, sb.String())

	case strings.HasPrefix(content, "!cron run "):
		name := strings.TrimSpace(strings.TrimPrefix(content, "!cron run "))
		job, ok := b.cronStore.FindByName(channelID, name)
		if !ok {
			ds.ChannelMessageSend(channelID, "❌ 找不到任務："+name)
			return
		}
		ds.ChannelMessageSend(channelID, fmt.Sprintf("⏰ 正在手動執行：**%s**", job.Name))
		job.NextRun = time.Now().Add(-time.Minute).Format(time.RFC3339)
		_ = b.cronStore.Update(job)

	default:
		ds.ChannelMessageSend(channelID, "Usage: `!cron list` | `!cron run <name>`")
	}
}

func init() {
	log.Println("[handler_cron] loaded")
}
