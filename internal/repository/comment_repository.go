package repository

import (
	"context"

	"post-comment-system/graph/model"
)

type CommentRepository interface {
	GetAllComments(ctx context.Context, limit, offset *int) ([]*model.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error)
	GetRepliesForComment(ctx context.Context, commentID string, limit, offset *int) ([]*model.Comment, error)
	CreateComment(ctx context.Context, input model.CreateComment) (*model.Comment, error)
}
