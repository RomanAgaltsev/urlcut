package repository

import (
	"context"
	"github.com/google/uuid"
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

	uid, _ := uuid.NewRandom()

	url := &model.URL{
		Long:   longURL,
		Base:   BaseURL,
		ID:     urlID,
		CorrID: "",
		UID:    uid,
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

	userURLs, err := inMemoryRepository.GetUserURLs(context.TODO(), uid)
	require.NoError(t, err)
	assert.IsType(t, []*model.URL{url}, userURLs)

	err = inMemoryRepository.DeleteURLs(context.TODO(), urls)
	require.NoError(t, err)

	err = inMemoryRepository.Close()
	require.NoError(t, err)
}
