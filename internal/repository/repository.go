package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/database"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
)

// Переменные ошибок.
var (
	// ErrInitRepositoryFailed ошибка инициации репозитория.
	ErrInitRepositoryFailed = fmt.Errorf("failed to init repository")

	// ErrConflict ошибка конфликта данных в БД.
	ErrConflict = fmt.Errorf("data conflict")
)

// NewRepository создает и возвращает новый репозиторий в соответствии с переданной конфигурацией приложения.
func NewRepository(cfg *config.Config) (interfaces.Repository, error) {
	if cfg.DatabaseDSN == "" {
		return NewInMemoryRepository(cfg.FileStoragePath), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Открываем новое соединение
	db, err := database.NewConnection(ctx, "pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, ErrInitRepositoryFailed
	}

	dbRepository, err := NewDBRepository(db)
	if err != nil {
		return nil, ErrInitRepositoryFailed
	}
	return dbRepository, nil
}
