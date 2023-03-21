package vgs

import "time"

type GatewayOptions struct {
	Currency        string          `json:"currency,omitempty"`
	ShippingAddress *ContactAddress `json:"shipping_address,omitempty"`
}

type VerificationsRequest struct {
	Card           *Card           `json:"card,omitempty"`
	GatewayOptions *GatewayOptions `json:"gateway_options,omitempty"`
}

type GatewayInfo struct {
	Type            string      `json:"type,omitempty"`
	ID              string      `json:"id,omitempty"`
	DefaultCurrency string      `json:"default_currency,omitempty"`
	DefaultGateway  bool        `json:"default_gateway,omitempty"`
	Config          interface{} `json:"config,omitempty"`
	CreatedAt       time.Time   `json:"created_at,omitempty"`
	UpdatedAt       time.Time   `json:"updated_at,omitempty"`
}

type GatewayResponse struct {
	ID          string `json:"id,omitempty"`
	Message     string `json:"message,omitempty"`
	State       string `json:"state,omitempty"`
	ErrorCode   string `json:"error_code,omitempty"`
	RawResponse string `json:"raw_response,omitempty"`
}

type AvsResult struct {
	Code        string `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	StreetMatch string `json:"street_match,omitempty"`
	PostalMatch string `json:"postal_match,omitempty"`
}

type Verficiation struct {
	ID              string           `json:"id,omitempty"`
	CreatedAt       time.Time        `json:"created_at,omitempty"`
	UpdatedAt       time.Time        `json:"updated_at,omitempty"`
	Type            string           `json:"type,omitempty"`
	Amount          int              `json:"amount,omitempty"`
	Fee             int              `json:"fee,omitempty"`
	Currency        string           `json:"currency,omitempty"`
	Gateway         *GatewayInfo     `json:"gateway,omitempty"`
	GatewayResponse *GatewayResponse `json:"gateway_response,omitempty"`
	Source          string           `json:"source,omitempty"`
	Destination     string           `json:"destination,omitempty"`
	State           string           `json:"state,omitempty"`
	AvsResult       *AvsResult       `json:"avs_result,omitempty"`
	SubAccountID    string           `json:"sub_account_id,omitempty"`
}

type VerficiationsResponse struct {
	Data Verficiation `json:"data,omitempty"`
}

func (c *Client) CreateVerifications(body *VerificationsRequest) (*VerficiationsResponse, error) {
	resp := &VerficiationsResponse{}
	_, err := c.Post("/verfications", body, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
