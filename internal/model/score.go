package model

import "time"

// Score はスコアレコードを表す
type Score struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Score     int       `json:"score" db:"score"`
	GameMode  string    `json:"game_mode" db:"game_mode"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// JOIN 用（ランキング表示時に使用）
	PlayerName string `json:"player_name,omitempty" db:"player_name"`
}

// PostScoreRequest はスコア登録リクエスト
type PostScoreRequest struct {
	Score    int    `json:"score" binding:"required,min=0"`
	GameMode string `json:"game_mode" binding:"omitempty,max=50"`
}

// RankingEntry はランキング表示用エントリ
type RankingEntry struct {
	Rank       int       `json:"rank"`
	PlayerName string    `json:"player_name"`
	Score      int       `json:"score"`
	GameMode   string    `json:"game_mode"`
	CreatedAt  time.Time `json:"created_at"`
}

// RankingResponse はランキングAPIのレスポンス
type RankingResponse struct {
	Rankings []RankingEntry `json:"rankings"`
	Cached   bool           `json:"cached"` // Redis キャッシュから返した場合 true
	Total    int            `json:"total"`
}
