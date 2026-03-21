# kiro-discord-bot

A Discord bot that bridges Discord channels to [kiro-cli](https://kiro.dev) AI agents via [acp-bridge](https://www.npmjs.com/package/acp-bridge). Each channel gets its own agent session and job queue.

**Created:** 2026-03-21 | **Language:** Go

---

## Deployment Guide

### Prerequisites

- Go 1.21+
- Node.js 20+ (for acp-bridge)
- [kiro-cli](https://cli.kiro.dev/install) installed and logged in
- A Discord bot token with the following:
  - Scopes: `bot`, `applications.commands`
  - Permissions: View Channels, Send Messages, Add Reactions, Read Message History
  - Privileged Intents: **Message Content Intent** enabled

---

### 1. Create a Discord Bot

1. Go to [Discord Developer Portal](https://discord.com/developers/applications) → New Application
2. **Bot** tab → Reset Token → copy the token
3. **Bot** tab → Privileged Gateway Intents → enable **Message Content Intent**
4. **OAuth2** tab → URL Generator → select scopes: `bot` + `applications.commands`
5. Select permissions: View Channels, Send Messages, Add Reactions, Read Message History
6. Open the generated URL to invite the bot to your server
7. Note your **Guild ID** (right-click server → Copy Server ID, requires Developer Mode)

---

### 2. Install acp-bridge

```bash
npm install -g acp-bridge
acp-bridged --version
```

---

### 3. Install kiro-cli and Log In

```bash
curl -fsSL https://cli.kiro.dev/install | bash
export PATH="$HOME/.local/bin:$PATH"
kiro-cli login
```

---

### 4. Clone and Configure

```bash
git clone https://github.com/nczz/kiro-discord-bot.git
cd kiro-discord-bot

cp .env.example .env
```

Edit `.env`:

```env
DISCORD_TOKEN=your-bot-token
DISCORD_GUILD_ID=your-guild-id
ACP_BRIDGE_URL=http://localhost:7800
KIRO_CLI_PATH=/home/user/.local/bin/kiro-cli
DEFAULT_CWD=/projects
DATA_DIR=/tmp/kiro-bot-data
ASK_TIMEOUT_SEC=300
STREAM_UPDATE_SEC=3
```

| Variable | Description | Default |
|----------|-------------|---------|
| `DISCORD_TOKEN` | Discord bot token | required |
| `DISCORD_GUILD_ID` | Guild ID for instant slash command registration | required |
| `ACP_BRIDGE_URL` | acp-bridge daemon URL | `http://localhost:7800` |
| `KIRO_CLI_PATH` | Full path to kiro-cli binary | `kiro-cli` |
| `DEFAULT_CWD` | Default working directory for agents | `/projects` |
| `DATA_DIR` | Directory for sessions.json | `./data` |
| `ASK_TIMEOUT_SEC` | Agent response timeout in seconds | `300` |
| `STREAM_UPDATE_SEC` | Discord message update interval during streaming | `3` |

---

### 5. Start (Local)

```bash
chmod +x start.sh
./start.sh
```

The script:
- Skips restart if bot is already running
- Starts acp-bridge with auto-restart watchdog
- Builds and starts the bot with auto-restart watchdog
- Reads all config from `.env`

To force restart:
```bash
pkill -f kiro-discord-bot && pkill -f acp-bridged
./start.sh
```

---

### 6. Deploy with Docker

```bash
docker compose up -d --build
```

`docker-compose.yml` uses `network_mode: host` and mounts `~/.kiro` so the bot inherits your kiro login and MCP settings.

---

### 7. Grant Channel Permissions

The bot needs explicit permission in each channel it should respond to:

1. Right-click the channel → Edit Channel → Permissions
2. Add the bot role/user
3. Enable: View Channel, Send Messages, Add Reactions, Read Message History

---

## Usage

### Commands

| Command | Description |
|---------|-------------|
| `/start <cwd>` | Bind channel to a project directory and start agent |
| `/reset` | Restart the agent for this channel |
| `/status` | Show agent state, queue length, session ID |
| `/cancel` | Cancel the currently running task |
| `/cwd` | Show current working directory |
| `/pause` | Switch to mention-only mode (bot ignores non-mention messages) |
| `/back` | Resume full-listen mode |

All commands also work with `!` prefix (e.g. `!status`, `!reset`).

### Sending Tasks

**Full-listen mode (default):** Any message in the channel is sent to the agent.

**Mention mode (after `/pause`):** Only `@BotName your message` triggers the agent.

### Status Indicators

| Reaction | Meaning |
|----------|---------|
| ⏳ | Queued |
| 🔄 | Processing |
| ⚙️ | Running a tool |
| ✅ | Done |
| ❌ | Error |
| ⚠️ | Timed out |

### Recovery

If a response is cut off, use `!resume` to re-post the agent's last output.

---

## Architecture

```
Discord User
    │ message / slash command
    ▼
Discord Bot (Go)
    ├── per-channel SessionStore   { agentName, sessionId, cwd }
    ├── per-channel JobQueue       buffered chan, FIFO
    └── per-channel Worker         goroutine, sequential execution
          │
          ▼
    AcpClient (HTTP)
          │
          ▼
acp-bridge daemon :7800
          │
          ▼
kiro-cli acp --trust-all-tools   (one process per channel)
          │
          ▼
    AWS Bedrock / Anthropic
```

---

## Project Structure

```
kiro-discord-bot/
├── main.go
├── config.go
├── start.sh              local start + watchdog script
├── bot/
│   ├── bot.go            Discord init, Ready handler, slash command registration
│   └── handler.go        message routing, slash command handlers
├── channel/
│   ├── manager.go        per-channel session + worker lifecycle
│   ├── session.go        session struct + JSON persistence
│   └── worker.go         job queue worker goroutine
├── acp/
│   ├── client.go         acp-bridge HTTP client + SSE stream parser
│   └── sse.go
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```

---

## Notes

- **Session persistence:** Sessions survive as long as the agent process is alive. Bot restart creates a new session (kiro-cli 1.28.1 does not support `session/load`).
- **MCP servers:** Inherited from `~/.kiro/settings/mcp.json` automatically.
- **Project steering:** Add `.kiro/steering/*.md` in the project directory to guide agent behavior.
- **Long responses:** Automatically split into multiple messages with `(1/N)` labels.

---

---

## 部署說明（中文）

### 前置需求

- Go 1.21+
- Node.js 20+（用於 acp-bridge）
- 已安裝並登入 [kiro-cli](https://cli.kiro.dev/install)
- Discord bot token，需具備：
  - Scopes：`bot`、`applications.commands`
  - 權限：查看頻道、發送訊息、新增反應、讀取訊息歷史
  - Privileged Intents：啟用 **Message Content Intent**

### 快速開始

```bash
# 1. 安裝 acp-bridge
npm install -g acp-bridge

# 2. 安裝並登入 kiro-cli
curl -fsSL https://cli.kiro.dev/install | bash
kiro-cli login

# 3. 設定環境變數
cp .env.example .env
# 編輯 .env，填入 DISCORD_TOKEN、DISCORD_GUILD_ID、KIRO_CLI_PATH 等

# 4. 啟動
chmod +x start.sh && ./start.sh
```

### 指令說明

| 指令 | 說明 |
|------|------|
| `/start <目錄>` | 綁定專案目錄並啟動 agent |
| `/reset` | 重啟此 channel 的 agent |
| `/status` | 查詢 agent 狀態、queue 長度 |
| `/cancel` | 取消目前執行中的任務 |
| `/cwd` | 查詢目前工作目錄 |
| `/pause` | 切換為 @mention 模式 |
| `/back` | 恢復完整監聽模式 |

所有指令也支援 `!` 前綴（如 `!status`、`!reset`）。

### 注意事項

- Bot 需要在各 channel 的權限設定中明確授予讀寫權限
- Session 在 agent process 存活期間持續，bot 重啟後建立新 session
- MCP 設定自動繼承 `~/.kiro/settings/mcp.json`
- 回應被截斷時可用 `!resume` 補完
