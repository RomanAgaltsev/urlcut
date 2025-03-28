package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/model"
)

// Неиспользуемая переменная для проверки реализации интерфейса хранилища in memory репозиторием.
var _ interfaces.Repository = (*InMemoryRepository)(nil)

// Переменные ошибок.
var (
	// ErrIDNotFound ошибка отсутствия URL в хранилище.
	ErrIDNotFound = fmt.Errorf("URL ID was not found in repository")

	// ErrStorageUnavailable ошибка недоступности хранилища.
	ErrStorageUnavailable = fmt.Errorf("storage unavailable")
)

// InMemoryRepository реализует in memory репозиторий.
type InMemoryRepository struct {
	m  map[string]*model.URL // мапа для хранения URL
	f  string                // адрес файлового хранилища - путь к файлу
	mu sync.RWMutex          // мьютекс хранилища
}

// NewInMemoryRepository создает новый in memory репозиторий.
func NewInMemoryRepository(fileStoragePath string) *InMemoryRepository {
	m, _ := readFromFile(fileStoragePath)

	return &InMemoryRepository{
		m: m,
		f: fileStoragePath,
	}
}

// Store сохраняет данные URL в in memory репозитории.
func (r *InMemoryRepository) Store(_ context.Context, urls []*model.URL) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, url := range urls {
		r.m[url.ID] = url
	}

	return nil, nil
}

// Get возвращает данные URL из in memory репозитория.
func (r *InMemoryRepository) Get(_ context.Context, id string) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if url, ok := r.m[id]; ok {
		return url, nil
	} else {
		return url, ErrIDNotFound
	}
}

// GetUserURLs возвращает URL пользователя из репозитория.
func (r *InMemoryRepository) GetUserURLs(_ context.Context, uid uuid.UUID) ([]*model.URL, error) {
	// Создаем слайс для возврата ссылок пользователя
	urls := make([]*model.URL, 0)

	// Отбираем URL пользователя простым перебором
	for _, url := range r.m {
		if url.UID == uid {
			urls = append(urls, url)
		}
	}

	return urls, nil
}

// DeleteURLs удаляет URL пользователя из репозитория.
func (r *InMemoryRepository) DeleteURLs(_ context.Context, urls []*model.URL) error {
	// Перебираем URL в цикле и устанавливаем пометки удаления
	for _, url := range urls {
		u, ok := r.m[url.ID]
		if ok && u.UID == url.UID {
			u.Deleted = true
		}
	}

	return nil
}

// Close сохраняет данные из in memory репозитория в файловое хранилище.
func (r *InMemoryRepository) Close() error {
	return writeToFile(r.f, r.m)
}
