package inmemory

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"post-comment-system/graph/model"
	"post-comment-system/internal/repository"
	"post-comment-system/internal/storage/inmemory"
)

type InMemoryPostRepo struct {
	storage *inmemory.InMemoryStorage
}

func NewInMemoryPostRepo(storage *inmemory.InMemoryStorage) repository.PostRepository {
	return &InMemoryPostRepo{
		storage: storage,
	}
}

func (r *InMemoryPostRepo) GetAllPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	r.storage.PostMutex.RLock()
	defer r.storage.PostMutex.RUnlock()

	var posts []*model.Post
	for _, post := range r.storage.Posts {
		posts = append(posts, post)
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt > posts[j].CreatedAt
	})

	start := *offset
	end := start + *limit

	if end > len(posts) {
		end = len(posts)
	}

	if start > len(posts) {
		return []*model.Post{}, nil
	}

	return posts[start:end], nil
}

func (r *InMemoryPostRepo) GetPostByID(ctx context.Context, id int) (*model.Post, error) {
	r.storage.PostMutex.RLock()
	defer r.storage.PostMutex.RUnlock()

	post, exist := r.storage.Posts[strconv.Itoa(id)]
	if !exist {
		return nil, errors.New("post not found")
	}

	r.storage.CommentMutex.RLock()
	// only post comments without replies
	comments := make([]*model.Comment, 0)
	for _, comment := range r.storage.Comments {
		if comment.PostID == post.ID && comment.ReplyTo == nil {
			comments = append(comments, comment)
		}
	}

	r.storage.CommentMutex.RUnlock()
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt > comments[j].CreatedAt
	})

	resPost := &model.Post{
		ID:            post.ID,
		Title:         post.Title,
		Content:       post.Content,
		CreatedAt:     post.CreatedAt,
		AllowComments: post.AllowComments,
		Author:        post.Author,
		Comments:      comments,
	}

	return resPost, nil
}

func (r *InMemoryPostRepo) CreatePost(ctx context.Context, input model.CreatePost) (*model.Post, error) {
	r.storage.UsersMutex.RLock()
	user, ok := r.storage.Users[input.AuthorID]
	r.storage.UsersMutex.RUnlock()
	if !ok {
		return nil, errors.New("user doesn't exists")
	}

	r.storage.PostMutex.Lock()
	defer r.storage.PostMutex.Unlock()

	r.storage.PostCounter++
	postID := strconv.Itoa(r.storage.PostCounter)
	newPost := model.Post{
		ID:            postID,
		Title:         input.Title,
		Content:       input.Content,
		AllowComments: input.AllowComments,
		Author:        user,
		CreatedAt:     time.Now().Format(time.RFC3339),
		Comments:      []*model.Comment{},
	}

	r.storage.Posts[postID] = &newPost

	return &newPost, nil
}
