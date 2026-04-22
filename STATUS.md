# game-score-api STATUS

## 最終更新
- 日時: 2026-04-22 11:30
- 環境: ターミナル
- トリガー: ポートフォリオ強化（README/Swagger/GitHub Actions）

---

## 現在のフェーズ

**Phase 2 完了・AWS EC2 公開中 🚀**
- 公開URL: `http://13.114.19.198:8080`
- ヘルスチェック: `http://13.114.19.198:8080/health`

---

## 背景・目的

クラウドワークス「フリーランススタート」の Goサーバーサイドエンジニア案件（80〜120万/月）を
ターゲットに技術習得するためのポートフォリオアプリ。

現状スキルセット（Unity / Web改修）とのギャップを埋めるために作成。
`sales/freelance-start/prospects/2026-04-20-go-game-server-engineer.md` に案件詳細あり。

---

## 完了済み

- 2026-04-22: ポートフォリオ強化 (PORT-1)
  - README.md 英語化・アーキテクチャ図・API表・技術スタック表・CI/CD Secrets一覧を追加
  - docs/swagger.yaml 作成（OpenAPI 3.0、全8エンドポイント・スキーマ定義・認証記述）
  - .github/workflows/deploy.yml 作成（Go 1.23ビルド→SCP転送→SSH再起動→ヘルスチェック）
- 2026-04-20: 設計書作成 (`dev/specs/2026-04-20-game-score-api.md`)
- 2026-04-20: フォルダ構成・スキャフォールド作成
- 2026-04-21: Phase 1 実装完了（全ファイル）
- 2026-04-21: **Phase 1 動作確認完了** ✅
  - `GET /health` → `{"status":"ok","version":"1.0.0"}`
  - `POST /api/v1/auth/register` → ユーザー登録・UUID発行
  - `POST /api/v1/auth/login` → JWT発行確認
  - `POST /api/v1/scores` → JWT認証・スコア登録確認
  - `GET /api/v1/rankings` → 1回目 `cached:false`（DB）、2回目 `cached:true`（Redis） ✅
  - `GET /api/v1/scores/me` → 認証済みユーザーのスコア一覧
  - `GET /api/v1/players` → プレイヤー一覧

---

## 実装ファイル一覧

| ファイル | 内容 | 状態 |
|---------|------|------|
| `cmd/server/main.go` | エントリーポイント・依存注入・ルーティング | ✅ |
| `internal/model/user.go` | User構造体・Request/Response型 | ✅ |
| `internal/model/score.go` | Score構造体・RankingResponse型 | ✅ |
| `internal/repository/user_repo.go` | Create / FindByEmail / FindByID / FindAll | ✅ |
| `internal/repository/score_repo.go` | Create / FindByUserID / GetRankings / CountAll | ✅ |
| `internal/service/auth_service.go` | Register（bcrypt）/ Login（JWT発行） | ✅ |
| `internal/service/score_service.go` | PostScore / GetMyScores / GetRankings（Redisキャッシュ） | ✅ |
| `internal/handler/auth.go` | POST /register / POST /login | ✅ |
| `internal/handler/score.go` | POST /scores / GET /scores/me / GET /rankings | ✅ |
| `internal/handler/player.go` | GET /players / GET /players/:id | ✅ |
| `pkg/auth/jwt.go` | JWT生成・検証 | ✅ |
| `pkg/database/postgres.go` | PostgreSQL接続プール | ✅ |
| `pkg/database/redis.go` | Redis接続 | ✅ |
| `pkg/middleware/auth.go` | JWTミドルウェア | ✅ |
| `migrations/001_create_users.sql` | usersテーブル定義 | ✅ |
| `migrations/002_create_scores.sql` | scoresテーブル定義 + インデックス | ✅ |
| `docker-compose.yml` | PostgreSQL + Redis（APIはローカル起動） | ✅ |
| `Dockerfile` | マルチステージビルド（alpine） | ✅ |
| `.env` | 環境変数（ASCII-only UTF-8） | ✅ |

---

## 進行中

なし

---

## 次のステップ

1. GitHub Secrets（EC2_HOST / EC2_USER / EC2_SSH_KEY）をリポジトリに設定し、GitHub Actions CI/CDを有効化
2. Phase 4: ユニットテスト・統合テスト追加
3. Phase 5: EKS移行（オプション）

