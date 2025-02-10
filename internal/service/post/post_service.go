package post

import (
	"context"
	"strconv"

	"post-comment-system/graph/model"
	"post-comment-system/internal/repository"
	"post-comment-system/internal/service/subscriber_manager"
)

type PostService interface {
	GetPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error)
	GetPostByID(ctx context.Context, id int) (*model.Post, error)
	CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error)
}

type Service struct {
	postRepo            repository.PostRepository
	commentRepo         repository.CommentRepository
	subscriptionManager *subscriber_manager.SubscriptionManager
}

func NewPostService(postRepo repository.PostRepository, commentsRepo repository.CommentRepository) *Service {
	return &Service{
		postRepo:    postRepo,
		commentRepo: commentsRepo,
	}
}

func (s *Service) GetPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	return s.postRepo.GetAllPosts(ctx, limit, offset)
}

func (s *Service) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	post, err := s.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByPostID(ctx, strconv.Itoa(id))

	if err != nil {
		return nil, err
	}
	post.Comments = comments
	return post, nil
}

func (s *Service) CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error) {
	return s.postRepo.CreatePost(ctx, input)
}
