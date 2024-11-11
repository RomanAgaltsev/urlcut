package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	application, err := New()
	require.NoError(t, err)

	go func() {
		_ = application.Run()
	}()
	time.Sleep(1 * time.Second)

	err = application.shortener.Close()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = application.server.Shutdown(ctx)
	require.NoError(t, err)
}
