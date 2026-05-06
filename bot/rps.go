package bot

import (
	"fmt"
	"math/rand"
	"strings"
)

type RPSChoice int

const (
	Rock RPSChoice = iota
	Paper
	Scissors
)

var rpsNames = map[RPSChoice]string{
	Rock:     "rock",
	Paper:    "paper",
	Scissors: "scissors",
}

var rpsEmojis = map[RPSChoice]string{
	Rock:     "🪨",
	Paper:    "📄",
	Scissors: "✂️",
}

type RPSResult struct {
	PlayerChoice RPSChoice
	BotChoice    RPSChoice
	Outcome      string
}

func parseRPSChoice(input string) (RPSChoice, bool) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "rock", "r":
		return Rock, true
	case "paper", "p":
		return Paper, true
	case "scissors", "s":
		return Scissors, true
	}
	return Rock, false
}

func determineRPSOutcome(player, bot RPSChoice) string {
	if player == bot {
		return "tie"
	}
	if (player == Rock && bot == Scissors) ||
		(player == Paper && bot == Rock) ||
		(player == Scissors && bot == Paper) {
		return "win"
	}
	return "lose"
}

func PlayRPS(playerInput string) (RPSResult, error) {
	if playerInput == "" {
		return RPSResult{}, fmt.Errorf("please provide a choice: rock, paper, or scissors")
	}
	player, ok := parseRPSChoice(playerInput)
	if !ok {
		return RPSResult{}, fmt.Errorf("invalid choice %q — use rock, paper, or scissors", playerInput)
	}
	bot := RPSChoice(rand.Intn(3))
	outcome := determineRPSOutcome(player, bot)
	return RPSResult{PlayerChoice: player, BotChoice: bot, Outcome: outcome}, nil
}

func FormatRPSResult(r RPSResult) string {
	pName := rpsNames[r.PlayerChoice]
	pEmoji := rpsEmojis[r.PlayerChoice]
	bName := rpsNames[r.BotChoice]
	bEmoji := rpsEmojis[r.BotChoice]

	switch r.Outcome {
	case "win":
		return fmt.Sprintf("You chose %s %s, I chose %s %s — **You win!** 🎉", pEmoji, pName, bEmoji, bName)
	case "lose":
		return fmt.Sprintf("You chose %s %s, I chose %s %s — **You lose!** 😈", pEmoji, pName, bEmoji, bName)
	default:
		return fmt.Sprintf("You chose %s %s, I chose %s %s — **It's a tie!** 🤝", pEmoji, pName, bEmoji, bName)
	}
}

func isRPSCommand(args []string) bool {
	return len(args) > 0 && strings.ToLower(args[0]) == "rps"
}

func handleRPSCommand(args []string) string {
	if len(args) < 2 {
		return "Usage: `!rps <rock|paper|scissors>` — Challenge me to a game!"
	}
	result, err := PlayRPS(args[1])
	if err != nil {
		return fmt.Sprintf("❌ %s", err.Error())
	}
	return FormatRPSResult(result)
}

func dispatchRPSCommand(args []string) string {
	if isRPSCommand(args) {
		return handleRPSCommand(args)
	}
	return ""
}
