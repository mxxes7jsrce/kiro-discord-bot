package bot

import (
	"fmt"
	"strconv"
	"strings"
)

// handleRandomJoke sends a random joke to the channel.
func handleRandomJoke(store *JokeStore) string {
	j, err := store.Random()
	if err != nil {
		return "😔 No jokes available yet. Add one with `!joke add <text>`!"
	}
	return fmt.Sprintf("😂 **Joke #%d:** %s", j.ID, j.Text)
}

// handleAddJoke adds a new joke to the store.
func handleAddJoke(store *JokeStore, args []string) string {
	if len(args) == 0 {
		return "Usage: `!joke add <joke text>`"
	}
	text := strings.Join(args, " ")
	id, err := store.Add(text)
	if err != nil {
		return fmt.Sprintf("❌ Error: %s", err.Error())
	}
	return fmt.Sprintf("✅ Joke added with ID **%d**!", id)
}

// handleRemoveJoke removes a joke by ID.
func handleRemoveJoke(store *JokeStore, args []string) string {
	if len(args) == 0 {
		return "Usage: `!joke remove <id>`"
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return "❌ Invalid ID. Please provide a numeric joke ID."
	}
	if err := store.Remove(id); err != nil {
		return fmt.Sprintf("❌ %s", err.Error())
	}
	return fmt.Sprintf("🗑️ Joke #%d removed.", id)
}

// handleListJokes returns a formatted list of all jokes.
func handleListJokes(store *JokeStore) string {
	jokes := store.List()
	if len(jokes) == 0 {
		return "📭 No jokes in the store."
	}
	var sb strings.Builder
	sb.WriteString("**Jokes:**\n")
	for _, j := range jokes {
		sb.WriteString(fmt.Sprintf("**#%d** %s\n", j.ID, j.Text))
	}
	return strings.TrimRight(sb.String(), "\n")
}
