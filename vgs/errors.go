package vgs

import (
	"fmt"
)

type VGSError struct {
	ErrorCode        string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func (e VGSError) Error() string {
	return fmt.Sprintf("API call error (%s): %s", e.ErrorCode, e.ErrorDescription)
}
