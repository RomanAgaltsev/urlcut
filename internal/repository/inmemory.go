package repository

import (
	"fmt"
	"github.com/google/uuid"
	"sync"

	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/model"
)

// Неиспользуемая переменная для проверки реализации интерфейса хранилища in memory репозиторием
var _ interfaces.Repository = (*InMemoryRepository)(nil)

var (
	// ErrIDNotFound ошибка отсутствия URL в хранилище.
	ErrIDNotFound = fmt.Errorf("URL ID was not found in repository")
	// ErrStorageUnavailable ошибка недоступности хранилища.
	ErrStorageUnavailable = fmt.Errorf("storage unavailable")
)

// InMemoryRepository реализует in memory репозиторий.
type InMemoryRepository struct {
	m map[string]*model.URL
	f string
	sync.RWMutex
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
func (r *InMemoryRepository) Store(urls []*model.URL) (*model.URL, error) {
	r.Lock()
	defer r.Unlock()

	for _, url := range urls {
		r.m[url.ID] = url
	}

	return nil, nil
}

// Get возвращает данные URL из in memory репозитория.
func (r *InMemoryRepository) Get(id string) (*model.URL, error) {
	r.Lock()
	defer r.Unlock()

	if url, ok := r.m[id]; ok {
		return url, nil
	} else {
		return url, ErrIDNotFound
	}
}

func (r *InMemoryRepository) GetUserURLs(uid uuid.UUID) ([]*model.URL, error) {
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

// Close сохраняет данные из in memory репозитория в файловое хранилище.
func (r *InMemoryRepository) Close() error {
	return writeToFile(r.f, r.m)
}
