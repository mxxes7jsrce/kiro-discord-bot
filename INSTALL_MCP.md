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

### Step 4: Restart the agent session

After completing the steps above, tell the user:

> MCP server installed. Please run `/reset` or `!reset` to restart the agent session. After restart, I'll have direct access to Discord.

## Available Tools (after enabled)

| Tool | Description |
|------|-------------|
| `discord_list_channels` | List text channels in a guild |
| `discord_read_messages` | Read recent messages from a channel |
| `discord_send_message` | Send a message to a channel |
| `discord_reply_message` | Reply to a specific message |
| `discord_add_reaction` | Add a reaction emoji to a message |
| `discord_list_members` | List members of a guild |
| `discord_search_messages` | Search recent messages by keyword |
| `discord_channel_info` | Get detailed info about a channel |
| `discord_send_file` | Upload a local file to a channel as an attachment |
| `discord_list_attachments` | List file attachments from recent messages |
| `discord_download_attachment` | Download a Discord attachment to a local file |
| `discord_edit_message` | Edit a message |
| `discord_delete_message` | Delete a message |
| `discord_get_message` | Get a single message by ID |
| `discord_send_embed` | Send a rich embed message |
| `discord_pin_message` | Pin or unpin a message |
| `discord_create_thread` | Create a thread from a message |
| `discord_list_threads`
