package repository

import (
	"context"
	"testing"

	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryRepository(t *testing.T) {
	const fileStoragePath = "test.json"

	const (
		longURL = "https://app.pachca.com"
		BaseURL = "http://localhost:8080"
		urlID   = "1q2w3e4r"
	)

	url := &model.URL{
		Long:   longURL,
		Base:   BaseURL,
		ID:     urlID,
		CorrID: "",
	}

	urls := []*model.URL{url}

	inMemoryRepository := NewInMemoryRepository(fileStoragePath)

	urlStore, err := inMemoryRepository.Store(context.TODO(), urls)
	require.NoError(t, err)
	assert.Nil(t, urlStore)

	urlGet, err := inMemoryRepository.Get(context.TODO(), urlID)
	require.NoError(t, err)
	assert.NotNil(t, urlGet)
	assert.Equal(t, url, urlGet)

	err = inMemoryRepository.Close()
	require.NoError(t, err)
}
