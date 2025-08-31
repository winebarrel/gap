package gap_test

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/gap"
	"google.golang.org/api/option"
)

func TestCallTokeninfo(t *testing.T) {
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

	client, err := gap.NewOauth2Client(option.WithHTTPClient(hc))
	require.NoError(err)
	ti, err := client.Tokeninfo("my-token")
	require.NoError(err)
	assert.Equal("scott@example.com", ti.Email)
}

func TestCallTokeninfo_InvalidRequest(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	hc := &http.Client{}
	httpmock.ActivateNonDefault(hc)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, "https://www.googleapis.com/oauth2/v2/tokeninfo", func(req *http.Request) (*http.Response, error) {
		token := req.URL.Query().Get("access_token")
		assert.Equal("my-token", token)
		body := `{"error":"invalid_request","error_description": "Invalid Credentials"}`
		return httpmock.NewStringResponse(http.StatusUnauthorized, body), nil
	})

	client, err := gap.NewOauth2Client(option.WithHTTPClient(hc))
	require.NoError(err)
	_, err = client.Tokeninfo("my-token")
	require.ErrorContains(err, `{"error":"invalid_request","error_description": "Invalid Credentials"}`)
}
