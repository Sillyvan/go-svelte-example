CREATE TABLE IF NOT EXISTS posts (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    coauthor TEXT,
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts (created_at DESC);
