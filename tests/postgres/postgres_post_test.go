package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"post-comment-system/graph/model"
	"post-comment-system/internal/repository/postgres"
	"post-comment-system/internal/service/post"
)

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postRepo := postgres.NewPostPostgresRepository(db)
	commentsRepo := postgres.NewPostgresCommentRepo(db)
	service := post.NewPostService(postRepo, commentsRepo)

	input := model.CreatePost{
		Title:         "Post 1",
		Content:       "Hello World",
		AuthorID:      "1",
		AllowComments: true,
	}

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "allow_comments", "author_id"}).
		AddRow(1, input.Title, input.Content, now, input.AllowComments, input.AuthorID)

	mock.ExpectQuery(`INSERT INTO posts`).
		WithArgs(input.Title, input.Content, sqlmock.AnyArg(), input.AllowComments, input.AuthorID).
		WillReturnRows(rows)

	postResult, err := service.CreatePost(context.Background(), input)
	require.NoError(t, err)

	require.Equal(t, input.Title, postResult.Title)
	require.Equal(t, input.Content, postResult.Content)
	require.Equal(t, "1", postResult.ID)
	require.Equal(t, input.AllowComments, postResult.AllowComments)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCreatePostError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postRepo := postgres.NewPostPostgresRepository(db)
	commentsRepo := postgres.NewPostgresCommentRepo(db)
	service := post.NewPostService(postRepo, commentsRepo)

	input := model.CreatePost{
		Title:         "Post 1",
		Content:       "Hello World",
		AuthorID:      "1",
		AllowComments: true,
	}

	mock.ExpectQuery(`INSERT INTO posts`).
		WithArgs(input.Title, input.Content, sqlmock.AnyArg(), input.AllowComments, input.AuthorID).
		WillReturnError(sql.ErrConnDone)

	postResult, err := service.CreatePost(context.Background(), input)
	require.Error(t, err)
	assert.Nil(t, postResult)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postRepo := postgres.NewPostPostgresRepository(db)
	commentsRepo := postgres.NewPostgresCommentRepo(db)
	service := post.NewPostService(postRepo, commentsRepo)

	limit := 10
	offset := 0

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "created_at", "allow_comments", "name", "author_id"}).
		AddRow(1, "Title 1", "Content 1", now, true, "Radmir", "1").
		AddRow(2, "Title 2", "Content 2", now, true, "Radmir", "1")

	mock.ExpectQuery(`SELECT (.+) FROM posts`).
		WithArgs(limit, offset).WillReturnRows(rows)

	posts, err := service.GetPosts(context.Background(), &limit, &offset)

	require.NoError(t, err)

	expectedPosts := []*model.Post{
		{
			ID:            "1",
			Title:         "Title 1",
			Content:       "Content 1",
			AllowComments: true,
			CreatedAt:     now.Format(time.RFC3339),
			Author: &model.User{
				ID:   "1",
				Name: "Radmir",
			},
			Comments: nil,
		},
		{
			ID:            "2",
			Title:         "Title 2",
			Content:       "Content 2",
			AllowComments: true,
			CreatedAt:     now.Format(time.RFC3339),
			Author: &model.User{
				ID:   "1",
				Name: "Radmir",
			},
			Comments: nil,
		},
	}

	assert.ElementsMatch(t, expectedPosts, posts)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
