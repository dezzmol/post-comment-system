package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"post-comment-system/graph/model"
)

type PostPostgresRepo struct {
	db *sql.DB
}

type postDB struct {
	ID            string    `db:"id"`
	Title         string    `db:"title"`
	Content       string    `db:"content"`
	CreatedAt     time.Time `db:"created_at"`
	AllowComments bool      `db:"allow_comments"`
	Username      *string   `db:"name"`
	AuthorId      *int      `db:"author_id"`
}

func NewPostPostgresRepository(db *sql.DB) *PostPostgresRepo {
	return &PostPostgresRepo{
		db: db,
	}
}

func (r *PostPostgresRepo) GetAllPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	query := `
		SELECT 
			posts.id, posts.title, posts.content, posts.created_at, posts.allow_comments,
			users.name, users.id AS author_id
		FROM posts
		JOIN users ON posts.author_id = users.id
		ORDER BY posts.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, *limit, *offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.Post
	for rows.Next() {
		var p postDB
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.AllowComments, &p.Username, &p.AuthorId); err != nil {
			return nil, err
		}

		post := &model.Post{
			ID:            p.ID,
			Title:         p.Title,
			Content:       p.Content,
			CreatedAt:     p.CreatedAt.Format(time.RFC3339),
			AllowComments: p.AllowComments,
			Author:        &model.User{},
		}
		if p.AuthorId != nil && p.Username != nil {
			post.Author.ID = strconv.Itoa(*p.AuthorId)
			post.Author.Name = *p.Username
		}
		results = append(results, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *PostPostgresRepo) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	query := `
		SELECT 
			posts.id, posts.title, posts.content, posts.created_at, posts.allow_comments,
			users.name, users.id AS author_id
		FROM posts
		JOIN users ON posts.author_id = users.id
		WHERE posts.id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)
	var p postDB
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.AllowComments, &p.Username, &p.AuthorId); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	var u model.User
	if p.AuthorId != nil && p.Username != nil {
		u.ID = strconv.Itoa(*p.AuthorId)
		u.Name = *p.Username
	}

	return &model.Post{
		ID:            p.ID,
		Title:         p.Title,
		Content:       p.Content,
		Author:        &u,
		CreatedAt:     p.CreatedAt.Format(time.RFC3339),
		AllowComments: p.AllowComments,
	}, nil
}

func (r *PostPostgresRepo) CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error) {
	query := `
		INSERT INTO posts (title, content, created_at, allow_comments, author_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, content, created_at, allow_comments, author_id
	`

	// Используем time.Now() для установки времени создания.
	var newPost postDB
	err := r.db.QueryRowContext(ctx, query, input.Title, input.Content, time.Now(), input.AllowComments, input.AuthorID).
		Scan(&newPost.ID, &newPost.Title, &newPost.Content, &newPost.CreatedAt, &newPost.AllowComments, &newPost.AuthorId)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:      newPost.ID,
		Title:   newPost.Title,
		Content: newPost.Content,
		Author: &model.User{
			ID: input.AuthorID,
		},
		CreatedAt:     newPost.CreatedAt.Format(time.RFC3339),
		AllowComments: newPost.AllowComments,
	}, nil
}
