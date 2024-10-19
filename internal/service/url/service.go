package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
)

type ShortenerService struct {
	repo     repository.URLRepository
	baseURL  string
	idLenght int
}

func NewShortener(repo repository.URLRepository, baseURL string, idLength int) *ShortenerService {
	return &ShortenerService{
		repo:     repo,
		baseURL:  baseURL,
		idLenght: idLength,
	}
}

func (s *ShortenerService) Shorten(ctx context.Context, longURL string) (*model.URL, error) {
	url := &model.URL{
		LongURL: longURL,
		ID:      urlID(s.idLenght),
	}
	err := s.repo.Store(ctx, url)
	if err != nil {
		return &model.URL{}, err
	}
	return url, nil
}

func (s *ShortenerService) Expand(ctx context.Context, id string) (*model.URL, error) {
	url, err := s.repo.Get(ctx, id)
	if err != nil {
		return &model.URL{}, fmt.Errorf("expanding URL failed: %w", err)
	}
	return url, nil
}

func urlID(lenght int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, lenght)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
