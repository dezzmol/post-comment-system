package repository

import (
	"context"

	"post-comment-system/graph/model"
)

type PostRepository interface {
	GetAllPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error)
	GetPostByID(ctx context.Context, id int) (*model.Post, error)
	CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error)
}
