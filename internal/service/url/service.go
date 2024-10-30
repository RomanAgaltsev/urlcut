package service

import (
	"fmt"

	"github.com/RomanAgaltsev/urlcut/internal/lib/random"
	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
	"github.com/RomanAgaltsev/urlcut/internal/service"
)

var _ service.URLService = (*ShortenerService)(nil)

type ShortenerService struct {
	repo     repository.URLRepository
	baseURL  string
	idLenght int
}

func New(repo repository.URLRepository, baseURL string, idLength int) *ShortenerService {
	return &ShortenerService{
		repo:     repo,
		baseURL:  baseURL,
		idLenght: idLength,
	}
}

func (s *ShortenerService) Shorten(longURL string) (*model.URL, error) {
	url := &model.URL{
		Long: longURL,
		Base: s.baseURL,
		ID:   random.String(s.idLenght),
	}

	err := s.repo.Store(url)
	if err != nil {
		return &model.URL{}, err
	}

	return url, nil
}

func (s *ShortenerService) Expand(id string) (*model.URL, error) {
	url, err := s.repo.Get(id)
	if err != nil {
		return &model.URL{}, fmt.Errorf("expanding URL failed: %w", err)
	}

	return url, nil
}
