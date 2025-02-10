package postgres

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"post-comment-system/graph/model"
	"post-comment-system/internal/repository/postgres"
	"post-comment-system/internal/service/comment"
	"post-comment-system/internal/service/subscriber_manager"
)

func TestCreateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postRepo := postgres.NewPostgresCommentRepo(db)
	sm := subscriber_manager.NewSubscriptionManager()
	service := comment.NewCommentService(postRepo, sm)

	input := &model.CreateComment{
		Text:     "test comment",
		AuthorID: "1",
		PostID:   "1",
		ReplyTo:  nil,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT allow_comments FROM posts WHERE id = $1`)).
		WithArgs(input.PostID).
		WillReturnRows(sqlmock.NewRows([]string{"AllowComments"}).AddRow(true))

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "post_id", "text", "author_id", "reply_to", "created_at"}).
		AddRow(1, input.PostID, input.Text, input.AuthorID, input.ReplyTo, now)

	mock.ExpectQuery(`INSERT INTO comments`).
		WithArgs(input.PostID, input.Text, input.ReplyTo, sqlmock.AnyArg(), input.AuthorID).
		WillReturnRows(rows)

	createdComment, err := service.CreateComment(context.Background(), *input)
	require.NoError(t, err)

	expected := &model.Comment{
		ID:        "1",
		PostID:    input.PostID,
		Text:      input.Text,
		Author:    &model.User{ID: input.AuthorID},
		ReplyTo:   nil,
		CreatedAt: now.Format(time.RFC3339),
		Replies:   nil,
	}

	require.Equal(t, expected.ID, createdComment.ID)
	require.Equal(t, expected.PostID, createdComment.PostID)
	require.Equal(t, expected.Text, createdComment.Text)
	require.Equal(t, expected.Author, createdComment.Author)
	require.Equal(t, expected.ReplyTo, createdComment.ReplyTo)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	postRepo := postgres.NewPostgresCommentRepo(db)
	sm := subscriber_manager.NewSubscriptionManager()
	service := comment.NewCommentService(postRepo, sm)

	limit := 10
	offset := 0

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "post_id", "text", "reply_to", "created_at", "id", "name"}).
		AddRow(1, 1, "Comment 1", nil, now, 1, "Radmir").
		AddRow(2, 1, "Comment 2", nil, now, 2, "Ivan")

	mock.ExpectQuery("SELECT c.id, c.post_id, c.text, c.reply_to, c.created_at, u.id AS user_id, u.name AS username").
		WithArgs(limit, offset).
		WillReturnRows(rows)

	comments, err := service.GetComments(context.Background(), &limit, &offset)

	require.NoError(t, err)

	expected := []*model.Comment{
		{
			ID:     "1",
			PostID: "1",
			Text:   "Comment 1",
			Author: &model.User{
				ID:   "1",
				Name: "Radmir",
			},
			ReplyTo: nil,
			Replies: []*model.Comment{},
		},
		{
			ID:     "2",
			PostID: "1",
			Text:   "Comment 2",
			Author: &model.User{
				ID:   "2",
				Name: "Ivan",
			},
			ReplyTo: nil,
			Replies: []*model.Comment{},
		},
	}

	opts := cmp.Options{
		cmpopts.IgnoreFields(model.Comment{}, "CreatedAt"),
	}

	require.True(t, cmp.Equal(expected, comments, opts), cmp.Diff(expected, comments, opts))
}
