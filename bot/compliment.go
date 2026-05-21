package bot

import "math/rand"

var defaultCompliments = []string{
	"You're doing an amazing job! 🌟",
	"You bring so much positivity to this server! ✨",
	"You're incredibly talented and creative! 🎨",
	"Your kindness makes a real difference! 💙",
	"You're stronger than you think! 💪",
	"The world is better with you in it! 🌍",
	"You have a great sense of humor! 😄",
	"You're a wonderful person! 🌺",
	"Your hard work is truly inspiring! 🚀",
	"You make this community awesome! 🎉",
}

// GetCompliment returns a random compliment, optionally targeted at a user.
func GetCompliment(target string) string {
	c := defaultCompliments[rand.Intn(len(defaultCompliments))]
	if target != "" {
		return target + ", " + c
	}
	return c
}

// isComplimentCommand reports whether the message content is a compliment command.
func isComplimentCommand(content string) bool {
	return len(content) >= 10 && content[:10] == "!compliment"
}

// handleComplimentCommand parses args and returns a compliment response.
func handleComplimentCommand(args []string) string {
	if len(args) > 0 {
		return GetCompliment(args[0])
	}
	return GetCompliment("")
}

// dispatchComplimentCommand is the top-level dispatcher for !compliment.
func dispatchComplimentCommand(content string) string {
	args := parseArgs(content, "!compliment")
	return handleComplimentCommand(args)
}
