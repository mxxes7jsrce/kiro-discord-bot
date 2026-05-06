package bot

import (
	"strings"
)

const jokeHelpText = `**Joke Commands:**
` + "```" + `
!joke          - Tell a random joke
!joke add      - Add a new joke: !joke add <text>
!joke remove   - Remove a joke by ID: !joke remove <id>
!joke list     - List all jokes
!joke help     - Show this help message
` + "```"

// isJokeCommand returns true if the message starts with the joke prefix.
func isJokeCommand(content string) bool {
	return strings.HasPrefix(strings.ToLower(content), "!joke")
}

// dispatchJokeCommand routes a joke command to the appropriate handler.
func dispatchJokeCommand(store *JokeStore, content string) string {
	parts := strings.Fields(content)
	if len(parts) < 2 {
		return handleRandomJoke(store)
	}
	subCmd := strings.ToLower(parts[1])
	args := parts[2:]
	switch subCmd {
	case "add":
		return handleAddJoke(store, args)
	case "remove", "rm", "delete":
		return handleRemoveJoke(store, args)
	case "list", "ls":
		return handleListJokes(store)
	case "help":
		return jokeHelpText
	default:
		return jokeHelpText
	}
}
