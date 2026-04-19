# Discord MCP Server — Install & Enable Guide

This project includes a Go-based Discord MCP Server (`cmd/mcp-discord/`) that gives the kiro agent direct access to Discord — read messages, send messages, list channels, search, add reactions, etc.

## Install Steps

> **Note:** Replace `<PROJECT_DIR>` with the absolute path of this project directory (run `pwd` to get it).

### Step 1: Build the binary

Run from the project root:

```bash
go build -o mcp-discord-server ./cmd/mcp-discord/
```

### Step 2: Install the steering file

Copy the steering file to the global kiro steering directory so the agent loads it regardless of working directory:

```bash
mkdir -p ~/.kiro/steering
cp .kiro/steering/discord-mcp.md ~/.kiro/steering/discord-mcp.md
```

### Step 3: Register the MCP Server

Edit `~/.kiro/settings/mcp.json` and add the following entry under `"mcpServers"` (keep existing entries intact):

```json
"mcp-discord": {
  "command": "sh",
  "args": [
    "-c",
    "set -a && . <PROJECT_DIR>/.env && exec <PROJECT_DIR>/mcp-discord-server"
  ]
}
```

Replace `<PROJECT_DIR>` with the actual absolute path.

For example, if the project is at `/home/user/kiro-discord-bot`:

```json
"mcp-discord": {
  "command": "sh",
  "args": [
    "-c",
    "set -a && . /home/user/kiro-discord-bot/.env && exec /home/user/kiro-discord-bot/mcp-discord-server"
  ]
}
```

This sources `DISCORD_TOKEN` from the project `.env` file at startup — no token duplication needed.

> **Tip (macOS):** If you're on macOS, the binary path might need to be absolute and the `.env` file must not contain any `export` keywords — just `KEY=VALUE` pairs.

> **Tip (Linux/WSL):** On Linux or WSL, make sure the binary is executable after building: `chmod +x mcp-discord-server`

> **Personal note:** I keep my `.env` file outside the repo directory (e.g. `~/.config/kiro-discord-bot/.env`) and symlink it here with `ln -s ~/.config/kiro-discord-bot/.env .env`. This avoids accidentally committing secrets even if `.gitignore` is misconfigured.

> **Personal note:** I also add `mcp-discord-server` to `.git/info/exclude` locally so the built binary never shows up as an untracked file without touching the shared `.gitignore`.

> **Personal note:** I created a small shell alias `alias build-mcp='go build -o mcp-discord-server ./cmd/mcp-discord/ && chmod +x mcp-discord-server'` in my `.bashrc` so I don't have to remember the full build command after pulling upstream changes.

> **Personal note:** If the agent fails to connect after restart, check that the MCP server process is actually running with `pgrep -a mcp-discord-server`. I've occasionally had stale processes from a previous session interfere — `pkill mcp-discord-server` clears it up before retrying.

> **Personal note:** I run the MCP server behind a simple health-check wrapper script that auto-restarts it if it crashes unexpectedly. Saved me a few times when the Discord gateway dropped the connection overnight. The wrapper is just a `while true; do ./mcp-discord-server; sleep 2; done` loop in a `tmux` session.

### Step 4: Restart the agent session

After completing the steps above, tell the user:

> MCP server installed. Please run `/reset` or `!reset` to restart the agent session. After restart, I'll have direct access to Discord.

## Available Tools (after e