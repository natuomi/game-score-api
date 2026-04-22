#!/bin/bash
# game-score-api 起動スクリプト
# 使い方: bash scripts/start.sh
# 環境変数は .env ファイルから読み込む（ps aux に秘密鍵が表示されない）

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# .env ファイルを読み込む
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
else
  echo "ERROR: .env file not found at $PROJECT_DIR/.env"
  exit 1
fi

# 既存プロセスを停止
pkill -f server-linux 2>/dev/null || true
sleep 1

# バックグラウンドで起動
nohup ./bin/server-linux > server.log 2>&1 &
echo "Server started (PID: $!)"
echo "Log: $PROJECT_DIR/server.log"
