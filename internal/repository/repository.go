package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/RomanAgaltsev/urlcut/internal/database"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
)

var (
	ErrInitRepositoryFailed = fmt.Errorf("failed to init repository")
	ErrConflict             = fmt.Errorf("data conflict")
)

func New(databaseDSN string, fileStoragePath string) (interfaces.Repository, error) {
	if databaseDSN == "" {
		return NewInMemoryRepository(fileStoragePath), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Открываем новое соединение
	db, err := database.NewConnection(ctx, "pgx", databaseDSN)
	if err != nil {
		return nil, ErrInitRepositoryFailed
	}

	dbRepository, err := NewDBRepository(db)
	if err != nil {
		return nil, ErrInitRepositoryFailed
	}
	return dbRepository, nil
}
