package vgs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Environment string

const (
	Sandbox Environment = "sandbox"
	Live    Environment = "live"
	LiveEU1 Environment = "live-eu-1"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Options struct {
	ClientID     string
	ClientSecret string
	VaultId      string
	RouteId      string
	Environment  Environment

	VaultURL      *url.URL
	PaymentURL    *url.URL
	HTTPClient    HTTPClient
	Authenticator Authenticator
}

func (o *Options) GetVaultUrl() (*url.URL, error) {
	if o.VaultId == "" {
		return nil, errors.New("VaultId is required")
	}
	u, err := url.Parse(fmt.Sprintf("https://%s.%s.verygoodproxy.com", o.VaultId, strings.ToLower(string(o.Environment))))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (o *Options) GetPaymentUrl() (*url.URL, error) {
	if o.VaultId == "" {
		return nil, errors.New("VaultId is required")
	}
	if o.RouteId == "" {
		return nil, errors.New("RouteId is required")
	}
	u, err := url.Parse(fmt.Sprintf("https://%s-%s.%s.verygoodproxy.com", o.VaultId, o.RouteId, strings.ToLower(string(o.Environment))))
	if err != nil {
		return nil, err
	}
	return u, nil
}

type Client struct {
	Options *Options
	Ctx     context.Context

	lastRequest  *http.Request
	lastResponse *http.Response
}

func NewClient(options *Options) (*Client, error) {
	return NewClientWithContext(context.Background(), options)
}

func NewClientWithContext(ctx context.Context, options *Options) (*Client, error) {
	if options.ClientID == "" || options.ClientSecret == "" {
		return nil, errors.New("client id and client secret is required")
	}
	if options.VaultId == "" {
		return nil, errors.New("vault id is required")
	}
	if options.RouteId == "" {
		return nil, errors.New("route id is required")
	}
	if options.HTTPClient == nil {
		options.HTTPClient = http.DefaultClient
	}
	if options.Environment == "" {
		options.Environment = Sandbox
	}
	if options.Authenticator == nil {
		options.Authenticator = NewOAuthAuthenticator(
			options.ClientID, options.ClientSecret,
		)
	}

	client := &Client{
		Options: options,
		Ctx:     ctx,
	}
	return client, nil
}

func (c *Client) NewRequest(request *Request) (*http.Request, error) {
	fullUrl, err := c.Options.GetPaymentUrl()
	if err != nil {
		log.Printf("Unable to parse required uri: %s", request.Uri)
		return nil, err
	}
	var buf io.Reader
	if request.hasBody() {
		buf, _ = request.jsonBody()
	}
	req, err := http.NewRequest(request.Method, fullUrl.String(), buf)
	if err != nil {
		log.Printf("Unable to create a new request: %s", err)
		return nil, err
	}
	// set request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// set authentication headers
	if err := c.Options.Authenticator.SetAuthentication(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(c.Ctx)
	c.lastRequest = req
	resp, err := c.Options.HTTPClient.Do(req)
	if err != nil {
		select {
		case <-c.Ctx.Done():
			return nil, c.Ctx.Err()
		default:
		}
		return nil, err
	}
	c.lastResponse = resp

	response := NewResponse(resp)
	defer resp.Body.Close()

	err = ValidateResponse(response)
	if err != nil {
		return response, err
	}

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, response.Body)
	default:
		if len(response.RawBody) > 0 {
			if err = json.Unmarshal(response.RawBody, &v); err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

func (c *Client) Get(uri string, v interface{}, options ...RequestOption) (*Response, error) {
	req, err := c.NewRequest(&Request{
		Method: http.MethodGet,
		Uri:    uri,
	})
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req, v)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Post(uri string, payload, v interface{}, options ...RequestOption) (*Response, error) {
	req, err := c.NewRequest(&Request{
		Method: http.MethodPost,
		Uri:    uri,
		Body:   payload,
	})
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req, v)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
