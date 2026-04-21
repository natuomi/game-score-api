# game-score-api

ゲームスコアランキングAPI — Go / PostgreSQL / Redis / AWS ECS

**目的:** Goサーバーサイド技術（Go・AWS・RDB/NoSQL）の習得とappworksポートフォリオ掲載

---

## 技術スタック

| レイヤー | 技術 |
|---------|------|
| 言語 | Go 1.22 |
| フレームワーク | Gin |
| RDB | PostgreSQL 16 |
| キャッシュ | Redis 7 |
| 認証 | JWT (golang-jwt) |
| コンテナ | Docker / Docker Compose |
| デプロイ | AWS ECS Fargate + RDS + ElastiCache |
| CI/CD | GitHub Actions |

---

## ローカル起動手順

### 1. 環境変数の設定

```bash
cp .env.example .env
# .env を編集（JWT_SECRET は必ず変更すること）
```

### 2. Docker Compose で起動

```bash
# PostgreSQL・Redis・APIサーバーをまとめて起動
docker compose up --build

# バックグラウンド起動
docker compose up -d --build
```

### 3. 動作確認

```bash
# ヘルスチェック
curl http://localhost:8080/health

# ユーザー登録
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"password123"}'

# ログイン → JWT取得
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"password123"}'

# ランキング取得（認証不要）
curl http://localhost:8080/api/v1/rankings
```

---

## APIエンドポイント一覧

| Method | Path | 認証 | 説明 |
|--------|------|------|------|
| POST | `/api/v1/auth/register` | 不要 | ユーザー登録 |
| POST | `/api/v1/auth/login` | 不要 | ログイン・JWT発行 |
| GET  | `/api/v1/rankings` | 不要 | ランキング取得（Redisキャッシュ） |
| GET  | `/api/v1/players` | 不要 | プレイヤー一覧 |
| GET  | `/api/v1/players/:id` | 不要 | プレイヤー詳細 |
| POST | `/api/v1/scores` | **必要** | スコア登録 |
| GET  | `/api/v1/scores/me` | **必要** | 自分のスコア履歴 |
| GET  | `/health` | 不要 | ヘルスチェック（ECS ALB用） |

---

## フォルダ構成

```
game-score-api/
├── cmd/server/main.go       # エントリーポイント
├── internal/
│   ├── handler/             # HTTPハンドラー
│   ├── service/             # ビジネスロジック
│   ├── repository/          # DB操作
│   └── model/               # データ構造（struct）
├── pkg/
│   ├── auth/jwt.go          # JWT生成・検証
│   ├── database/            # PostgreSQL・Redis接続
│   └── middleware/auth.go   # JWTミドルウェア
├── migrations/              # SQLマイグレーション
├── docker-compose.yml
├── Dockerfile               # マルチステージビルド
└── .env.example
```

---

## 学習フェーズ

- **Phase 1** ローカル開発 — Go + Gin + PostgreSQL + Redis + JWT
- **Phase 2** AWSデプロイ — ECS Fargate + RDS + ElastiCache + GitHub Actions
- **Phase 3** 品質向上 — テスト・Swagger・レートリミット
- **Phase 4** Kubernetes — EKS移行（オプション）

---

## 設計書

`dev/specs/2026-04-20-game-score-api.md` を参照。
