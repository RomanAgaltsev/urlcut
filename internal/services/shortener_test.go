package services

import (
	"testing"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/mocks"
	"github.com/RomanAgaltsev/urlcut/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestShortener(t *testing.T) {
	const (
		serverPort = "localhost:8080"
		baseURL    = "http://localhost:8080"
		longURL    = "https://practicum.yandex.ru/"
		urlID      = "h7Ds18sD"
		idLength   = 8
	)

	uid := uuid.New()

	cfg := &config.Config{
		ServerPort:      serverPort,
		BaseURL:         baseURL,
		FileStoragePath: "",
		DatabaseDSN:     "",
		IDlength:        idLength,
	}

	url := &model.URL{
		Long: longURL,
		Base: baseURL,
		ID:   urlID,
		UID:  uid,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		Store(gomock.Any()).
		Return(nil, nil).
		Times(1)

	shortener, err := NewShortener(mockRepo, cfg)
	require.NoError(t, err)

	urlS, err := shortener.Shorten(longURL, uid)
	require.NoError(t, err)
	assert.Equal(t, url.Long, urlS.Long)
	assert.Equal(t, url.Base, urlS.Base)
	assert.IsType(t, url.ID, urlS.ID)
	assert.Len(t, urlS.ID, idLength)

	mockRepo.
		EXPECT().
		Get(urlS.ID).
		Return(&model.URL{
			Long: urlS.Long,
			Base: urlS.Base,
			ID:   urlS.ID,
			UID:  uid,
		}, nil).
		Times(1)

	urlE, err := shortener.Expand(urlS.ID)
	require.NoError(t, err)
	assert.Equal(t, urlS, urlE)

	mockRepo.
		EXPECT().
		Close().
		Return(nil).
		Times(1)

	err = shortener.Close()
	require.NoError(t, err)
}
