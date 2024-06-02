package api

import (
	"GraphQLPostComments/api/model"
	"GraphQLPostComments/internal/storage"
	"sync"
)

type Resolver struct {
	Store     storage.Storage
	observers map[string]map[string]chan *model.Comment
	mu        sync.Mutex
}

func NewResolver(store storage.Storage) *Resolver {
	return &Resolver{
		Store:     store,
		observers: make(map[string]map[string]chan *model.Comment),
	}
}
