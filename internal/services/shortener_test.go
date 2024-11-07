package services

import (
	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/RomanAgaltsev/urlcut/internal/mocks"
)

func TestShortener(t *testing.T) {
	const (
		baseURL  = "http://localhost:8080"
		longURL  = "https://practicum.yandex.ru/"
		urlID    = "h7Ds18sD"
		idLength = 8
	)
	url := &model.URL{
		Long: longURL,
		Base: baseURL,
		ID:   urlID,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		Store(gomock.Any()).
		Return(nil).
		Times(1)

	shortener, err := NewShortener(mockRepo, baseURL, idLength)
	require.NoError(t, err)

	urlS, err := shortener.Shorten(longURL)
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
		}, nil).
		Times(1)

	urlE, err := shortener.Expand(urlS.ID)
	require.NoError(t, err)
	assert.Equal(t, urlS, urlE)

	mockRepo.
		EXPECT().
		Check().
		Return(nil).
		Times(1)

	err = shortener.Check()
	require.NoError(t, err)

	mockRepo.
		EXPECT().
		Close().
		Return(nil).
		Times(1)

	err = shortener.Close()
	require.NoError(t, err)
}
