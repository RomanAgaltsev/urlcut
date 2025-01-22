package interfaces

import (
	"context"

	"github.com/google/uuid"

	"github.com/RomanAgaltsev/urlcut/internal/model"
)

// Service интерфейс сервиса сокращения URL.
type Service interface {
	Shorten(ctx context.Context, longURL string, uid uuid.UUID) (*model.URL, error)
	ShortenBatch(ctx context.Context, batch []model.IncomingBatchDTO, uid uuid.UUID) ([]model.OutgoingBatchDTO, error)
	Expand(ctx context.Context, id string) (*model.URL, error)
	UserURLs(ctx context.Context, uid uuid.UUID) ([]model.UserURLDTO, error)
	DeleteUserURLs(ctx context.Context, uid uuid.UUID, shortURLs *model.ShortURLsDTO) error
	Close() error
}

// Repository интерфейс хранилища сокращенных URL.
type Repository interface {
	Store(ctx context.Context, urls []*model.URL) (*model.URL, error)
	Get(ctx context.Context, id string) (*model.URL, error)
	GetUserURLs(ctx context.Context, uid uuid.UUID) ([]*model.URL, error)
	DeleteURLs(ctx context.Context, urls []*model.URL) error
	Close() error
}
