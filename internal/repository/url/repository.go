package url

import (
	"fmt"
	"sync"

	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
)

var _ repository.URLRepository = (*InMemoryRepository)(nil)

var ErrIDNotFound = fmt.Errorf("URL ID was not found in repository")

func New() *InMemoryRepository {
	return &InMemoryRepository{
		m: make(map[string]*model.URL),
	}
}

type InMemoryRepository struct {
	m map[string]*model.URL
	sync.RWMutex
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
