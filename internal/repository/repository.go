package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/database"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
)

var (
	ErrInitRepositoryFailed = fmt.Errorf("failed to init repository")
	ErrConflict             = fmt.Errorf("data conflict")
)

func New(cfg *config.Config) (interfaces.Repository, error) {
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
