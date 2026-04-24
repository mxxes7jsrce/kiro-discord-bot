package bot

import (
	"fmt"
	"strings"
)

// triviaHelpText describes how to use trivia commands.
const triviaHelpText = `**Trivia Commands**
` + "`!trivia start`" + ` — Start a trivia question in this channel
` + "`!trivia answer <1-4>`" + ` — Submit your answer
` + "`!trivia stop`" + ` — Cancel the current trivia session`

// handleTriviaCommand routes trivia sub-commands.
func (h *Handler) handleTriviaCommand(channelID, content string) string {
	parts := strings.Fields(content)
	if len(parts) < 2 {
		return triviaHelpText
	}

	switch strings.ToLower(parts[1]) {
	case "start":
		return h.handleTriviaStart(channelID)
	case "answer":
		if len(parts) < 3 {
			return "Usage: `!trivia answer <1-4>`"
		}
		return h.handleTriviaAnswer(channelID, parts[2])
	case "stop":
		return h.handleTriviaStop(channelID)
	default:
		return triviaHelpText
	}
}

func (h *Handler) handleTriviaStart(channelID string) string {
	q, ok := h.trivia.Start(channelID)
	if !ok {
		return "A trivia session is already active! Use `!trivia answer <1-4>` to answer."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🎯 **Trivia Time!**\n%s\n", q.Question))
	for i, opt := range q.Options {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, opt))
	}
	sb.WriteString("Use `!trivia answer <1-4>` to submit your answer!")
	return sb.String()
}

func (h *Handler) handleTriviaAnswer(channelID, raw string) string {
	var idx int
	_, err := fmt.Sscanf(raw, "%d", &idx)
	if err != nil || idx < 1 || idx > 4 {
		return "Please provide a valid answer number between 1 and 4."
	}

	correct, found := h.trivia.Answer(channelID, idx-1)
	if !found {
		return "No active trivia session. Start one with `!trivia start`."
	}
	if correct {
		return "✅ Correct! Well done!"
	}
	return "❌ Wrong answer! Try again with `!trivia answer <1-4>`."
}

func (h *Handler) handleTriviaStop(channelID string) string {
	q, ok := h.trivia.Stop(channelID)
	if !ok {
		return "No active trivia session to stop."
	}
	return fmt.Sprintf("🛑 Trivia stopped. The correct answer was: **%s**", q.Options[q.Answer])
}
