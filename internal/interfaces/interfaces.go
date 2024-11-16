package interfaces

import (
	"github.com/RomanAgaltsev/urlcut/internal/model"

	"github.com/google/uuid"
)

// Service интерфейс сервиса сокращения URL.
type Service interface {
	Shorten(longURL string, uid uuid.UUID) (*model.URL, error)
	ShortenBatch(batch []model.IncomingBatchDTO, uid uuid.UUID) ([]model.OutgoingBatchDTO, error)
	Expand(id string) (*model.URL, error)
	UserURLs(uid uuid.UUID) ([]model.UserURLDTO, error)
	Close() error
}

// Repository интерфейс хранилища сокращенных URL.
type Repository interface {
	Store(urls []*model.URL) (*model.URL, error)
	Get(id string) (*model.URL, error)
	GetUserURLs(uid uuid.UUID) ([]*model.URL, error)
	Close() error
}
