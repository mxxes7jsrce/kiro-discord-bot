---
name: discord-mcp
description: Discord MCP server that gives this agent the ability to interact with Discord channels. If the discord_* tools are not available, read INSTALL_MCP.md and follow the steps to install.
---

# Discord MCP Server

## 檢查是否已啟用

如果你的 tool list 中有 `discord_list_channels`、`discord_send_message` 等 `discord_*` 開頭的 tools，表示已啟用，可直接使用。

## 尚未啟用時

1. 讀取專案根目錄的 `INSTALL_MCP.md`
2. 用 `pwd` 取得當前專案的絕對路徑，作為文件中 `<PROJECT_DIR>` 的值
3. 依照文件步驟執行 build 和註冊，將所有 `<PROJECT_DIR>` 替換為實際路徑

## 使用方式

每則訊息的 prompt 開頭帶有 `[Discord context] channel_id=... guild_id=...`，直接用這些 ID 呼叫 discord_* tools。
