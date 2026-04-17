// Package config handles loading and validating bot configuration from environment variables.
package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values for the bot.
type Config struct {
	// DiscordToken is the bot token used to authenticate with Discord.
	DiscordToken string

	// CommandPrefix is the prefix used to trigger bot commands (e.g. "!").
	CommandPrefix string

	// GuildID is the optional Discord server ID for guild-specific commands.
	// If empty, commands are registered globally.
	GuildID string

	// OpenAIAPIKey is used for AI-powered features, if enabled.
	OpenAIAPIKey string
}

// Load reads configuration from a .env file and environment variables.
// Environment variables take precedence over .env file values.
func Load() (*Config, error) {
	// Attempt to load .env file; ignore error if it doesn't exist
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, errors.New("DISCORD_TOKEN environment variable is required")
	}

	prefix := os.Getenv("COMMAND_PREFIX")
	if prefix == "" {
		prefix = "!"
	}

	return &Config{
		DiscordToken:  token,
		CommandPrefix: prefix,
		GuildID:       os.Getenv("GUILD_ID"),
		OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),
	}, nil
}
