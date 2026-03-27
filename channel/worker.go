package channel

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/nczz/kiro-discord-bot/acp"
	L "github.com/nczz/kiro-discord-bot/locale"
)

// Job represents a single user message to be processed.
type Job struct {
	ChannelID   string
	MessageID   string
	Prompt      string
	Session     *discordgo.Session
	UserID      string
	Username    string
	Attachments []string
}

// Worker manages a per-channel job queue and executes jobs sequentially.
type Worker struct {
	channelID       string
	agent           *acp.Agent
	queue           chan *Job
	askTimeoutSec   int
	streamUpdateSec int
	stopCh          chan struct{}
	stopped         sync.Once
	started         sync.Once
	logger          *ChatLogger
	model           string

	cancelMu sync.Mutex
	cancelFn context.CancelFunc // cancel the currently running job
}

func NewWorker(channelID string, agent *acp.Agent, bufSize, askTimeoutSec, streamUpdateSec int, logger *ChatLogger, model string) *Worker {
	return &Worker{
		channelID:       channelID,
		agent:           agent,
		queue:           make(chan *Job, bufSize),
		askTimeoutSec:   askTimeoutSec,
		streamUpdateSec: streamUpdateSec,
		stopCh:          make(chan struct{}),
		logger:          logger,
		model:           model,
	}
}

func (w *Worker) Enqueue(job *Job) error {
	select {
	case w.queue <- job:
		return nil
	default:
		return fmt.Errorf("queue full")
	}
}

func (w *Worker) QueueLen() int {
	return len(w.queue)
}

func (w *Worker) Start() {
	w.started.Do(func() {
		go w.run()
	})
}

func (w *Worker) Stop() {
	w.stopped.Do(func() {
		close(w.stopCh)
	})
}

// CancelCurrent cancels the currently running job, if any.
func (w *Worker) CancelCurrent() {
	w.cancelMu.Lock()
	fn := w.cancelFn
	w.cancelMu.Unlock()
	if fn != nil {
		fn()
	}
}

func (w *Worker) run() {
	for {
		select {
		case <-w.stopCh:
			return
		case job := <-w.queue:
			w.execute(job)
		}
	}
}

func (w *Worker) execute(job *Job) {
	ds := job.Session

	if w.logger != nil {
		w.logger.Log(w.channelID, ChatEntry{
			Role:        "user",
			UserID:      job.UserID,
			Username:    job.Username,
			MessageID:   job.MessageID,
			Content:     job.Prompt,
			Attachments: job.Attachments,
		})
	}

	swapReaction(ds, job.ChannelID, job.MessageID, "⏳", "🔄")

	replyMsg, err := ds.ChannelMessageSendReply(job.ChannelID, L.Get("worker.processing"), &discordgo.MessageReference{
		MessageID: job.MessageID,
		ChannelID: job.ChannelID,
	})
	if err != nil {
		log.Printf("[worker %s] send placeholder: %v", w.channelID, err)
		swapReaction(ds, job.ChannelID, job.MessageID, "🔄", "❌")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(w.askTimeoutSec)*time.Second)
	defer cancel()
	w.cancelMu.Lock()
	w.cancelFn = cancel
	w.cancelMu.Unlock()
	defer func() {
		w.cancelMu.Lock()
		w.cancelFn = nil
		w.cancelMu.Unlock()
	}()

	var mu sync.Mutex
	accumulated := ""
	lastUpdate := time.Now()
	statusLine := L.Get("worker.processing")

	w.agent.OnToolUseFunc(func(started bool) {
		mu.Lock()
		if started {
			statusLine = L.Get("worker.tool_running")
			swapReaction(ds, job.ChannelID, job.MessageID, "🔄", "⚙️")
		} else {
			statusLine = L.Get("worker.processing")
			swapReaction(ds, job.ChannelID, job.MessageID, "⚙️", "🔄")
		}
		mu.Unlock()
	})

	onChunk := func(chunk string) {
		mu.Lock()
		accumulated += chunk
		shouldUpdate := time.Since(lastUpdate) >= time.Duration(w.streamUpdateSec)*time.Second
		snap := accumulated
		status := statusLine
		mu.Unlock()

		if shouldUpdate {
			mu.Lock()
			lastUpdate = time.Now()
			mu.Unlock()
			editMessage(ds, job.ChannelID, replyMsg.ID, status+"\n\n"+snap)
		}
	}

	response, err := w.agent.Ask(ctx, job.Prompt, onChunk)

	if err != nil {
		errMsg := err.Error()
		if ctx.Err() == context.DeadlineExceeded {
			errMsg = L.Getf("worker.timeout", w.askTimeoutSec)
			swapReaction(ds, job.ChannelID, job.MessageID, "🔄", "⚠️")
		} else if ctx.Err() == context.Canceled {
			errMsg = L.Get("cancel.success")
			swapReaction(ds, job.ChannelID, job.MessageID, "🔄", "⚠️")
		} else {
			swapReaction(ds, job.ChannelID, job.MessageID, "🔄", "❌")
		}
		editMessage(ds, job.ChannelID, replyMsg.ID, "❌ "+errMsg)
		if w.logger != nil {
			w.logger.Log(w.channelID, ChatEntry{
				Role:    "assistant",
				Content: "❌ " + errMsg,
				Model:   w.model,
			})
		}
		return
	}

	if response == "" {
		mu.Lock()
		response = accumulated
		mu.Unlock()
	}

	swapReaction(ds, job.ChannelID, job.MessageID, "🔄", "✅")
	sendLong(ds, job.ChannelID, replyMsg.ID, response)

	if w.logger != nil {
		w.logger.Log(w.channelID, ChatEntry{
			Role:    "assistant",
			Content: response,
			Model:   w.model,
		})
	}
}

func editMessage(ds *discordgo.Session, channelID, msgID, content string) {
	if len(content) > 2000 {
		content = truncateUTF8(content, 1997) + "..."
	}
	_, _ = ds.ChannelMessageEdit(channelID, msgID, content)
}

func sendLong(ds *discordgo.Session, channelID, placeholderID, content string) {
	const limit = 1990
	parts := splitMessage(content, limit)
	if len(parts) == 0 {
		editMessage(ds, channelID, placeholderID, L.Get("worker.empty_response"))
		return
	}

	prefix := ""
	if len(parts) > 1 {
		prefix = fmt.Sprintf("(1/%d) ", len(parts))
	}
	editMessage(ds, channelID, placeholderID, prefix+parts[0])

	for i := 1; i < len(parts); i++ {
		label := fmt.Sprintf("(%d/%d) ", i+1, len(parts))
		_, _ = ds.ChannelMessageSend(channelID, label+parts[i])
	}
}

func splitMessage(s string, limit int) []string {
	var parts []string
	for len(s) > limit {
		idx := strings.LastIndex(s[:limit], "\n")
		if idx < limit/2 {
			// No good newline break — find a valid UTF-8 boundary near limit
			idx = limit
			for idx > 0 && !utf8.RuneStart(s[idx]) {
				idx--
			}
		}
		parts = append(parts, s[:idx])
		s = s[idx:]
		if len(s) > 0 && s[0] == '\n' {
			s = s[1:]
		}
	}
	if s != "" {
		parts = append(parts, s)
	}
	return parts
}

// truncateUTF8 truncates s to at most maxBytes without breaking a multi-byte character.
func truncateUTF8(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	for maxBytes > 0 && !utf8.RuneStart(s[maxBytes]) {
		maxBytes--
	}
	return s[:maxBytes]
}
