package repository

import (
    //    "bufio"
    //    "encoding/json"
    //    "os"
    "testing"
    //
    //    "github.com/RomanAgaltsev/urlcut/internal/model"
    //
    //    "github.com/stretchr/testify/assert"
    //    "github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
    //	t.Run("Saver test", func(t *testing.T) {
    //		urlS := &model.URL{
    //			Long: "https://app.pachca.com",
    //			Base: "http://localhost:8080",
    //			ID:   "1q2w3e4r",
    //		}
    //		state := map[string]*model.URL{
    //			"1q2w3e4r": urlS,
    //		}
    //
    //		inMemoryRepository := repository.NewInMemoryRepository("")
    //		stateSaver := NewStateSaver("test.json")
    //
    //		err := stateSaver.SaveState(state)
    //		require.NoError(t, err)
    //
    //		restoredState, err := stateSaver.RestoreState()
    //		require.NoError(t, err)
    //		assert.Equal(t, state, restoredState)
    //
    //		err = inMemoryRepository.SetState(restoredState)
    //		require.NoError(t, err)
    //
    //		fromMemoryState := inMemoryRepository.GetState()
    //		assert.Equal(t, state, fromMemoryState)
    //
    //		err = stateSaver.SaveState(state)
    //		require.NoError(t, err)
    //
    //		file, err := os.OpenFile("test.json", os.O_RDONLY|os.O_CREATE, 0666)
    //		require.NoError(t, err)
    //
    //		scanner := bufio.NewScanner(file)
    //		scanner.Scan()
    //
    //		data := scanner.Bytes()
    //
    //		var urlU model.URL
    //		err = json.Unmarshal(data, &urlU)
    //		require.NoError(t, err)
    //		assert.Equal(t, *urlS, urlU)
    //
    //		err = file.Close()
    //		require.NoError(t, err)
    //
    //		err = os.Remove("test.json")
    //		require.NoError(t, err)
    //	})
}
