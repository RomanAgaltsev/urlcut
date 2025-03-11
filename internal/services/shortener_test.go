package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/mocks"
	"github.com/RomanAgaltsev/urlcut/internal/model"
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
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	shortener, err := NewShortener(mockRepo, cfg)
	require.NoError(t, err)

	mockRepo.EXPECT().
		Store(context.TODO(), gomock.Any()).
		Return(nil, nil).
		Times(1)

	urlS, err := shortener.Shorten(context.TODO(), longURL, uid)
	require.NoError(t, err)
	assert.Equal(t, url.Long, urlS.Long)
	assert.Equal(t, url.Base, urlS.Base)
	assert.IsType(t, url.ID, urlS.ID)
	assert.Len(t, urlS.ID, idLength)

	mockRepo.
		EXPECT().
		Get(context.TODO(), urlS.ID).
		Return(&model.URL{
			Long: urlS.Long,
			Base: urlS.Base,
			ID:   urlS.ID,
			UID:  uid,
		}, nil).
		Times(1)

	urlE, err := shortener.Expand(context.TODO(), urlS.ID)
	require.NoError(t, err)
	assert.Equal(t, urlS, urlE)

	mockRepo.
		EXPECT().
		Close().
		Return(nil).
		Times(1)

	mockRepo.EXPECT().
		Store(context.TODO(), gomock.Any()).
		Return(nil, nil).
		Times(1)

	inBatch := []model.IncomingBatchDTO{
		{CorrelationID: urlID, OriginalURL: longURL},
	}

	outBatch, err := shortener.ShortenBatch(context.TODO(), inBatch, uid)
	require.NoError(t, err)
	assert.Equal(t, len(inBatch), len(outBatch))

	mockRepo.EXPECT().
		GetUserURLs(context.TODO(), gomock.Any()).
		Return(nil, nil).
		Times(1)

	_, err = shortener.UserURLs(context.TODO(), uid)
	require.NoError(t, err)

	mockRepo.EXPECT().
		DeleteURLs(context.TODO(), gomock.Any()).
		Return(nil).
		Times(1)

	shortURLs := model.ShortURLsDTO{IDs: []string{urlID}}

	err = shortener.DeleteUserURLs(context.TODO(), uid, &shortURLs)
	require.NoError(t, err)

	mockRepo.EXPECT().
		GetStats(context.TODO()).
		Return(&model.Stats{
			Urls:  0,
			Users: 0,
		}, nil).
		Times(1)

	_, err = shortener.Stats(context.TODO())
	require.NoError(t, err)

	time.Sleep(11 * time.Second)

	err = shortener.Close()
	require.NoError(t, err)
}
