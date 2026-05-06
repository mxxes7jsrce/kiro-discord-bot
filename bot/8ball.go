package bot

import (
	"math/rand"
)

// eightBallResponses contains classic Magic 8-Ball responses.
var eightBallResponses = []string{
	"It is certain.",
	"It is decidedly so.",
	"Without a doubt.",
	"Yes, definitely.",
	"You may rely on it.",
	"As I see it, yes.",
	"Most likely.",
	"Outlook good.",
	"Yes.",
	"Signs point to yes.",
	"Reply hazy, try again.",
	"Ask again later.",
	"Better not tell you now.",
	"Cannot predict now.",
	"Concentrate and ask again.",
	"Don't count on it.",
	"My reply is no.",
	"My sources say no.",
	"Outlook not so good.",
	"Very doubtful.",
}

// AskEightBall returns a random Magic 8-Ball response for the given question.
// Returns an empty string if the question is empty.
func AskEightBall(question string) string {
	if question == "" {
		return ""
	}
	return eightBallResponses[rand.Intn(len(eightBallResponses))]
}

// isEightBallCommand reports whether the command string is a valid 8ball command.
func isEightBallCommand(cmd string) bool {
	return cmd == "8ball" || cmd == "8b"
}

// handleEightBall processes a Magic 8-Ball question and returns a formatted reply.
func handleEightBall(question string) string {
	if question == "" {
		return "🎱 Please ask a question! Usage: `!8ball <question>`"
	}
	response := AskEightBall(question)
	return "🎱 **" + response + "**"
}
