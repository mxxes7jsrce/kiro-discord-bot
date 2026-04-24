package bot

import (
	"strings"
	"testing"
)

// newTestHandler builds a minimal Handler wired with a fresh TriviaStore.
func newTestHandlerForTrivia() *Handler {
	h := &Handler{
		commands: make(map[string]commandFunc),
		trivia:   NewTriviaStore(),
	}
	h.registerTriviaHandlers()
	return h
}

func TestHandlerTrivia_StartAndAnswer(t *testing.T) {
	h := newTestHandlerForTrivia()
	channelID := "integration-ch-001"

	resp := h.handleTriviaCommand(channelID, "!trivia start")
	if !strings.Contains(resp, "Trivia Time") {
		t.Errorf("expected trivia start message, got: %s", resp)
	}

	// Duplicate start
	resp = h.handleTriviaCommand(channelID, "!trivia start")
	if !strings.Contains(resp, "already active") {
		t.Errorf("expected already active message, got: %s", resp)
	}

	// Invalid answer format
	resp = h.handleTriviaCommand(channelID, "!trivia answer abc")
	if !strings.Contains(resp, "valid answer") {
		t.Errorf("expected validation message, got: %s", resp)
	}

	// Out-of-range answer
	resp = h.handleTriviaCommand(channelID, "!trivia answer 9")
	if !strings.Contains(resp, "valid answer") {
		t.Errorf("expected validation message for out-of-range, got: %s", resp)
	}
}

func TestHandlerTrivia_Stop(t *testing.T) {
	h := newTestHandlerForTrivia()
	channelID := "integration-ch-002"

	resp := h.handleTriviaCommand(channelID, "!trivia stop")
	if !strings.Contains(resp, "No active") {
		t.Errorf("expected no active session message, got: %s", resp)
	}

	h.handleTriviaCommand(channelID, "!trivia start")
	resp = h.handleTriviaCommand(channelID, "!trivia stop")
	if !strings.Contains(resp, "correct answer") {
		t.Errorf("expected correct answer reveal on stop, got: %s", resp)
	}
}

func TestHandlerTrivia_HelpFallback(t *testing.T) {
	h := newTestHandlerForTrivia()
	channelID := "integration-ch-003"

	resp := h.handleTriviaCommand(channelID, "!trivia")
	if !strings.Contains(resp, "Trivia Commands") {
		t.Errorf("expected help text, got: %s", resp)
	}

	resp = h.handleTriviaCommand(channelID, "!trivia unknown")
	if !strings.Contains(resp, "Trivia Commands") {
		t.Errorf("expected help text for unknown sub-command, got: %s", resp)
	}
}
