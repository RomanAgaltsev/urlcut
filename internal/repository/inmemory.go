package repository

import (
	"fmt"
	"sync"

	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/model"
)

var _ interfaces.URLStoreGetter = (*InMemoryRepository)(nil)

var (
	ErrIDNotFound         = fmt.Errorf("URL ID was not found in repository")
	ErrStorageUnavailable = fmt.Errorf("storage unavailable")
)

type InMemoryRepository struct {
	m map[string]*model.URL
	sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		m: make(map[string]*model.URL),
	}
}

func (r *InMemoryRepository) Store(url *model.URL) error {
	r.Lock()
	defer r.Unlock()

	r.m[url.ID] = url

	return nil
}

func (r *InMemoryRepository) Get(id string) (*model.URL, error) {
	r.Lock()
	defer r.Unlock()

	if url, ok := r.m[id]; ok {
		return url, nil
	} else {
		return url, ErrIDNotFound
	}
}

func (r *InMemoryRepository) SetState(state map[string]*model.URL) error {
	r.m = state
	return nil
}

func (r *InMemoryRepository) GetState() map[string]*model.URL {
	return r.m
}

func (r *InMemoryRepository) Check() error {
	if r.m == nil {
		return ErrStorageUnavailable
	}
	return nil
}
