package bot

import (
	"strings"
)

// isWeatherCommand returns true if the message is a weather command.
func isWeatherCommand(content string) bool {
	return strings.HasPrefix(content, "!weather")
}

// handleWeatherCommand processes the !weather command and returns a response.
func handleWeatherCommand(content string) string {
	parts := strings.SplitN(content, " ", 2)
	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "Usage: `!weather <location>` — e.g. `!weather Tokyo`"
	}

	location := strings.TrimSpace(parts[1])
	w, err := GetWeather(location)
	if err != nil {
		return "❌ Could not retrieve weather: " + err.Error()
	}

	return FormatWeatherReport(location, w)
}

// dispatchWeatherCommand routes weather subcommands.
func dispatchWeatherCommand(content string) string {
	if !isWeatherCommand(content) {
		return ""
	}
	return handleWeatherCommand(content)
}
