package url

import (
	"fmt"
	"sync"

	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
)

var _ repository.URLRepository = (*InMemoryRepository)(nil)

var ErrIDNotFound = fmt.Errorf("URL ID was not found in repository")

type InMemoryRepository struct {
	m map[string]*model.URL
	f string
	sync.RWMutex
}

func New(filename string) *InMemoryRepository {
	return &InMemoryRepository{
		m: make(map[string]*model.URL),
		f: filename,
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

func (r *InMemoryRepository) SaveState() error {
	strg, err := newStorage(r.f)
	if err != nil {
		return err
	}
	strg.m = r.m

	if err = strg.save(); err != nil {
		return err
	}

	return nil
}

func (r *InMemoryRepository) RestoreState() error {
	strg, err := newStorage(r.f)
	if err != nil {
		return err
	}

	if err = strg.restore()
		err != nil {
		return err
	}

	if len(strg.m) > 0 {
		r.m = strg.m
	}

	return nil
}
