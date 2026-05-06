package bot

import (
	"fmt"
	"strconv"
	"strings"
)

// CountdownResult holds the result of a countdown formatting operation.
type CountdownResult struct {
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

// ParseCountdownSeconds converts a raw seconds value into a CountdownResult.
func ParseCountdownSeconds(totalSeconds int) (CountdownResult, error) {
	if totalSeconds < 0 {
		return CountdownResult{}, fmt.Errorf("seconds must be non-negative")
	}
	days := totalSeconds / 86400
	remaining := totalSeconds % 86400
	hours := remaining / 3600
	remaining = remaining % 3600
	minutes := remaining / 60
	seconds := remaining % 60
	return CountdownResult{
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}, nil
}

// FormatCountdown returns a human-readable countdown string.
func FormatCountdown(r CountdownResult) string {
	parts := []string{}
	if r.Days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", r.Days))
	}
	if r.Hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", r.Hours))
	}
	if r.Minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", r.Minutes))
	}
	parts = append(parts, fmt.Sprintf("%ds", r.Seconds))
	return strings.Join(parts, " ")
}

// isCountdownCommand returns true if the message starts with !countdown.
func isCountdownCommand(content string) bool {
	return strings.HasPrefix(strings.ToLower(content), "!countdown")
}

// handleCountdownCommand parses the seconds argument and returns formatted output.
func handleCountdownCommand(args []string) string {
	if len(args) == 0 {
		return "Usage: !countdown <seconds>"
	}
	sec, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Sprintf("Invalid seconds value: %s", args[0])
	}
	result, err := ParseCountdownSeconds(sec)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}
	return fmt.Sprintf("⏳ Countdown: %s", FormatCountdown(result))
}

// dispatchCountdownCommand routes the !countdown command.
func dispatchCountdownCommand(parts []string) string {
	if len(parts) < 2 {
		return "Usage: !countdown <seconds>"
	}
	return handleCountdownCommand(parts[1:])
}
