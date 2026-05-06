package bot

import (
	"math/rand"
)

// CoinFlipResult represents the outcome of a coin flip.
type CoinFlipResult struct {
	Side  string
	Emoji string
}

// FlipCoin simulates a coin flip and returns the result.
func FlipCoin() CoinFlipResult {
	if rand.Intn(2) == 0 {
		return CoinFlipResult{Side: "Heads", Emoji: "🪙"}
	}
	return CoinFlipResult{Side: "Tails", Emoji: "🌕"}
}

// isCoinFlipCommand returns true if the message is a coinflip command.
func isCoinFlipCommand(content string) bool {
	return content == "!coinflip" || content == "!flip" || content == "!coin"
}

// handleCoinFlip processes a coin flip command and returns the response string.
func handleCoinFlip() string {
	result := FlipCoin()
	return result.Emoji + " It's **" + result.Side + "**!"
}

// dispatchCoinFlipCommand handles routing for coinflip-related commands.
func dispatchCoinFlipCommand(content string) (string, bool) {
	if !isCoinFlipCommand(content) {
		return "", false
	}
	return handleCoinFlip(), true
}
