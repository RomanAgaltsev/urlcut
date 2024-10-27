package url

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/RomanAgaltsev/urlcut/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("Storage test", func(t *testing.T) {
		urlS := &model.URL{
			Long: "https://app.pachca.com",
			Base: "http://localhost:8080",
			ID:   "1q2w3e4r",
		}

		InMemRepo := New("test.json")

		err := InMemRepo.RestoreState()
		require.NoError(t, err)

		err = InMemRepo.Store(urlS)
		require.NoError(t, err)

		urlG, err := InMemRepo.Get("1q2w3e4r")
		require.NoError(t, err)
		assert.Equal(t, *urlS, *urlG)

		err = InMemRepo.SaveState()
		require.NoError(t, err)

		file, err := os.OpenFile("test.json", os.O_RDONLY|os.O_CREATE, 0666)
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
	})
}
