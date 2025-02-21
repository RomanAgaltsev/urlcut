package cert

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateCertificate(t *testing.T) {
	const (
		certFileName = "cert.pem"
		keyFileName  = "privatekey.pem"
	)

	err := CreateCertificate(certFileName, keyFileName)

	require.NoError(t, err)
	require.FileExists(t, certFileName)
	require.FileExists(t, keyFileName)
}
