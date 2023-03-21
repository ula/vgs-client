package vgs

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var (
	VGSRequestId = "VGS-Request-Id"
	TraceId      = "Trace-Id"
)

type ResponseLinks struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Self  string `json:"self,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}
type ResponseMeta struct {
	TotalElements int `json:"total_elements,omitempty"`
	TotalPages    int `json:"total_pages,omitempty"`
}

type Response struct {
	*http.Response

	RawBody []byte

	Data  []interface{}
	Links ResponseLinks
	Meta  ResponseMeta

	VGSRequestId string
	TraceId      string
}

func NewResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	response.parseHeaders()
	response.readBody()
	return response
}

func (r *Response) parseHeaders() {
	if requestId := r.Header.Get(VGSRequestId); requestId != "" {
		r.VGSRequestId = requestId
	}
	if traceId := r.Header.Get(TraceId); traceId != "" {
		r.TraceId = traceId
	}
}

func (r *Response) readBody() {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Unable to read response body: %v\n", err)
		return
	}
	r.RawBody = body
}

func ValidateResponse(r *Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}
	var vgsError VGSError
	if err := json.Unmarshal(r.RawBody, &vgsError); err != nil {
		return err
	}
	return vgsError
}
