package repository

import (
	"fmt"
)

var ErrIDNotFound = fmt.Errorf("URL ID was not found in repository")

// New - конструктор хранилища-мапы
func New() *InMemoryRepository {
	return &InMemoryRepository{
		m: make(map[string]string),
	}
}

// Repository - интерфейс хранилища сокращенных URL
type Repository interface {
	Store(id string, url string) error
	GetURL(id string) (string, error)
}

// InMemoryRepository - хранилище-мапа
type InMemoryRepository struct {
	m map[string]string
}

// Store - сохраняет пару ID/URL в хранилище-мапе
func (r *InMemoryRepository) Store(id string, url string) error {
	r.m[id] = url
	return nil
}

// GetURL - возвращает из хранилища URL по его ID
func (r *InMemoryRepository) GetURL(id string) (string, error) {
	if url, ok := r.m[id]; ok {
		return url, nil
	}
	return "", ErrIDNotFound
}
