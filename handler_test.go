package gap_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/gap"
	"google.golang.org/api/option"
)

func TestHandlePing(t *testing.T) {
	assert := assert.New(t)
	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	w := httptest.NewRecorder()
	gap.HandlePing(w, r)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("pong", w.Body.String())
}

func TestAuthHandler(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	hc := &http.Client{}
	httpmock.ActivateNonDefault(hc)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, "https://www.googleapis.com/oauth2/v2/tokeninfo", func(req *http.Request) (*http.Response, error) {
		token := req.URL.Query().Get("access_token")
		assert.Equal("my-token", token)
		body := `{"email":"scott@example.com"}`
		return httpmock.NewStringResponse(http.StatusOK, body), nil
	})

	oc, err := gap.NewOauth2Client(option.WithHTTPClient(hc))
	require.NoError(err)

	handler := gap.AuthHandler{
		Options: &gap.Options{
			HeaderName: "my-header",
			AllowList:  []string{"scott@example.com"},
		},
		Oauth2: oc,
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Header.Add("my-header", "my-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("proxied", w.Body.String())
}

func TestAuthHandlerWithMultiEmail(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	hc := &http.Client{}
	httpmock.ActivateNonDefault(hc)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, "https://www.googleapis.com/oauth2/v2/tokeninfo", func(req *http.Request) (*http.Response, error) {
		token := req.URL.Query().Get("access_token")
		assert.Equal("my-token", token)
		body := `{"email":"scott@example.com"}`
		return httpmock.NewStringResponse(http.StatusOK, body), nil
	})

	oc, err := gap.NewOauth2Client(option.WithHTTPClient(hc))
	require.NoError(err)

	handler := gap.AuthHandler{
		Options: &gap.Options{
			HeaderName: "my-header",
			AllowList:  []string{"scott@example.com", "tiger@example.com"},
		},
		Oauth2: oc,
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Header.Add("my-header", "my-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("proxied", w.Body.String())
}

func TestAuthHandlerWithWildcard(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	hc := &http.Client{}
	httpmock.ActivateNonDefault(hc)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, "https://www.googleapis.com/oauth2/v2/tokeninfo", func(req *http.Request) (*http.Response, error) {
		token := req.URL.Query().Get("access_token")
		assert.Equal("my-token", token)
		body := `{"email":"scott@example.com"}`
		return httpmock.NewStringResponse(http.StatusOK, body), nil
	})

	oc, err := gap.NewOauth2Client(option.WithHTTPClient(hc))
	require.NoError(err)

	handler := gap.AuthHandler{
		Options: &gap.Options{
			HeaderName: "my-header",
			AllowList:  []string{"*@example.com"},
		},
		Oauth2: oc,
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Header.Add("my-header", "my-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("proxied", w.Body.String())
}

func TestAuthHandler_InvalidRequest(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	hc := &http.Client{}
	httpmock.ActivateNonDefault(hc)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, "https://www.googleapis.com/oauth2/v2/tokeninfo", func(req *http.Request) (*http.Response, error) {
		token := req.URL.Query().Get("access_token")
		assert.Equal("my-token", token)
		body := `{"error":"invalid_request","error_description": "Invalid Credentials"}`
		return httpmock.NewStringResponse(http.StatusForbidden, body), nil
	})

	oc, err := gap.NewOauth2Client(option.WithHTTPClient(hc))
	require.NoError(err)

	handler := gap.AuthHandler{
		Options: &gap.Options{
			HeaderName: "my-header",
			AllowList:  []string{"*@example.com"},
		},
		Oauth2: oc,
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Header.Add("my-header", "my-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	assert.Equal(http.StatusForbidden, w.Code)
	assert.Equal("forbidden\n", w.Body.String())
}

func TestAuthHandler_EmptyToken(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	hc := &http.Client{}
	httpmock.ActivateNonDefault(hc)
	defer httpmock.DeactivateAndReset()

	oc, err := gap.NewOauth2Client(option.WithHTTPClient(hc))
	require.NoError(err)

	handler := gap.AuthHandler{
		Options: &gap.Options{
			HeaderName: "my-header",
			AllowList:  []string{"*@example.com"},
		},
		Oauth2: oc,
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	assert.Equal(http.StatusForbidden, w.Code)
	assert.Equal("forbidden\n", w.Body.String())
}

func TestAuthHandler_EmptyEmail(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	hc := &http.Client{}
	httpmock.ActivateNonDefault(hc)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, "https://www.googleapis.com/oauth2/v2/tokeninfo", func(req *http.Request) (*http.Response, error) {
		token := req.URL.Query().Get("access_token")
		assert.Equal("my-token", token)
		body := `{}`
		return httpmock.NewStringResponse(http.StatusOK, body), nil
	})

	oc, err := gap.NewOauth2Client(option.WithHTTPClient(hc))
	require.NoError(err)

	handler := gap.AuthHandler{
		Options: &gap.Options{
			HeaderName: "my-header",
			AllowList:  []string{"*@example.com"},
		},
		Oauth2: oc,
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Header.Add("my-header", "my-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	assert.Equal(http.StatusForbidden, w.Code)
	assert.Equal("cannot get email\n", w.Body.String())
}

func TestAuthHandler_Disallowed(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	hc := &http.Client{}
	httpmock.ActivateNonDefault(hc)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, "https://www.googleapis.com/oauth2/v2/tokeninfo", func(req *http.Request) (*http.Response, error) {
		token := req.URL.Query().Get("access_token")
		assert.Equal("my-token", token)
		body := `{"email":"scott@example.com"}`
		return httpmock.NewStringResponse(http.StatusOK, body), nil
	})

	oc, err := gap.NewOauth2Client(option.WithHTTPClient(hc))
	require.NoError(err)

	handler := gap.AuthHandler{
		Options: &gap.Options{
			HeaderName: "my-header",
			AllowList:  []string{"x@example.com"},
		},
		Oauth2: oc,
		Proxy: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "proxied")
		},
	}

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Header.Add("my-header", "my-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	assert.Equal(http.StatusForbidden, w.Code)
	assert.Equal("not allowed\n", w.Body.String())
}
