package repository

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/RomanAgaltsev/urlcut/internal/model"
)

func TestFileStorage(t *testing.T) {
	const fileStoragePath = "test.json"

	urlS := &model.URL{
		Long: "https://app.pachca.com",
		Base: "http://localhost:8080",
		ID:   "1q2w3e4r",
	}
	state := map[string]*model.URL{
		"1q2w3e4r": urlS,
	}

	err := writeToFile(fileStoragePath, state)
	require.NoError(t, err)

	restoredState, err := readFromFile(fileStoragePath)
	require.NoError(t, err)
	assert.Equal(t, state, restoredState)

	file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	require.NoError(t, err)

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	data := scanner.Bytes()

	var urlU model.URL
	err = json.Unmarshal(data, &urlU)
	require.NoError(t, err)
	assert.Equal(t, *urlS, urlU)

	err = file.Close()
	require.NoError(t, err)

	err = os.Remove(fileStoragePath)
	require.NoError(t, err)

}
