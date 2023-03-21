package vgs

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFinancialInstruments(t *testing.T) {
	t.Parallel()
	c := NewMockClientWithHandler(newMockHandler(http.StatusOK, `{"meta": {}, "links": {}, "data": [{"id": "dummy"}]}`, nil))
	gateways, err := c.GetFinancialInstruments()
	assert.Nil(t, err)
	assert.NotNil(t, gateways)
	assert.NotNil(t, gateways.Data)
	assert.Len(t, gateways.Data, 1)
	assert.Equal(t, gateways.Data[0].ID, "dummy")
	assert.NotNil(t, gateways.Meta)
	assert.NotNil(t, gateways.Links)
}
func TestGetFinancialInstrumentsBadRequest(t *testing.T) {
	t.Parallel()
	c := NewMockClientWithHandler(newMockHandler(http.StatusBadRequest, `{"error": "error", "error_description": "dummy error"}`, nil))
	gateways, err := c.GetFinancialInstruments()
	assert.NotNil(t, err)
	assert.Nil(t, gateways)
	assert.ErrorContains(t, err, "dummy")
}
