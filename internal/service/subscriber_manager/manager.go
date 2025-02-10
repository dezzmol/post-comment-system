package subscriber_manager

import (
	"sync"

	"post-comment-system/graph/model"
)

type SubscriptionManager struct {
	mu          sync.Mutex
	subscribers map[string][]chan *model.Comment // Ключ - идентификатор поста
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscribers: make(map[string][]chan *model.Comment),
	}
}

func (sm *SubscriptionManager) Subscribe(postID string, ch chan *model.Comment) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.subscribers[postID] = append(sm.subscribers[postID], ch)
}

// Метод для публикации нового комментария
func (sm *SubscriptionManager) PublishComment(postID string, comment *model.Comment) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if subs, ok := sm.subscribers[postID]; ok {
		for _, ch := range subs {
			// Используем неблокирующую отправку
			select {
			case ch <- comment:
			default:
			}
		}
	}
}

func (sm *SubscriptionManager) Unsubscribe(postID string, ch chan *model.Comment) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	subs := sm.subscribers[postID]
	// Ищем и удаляем канал подписчика из среза
	for i, subscriber := range subs {
		if subscriber == ch {
			subs = append(subs[:i], subs[i+1:]...)
			break
		}
	}

	// Если подписчиков осталось, обновляем срез, иначе удаляем ключ из карты
	if len(subs) > 0 {
		sm.subscribers[postID] = subs
	} else {
		delete(sm.subscribers, postID)
	}
}
