package vgs

import "time"

type Gateway struct {
	Type_ string `json:"type"`
	// Unique identifier for this gateway.  Used to refer to gateway in rules.
	Id string `json:"id"`
	// ISO 4217 currency code. Defaults to USD
	DefaultCurrency string `json:"default_currency"`
	// Is this gateway the default gateway or not.  A default gateway is needed and will be used when a transfer without matching any routing rule created. There could be only one default gateway at the same time. When a new gateway created as the default gateway, the old default gateway will no longer be the default gateway anymore.
	DefaultGateway bool `json:"default_gateway,omitempty"`
	// Any specific keys passed through to the gateway configuration. Refer to docs
	Config *interface{} `json:"config"`
	// Creation time, in UTC.
	CreatedAt time.Time `json:"created_at"`
	// Last time psp token was updated, in UTC.
	UpdatedAt time.Time `json:"updated_at"`
}

type Gateways struct {
	Response
	Data []Gateway `json:"data"`
}

func (c *Client) GetGateways() (*Gateways, error) {
	gateways := &Gateways{}
	_, err := c.Get("/gateways", gateways)
	if err != nil {
		return nil, err
	}
	return gateways, nil
}
