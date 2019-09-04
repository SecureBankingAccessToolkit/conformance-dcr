package step

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/client"
	"bitbucket.org/openbankingteam/conformance-dcr/pkg/compliant/openid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	clientID     = "foo"
	clientSecret = "bar"
)

func TestNewClientDelete(t *testing.T) {
	softClient := client.NewClient(clientID, clientSecret)
	// creating a stub server that expects a JWT body posted
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s", clientID), r.URL.EscapedPath())
		_, err := w.Write([]byte(`OK`))
		require.NoError(t, err)
	}))
	defer server.Close()

	ctx := NewContext()
	ctx.SetClient("clientKey", softClient)
	step := NewClientDelete("responseCtxKey", server.URL, "clientKey", server.Client())

	result := step.Run(ctx)

	assert.True(t, result.Pass)
	assert.Equal(t, "Software client delete", result.Name)
	assert.Equal(t, "", result.FailReason)

	// assert that response in now in ctx
	_, err := ctx.GetResponse("responseCtxKey")
	assert.NoError(t, err)
}

func TestNewClientDelete_HandlesCreateRequestError(t *testing.T) {
	softClient := client.NewClient(clientID, clientSecret)
	ctx := NewContext()
	ctx.SetClient("clientKey", softClient)
	step := NewClientDelete("responseCtxKey", string(0x7f), "clientKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to create request \u007f/foo: parse \u007f/foo: net/url: invalid control character in URL",
		result.FailReason,
	)
}

func TestNewClientDelete_HandlesExecuteRequestError(t *testing.T) {
	softClient := client.NewClient(clientID, clientSecret)
	ctx := NewContext()
	ctx.SetClient("clientKey", softClient)
	step := NewClientDelete("responseCtxKey", "localhost", "clientKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to call endpoint localhost/foo: Delete localhost/foo: unsupported protocol scheme \"\"",
		result.FailReason,
	)
}

func TestNewClientDelete_HandlesErrorForClientNotFound(t *testing.T) {
	ctx := NewContext()
	registrationEndpoint := string(0x7f)
	ctx.SetOpenIdConfig("openIdConfigCtxKey", openid.Configuration{
		RegistrationEndpoint: &registrationEndpoint,
		TokenEndpoint:        "",
	})
	step := NewClientDelete("responseCtxKey", "openIdConfigCtxKey", "clientKey", &http.Client{})

	result := step.Run(ctx)

	assert.False(t, result.Pass)
	assert.Equal(
		t,
		"unable to find client clientKey in context: key not found in context",
		result.FailReason,
	)
}
