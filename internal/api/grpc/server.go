package url

import (
	"context"

	"github.com/google/uuid"

	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/services"
	pb "github.com/RomanAgaltsev/urlcut/pkg/shortener/v1"
)

// ShortenerService интерфейс сервиса сокращения URL.
type ShortenerService interface {
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

	// Stats возвращает статистику по количеству сокращенных URL и пользователей.
	Stats(ctx context.Context) (*model.StatsDTO, error)
}

var _ ShortenerService = (*services.Shortener)(nil)

type ShortenerServer struct {
	pb.UnimplementedURLShortenerServiceServer
	shortener ShortenerService
}

func NewShortenerServer(shortener ShortenerService) *ShortenerServer {
	return &ShortenerServer{
		shortener: shortener,
	}
}
