# ── Stage 1: ビルド ──────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 依存解決（キャッシュ活用のため先にコピー）
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# ── Stage 2: 実行 ─────────────────────────────────────
# alpineイメージで最小サイズに（ECS Fargate のコールドスタート短縮）
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Tokyo

WORKDIR /app
COPY --from=builder /app/server .

# ECS のヘルスチェックポート
EXPOSE 8080

CMD ["/app/server"]
