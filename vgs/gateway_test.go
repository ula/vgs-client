package vgs

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGateways(t *testing.T) {
	t.Parallel()
	c := NewMockClientWithHandler(newMockHandler(http.StatusOK, `{"meta": {}, "links": {}}`, nil))
	gateways, err := c.GetGateways()
	assert.Nil(t, err)
	assert.NotNil(t, gateways)
	assert.Empty(t, gateways.Data)
	assert.NotNil(t, gateways.Meta)
	assert.NotNil(t, gateways.Links)
}
func TestGetGatewaysBadRequest(t *testing.T) {
	t.Parallel()
	c := NewMockClientWithHandler(newMockHandler(http.StatusBadRequest, `{"error": "error", "error_description": "dummy error"}`, nil))
	gateways, err := c.GetGateways()
	assert.NotNil(t, err)
	assert.Nil(t, gateways)
	assert.ErrorContains(t, err, "dummy")
}
