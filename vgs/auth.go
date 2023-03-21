package vgs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var AuthEndpoint = "https://auth.verygoodsecurity.com/auth/realms/vgs/protocol/openid-connect/token"

type Authenticator interface {
	Authenticate() (*OAuthToken, error)
	SetAuthentication(r *http.Request) error
}

type OAuthToken struct {
	AccessToken      string `json:"access_token,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	RefreshExpiresIn int    `json:"refresh_expires_in,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	NotBeforePolicy  int    `json:"not-before-policy,omitempty"`
	Scope            string `json:"scope,omitempty"`
	CreatedAt        time.Time
}

func (o *OAuthToken) IsValid() bool {
	if o == nil || o.AccessToken == "" {
		return false
	}
	now := time.Now()
	return now.Before(o.CreatedAt.Add(time.Second * time.Duration(o.ExpiresIn)))
}

type OauthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type OAuthAuthenticator struct {
	OAuthURL   string
	Config     *OauthConfig
	HTTPClient HTTPClient
	Token      *OAuthToken
}

func NewOAuthAuthenticator(clientId, clientSecret string) *OAuthAuthenticator {
	return &OAuthAuthenticator{
		OAuthURL: AuthEndpoint,
		Config: &OauthConfig{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			GrantType:    "client_credentials",
		},
		HTTPClient: http.DefaultClient,
	}
}

func (o *OAuthAuthenticator) SetAuthentication(r *http.Request) error {
	token, err := o.Authenticate()
	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("token is not set")
	}
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	return nil
}

func (o *OAuthAuthenticator) Authenticate() (*OAuthToken, error) {
	if o.Token.IsValid() {
		return o.Token, nil
	}
	token, err := o.FetchToken()
	if err != nil {
		return nil, err
	}
	// set token value to authenticator
	o.Token = token

	return token, nil
}

func (o *OAuthAuthenticator) FetchToken() (*OAuthToken, error) {
	data := url.Values{}
	data.Add("grant_type", o.Config.GrantType)
	data.Add("client_id", o.Config.ClientID)
	data.Add("client_secret", o.Config.ClientSecret)

	buf := strings.NewReader(data.Encode())
	req, err := http.NewRequest(http.MethodPost, o.OAuthURL, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := o.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var token OAuthToken
		err = json.NewDecoder(resp.Body).Decode(&token)
		if err != nil {
			return nil, err
		}
		token.CreatedAt = time.Now()
		return &token, nil
	}

	var vgsError VGSError
	err = json.NewDecoder(resp.Body).Decode(&vgsError)
	return nil, err
}
