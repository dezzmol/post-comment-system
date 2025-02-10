package inmemory

import (
	"sync"

	"post-comment-system/graph/model"
)

// Лок только в таком порядке: UsersLock -> PostsLock -> CommentsLock чтобы не допустить дедлоков

type InMemoryStorage struct {
	Users        map[string]*model.User
	UsersMutex   sync.RWMutex
	UsersCounter int

	Posts       map[string]*model.Post
	PostMutex   sync.RWMutex
	PostCounter int

	Comments        map[string]*model.Comment
	CommentMutex    sync.RWMutex
	CommentsCounter int
}

func NewInMemoryStorage() *InMemoryStorage {
	storage := &InMemoryStorage{
		Users:           make(map[string]*model.User),
		UsersCounter:    3,
		Posts:           make(map[string]*model.Post),
		PostCounter:     0,
		Comments:        make(map[string]*model.Comment),
		CommentsCounter: 0,
	}

	user1 := &model.User{
		ID:   "1",
		Name: "Радмир",
	}

	user2 := &model.User{
		ID:   "2",
		Name: "Иван",
	}

	user3 := &model.User{
		ID:   "3",
		Name: "Петя",
	}

	storage.Users["1"] = user1
	storage.Users["2"] = user2
	storage.Users["3"] = user3

	return storage
}
