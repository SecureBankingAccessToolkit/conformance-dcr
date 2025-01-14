package auth

import (
	"testing"
	"time"

	"github.com/OpenBankingUK/conformance-dcr/pkg/certs"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientSecretJWT_Claims(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	auther := NewClientSecretJWT(
		"tokenEndpoint",
		NewJwtSigner(
			jwt.SigningMethodRS256,
			"ssa",
			"issuer",
			"aud",
			"kid",
			"client_secret_jwt",
			"none",
			[]string{},
			nil,
			privateKey,
			time.Hour,
			nil,
			"",
			"",
			"",
		),
	)

	claims, err := auther.Claims()

	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
}

func TestClientSecretJWT_Client_ReturnsAClient(t *testing.T) {
	privateKey, err := certs.ParseRsaPrivateKeyFromPemFile("testdata/private-sign.key")
	require.NoError(t, err)
	auther := NewClientSecretJWT(
		"tokenEndpoint",
		NewJwtSigner(
			jwt.SigningMethodRS256,
			"ssa",
			"issuer",
			"aud",
			"kid",
			"client_secret_jwt",
			"none",
			[]string{},
			nil,
			privateKey,
			time.Hour,
			nil,
			"",
			"",
			"",
		),
	)

	client, err := auther.Client([]byte(`{"client_id": "12345", "registration_access_token": "abcdef", "client_secret": "54321"}`))

	require.NoError(t, err)
	assert.Equal(t, "12345", client.Id())
	assert.Equal(t, "abcdef", client.RegistrationAccessToken())
}
