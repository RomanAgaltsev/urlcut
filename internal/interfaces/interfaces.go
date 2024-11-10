package interfaces

import (
	"github.com/RomanAgaltsev/urlcut/internal/model"
)

type Service interface {
	Shorten(longURL string) (*model.URL, error)
	ShortenBatch(batch []model.BatchRequest) ([]model.BatchResponse, error)
	Expand(id string) (*model.URL, error)
	Close() error
	Check() error
}

type Repository interface {
	Store(urls []*model.URL) (*model.URL, error)
	Get(id string) (*model.URL, error)
	Close() error
	Check() error
}
