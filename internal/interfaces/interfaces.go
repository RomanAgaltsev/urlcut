package interfaces

import (
	"context"

	"github.com/google/uuid"

	"github.com/RomanAgaltsev/urlcut/internal/model"
)

// Service интерфейс сервиса сокращения URL.
type Service interface {
	// Shorten сокращает переданный оригинальный URL и возвращает сокращенный URL.
	Shorten(ctx context.Context, longURL string, uid uuid.UUID) (*model.URL, error)

	// ShortenBatch сокращает переданный слайс оригинальных URL и возвращает слайс сокращенных URL.
	ShortenBatch(ctx context.Context, batch []model.IncomingBatchDTO, uid uuid.UUID) ([]model.OutgoingBatchDTO, error)

	// Expand возвращает оригинальный URL по переданному идентификатору сокращенного URL.
	Expand(ctx context.Context, id string) (*model.URL, error)

	// UserURLs возвращает слайс URL по переданному идентификатору пользователя.
	UserURLs(ctx context.Context, uid uuid.UUID) ([]model.UserURLDTO, error)

	// DeleteUserURLs удаляет URL пользователя по переданным идентификаторам сокращенных URL.
	DeleteUserURLs(ctx context.Context, uid uuid.UUID, shortURLs *model.ShortURLsDTO) error

	// Close закрывает сервис. Используется в текущей реализации graceful shutdown.
	Close() error
}

// Repository интерфейс хранилища сокращенных URL.
type Repository interface {
	// Store сохраняет переданные URL в базе данных.
	Store(ctx context.Context, urls []*model.URL) (*model.URL, error)

	// Get возвращает из базы данных URL по переданному идентификатору сокращенного URL.
	Get(ctx context.Context, id string) (*model.URL, error)

	// GetUserURLs возвращает из БД URL пользователя по переданному идентификатору.
	GetUserURLs(ctx context.Context, uid uuid.UUID) ([]*model.URL, error)

	// DeleteURLs удаляет в БД переданные URL.
	DeleteURLs(ctx context.Context, urls []*model.URL) error

	// Close - закрывает соединения с БД. Используется в текущей реализации graceful shutdown.
	Close() error
}
