package service

import (
	"fmt"
	"math/rand"

	"github.com/RomanAgaltsev/urlcut/internal/repository"
)

type Service interface {
	ShortenURL(url string) (string, error)
	ExpandURL(id string) (string, error)
}

// ShortenerService - структура сервиса сокращения URL
type ShortenerService struct {
	repo     repository.Repository // Репозиторий
	baseURL  string                // Базовый URL для сокращенного URL
	idLenght int                   // Длина идентификатора сокращенного URL
}

// NewShortener - конструктор нового сервиса сокращения URL
func NewShortener(repo repository.Repository, baseURL string, idLength int) *ShortenerService {
	return &ShortenerService{
		repo:     repo,
		baseURL:  baseURL,
		idLenght: idLength,
	}
}

// ShortenURL - сокращает URL, сохраняет пару ID/URL и возвращает сокращенный URL
func (s *ShortenerService) ShortenURL(url string) (string, error) {
	// Получаем новый произвольный идентификатор заданной длины
	id := urlID(s.idLenght)
	// Сохраняем пару ID/URL в репозитории
	err := s.repo.Store(id, url)
	// Проверяем на ошибку
	if err != nil {
		// Была ошибка
		return "", err
	}
	// Возвращаем сокращенный URL без ошибки
	return fmt.Sprintf("%s/%s", s.baseURL, id), nil
}

// ExpandURL - вовзращает оригинальный URL по переданному ID
func (s *ShortenerService) ExpandURL(id string) (string, error) {
	// Получаем оригинальный URL из хранилища по ID
	url, err := s.repo.GetURL(id)
	// Проверяем наличие ошибки
	if err != nil {
		return "", fmt.Errorf("expanding URL failed: %w", err)
	}
	// Возвращаем оригинальный URL
	return url, nil
}

// urlID - возвращает идентификатор сокращенного URL
func urlID(lenght int) string {
	// Символы, которые могут входить в идентификатор
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Инициируем слайс байт с длиной, равной длине идентификатора
	b := make([]byte, lenght)
	// Заполняем слайс произвольными символами из доступных
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	// Возвращаем получившуюся строку - идентификатор
	return string(b)
}
