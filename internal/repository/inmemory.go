package repository

import (
	"fmt"
	"sync"

	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/model"
)

var _ interfaces.Repository = (*InMemoryRepository)(nil)

var (
	ErrIDNotFound         = fmt.Errorf("URL ID was not found in repository")
	ErrStorageUnavailable = fmt.Errorf("storage unavailable")
)

type InMemoryRepository struct {
	m map[string]*model.URL
	f string
	sync.RWMutex
}

func NewInMemoryRepository(fileStoragePath string) *InMemoryRepository {
	m, _ := readFromFile(fileStoragePath)

	return &InMemoryRepository{
		m: m,
		f: fileStoragePath,
	}
}

func (r *InMemoryRepository) Store(urls []*model.URL) (*model.URL, error) {
	r.Lock()
	defer r.Unlock()

	for _, url := range urls {
		r.m[url.ID] = url
	}

	return nil, nil
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

func (r *InMemoryRepository) Close() error {
	return writeToFile(r.f, r.m)
}

func (r *InMemoryRepository) Check() error {
	if r.m == nil {
		return ErrStorageUnavailable
	}
	return nil
}
