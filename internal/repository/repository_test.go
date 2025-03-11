package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RomanAgaltsev/urlcut/internal/config"
)

func TestNewRepository(t *testing.T) {
	cfg := &config.Config{}

	repo, err := NewRepository(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, repo, nil)

	cfg.DatabaseDSN = "postgres://postgres:12345@localhost:5432/praktikum?sslmode=disable"

	repo, err = NewRepository(cfg)
	assert.Error(t, err)
	assert.Equal(t, repo, nil)
}
