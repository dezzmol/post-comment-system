package graph

//go:generate go run github.com/99designs/gqlgen generate
import (
	"post-comment-system/internal/service/comment"
	"post-comment-system/internal/service/post"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostService    *post.Service
	CommentService *comment.Service
}

func NewResolver(postService *post.Service, commentService *comment.Service) *Resolver {
	return &Resolver{
		PostService:    postService,
		CommentService: commentService,
	}
}
