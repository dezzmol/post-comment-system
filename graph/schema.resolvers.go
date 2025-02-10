package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"context"

	"post-comment-system/graph/model"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error) {
	return r.PostService.CreatePost(ctx, input)
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateComment) (*model.Comment, error) {
	return r.CommentService.CreateComment(ctx, input)
}

// GetPosts is the resolver for the getPosts field.
func (r *queryResolver) GetPosts(ctx context.Context, limit *int, offset *int) ([]*model.Post, error) {
	return r.PostService.GetPosts(ctx, limit, offset)
}

// GetPostByID is the resolver for the getPostByID field.
func (r *queryResolver) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	return r.PostService.GetPostByID(ctx, id)
}

// GetComments is the resolver for the getComments field.
func (r *queryResolver) GetComments(ctx context.Context, limit *int, offset *int) ([]*model.Comment, error) {
	return r.CommentService.GetComments(ctx, limit, offset)
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	commentChan := make(chan *model.Comment, 1)

	// Регистрируем подписчика для указанного postID
	r.CommentService.SubscribeToPost(ctx, postID, commentChan)

	// При завершении контекста отписываемся
	go func() {
		<-ctx.Done()
		r.CommentService.UnsubscribeFromPost(ctx, postID, commentChan)
	}()

	return commentChan, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
