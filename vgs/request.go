package vgs

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
)

type Request struct {
	Method string      `json:"method"`
	Uri    string      `json:"uri"`
	Body   interface{} `json:"body"`
	Values url.Values  `json:"data"`
}

type RequestOption func(*Request)

func NewJsonRequest(method, uri string, body interface{}, options ...RequestOption) *Request {
	request := &Request{
		Method: method,
		Uri:    uri,
		Body:   body,
	}
	for _, opt := range options {
		opt(request)
	}
	return request
}

func (r *Request) jsonBody() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(r.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (r *Request) hasBody() bool {
	return r.Body != ""
}

func (r *Request) BuildURL(baseUrl *url.URL) (*url.URL, error) {
	u, err := baseUrl.Parse(r.Uri)
	if err != nil {
		log.Printf("Unable to parse required uri: %s", r.Uri)
		return nil, err
	}
	return u, nil
}
