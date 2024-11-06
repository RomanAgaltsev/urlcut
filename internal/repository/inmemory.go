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
	sync.RWMutex
}

func NewInMemoryRepository(fileStoragePath string) *InMemoryRepository {
	//	saver := services.NewStateSaver(a.config.FileStoragePath)
	//	state, err := saver.RestoreState()
	//	if err == nil {
	//		if err := inMemoryRepository.SetState(state); err != nil {
	//			slog.Error(
	//				"failed to restore url storage from file",
	//				slog.String("error", err.Error()),
	//			)
	//		}
	//	}
	//	a.repository = inMemoryRepository
	//	a.stater = inMemoryRepository

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
