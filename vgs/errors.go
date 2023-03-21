package vgs

import (
	"fmt"
	"strings"
)

type ErrorDetail struct {
	Code   string `json:"code,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type VGSError struct {
	Errors           []ErrorDetail `json:"errors,omitempty"`
	ErrorCode        string        `json:"error,omitempty"`
	ErrorDescription string        `json:"error_description,omitempty"`
}

func (e VGSError) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("API call error (%s): %s", e.ErrorCode, e.ErrorDescription)
	}
	msgs := []string{}
	for i := 0; i < len(e.Errors); i++ {
		msgs = append(msgs, fmt.Sprintf("#%v. Code: %s; Details: %s", i+1, e.Errors[i].Code, e.Errors[i].Detail))
	}
	return fmt.Sprintf("Response errors:\n %v", strings.Join(msgs, "\n"))
}
