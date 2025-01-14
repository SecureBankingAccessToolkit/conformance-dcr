package step

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/auth"
	"github.com/OpenBankingUK/conformance-dcr/pkg/compliant/client"
	"github.com/stretchr/testify/assert"
)

const (
	clientID                = "foo"
	clientSecret            = "bar"
	registrationAccessToken = "accessToken123"
)

func TestNewClientDelete(t *testing.T) {
	// creating a stub server that expects a JWT body posted
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s", clientID), r.URL.EscapedPath())
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	softClient := client.NewClientSecretBasic(clientID, registrationAccessToken, server.URL, clientSecret)
	ctx := NewContext()
	ctx.SetClient("clientKey", softClient)
	ctx.SetGrantToken("clientGrantKey", auth.GrantToken{})
	step := NewClientDelete(server.URL, "clientKey", "clientGrantKey", server.Client())

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Software client delete", result.Name)
	assert.Equal(t, "", result.FailReason)
}

func TestNewClientDelete_Expects204(t *testing.T) {
	// creating a stub server that expects a JWT body posted
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s", clientID), r.URL.EscapedPath())
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	softClient := client.NewClientSecretBasic(clientID, registrationAccessToken, server.URL, clientSecret)
	ctx := NewContext()
	ctx.SetClient("clientKey", softClient)
	ctx.SetGrantToken("clientGrantKey", auth.GrantToken{})
	step := NewClientDelete(server.URL, "clientKey", "clientGrantKey", server.Client())

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Contains(t, result.FailReason, "unexpected status code 200, should be 204")
}

func TestNewClientDelete_HandlesCreateRequestError(t *testing.T) {
	softClient := client.NewClientSecretBasic(clientID, registrationAccessToken, clientSecret, "")
	ctx := NewContext()
	ctx.SetClient("clientKey", softClient)
	ctx.SetGrantToken("clientGrantKey", auth.GrantToken{})
	step := NewClientDelete(string(rune(0x7f)), "clientKey", "clientGrantKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to create request \u007f/foo: parse \"\\u007f/foo\": net/url: invalid control character in URL",
		result.FailReason,
	)
}

func TestNewClientDelete_HandlesExecuteRequestError(t *testing.T) {
	softClient := client.NewClientSecretBasic(clientID, registrationAccessToken, clientSecret, "")
	ctx := NewContext()
	ctx.SetClient("clientKey", softClient)
	ctx.SetGrantToken("clientGrantKey", auth.GrantToken{})
	step := NewClientDelete("localhost", "clientKey", "clientGrantKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to call endpoint localhost/foo: Delete \"localhost/foo\": unsupported protocol scheme \"\"",
		result.FailReason,
	)
}

func TestNewClientDelete_HandlesErrorForClientNotFound(t *testing.T) {
	ctx := NewContext()
	ctx.SetGrantToken("clientGrantKey", auth.GrantToken{})
	step := NewClientDelete("localhost", "clientKey", "clientGrantKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to find client clientKey in context: key not found in context",
		result.FailReason,
	)
}
