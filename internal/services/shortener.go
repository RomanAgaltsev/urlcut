package services

import (
	"fmt"

	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/lib/random"
	"github.com/RomanAgaltsev/urlcut/internal/model"
)

var _ interfaces.URLShortExpander = (*Shortener)(nil)

type Shortener struct {
	repository interfaces.URLStoreGetter
	baseURL    string
	idLenght   int
}

func NewShortener(repository interfaces.URLStoreGetter, baseURL string, idLength int) *Shortener {
	return &Shortener{
		repository: repository,
		baseURL:    baseURL,
		idLenght:   idLength,
	}
}

func (s *Shortener) Shorten(longURL string) (*model.URL, error) {
	url := &model.URL{
		Long: longURL,
		Base: s.baseURL,
		ID:   random.String(s.idLenght),
	}

	err := s.repository.Store(url)
	if err != nil {
		return &model.URL{}, err
	}

	return url, nil
}

func (s *Shortener) Expand(id string) (*model.URL, error) {
	url, err := s.repository.Get(id)
	if err != nil {
		return &model.URL{}, fmt.Errorf("expanding URL failed: %w", err)
	}

	return url, nil
}

func (s *Shortener) Check() error {
	return s.repository.Check()
}
