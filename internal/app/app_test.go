package app

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	t.Run("App test", func(t *testing.T) {
		application, err := New()
		require.NoError(t, err)

		go func() {
			_ = application.Run()
		}()
		time.Sleep(1 * time.Second)

		p, err := os.FindProcess(os.Getpid())
		require.NoError(t, err)

		err = p.Signal(syscall.SIGTERM)
		require.NoError(t, err)
	})
}
