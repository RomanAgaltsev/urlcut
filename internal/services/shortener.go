// Пакет services реализует сервис сокращения URL.
package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/pkg/random"
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

	urlDelChan chan *model.URL // Канал для сбора URL к отложенному удалению
}

// NewShortener создает новый сокращатель ссылок.
func NewShortener(repository interfaces.Repository, cfg *config.Config) (*Shortener, error) {
	// Считаем, что без базового адреса и без идентификатора сокращенной ссылки быть не может
	if cfg.BaseURL == "" || cfg.IDlength == 0 {
		return nil, ErrInitServiceFailed
	}

	shortener := &Shortener{
		repository: repository,
		cfg:        cfg,
		urlDelChan: make(chan *model.URL, 1024),
	}

	// Запускаем горутину фоновой пометки на удаление URL
	go shortener.deleteURLs()

	return shortener, nil
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
	// Перекладываем идентификаторы и uid в слайс URL
	for _, id := range shortURLs.IDs {
		s.urlDelChan <- &model.URL{
			ID:  id,
			UID: uid,
		}
	}

	return nil
}

// deleteURLs устанавливаем пометку на удаление URL с определенным интервалом.
func (s *Shortener) deleteURLs() {
	// Сохраняем URL, накопленные за последние 10 секунд
	ticker := time.NewTicker(10 * time.Second)

	// Накапливаем URL к удалению в слайсе
	var urls []*model.URL

	for {
		select {
		// Пробуем получить URL из канала
		case url := <-s.urlDelChan:
			// Полученный URL добавляем в слайс
			urls = append(urls, url)
		case <-ticker.C:
			// Очередной интервал прошел, проверяем наличие URL к удалению
			if len(urls) == 0 {
				continue
			}
			// Устанавливаем пометку на удаление для полученных URL
			err := s.repository.DeleteURLs(context.TODO(), urls)
			if err != nil {
				slog.Info("failed to delete URLs", "error", err.Error())
				continue
			}
			// Обнуляем слайс URL
			urls = nil
		}
	}
}

func (s *Shortener) Stats(ctx context.Context) (*model.StatsDTO, error) {
	return nil, nil
}

// Close закрывает репозиторий ссылок сокращателя.
func (s *Shortener) Close() error {
	// Закрываем канал сбора URL к удалению
	close(s.urlDelChan)
	// Закрываем соединения
	return s.repository.Close()
}
