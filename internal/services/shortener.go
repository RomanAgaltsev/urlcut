package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/pkg/random"
	"github.com/RomanAgaltsev/urlcut/internal/repository"

	"github.com/google/uuid"
)

// Неиспользуемая переменная для проверки соответствия сокращателя интерфейсу сервиса
var _ interfaces.Service = (*Shortener)(nil)

// ErrInitServiceFailed ошибка инициализации сервиса сокращателя.
var ErrInitServiceFailed = fmt.Errorf("failed to init service")

// Shortener реализует сервис сокращателя ссылок.
type Shortener struct {
	repository interfaces.Repository // Репозиторий (интерфейс) для хранения URL
	cfg        *config.Config        // Конфигурация приложения
}

// NewShortener создает новый сокращатель ссылок.
func NewShortener(repository interfaces.Repository, cfg *config.Config) (*Shortener, error) {
	// Считаем, что без базового адреса и без идентификатора сокращенной ссылки быть не может
	if cfg.BaseURL == "" || cfg.IDlength == 0 {
		return nil, ErrInitServiceFailed
	}

	return &Shortener{
		repository: repository,
		cfg:        cfg,
	}, nil
}

// Shorten сокращает переданную ссылку.
func (s *Shortener) Shorten(ctx context.Context, longURL string, uid uuid.UUID) (*model.URL, error) {
	// Создаем структуру URL
	url := &model.URL{
		Long: longURL,
		Base: s.cfg.BaseURL,
		ID:   random.String(s.cfg.IDlength),
		UID:  uid,
	}

	// Сохраняем структуру URL в репозитории
	duplicatedURL, err := s.repository.Store(ctx, []*model.URL{url})
	if errors.Is(err, repository.ErrConflict) {
		// При наличии конфликта должна вернуться ранее сохраненная сокращенная ссылка
		return duplicatedURL, err
	}
	if err != nil {
		return &model.URL{}, err
	}

	return url, nil
}

// ShortenBatch сокращает переданный батч ссылок.
func (s *Shortener) ShortenBatch(ctx context.Context, batch []model.IncomingBatchDTO, uid uuid.UUID) ([]model.OutgoingBatchDTO, error) {
	// Создаем слайс для хранения сокращенных ссылок батча
	batchShortened := make([]model.OutgoingBatchDTO, 0, len(batch))
	// Создаем слайс URL
	urls := make([]*model.URL, 0, len(batch))

	// Обходим батч, создаем сокращенные ссылки и сохраняем в слайс URL
	for _, batchReq := range batch {
		urls = append(urls, &model.URL{
			Long:   batchReq.OriginalURL,
			Base:   s.cfg.BaseURL,
			ID:     random.String(s.cfg.IDlength),
			CorrID: batchReq.CorrelationID,
			UID:    uid,
		})
	}

	// Сохраняем слайс URL в БД
	duplicatedURL, err := s.repository.Store(ctx, urls)
	if err != nil && !errors.Is(err, repository.ErrConflict) {
		return batchShortened, err
	}

	// Если был конфликт, вернем дубль
	if errors.Is(err, repository.ErrConflict) {
		// При наличии конфликта должна вернуться ранее сохраненная сокращенная ссылка
		batchShortened = append(batchShortened, model.OutgoingBatchDTO{
			CorrelationID: duplicatedURL.CorrID,
			ShortURL:      duplicatedURL.Short(),
		})

		return batchShortened, nil
	}

	// Перекладываем сокращенные ссылки в слайс батча для возврата
	for _, url := range urls {
		batchShortened = append(batchShortened, model.OutgoingBatchDTO{
			CorrelationID: url.CorrID,
			ShortURL:      url.Short(),
		})
	}

	return batchShortened, nil
}

// Expand возвращает оригинальную ссылку по переданному идентификатору.
func (s *Shortener) Expand(ctx context.Context, id string) (*model.URL, error) {
	url, err := s.repository.Get(ctx, id)
	if err != nil {
		return &model.URL{}, fmt.Errorf("expanding URL failed: %w", err)
	}

	return url, nil
}

// UserURLs возвращает слайс URL пользователя с переданным uid пользователя.
func (s *Shortener) UserURLs(ctx context.Context, uid uuid.UUID) ([]model.UserURLDTO, error) {
	urls, err := s.repository.GetUserURLs(ctx, uid)
	if err != nil {
		return nil, err
	}

	// Создаем слайс для возврата ссылок пользователя
	userURLs := make([]model.UserURLDTO, 0, len(urls))

	// Перекладываем URL пользователя в слайс структур DTO
	for _, url := range urls {
		userURLs = append(userURLs, model.UserURLDTO{
			ShortURL:    url.Short(),
			OriginalURL: url.Long,
		})
	}

	return userURLs, nil
}

// DeleteUserURLs устанавливает пометку на удаление у всех URL с переданным uid пользователя и идентификаторами URL.
func (s *Shortener) DeleteUserURLs(ctx context.Context, uid uuid.UUID, shortURLs *model.ShortURLsDTO) error {
	// Создаем слайс URL для передачи в репозиторий
	urls := make([]*model.URL, 0, len(shortURLs.IDs))

	// Перекладываем идентификаторы и uid в слайс URL
	for _, id := range shortURLs.IDs {
		urls = append(urls, &model.URL{
			ID:  id,
			UID: uid,
		})
	}

	return s.repository.DeleteURLs(ctx, urls)
}

// Close закрывает репозиторий ссылок сокращателя.
func (s *Shortener) Close() error {
	return s.repository.Close()
}
