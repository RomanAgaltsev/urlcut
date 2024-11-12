package interfaces

import "github.com/RomanAgaltsev/urlcut/internal/model"

// Service интерфейс сервиса сокращения URL.
type Service interface {
	Shorten(longURL string) (*model.URL, error)
	ShortenBatch(batch []model.BatchRequest) ([]model.BatchResponse, error)
	Expand(id string) (*model.URL, error)
	Close() error
}

// Repository интерфейс хранилища сокращенных URL.
type Repository interface {
	Store(urls []*model.URL) (*model.URL, error)
	Get(id string) (*model.URL, error)
	Close() error
}
