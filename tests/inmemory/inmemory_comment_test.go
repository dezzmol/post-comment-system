package inmemory

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"post-comment-system/graph/model"
	inmemory2 "post-comment-system/internal/repository/inmemory"
	"post-comment-system/internal/service/comment"
	"post-comment-system/internal/service/subscriber_manager"
	"post-comment-system/internal/storage/inmemory"
)

func TestGetComments(t *testing.T) {
	t.Parallel()
	storage := inmemory.NewInMemoryStorage()

	storage.Comments["1"] = &model.Comment{
		ID:        "1",
		PostID:    "post1",
		Author:    &model.User{ID: "1"},
		ReplyTo:   nil,
		Text:      "Comment 1",
		CreatedAt: "2024-2-7T14:00:00Z",
	}
	storage.Comments["2"] = &model.Comment{
		ID:        "2",
		PostID:    "post1",
		Author:    &model.User{ID: "2"},
		ReplyTo:   nil,
		Text:      "Comment 2",
		CreatedAt: "2024-2-7T15:00:00Z",
	}

	repo := inmemory2.NewInMemoryCommentRepo(storage)
	sm := subscriber_manager.NewSubscriptionManager()
	service := comment.NewCommentService(repo, sm)

	limit := 2
	offset := 0

	comments, err := service.GetComments(context.Background(), &limit, &offset)
	require.NoError(t, err)

	expected := []*model.Comment{
		{
			ID:        "2",
			PostID:    "post1",
			Author:    &model.User{ID: "2"},
			ReplyTo:   nil,
			Text:      "Comment 2",
			CreatedAt: "2024-2-7T15:00:00Z",
		},
		{
			ID:        "1",
			PostID:    "post1",
			Author:    &model.User{ID: "1"},
			ReplyTo:   nil,
			Text:      "Comment 1",
			CreatedAt: "2024-2-7T14:00:00Z",
		},
	}
	require.Equal(t, expected, comments)
}

func TestGetReplies(t *testing.T) {
	t.Parallel()
	storage := inmemory.NewInMemoryStorage()

	storage.Comments["1"] = &model.Comment{
		ID:        "1",
		PostID:    "post1",
		Author:    &model.User{ID: "1"},
		ReplyTo:   nil,
		Text:      "Comment 1",
		CreatedAt: "2024-2-7T15:00:00Z",
	}
	storage.Comments["2"] = &model.Comment{
		ID:        "2",
		PostID:    "post1",
		Author:    &model.User{ID: "2"},
		ReplyTo:   storage.Comments["1"],
		Text:      "Comment 2",
		CreatedAt: "2024-2-7T16:00:00Z",
	}

	repo := inmemory2.NewInMemoryCommentRepo(storage)
	sm := subscriber_manager.NewSubscriptionManager()
	service := comment.NewCommentService(repo, sm)

	limit := 2
	offset := 0

	comments, err := service.GetRepliesForComment(context.Background(), "1", &limit, &offset)
	require.NoError(t, err)
	expected := []*model.Comment{
		{
			ID:        "2",
			PostID:    "post1",
			Author:    &model.User{ID: "2"},
			ReplyTo:   storage.Comments["1"],
			Text:      "Comment 2",
			CreatedAt: "2024-2-7T16:00:00Z",
		},
	}

	require.Equal(t, expected, comments)
}

func TestGetRepliesNotFoundError(t *testing.T) {
	t.Parallel()
	storage := inmemory.NewInMemoryStorage()
	repo := inmemory2.NewInMemoryCommentRepo(storage)
	sm := subscriber_manager.NewSubscriptionManager()
	service := comment.NewCommentService(repo, sm)

	limit := 2
	offset := 0

	comments, err := service.GetRepliesForComment(context.Background(), "1", &limit, &offset)
	require.Error(t, err)
	assert.Nil(t, comments)
}

