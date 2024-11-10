package services

import (
	"errors"
	"fmt"

	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/lib/random"
	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
)

var _ interfaces.Service = (*Shortener)(nil)
var ErrInitServiceFailed = fmt.Errorf("failed to init service")

type Shortener struct {
	repository interfaces.Repository
	baseURL    string
	idLenght   int
}

func NewShortener(repository interfaces.Repository, baseURL string, idLength int) (*Shortener, error) {
	if baseURL == "" || idLength == 0 {
		return nil, ErrInitServiceFailed
	}

	return &Shortener{
		repository: repository,
		baseURL:    baseURL,
		idLenght:   idLength,
	}, nil
}

func (s *Shortener) Shorten(longURL string) (*model.URL, error) {
	url := &model.URL{
		Long: longURL,
		Base: s.baseURL,
		ID:   random.String(s.idLenght),
	}

	duplicatedURL, err := s.repository.Store([]*model.URL{url})
	if errors.Is(err, repository.ErrConflict) {
		return duplicatedURL, err
	}
	if err != nil {
		return &model.URL{}, err
	}

	return url, nil
}

func (s *Shortener) ShortenBatch(batch []model.BatchRequest) ([]model.BatchResponse, error) {
	batchShortened := make([]model.BatchResponse, 0, len(batch))
	urls := make([]*model.URL, 0, len(batch))

	for _, batchReq := range batch {
		urls = append(urls, &model.URL{
			Long:   batchReq.OriginalURL,
			Base:   s.baseURL,
			ID:     random.String(s.idLenght),
			CorrID: batchReq.CorrelationID,
		})
	}

	_, err := s.repository.Store(urls)
	if err != nil {
		return batchShortened, err
	}

	for _, url := range urls {
		batchShortened = append(batchShortened, model.BatchResponse{
			CorrelationID: url.CorrID,
			ShortURL:      url.Short(),
		})
	}

	return batchShortened, nil
}

func (s *Shortener) Expand(id string) (*model.URL, error) {
	url, err := s.repository.Get(id)
	if err != nil {
		return &model.URL{}, fmt.Errorf("expanding URL failed: %w", err)
	}

	return url, nil
}

func (s *Shortener) Close() error {
	return s.repository.Close()
}

func (s *Shortener) Check() error {
	return s.repository.Check()
}
