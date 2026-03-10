package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "turso.tech/database/tursogo"
)

var ErrPostNotFound = errors.New("post not found")

type Post struct {
	ID        string    `json:"id" binding:"required"`
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
		SELECT uuid_str(id), title, content, coauthor, created_at
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

func (s *Store) GetPost(ctx context.Context, id string) (Post, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT uuid_str(id), title, content, coauthor, created_at
		FROM posts
		WHERE id = uuid_blob(?)
	`, id)

	post, err := scanPost(row.Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return Post{}, ErrPostNotFound
	}

	return post, err
}

func (s *Store) CreatePost(ctx context.Context, input CreatePostInput) (Post, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return Post{}, err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `BEGIN CONCURRENT`); err != nil {
		return Post{}, err
	}

	committed := false
	defer func() {
		if committed {
			return
		}

		_, _ = conn.ExecContext(context.Background(), `ROLLBACK`)
	}()

	createdAt := time.Now().UTC()
	var id string
	err = conn.QueryRowContext(ctx, `
		INSERT INTO posts (title, content, coauthor, created_at)
		VALUES (?, ?, ?, ?)
		RETURNING uuid_str(id)
	`, input.Title, input.Content, input.Coauthor, createdAt.Format(time.RFC3339Nano)).Scan(&id)
	if err != nil {
		return Post{}, err
	}

	if _, err := conn.ExecContext(ctx, `COMMIT`); err != nil {
		return Post{}, err
	}
	committed = true

	return Post{
		ID:        id,
		Title:     input.Title,
		Content:   input.Content,
		Coauthor:  input.Coauthor,
		CreatedAt: createdAt,
	}, nil
}

func (s *Store) init(ctx context.Context) error {
	if err := s.enableMVCC(ctx); err != nil {
		return err
	}

	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id UUID PRIMARY KEY DEFAULT (uuid7()),
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			coauthor TEXT,
			created_at TEXT NOT NULL
		)
	`)
	return err
}

func (s *Store) enableMVCC(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `PRAGMA journal_mode = 'mvcc'`)
	return err
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
