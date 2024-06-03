package dataloader

import (
	"context"
	"sync"

	"GraphQLPostComments/api/model"
	"GraphQLPostComments/internal/storage"
)

type PostLoader struct {
	store storage.Storage
	mu    sync.Mutex
	cache map[string]*model.Post
}

func NewPostLoader(store storage.Storage) *PostLoader {
	return &PostLoader{
		store: store,
		cache: make(map[string]*model.Post),
	}
}

func (l *PostLoader) Load(ctx context.Context, key string) (*model.Post, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if post, ok := l.cache[key]; ok {
		return post, nil
	}

	post, err := l.store.GetPostByID(key)
	if err != nil {
		return nil, err
	}

	l.cache[key] = post
	return post, nil
}

func (l *PostLoader) LoadMany(ctx context.Context, keys []string) ([]*model.Post, error) {
	posts := make([]*model.Post, len(keys))

	for i, key := range keys {
		post, err := l.Load(ctx, key)
		if err != nil {
			return nil, err
		}

		posts[i] = post
	}

	return posts, nil
}
