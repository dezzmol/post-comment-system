package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"post-comment-system/graph/model"
)

type PostgresCommentRepo struct {
	db *sql.DB
}

func NewPostgresCommentRepo(db *sql.DB) *PostgresCommentRepo {
	return &PostgresCommentRepo{db: db}
}

type mappingCommentDB struct {
	ID        int     `db:"id"`
	PostID    int     `db:"post_id"`
	Text      string  `db:"text"`
	ReplyTo   *int    `db:"reply_to"`
	CreatedAt string  `db:"created_at"`
	UserID    *int    `db:"id"`
	Username  *string `db:"name"`
}

type commentDB struct {
	ID        int    `db:"id"`
	PostID    int    `db:"post_id"`
	Text      string `db:"text"`
	AuthorID  *int   `db:"author_id"`
	ReplyTo   *int   `db:"reply_to"`
	CreatedAt string `db:"created_at"`
}

func (r *PostgresCommentRepo) GetAllComments(ctx context.Context, limit, offset *int) ([]*model.Comment, error) {
	query := `
		SELECT 
			c.id, c.post_id, c.text, c.reply_to, c.created_at,
			u.id AS user_id, u.name AS username
		FROM comments c
		JOIN users u ON c.author_id = u.id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, *limit, *offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentsDB []*mappingCommentDB
	for rows.Next() {
		var m mappingCommentDB
		if err := rows.Scan(&m.ID, &m.PostID, &m.Text, &m.ReplyTo, &m.CreatedAt, &m.UserID, &m.Username); err != nil {
			return nil, err
		}

		commentsDB = append(commentsDB, &m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return buildCommentsTree(commentsDB)
}

func (r *PostgresCommentRepo) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error) {
	query := `
		SELECT 
			comments.id, comments.post_id, comments.text, comments.reply_to, comments.created_at,
			users.id AS user_id, users.name AS username
		FROM comments
		JOIN users ON comments.author_id = users.id
		WHERE comments.post_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentsDB []*mappingCommentDB
	for rows.Next() {
		var m mappingCommentDB
		if err := rows.Scan(&m.ID, &m.PostID, &m.Text, &m.ReplyTo, &m.CreatedAt, &m.UserID, &m.Username); err != nil {
			return nil, err
		}

		commentsDB = append(commentsDB, &m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return buildCommentsTree(commentsDB)
}

// Функция для построения древоидной структуры комментариев
func buildCommentsTree(dbComments []*mappingCommentDB) ([]*model.Comment, error) {
	allComments := make(map[int]*model.Comment)

	// Заполняем мапу, чтобы получать быстрый доступ к комментариям
	for _, comment := range dbComments {
		var u model.User
		if comment.UserID != nil && comment.Username != nil {
			u.ID = strconv.Itoa(comment.ID)
			u.Name = *comment.Username
		}
		allComments[comment.ID] = &model.Comment{
			ID:        strconv.Itoa(comment.ID),
			PostID:    strconv.Itoa(comment.PostID),
			Text:      comment.Text,
			CreatedAt: comment.CreatedAt,
			Author:    &u,
			Replies:   []*model.Comment{},
		}
	}

	// Заполняем комментарии к посту, комментарии ответы подвешиваем к нужным комментариям
	var roots []*model.Comment
	for _, dbComment := range dbComments {
		comment := allComments[dbComment.ID]
		if dbComment.ReplyTo != nil {
			parent, ok := allComments[*dbComment.ReplyTo]
			if ok {
				comment.ReplyTo = parent
				parent.Replies = append(parent.Replies, comment)
			} else {
				// Если в мапе не оказалось сообщения, но есть replyid - значит, что комментарий для ответа не был прогружен из БД
				comment.ReplyTo = &model.Comment{
					ID: strconv.Itoa(*dbComment.ReplyTo),
				}
			}
		} else {
			roots = append(roots, comment)
		}
	}

	return roots, nil
}

func (r *PostgresCommentRepo) GetRepliesForComment(ctx context.Context, commentID string, limit, offset *int) ([]*model.Comment, error) {
	query := `
		SELECT 
			c.id, c.post_id, c.text, c.reply_to, c.created_at,
			u.id AS user_id, u.name AS username
		FROM comments c
		JOIN users u ON c.author_id = u.id
		WHERE c.reply_to = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, commentID, *limit, *offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var m mappingCommentDB
		if err := rows.Scan(&m.ID, &m.PostID, &m.Text, &m.ReplyTo, &m.CreatedAt, &m.UserID, &m.Username); err != nil {
			return nil, err
		}

		var u model.User
		if m.UserID != nil && m.Username != nil {
			u.ID = strconv.Itoa(*m.UserID)
			u.Name = *m.Username
		}

		comment := &model.Comment{
			ID:        strconv.Itoa(m.ID),
			PostID:    strconv.Itoa(m.PostID),
			Text:      m.Text,
			Author:    &u,
			CreatedAt: m.CreatedAt,
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *PostgresCommentRepo) CreateComment(ctx context.Context, input model.CreateComment) (*model.Comment, error) {
	query := `SELECT allow_comments FROM posts WHERE id = $1`
	var allowComments bool
	err := r.db.QueryRowContext(ctx, query, input.PostID).Scan(&allowComments)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	if !allowComments {
		return nil, errors.New("commenting is not allowed")
	}

	insertQuery := `
		INSERT INTO comments (post_id, text, reply_to, created_at, author_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, post_id, text, author_id, reply_to, created_at
	`
	var c commentDB
	err = r.db.QueryRowContext(ctx, insertQuery, input.PostID, input.Text, input.ReplyTo, time.Now(), input.AuthorID).
		Scan(&c.ID, &c.PostID, &c.Text, &c.AuthorID, &c.ReplyTo, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	var u model.User
	if c.AuthorID != nil {
		u.ID = strconv.Itoa(*c.AuthorID)
	}

	comment := &model.Comment{
		ID:        strconv.Itoa(c.ID),
		PostID:    strconv.Itoa(c.PostID),
		Text:      c.Text,
		Author:    &u,
		CreatedAt: c.CreatedAt,
	}

	if c.ReplyTo != nil {
		comment.ReplyTo = &model.Comment{
			ID: strconv.Itoa(*c.ReplyTo),
		}
	}

	return comment, nil
}
