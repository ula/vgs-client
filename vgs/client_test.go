package vgs

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockAuthenticator struct {
	Token *OAuthToken
}

func (m *MockAuthenticator) Authenticate() (*OAuthToken, error) {
	return &OAuthToken{
		AccessToken: "test-token",
		ExpiresIn:   3600,
		CreatedAt:   time.Now(),
	}, nil
}
func (m *MockAuthenticator) SetAuthentication(r *http.Request) error {
	r.Header.Set("Authentication", "Bearer test")
	return nil
}

type mockHTTPClient struct {
	mockHandler http.HandlerFunc
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(m.mockHandler)
	handler.ServeHTTP(resp, req)
	return resp.Result(), nil
}

type badMockHTTPClient struct {
	mockHandler http.HandlerFunc
}

func (mtc *badMockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("Request failed")
}

func NewMockClient(mockHttpClient HTTPClient) (*Client, error) {
	options := &Options{
		ClientID:      "test-client",
		ClientSecret:  "test-secret",
		VaultId:       "test-vault",
		RouteId:       "test-route",
		HTTPClient:    mockHttpClient,
		Authenticator: &MockAuthenticator{},
	}
	return NewClient(options)
}

func NewMockClientWithHandler(mockHandler http.HandlerFunc) *Client {
	c, err := NewMockClient(&mockHTTPClient{mockHandler: mockHandler})
	if err != nil {
		log.Fatalf("cannot init client")
	}
	return c
}

func newMockHandler(statusCode int, json string, headers map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(headers) > 0 {
			for key, value := range headers {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(statusCode)
		w.Write([]byte(json))
	}
}

func TestNewClient(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		expectedError bool
		options       *Options
	}{
		{
			true,
			&Options{},
		},
		{
			expectedError: true,
			options:       &Options{ClientID: "", ClientSecret: ""},
		},
		{
			expectedError: true,
			options:       &Options{ClientID: "test", ClientSecret: "test"},
		},
		{
			expectedError: false,
			options:       &Options{ClientID: "test", ClientSecret: "test", VaultId: "test", RouteId: "test"},
		},
	}

	for _, testCase := range testCases {
		c, err := NewClient(testCase.options)

		if testCase.expectedError {
			assert.NotNil(t, err)
			assert.Error(t, err)
			assert.Nil(t, c)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, c)
			assert.NotNil(t, c.Ctx)
			assert.Equal(t, c.Options.HTTPClient, http.DefaultClient)
		}
	}
}

func TestNewClientContextExists(t *testing.T) {
	t.Parallel()

	mockHTTPClient := &mockHTTPClient{mockHandler: func(w http.ResponseWriter, r *http.Request) {
		assert.NotNil(t, r.Context())
	}}

	c, err := NewMockClient(mockHTTPClient)

	assert.Nil(t, err)
	assert.NotNil(t, c)

}

func TestNewClientContextEvent(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wasCanceled := false
	mockHTTPClient := &mockHTTPClient{func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
		wasCanceled = true
	}}

	c, err := NewClientWithContext(
		ctx,
		&Options{
			ClientID:      "test",
			ClientSecret:  "test",
			VaultId:       "test",
			RouteId:       "test",
			HTTPClient:    mockHTTPClient,
			Authenticator: &MockAuthenticator{},
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	cancel()
	account, err := c.GetGateways()
	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.True(t, wasCanceled)
}

func TestHTTPClientFailedDo(t *testing.T) {
	t.Parallel()
	badClient := &badMockHTTPClient{
		mockHandler: newMockHandler(0, "", nil),
	}
	c, err := NewMockClient(badClient)
	assert.Nil(t, err)

	account, err := c.GetGateways()
	assert.NotNil(t, err)
	assert.Error(t, err, "Request failed")
	assert.Nil(t, account)
}

func TestHTTPClientBadJSON(t *testing.T) {
	t.Parallel()
	c := NewMockClientWithHandler(newMockHandler(http.StatusOK, `{"data": [], "meta": }`, nil))
	account, err := c.GetGateways()
	assert.NotNil(t, err)
	assert.Nil(t, account)
	assert.ErrorContains(t, err, "invalid character")
}
