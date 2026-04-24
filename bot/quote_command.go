package bot

import (
	"fmt"
	"strings"
)

// quoteCommandResult holds the response text for a quote command.
type quoteCommandResult struct {
	Message string
}

// handleAddQuote processes "!quote add <text> -- <author>" style input.
func handleAddQuote(store *QuoteStore, guildID, args string) quoteCommandResult {
	parts := strings.SplitN(args, "--", 2)
	text := strings.TrimSpace(parts[0])
	author := "unknown"
	if len(parts) == 2 {
		author = strings.TrimSpace(parts[1])
	}
	if text == "" {
		return quoteCommandResult{Message: "Usage: !quote add <text> [-- <author>]"}
	}
	id, err := store.Add(guildID, text, author)
	if err != nil {
		return quoteCommandResult{Message: "Error: " + err.Error()}
	}
	return quoteCommandResult{Message: fmt.Sprintf("Quote #%d added.", id)}
}

// handleRemoveQuote processes "!quote remove <id>" commands.
func handleRemoveQuote(store *QuoteStore, guildID, args string) quoteCommandResult {
	var id int
	if _, err := fmt.Sscanf(strings.TrimSpace(args), "%d", &id); err != nil {
		return quoteCommandResult{Message: "Usage: !quote remove <id>"}
	}
	if err := store.Remove(guildID, id); err != nil {
		return quoteCommandResult{Message: "Error: " + err.Error()}
	}
	return quoteCommandResult{Message: fmt.Sprintf("Quote #%d removed.", id)}
}

// handleRandomQuote returns a random quote for the guild.
func handleRandomQuote(store *QuoteStore, guildID string) quoteCommandResult {
	q, err := store.Random(guildID)
	if err != nil {
		return quoteCommandResult{Message: "No quotes available. Add one with !quote add <text>."}
	}
	return quoteCommandResult{Message: fmt.Sprintf("📜 #%d: \"%s\" — %s", q.ID, q.Text, q.Author)}
}

// handleListQuotes returns a formatted list of all quotes.
func handleListQuotes(store *QuoteStore, guildID string) quoteCommandResult {
	list := store.List(guildID)
	if len(list) == 0 {
		return quoteCommandResult{Message: "No quotes stored yet."}
	}
	var sb strings.Builder
	sb.WriteString("**Quotes:**\n")
	for _, q := range list {
		sb.WriteString(fmt.Sprintf("#%d: \"%s\" — %s\n", q.ID, q.Text, q.Author))
	}
	return quoteCommandResult{Message: sb.String()}
}
