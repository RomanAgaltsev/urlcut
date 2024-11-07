package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	hlp := newHelper(t)

	server, err := NewServer(hlp.shortener, hlp.serverPort)
	require.NoError(t, err)
	assert.Equal(t, hlp.serverPort, server.Addr)
}
