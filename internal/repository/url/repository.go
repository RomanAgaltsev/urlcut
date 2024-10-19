package url

import (
	"context"
	"fmt"

	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
)

var _ repository.URLRepository = (*InMemoryRepository)(nil)

var ErrIDNotFound = fmt.Errorf("URL ID was not found in repository")

func New() *InMemoryRepository {
	return &InMemoryRepository{
		m: make(map[string]string),
	}
}

type InMemoryRepository struct {
	m map[string]string
}

func (r *InMemoryRepository) Store(_ context.Context, url *model.URL) error {
	r.m[url.ID] = url.LongURL
	return nil
}

func (r *InMemoryRepository) Get(_ context.Context, id string) (*model.URL, error) {
	if longURL, ok := r.m[id]; ok {
		return &model.URL{LongURL: longURL, ID: id}, nil
	}
	return &model.URL{}, ErrIDNotFound
}
