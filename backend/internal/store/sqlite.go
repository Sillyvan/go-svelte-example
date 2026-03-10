package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "modernc.org/sqlite"
)

var ErrPostNotFound = errors.New("post not found")

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePostInput struct {
	Title   string
	Content string
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.init(context.Background()); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) ListPosts(ctx context.Context) ([]Post, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, content, created_at
		FROM posts
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]Post, 0)
	for rows.Next() {
		post, err := scanPost(rows.Scan)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func (s *SQLiteStore) GetPost(ctx context.Context, id int64) (Post, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, title, content, created_at
		FROM posts
		WHERE id = ?
	`, id)

	post, err := scanPost(row.Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return Post{}, ErrPostNotFound
	}

	return post, err
}

func (s *SQLiteStore) CreatePost(ctx context.Context, input CreatePostInput) (Post, error) {
	createdAt := time.Now().UTC()
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO posts (title, content, created_at)
		VALUES (?, ?, ?)
	`, input.Title, input.Content, createdAt.Format(time.RFC3339Nano))
	if err != nil {
		return Post{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Post{}, err
	}

	return Post{
		ID:        id,
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: createdAt,
	}, nil
}

func (s *SQLiteStore) init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TEXT NOT NULL
		)
	`)
	return err
}

type scanner func(dest ...any) error

func scanPost(scan scanner) (Post, error) {
	var post Post
	var createdAt string
	if err := scan(&post.ID, &post.Title, &post.Content, &createdAt); err != nil {
		return Post{}, err
	}

	parsedCreatedAt, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return Post{}, err
	}

	post.CreatedAt = parsedCreatedAt
	return post, nil
}
