package repository

import (
	"fmt"

	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
)

var ErrInitRepositoryFailed = fmt.Errorf("failed to init repository")

func New(databaseDSN string, fileStoragePath string) (interfaces.Repository, error) {
	if databaseDSN == "" {
		return NewInMemoryRepository(fileStoragePath), nil
	}
	dbRepository, err := NewDBRepository(databaseDSN)
	if err != nil {
		return nil, ErrInitRepositoryFailed
	}
	return dbRepository, nil
}
