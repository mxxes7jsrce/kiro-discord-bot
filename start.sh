#!/bin/bash

cd "$(dirname "$0")"

# Load .env if exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# If already running, do nothing
if pgrep -f "kiro-discord-bot" > /dev/null 2>&1; then
  echo "bot already running (PID: $(pgrep -f kiro-discord-bot))"
  exit 0
fi

# Kill stale acp-bridge if any
pkill -f "acp-bridged" 2>/dev/null || true
sleep 1

# Build
go build -o /tmp/kiro-discord-bot .

# Start acp-bridge with watchdog
(while true; do
  acp-bridged >> /tmp/acp-bridge.log 2>&1
  echo "[watchdog] acp-bridge exited, restarting in 3s..."
  sleep 3
done) &

# Wait for acp-bridge
for i in $(seq 1 10); do
  curl -sf http://localhost:7800/health > /dev/null && break
  sleep 1
done
echo "acp-bridge ready"

# Start bot with watchdog
(while true; do
  echo "[watchdog] starting bot..."
  /tmp/kiro-discord-bot >> /tmp/kiro-bot.log 2>&1
  echo "[watchdog] bot exited, restarting in 3s..."
  sleep 3
done) &

sleep 5
tail -3 /tmp/kiro-bot.log
