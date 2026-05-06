package bot

import "strings"

// handleMagic8BallCommand processes the !8ball command and returns a response.
// It delegates to the AskEightBall and handleEightBall functions in 8ball.go.
func handleMagic8BallCommand(args []string) string {
	question := strings.Join(args, " ")
	return handleEightBall(question)
}

// dispatchMagic8BallCommand routes the !8ball command with its arguments.
func dispatchMagic8BallCommand(parts []string) string {
	if len(parts) < 2 {
		return handleEightBall("")
	}
	return handleMagic8BallCommand(parts[1:])
}
