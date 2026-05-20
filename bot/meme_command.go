package bot

import (
	"strings"
)

func isMemeCommand(content string) bool {
	lower := strings.ToLower(strings.TrimSpace(content))
	return lower == "!meme" ||
		strings.HasPrefix(lower, "!meme ") ||
		strings.HasPrefix(lower, "!meme\t")
}

// handleMemeCommand parses the meme command and returns a response string.
// Usage:
//
//	!meme                          — random meme, no args
//	!meme list                     — list available templates
//	!meme <template> <arg1> [arg2] — generate a specific meme
//	!meme random <arg1> [arg2]     — generate a random meme with args
func handleMemeCommand(args []string) string {
	if len(args) == 0 {
		result, err := GenerateMeme("random", []string{"something", "something else"})
		if err != nil {
			return "🖼️ Could not generate meme."
		}
		return "🖼️ " + result
	}

	subcmd := strings.ToLower(args[0])

	if subcmd == "list" {
		names := GetMemeTemplateNames()
		return "🖼️ Available meme templates: " + strings.Join(names, ", ")
	}

	// Treat args[0] as template name, rest as meme args
	memeArgs := args[1:]
	result, err := GenerateMeme(subcmd, memeArgs)
	if err != nil {
		return "❌ " + err.Error() + ". Use `!meme list` to see available templates."
	}
	return "🖼️ " + result
}

func dispatchMemeCommand(content string) string {
	trimmed := strings.TrimSpace(content)
	parts := strings.Fields(trimmed)
	// parts[0] == "!meme"
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}
	return handleMemeCommand(args)
}
