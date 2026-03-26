package heartbeat

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

// CronHistory is a single execution record.
type CronHistory struct {
	Timestamp   string `json:"ts"`
	Prompt      string `json:"prompt"`
	Response    string `json:"response"`
	Status      string `json:"status"` // "ok" or "error"
	DurationSec int    `json:"duration_sec"`
}

// CronDeps abstracts dependencies for the cron task.
type CronDeps interface {
	StartTempAgent(name, cwd string) error
	StopTempAgent(name string)
	AskAgent(ctx context.Context, name, prompt string) (string, error)
	Notify(channelID, msg string)
}

// CronTask checks and executes due cron jobs.
type CronTask struct {
	store    *CronStore
	deps     CronDeps
	dataDir  string
	location *time.Location
	parser   cron.Parser
}

func NewCronTask(store *CronStore, deps CronDeps, dataDir string, tz string) *CronTask {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.Now().Location()
	}
	return &CronTask{
		store:    store,
		deps:     deps,
		dataDir:  dataDir,
		location: loc,
		parser:   cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
	}
}

func (c *CronTask) Name() string { return "cron" }

func (c *CronTask) ShouldRun(_ time.Time) bool {
	return len(c.store.All()) > 0
}

func (c *CronTask) Run() error {
	now := time.Now().In(c.location)
	for _, job := range c.store.All() {
		if !job.Enabled || job.Running {
			continue
		}
		if !c.isDue(job, now) {
			continue
		}
		job.Running = true
		go c.execute(job, now)
	}
	return nil
}

func (c *CronTask) isDue(job *CronJob, now time.Time) bool {
	if job.NextRun == "" {
		// First run — compute next from schedule
		next, err := c.computeNext(job.Schedule, job.CreatedAt)
		if err != nil {
			return false
		}
		job.NextRun = next.Format(time.RFC3339)
		_ = c.store.Update(job)
	}
	nextRun, err := time.ParseInLocation(time.RFC3339, job.NextRun, c.location)
	if err != nil {
		return false
	}
	return !now.Before(nextRun)
}

func (c *CronTask) computeNext(schedule string, afterStr string) (time.Time, error) {
	sched, err := c.parser.Parse(schedule)
	if err != nil {
		return time.Time{}, err
	}
	after := time.Now().In(c.location)
	if afterStr != "" {
		if t, err := time.ParseInLocation(time.RFC3339, afterStr, c.location); err == nil {
			after = t
		}
	}
	return sched.Next(after), nil
}

func (c *CronTask) execute(job *CronJob, now time.Time) {
	defer func() { job.Running = false }()

	agentName := "cron-" + job.ID
	start := time.Now()

	log.Printf("[cron] executing job %s (%s)", job.ID, job.Name)
	c.deps.Notify(job.ChannelID, fmt.Sprintf("⏰ 排程任務執行中：**%s**", job.Name))

	// Load history
	history := c.loadHistory(job.ID, job.HistoryLimit)
	prompt := c.buildPrompt(job, history)

	// Start temp agent
	cwd := job.CWD
	if cwd == "" {
		cwd = "/tmp"
	}
	if err := c.deps.StartTempAgent(agentName, cwd); err != nil {
		log.Printf("[cron] start agent for %s failed: %v", job.ID, err)
		c.deps.Notify(job.ChannelID, fmt.Sprintf("❌ 排程任務 **%s** 啟動失敗：%s", job.Name, err.Error()))
		c.saveHistory(job.ID, CronHistory{
			Timestamp: now.Format(time.RFC3339), Prompt: job.Prompt, Response: err.Error(), Status: "error",
			DurationSec: int(time.Since(start).Seconds()),
		})
		c.advanceNextRun(job, now)
		return
	}
	defer c.deps.StopTempAgent(agentName)

	// Ask
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	response, err := c.deps.AskAgent(ctx, agentName, prompt)
	duration := int(time.Since(start).Seconds())
	status := "ok"

	if err != nil {
		response = err.Error()
		status = "error"
		c.deps.Notify(job.ChannelID, fmt.Sprintf("❌ 排程任務 **%s** 失敗：%s", job.Name, err.Error()))
	} else {
		// Truncate for Discord (2000 char limit)
		display := response
		if len(display) > 1900 {
			display = display[:1900] + "\n...(truncated)"
		}
		c.deps.Notify(job.ChannelID, fmt.Sprintf("⏰ **%s** 執行完成：\n%s", job.Name, display))
	}

	c.saveHistory(job.ID, CronHistory{
		Timestamp: now.Format(time.RFC3339), Prompt: job.Prompt, Response: response, Status: status, DurationSec: duration,
	})
	c.advanceNextRun(job, now)
}

func (c *CronTask) buildPrompt(job *CronJob, history []CronHistory) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[排程任務: %s]\n", job.Name))
	if len(history) > 0 {
		sb.WriteString("\n以下是過去的執行紀錄：\n---\n")
		for _, h := range history {
			ts, _ := time.Parse(time.RFC3339, h.Timestamp)
			sb.WriteString(fmt.Sprintf("[%s] (%s) %s\n", ts.Format("01/02 15:04"), h.Status, h.Response))
		}
		sb.WriteString("---\n\n")
	}
	sb.WriteString("請執行：")
	sb.WriteString(job.Prompt)
	return sb.String()
}

func (c *CronTask) advanceNextRun(job *CronJob, after time.Time) {
	sched, err := c.parser.Parse(job.Schedule)
	if err != nil {
		return
	}
	job.LastRun = after.Format(time.RFC3339)
	job.NextRun = sched.Next(after).Format(time.RFC3339)
	_ = c.store.Update(job)
}

func (c *CronTask) historyPath(jobID string) string {
	dir := filepath.Join(c.dataDir, "cron", jobID)
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "history.jsonl")
}

func (c *CronTask) loadHistory(jobID string, limit int) []CronHistory {
	f, err := os.Open(c.historyPath(jobID))
	if err != nil {
		return nil
	}
	defer f.Close()

	var all []CronHistory
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var h CronHistory
		if json.Unmarshal(scanner.Bytes(), &h) == nil {
			all = append(all, h)
		}
	}
	if len(all) > limit {
		all = all[len(all)-limit:]
	}
	return all
}

func (c *CronTask) saveHistory(jobID string, entry CronHistory) {
	data, err := json.Marshal(entry)
	if err != nil {
		return
	}
	data = append(data, '\n')
	f, err := os.OpenFile(c.historyPath(jobID), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write(data)
}
