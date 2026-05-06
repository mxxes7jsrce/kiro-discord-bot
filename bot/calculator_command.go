package bot

import (
	"strings"
)

// isCalcCommand returns true if the message starts with !calc or !calculate.
func isCalcCommand(content string) bool {
	lower := strings.ToLower(content)
	return strings.HasPrefix(lower, "!calc ") ||
		strings.HasPrefix(lower, "!calculate ") ||
		lower == "!calc" ||
		lower == "!calculate"
}

// handleCalcCommand processes a !calc command and returns a response string.
func handleCalcCommand(content string) string {
	parts := strings.SplitN(content, " ", 2)
	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "Usage: `!calc <expression>` — e.g. `!calc 3 + 4`, `!calc 10 / 2`, `!calc 2 ^ 8`"
	}

	expr := strings.TrimSpace(parts[1])
	result, err := EvaluateSimple(expr)
	if err != nil {
		return "❌ Could not evaluate expression: " + err.Error()
	}
	return FormatCalcResult(result)
}

// dispatchCalcCommand is the entry-point dispatcher for calculator commands.
func dispatchCalcCommand(content string) string {
	return handleCalcCommand(content)
}
