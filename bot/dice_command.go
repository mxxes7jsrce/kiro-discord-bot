package bot

import (
	"fmt"
	"strings"
)

const diceHelpText = `**Dice Commands**
` + "`!roll <NdM>` — Roll N dice with M sides (e.g. `!roll 2d6`)" + `
` + "`!roll d20` — Roll a single 20-sided die" + `
` + "`!roll` — Roll a single d6 by default"

// handleRollDice processes the !roll command
func handleRollDice(args []string) string {
	notation := "1d6"
	if len(args) > 0 {
		notation = strings.TrimSpace(args[0])
	}

	dice, sides, err := ParseDiceNotation(notation)
	if err != nil {
		return fmt.Sprintf("❌ %s\n%s", err.Error(), diceHelpText)
	}

	result, err := RollDice(dice, sides)
	if err != nil {
		return fmt.Sprintf("❌ %s", err.Error())
	}

	return FormatDiceResult(result)
}

// isDiceCommand returns true if the message is a dice command
func isDiceCommand(content string) bool {
	return strings.HasPrefix(content, "!roll")
}

// dispatchDiceCommand routes dice-related commands
func dispatchDiceCommand(content string) string {
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return diceHelpText
	}

	switch parts[0] {
	case "!roll":
		return handleRollDice(parts[1:])
	default:
		return diceHelpText
	}
}
