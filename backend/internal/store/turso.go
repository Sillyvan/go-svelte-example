package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "turso.tech/database/tursogo"
)

var ErrPostNotFound = errors.New("post not found")

type Post struct {
	ID        int64     `json:"id" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	Coauthor  *string   `json:"coauthor,omitempty"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
}

type CreatePostInput struct {
	Title    string
	Content  string
	Coauthor *string
}

type Store struct {
	db *sql.DB
}

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("turso", path)
	if err != nil {
		return nil, err
	}

	store := &Store{db: db}
	if err := store.init(context.Background()); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) ListPosts(ctx context.Context) ([]Post, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, content, coauthor, created_at
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

func (s *Store) GetPost(ctx context.Context, id int64) (Post, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, title, content, coauthor, created_at
		FROM posts
		WHERE id = ?
	`, id)

	post, err := scanPost(row.Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return Post{}, ErrPostNotFound
	}

	return post, err
}

func (s *Store) CreatePost(ctx context.Context, input CreatePostInput) (Post, error) {
	createdAt := time.Now().UTC()
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO posts (title, content, coauthor, created_at)
		VALUES (?, ?, ?, ?)
	`, input.Title, input.Content, input.Coauthor, createdAt.Format(time.RFC3339Nano))
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
		Coauthor:  input.Coauthor,
		CreatedAt: createdAt,
	}, nil
}

func (s *Store) init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			coauthor TEXT,
			created_at TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `
		ALTER TABLE posts
		ADD COLUMN coauthor TEXT
	`)
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		return err
	}

	return nil
}

type scanner func(dest ...any) error

func scanPost(scan scanner) (Post, error) {
	var post Post
	var coauthor sql.NullString
	var createdAt string
	if err := scan(&post.ID, &post.Title, &post.Content, &coauthor, &createdAt); err != nil {
		return Post{}, err
	}

	if coauthor.Valid {
		post.Coauthor = &coauthor.String
	}

	parsedCreatedAt, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return Post{}, err
	}

	post.CreatedAt = parsedCreatedAt
	return post, nil
}
