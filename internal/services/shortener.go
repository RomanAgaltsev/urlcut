package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/lib/random"
	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
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
func (s *Shortener) Shorten(longURL string) (*model.URL, error) {
	// Создаем структуру URL
	url := &model.URL{
		Long: longURL,
		Base: s.cfg.BaseURL,
		ID:   random.String(s.cfg.IDlength),
	}

	// Сохраняем структуру URL в репозитории
	duplicatedURL, err := s.repository.Store([]*model.URL{url})
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
func (s *Shortener) ShortenBatch(batch []model.BatchRequest) ([]model.BatchResponse, error) {
	// Создаем слайс для хранения сокращенных ссылок батча
	batchShortened := make([]model.BatchResponse, 0, len(batch))
	// Создаем слайс URL
	urls := make([]*model.URL, 0, len(batch))

	// Обходим батч, создаем сокращенные ссылки и сохраняем в слайс URL
	for _, batchReq := range batch {
		urls = append(urls, &model.URL{
			Long:   batchReq.OriginalURL,
			Base:   s.cfg.BaseURL,
			ID:     random.String(s.cfg.IDlength),
			CorrID: batchReq.CorrelationID,
		})
	}

	// Сохраняем слайс URL в БД
	_, err := s.repository.Store(urls)
	if err != nil {
		return batchShortened, err
	}

	// Перекладываем сокращенные ссылки в слайс батча для возврата
	for _, url := range urls {
		batchShortened = append(batchShortened, model.BatchResponse{
			CorrelationID: url.CorrID,
			ShortURL:      url.Short(),
		})
	}

	return batchShortened, nil
}

// Expand возвращает оригинальную ссылку по переданному идентификатору.
func (s *Shortener) Expand(id string) (*model.URL, error) {
	url, err := s.repository.Get(id)
	if err != nil {
		return &model.URL{}, fmt.Errorf("expanding URL failed: %w", err)
	}

	return url, nil
}

func (s *Shortener) UserURLs(uid uuid.UUID) ([]model.UserURL, error) {
	return nil, nil
}

// Close закрывает репозиторий ссылок сокращателя.
func (s *Shortener) Close() error {
	return s.repository.Close()
}
