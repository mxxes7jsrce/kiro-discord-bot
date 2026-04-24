package bot

import "strings"

// registerTriviaHandlers wires the trivia command into the handler dispatch table.
func (h *Handler) registerTriviaHandlers() {
	h.commands["trivia"] = h.dispatchTrivia
}

// dispatchTrivia is the top-level entry point called by the command router.
func (h *Handler) dispatchTrivia(channelID, content string) string {
	return h.handleTriviaCommand(channelID, content)
}

// isTriviaCommand returns true when the message starts with the trivia prefix.
func isTriviaCommand(content string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(content)), "!trivia")
}