func TestCreateComment(t *testing.T) {
	t.Parallel()
	storage := inmemory.NewInMemoryStorage()
	repo := inmemory2.NewInMemoryCommentRepo(storage)
	sm := subscriber_manager.NewSubscriptionManager()
	service := comment.NewCommentService(repo, sm)

	storage.Posts["1"] = &model.Post{
		ID:            "1",
		Title:         "Post 1",
		Content:       "Content",
		Author:        &model.User{ID: "1"},
		CreatedAt:     "2024-2-7T16:00:00Z",
		AllowComments: true,
		Comments:      []*model.Comment{},
	}

	input := &model.CreateComment{
		Text:     "Comment 1",
		PostID:   "1",
		AuthorID: "1",
		ReplyTo:  nil,
	}

	comment, err := service.CreateComment(context.Background(), *input)

	expected := &model.Comment{
		ID:        "1",
		PostID:    "1",
		Text:      "Comment 1",
		ReplyTo:   nil,
		CreatedAt: comment.CreatedAt,
	}

	require.NoError(t, err)
	require.Equal(t, expected.ID, comment.ID)
	require.Equal(t, expected.PostID, comment.PostID)
	require.Equal(t, expected.Text, comment.Text)
	assert.True(t, len(expected.CreatedAt) > 0)
}

func TestCreateCommentPostNotFoundError(t *testing.T) {
	t.Parallel()

	storage := inmemory.NewInMemoryStorage()
	repo := inmemory2.NewInMemoryCommentRepo(storage)
	sm := subscriber_manager.NewSubscriptionManager()
	service := comment.NewCommentService(repo, sm)

	input := &model.CreateComment{
		Text:   "Comment 1",
		PostID: "1",
	}
	expected, err := service.CreateComment(context.Background(), *input)
	require.Error(t, err)
	assert.Nil(t, expected)
}

func TestCreateCommentNotAllowedCommentingError(t *testing.T) {
	t.Parallel()
	storage := inmemory.NewInMemoryStorage()
	repo := inmemory2.NewInMemoryCommentRepo(storage)
	sm := subscriber_manager.NewSubscriptionManager()
	service := comment.NewCommentService(repo, sm)

	storage.Posts["1"] = &model.Post{
		ID:            "1",
		Title:         "Post 1",
		Content:       "Content",
		Author:        &model.User{ID: "1"},
		CreatedAt:     "2024-2-7T15:00:00Z",
		AllowComments: false,
	}

	input := &model.CreateComment{
		Text:     "Comment 1",
		PostID:   "1",
		AuthorID: "1",
		ReplyTo:  nil,
	}

	expected, err := service.CreateComment(context.Background(), *input)
	require.Error(t, err)
	assert.Nil(t, expected)
}

func TestCreateCommentTooLongMessageError(t *testing.T) {
	t.Parallel()
	storage := inmemory.NewInMemoryStorage()
	repo := inmemory2.NewInMemoryCommentRepo(storage)
	sm := subscriber_manager.NewSubscriptionManager()
	service := comment.NewCommentService(repo, sm)

	input := &model.CreateComment{
		Text:   strings.Repeat("a", 2001),
		PostID: "1",
	}

	expected, err := service.CreateComment(context.Background(), *input)
	require.Error(t, err)
	assert.Nil(t, expected)
}

func TestCreateCommentAuthorNotFoundError(t *testing.T) {
	t.Parallel()
	storage := inmemory.NewInMemoryStorage()
	repo := inmemory2.NewInMemoryCommentRepo(storage)
	sm := subscriber_manager.NewSubscriptionManager()
	service := comment.NewCommentService(repo, sm)

	input := &model.CreateComment{
		Text:     "Comment 1",
		PostID:   "1",
		AuthorID: "4",
		ReplyTo:  nil,
	}

	expected, err := service.CreateComment(context.Background(), *input)
	require.Error(t, err)
	assert.Nil(t, expected)
}
