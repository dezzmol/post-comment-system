package inmemory

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"post-comment-system/graph/model"
	"post-comment-system/internal/storage/inmemory"
)

type InMemoryCommentRepo struct {
	s *inmemory.InMemoryStorage
}

func NewInMemoryCommentRepo(s *inmemory.InMemoryStorage) *InMemoryCommentRepo {
	return &InMemoryCommentRepo{s: s}
}

func (r *InMemoryCommentRepo) GetAllComments(ctx context.Context, limit, offset *int) ([]*model.Comment, error) {
	r.s.CommentMutex.RLock()
	defer r.s.CommentMutex.RUnlock()

	comments := make([]*model.Comment, 0)
	for _, comment := range r.s.Comments {
		comments = append(comments, comment)
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt > comments[j].CreatedAt
	})

	start := *offset
	end := start + *limit

	if end > len(comments) {
		end = len(comments)
	}

	if start > len(comments) {
		return []*model.Comment{}, nil
	}

	return comments[start:end], nil
}

func (r *InMemoryCommentRepo) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error) {
	r.s.CommentMutex.RLock()
	defer r.s.CommentMutex.RUnlock()

	comments := make([]*model.Comment, 0)
	for _, comment := range r.s.Comments {
		if comment.PostID == postID && comment.ReplyTo == nil {
			comments = append(comments, comment)
		}
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt > comments[j].CreatedAt
	})

	return comments, nil
}

func (r *InMemoryCommentRepo) GetRepliesForComment(ctx context.Context, commentID string, limit, offset *int) ([]*model.Comment, error) {
	r.s.CommentMutex.RLock()
	defer r.s.CommentMutex.RUnlock()
	baseComment, ok := r.s.Comments[commentID]
	if !ok {
		return nil, errors.New("comment not found")
	}

	comments := make([]*model.Comment, 0)
	for _, comment := range r.s.Comments {
		if comment.ReplyTo == baseComment {
			comments = append(comments, comment)
		}
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt > comments[j].CreatedAt
	})

	start := *offset
	end := start + *limit
	if end > len(comments) {
		end = len(comments)
	}

	if start > len(comments) {
		return []*model.Comment{}, nil
	}

	return comments[start:end], nil
}

func (r *InMemoryCommentRepo) CreateComment(ctx context.Context, input model.CreateComment) (*model.Comment, error) {
	// lock users -> posts -> comments
	r.s.UsersMutex.RLock()
	user, ok := r.s.Users[input.AuthorID]
	r.s.UsersMutex.RUnlock()
	if !ok {
		return nil, errors.New("user not found")
	}

	r.s.PostMutex.RLock()
	post, ok := r.s.Posts[input.PostID]
	r.s.PostMutex.RUnlock()
	if !ok {
		return nil, errors.New("post not found")
	}

	if !post.AllowComments {
		return nil, errors.New("comments are not allowed in this post")
	}

	r.s.CommentMutex.Lock()
	defer r.s.CommentMutex.Unlock()

	r.s.CommentsCounter++
	id := r.s.CommentsCounter
	comment := model.Comment{
		ID:        strconv.Itoa(id),
		PostID:    input.PostID,
		Text:      input.Text,
		Author:    user,
		CreatedAt: time.Now().Format(time.RFC3339),
		Replies:   []*model.Comment{},
	}

	if input.ReplyTo != nil {
		replyTo, ok := r.s.Comments[*input.ReplyTo]
		if !ok {
			return nil, errors.New("comment to reply not found")
		}
		comment.ReplyTo = replyTo
		replyTo.Replies = append(replyTo.Replies, &comment)
	}

	r.s.Comments[comment.ID] = &comment
	return &comment, nil
}
