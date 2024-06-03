package dataloader

import (
	"context"
	"sync"

	"GraphQLPostComments/api/model"
	"GraphQLPostComments/internal/storage"
)

type CommentLoader struct {
	postLoader *PostLoader
	store      storage.Storage
	mu         sync.Mutex
	cache      map[string]*model.Comment
}

func NewCommentLoader(postLoader *PostLoader, store storage.Storage) *CommentLoader {
	return &CommentLoader{
		postLoader: postLoader,
		store:      store,
		cache:      make(map[string]*model.Comment),
	}
}

func (l *CommentLoader) Load(ctx context.Context, key string) (*model.Comment, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if comment, ok := l.cache[key]; ok {
		return comment, nil
	}

	comment, err := l.store.GetCommentByID(key)
	if err != nil {
		return nil, err
	}

	post, err := l.postLoader.Load(ctx, comment.Post.ID)
	if err != nil {
		return nil, err
	}

	comment.Post = post

	l.cache[key] = comment
	return comment, nil
}

func (l *CommentLoader) LoadMany(ctx context.Context, keys []string) ([]*model.Comment, error) {
	comments := make([]*model.Comment, len(keys))

	wg := sync.WaitGroup{}
	ch := make(chan error, len(keys))
	wg.Add(len(keys))

	for i, key := range keys {
		go func(i int, key string) {
			defer wg.Done()

			comment, err := l.Load(ctx, key)
			if err != nil {
				ch <- err
				return
			}

			comments[i] = comment
		}(i, key)
	}
	wg.Wait()
	close(ch)

	err := <-ch
	if err != nil {
		return nil, err
	}
	return comments, nil
}
