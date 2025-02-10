package comment

import (
	"context"
	"errors"

	"post-comment-system/graph/model"
	"post-comment-system/internal/repository"
	"post-comment-system/internal/service/subscriber_manager"
)

type CommentService interface {
	GetComments(ctx context.Context, limit, offset *int) ([]*model.Comment, error)
	GetRepliesForComment(ctx context.Context, commentID string, limit, offset *int) ([]*model.Comment, error)
	CreateComment(ctx context.Context, input model.CreateComment) (*model.Comment, error)
	SubscribeToPost(ctx context.Context, postID string)
	UnsubscribeFromPost(ctx context.Context, postID string, ch chan *model.Comment)
}

type Service struct {
	repo                repository.CommentRepository
	subscriptionManager *subscriber_manager.SubscriptionManager
}

func NewCommentService(repo repository.CommentRepository, sm *subscriber_manager.SubscriptionManager) *Service {
	return &Service{
		repo:                repo,
		subscriptionManager: sm,
	}
}

func (s *Service) GetComments(ctx context.Context, limit, offset *int) ([]*model.Comment, error) {
	return s.repo.GetAllComments(ctx, limit, offset)
}

func (s *Service) GetRepliesForComment(ctx context.Context, commentID string, limit, offset *int) ([]*model.Comment, error) {
	return s.repo.GetRepliesForComment(ctx, commentID, limit, offset)
}

func (s *Service) CreateComment(ctx context.Context, input model.CreateComment) (*model.Comment, error) {
	if len([]rune(input.Text)) > 2000 {
		return nil, errors.New("text too long")
	}

	comment, err := s.repo.CreateComment(ctx, input)
	if err != nil {
		return nil, err
	}
	s.subscriptionManager.PublishComment(comment.PostID, comment)
	return comment, nil
}

func (s *Service) SubscribeToPost(ctx context.Context, postID string, ch chan *model.Comment) {
	s.subscriptionManager.Subscribe(postID, ch)
}

func (s *Service) UnsubscribeFromPost(ctx context.Context, postID string, ch chan *model.Comment) {
	s.subscriptionManager.Unsubscribe(postID, ch)
}
