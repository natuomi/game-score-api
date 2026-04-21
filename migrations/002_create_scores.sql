-- スコアテーブル作成
CREATE TABLE IF NOT EXISTS scores (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    score      INTEGER     NOT NULL CHECK (score >= 0),
    game_mode  VARCHAR(50) NOT NULL DEFAULT 'classic',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ランキング取得クエリ用インデックス（高スコア順）
CREATE INDEX IF NOT EXISTS idx_scores_score_desc ON scores(score DESC);

-- ユーザー別スコア一覧用インデックス
CREATE INDEX IF NOT EXISTS idx_scores_user_id ON scores(user_id);

-- ゲームモード別フィルタ用インデックス
CREATE INDEX IF NOT EXISTS idx_scores_game_mode ON scores(game_mode);
