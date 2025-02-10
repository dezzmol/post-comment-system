package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"post-comment-system/graph/model"
	inmemory2 "post-comment-system/internal/repository/inmemory"
	"post-comment-system/internal/service/post"
	"post-comment-system/internal/storage/inmemory"
)

func TestGetAllPosts(t *testing.T) {
	t.Parallel()

	storage := inmemory.NewInMemoryStorage()
	repo := inmemory2.NewInMemoryPostRepo(storage)
	commentRepo := inmemory2.NewInMemoryCommentRepo(storage)
	service := post.NewPostService(repo, commentRepo)

	storage.Posts["1"] = &model.Post{
		ID:            "1",
		Title:         "Post 1",
		Content:       "Content 1",
		Author:        &model.User{ID: "1"},
		AllowComments: true,
		CreatedAt:     "2024-2-7T15:00:00Z",
		Comments:      []*model.Comment{},
	}
	storage.Posts["2"] = &model.Post{
		ID:            "2",
		Title:         "Post 2",
		Content:       "Content 2",
		Author:        &model.User{ID: "2"},
		AllowComments: true,
		CreatedAt:     "2024-2-7T16:00:00Z",
		Comments:      []*model.Comment{},
	}

	limit := 2
	offset := 0

	posts, err := service.GetPosts(context.Background(), &limit, &offset)

	expected := []*model.Post{
		{
			ID:            "2",
			Title:         "Post 2",
			Content:       "Content 2",
			Author:        &model.User{ID: "2"},
			AllowComments: true,
			CreatedAt:     "2024-2-7T16:00:00Z",
			Comments:      []*model.Comment{},
		},
		{
			ID:            "1",
			Title:         "Post 1",
			Content:       "Content 1",
			Author:        &model.User{ID: "1"},
			AllowComments: true,
			CreatedAt:     "2024-2-7T15:00:00Z",
			Comments:      []*model.Comment{},
		},
	}

	require.NoError(t, err)
	require.Equal(t, expected, posts)

}

func TestGetPostByID(t *testing.T) {
	t.Parallel()

	storage := inmemory.NewInMemoryStorage()
	repo := inmemory2.NewInMemoryPostRepo(storage)
	commentRepo := inmemory2.NewInMemoryCommentRepo(storage)
	service := post.NewPostService(repo, commentRepo)
	storage.Posts["1"] = &model.Post{
		ID:            "1",
		Title:         "Post 1",
		Content:       "Content 1",
		Author:        &model.User{ID: "1"},
		AllowComments: true,
		CreatedAt:     "2024-2-7T15:00:00Z",
		Comments:      []*model.Comment{},
	}

	createdPost, err := service.GetPostByID(context.Background(), 1)
	expected := &model.Post{
		ID:            "1",
		Title:         "Post 1",
		Content:       "Content 1",
		Author:        &model.User{ID: "1"},
		AllowComments: true,
		CreatedAt:     "2024-2-7T15:00:00Z",
		Comments:      []*model.Comment{},
	}
	require.NoError(t, err)
	require.Equal(t, expected, createdPost)
}

func TestGetPostByIDNotFoundError(t *testing.T) {
	t.Parallel()

	storage := inmemory.NewInMemoryStorage()
	repo := inmemory2.NewInMemoryPostRepo(storage)
	commentRepo := inmemory2.NewInMemoryCommentRepo(storage)
	service := post.NewPostService(repo, commentRepo)

	createdPost, err := service.GetPostByID(context.Background(), 1)
	require.Nil(t, createdPost)
	require.Error(t, err)
}

func TestCreatePost(t *testing.T) {
	t.Parallel()

	storage := inmemory.NewInMemoryStorage()
	repo := inmemory2.NewInMemoryPostRepo(storage)
	commentRepo := inmemory2.NewInMemoryCommentRepo(storage)
	service := post.NewPostService(repo, commentRepo)

	newPost := &model.CreatePost{
		Title:         "1",
		Content:       "Content 1",
		AuthorID:      "1",
		AllowComments: true,
	}

	expected, err := service.CreatePost(context.Background(), *newPost)
	require.NoError(t, err)
	require.Equal(t, expected.Title, newPost.Title)
	require.Equal(t, expected.Content, newPost.Content)
	require.Equal(t, expected.Author.ID, newPost.AuthorID)
}

func TestCreatePostUserNotFoundError(t *testing.T) {
	t.Parallel()

	storage := inmemory.NewInMemoryStorage()
	repo := inmemory2.NewInMemoryPostRepo(storage)
	commentRepo := inmemory2.NewInMemoryCommentRepo(storage)
	service := post.NewPostService(repo, commentRepo)
	newPost := &model.CreatePost{
		Title:         "title 1",
		Content:       "Content 1",
		AuthorID:      "4",
		AllowComments: true,
	}

	expected, err := service.CreatePost(context.Background(), *newPost)
	require.Error(t, err)
	require.Nil(t, expected)
}
