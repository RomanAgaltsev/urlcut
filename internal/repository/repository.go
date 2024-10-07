package repository

import "errors"

// Repository - интерфейс хранилища сокращенных URL
type Repository interface {
    Store(id string, url string) error
    GetUrl(id string) (string, error)
}

// MapRepository - хранилище-мапа
type MapRepository struct {
    m map[string]string
}

// Store - сохраняет пару ID/URL в хранилище-мапе
func (r *MapRepository) Store(id string, url string) error {
    r.m[id] = url
    return nil
}

// GetUrl - возвращает из хранилища URL по его ID
func (r *MapRepository) GetUrl(id string) (string, error) {
    if url, ok := r.m[id]; ok {
        return url, nil
    }
    return "", errors.New("URL ID was not found in repository")
}

// New - конструктор хранилища-мапы
func New() *MapRepository {
    return &MapRepository{
        m: make(map[string]string),
    }
}
