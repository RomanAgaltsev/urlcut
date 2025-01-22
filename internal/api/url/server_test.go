package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/RomanAgaltsev/urlcut/internal/config"
)

func TestServer(t *testing.T) {
	hlp := newHelper(t)

	cfg := &config.Config{}

	_, err := NewServer(hlp.shortener, cfg)
	assert.Equal(t, ErrInitServerFailed, err)

	server, err := NewServer(hlp.shortener, hlp.cfg)
	require.NoError(t, err)
	assert.Equal(t, hlp.cfg.ServerPort, server.Addr)
}