---

## ローカル起動手順（再開時）

```powershell
# ① Docker コンテナ起動（PostgreSQL + Redis）
cd "M:\Aiplay\cursor\.company\test company 1\.company\webapp\projects\game-score-api"
docker compose up -d postgres redis

# ② APIサーバー起動（環境変数を明示的にセット）
$env:PORT="8080"; $env:GIN_MODE="debug"; $env:JWT_SECRET="game-score-dev-secret-2026"
$env:DB_HOST="localhost"; $env:DB_PORT="5432"; $env:DB_USER="postgres"
$env:DB_PASSWORD="postgres"; $env:DB_NAME="game_score"
$env:REDIS_HOST="localhost"; $env:REDIS_PORT="6379"; $env:REDIS_PASSWORD=""
$env:RANKING_CACHE_TTL="60"; $env:RANKING_LIMIT="100"
& "C:\Program Files\Go\bin\go.exe" run ./cmd/server

# ③ ヘルスチェック
Invoke-RestMethod -Uri "http://localhost:8080/health"
```

> **注意**: Docker APIコンテナは runc の相対パス問題で起動不可。APIは `go run` でローカル実行する。

---

## 技術スタック

| 技術 | バージョン | 用途 |
|------|----------|------|
| Go | 1.23 | アプリケーション本体 |
| Gin | v1.9.1 | HTTPフレームワーク |
| PostgreSQL | 16 | ユーザー・スコアデータ |
| Redis | 7 | ランキングキャッシュ |
| golang-jwt | v5.2.1 | JWT認証 |
| pgx | v5.5.5 | PostgreSQLドライバ |
| go-redis | v9.5.1 | Redisクライアント |
| bcrypt | golang.org/x/crypto | パスワードハッシュ |

---

## AWS EC2 情報
- インスタンス: t3.micro（東京リージョン ap-northeast-1）
- パブリックIP: 13.114.19.198
- ポート: 8080
- 起動方法: `nohup ./bin/server-linux > server.log 2>&1 &`
- DB/Redis: `sudo docker-compose up -d postgres redis`

## Git状態
- ブランチ: main
- リポジトリ: https://github.com/natuomi/game-score-api
- 最終コミット: feat: add pre-built Linux binary
- 最終プッシュ: 2026-04-22 00:10

## 関連ファイル
- 設計書: `.company/dev/specs/2026-04-20-game-score-api.md`
- 案件メモ: `.company/sales/freelance-start/prospects/2026-04-20-go-game-server-engineer.md`
- パイプライン: `.company/sales/freelance-start/pipeline.md`

## 作業ログ（自動追記）

- 2026-04-20 21:47 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md
- 2026-04-21 07:01 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md
- 2026-04-21 22:30 / ターミナル / Phase 1 動作確認完了。全エンドポイント・Redisキャッシュ動作確認済み。
- 2026-04-21 22:30 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go
- 2026-04-21 22:36 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go
- 2026-04-21 22:37 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go
- 2026-04-21 22:38 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go
- 2026-04-21 22:41 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go
- 2026-04-21 22:44 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 22:53 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 22:54 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:14 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:15 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:16 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:17 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:17 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:18 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:19 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:19 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:22 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:22 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:23 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:24 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:24 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:24 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:25 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:26 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:26 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:28 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:28 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:29 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:29 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:30 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:32 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:32 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:35 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-21 23:44 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 07:12 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 07:13 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 07:13 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 07:14 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 07:16 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 07:16 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 07:19 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 07:20 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 07:20 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 09:03 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 09:05 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 09:08 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 09:08 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 09:08 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 09:09 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 09:09 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 09:09 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 09:10 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
- 2026-04-22 09:12 / ターミナル / edited: go.mod, cmd/server/main.go, internal/model/user.go, internal/model/score.go, pkg/database/postgres.go, pkg/database/redis.go, pkg/auth/jwt.go, docker-compose.yml, Dockerfile, .env.example, migrations/001_create_users.sql, migrations/002_create_scores.sql, pkg/middleware/auth.go, README.md, internal/repository/user_repo.go, internal/repository/score_repo.go, internal/service/auth_service.go, internal/service/score_service.go, internal/handler/auth.go, internal/handler/score.go, internal/handler/player.go, .gitignore
