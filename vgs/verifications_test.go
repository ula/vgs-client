package vgs

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateVerifications(t *testing.T) {
	t.Parallel()
	c := NewMockClientWithHandler(newMockHandler(http.StatusOK, `{"data": {}}`, nil))
	verification, err := c.CreateVerifications(&VerificationsRequest{
		Card:           &Card{},
		GatewayOptions: &GatewayOptions{},
	})
	assert.Nil(t, err)
	assert.NotNil(t, verification)
	assert.NotNil(t, verification.Data)
}
