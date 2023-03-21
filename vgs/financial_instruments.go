package vgs

import "time"

type ContactAddress struct {
	Name       string `json:"name,omitempty"`
	Company    string `json:"company,omitempty"`
	Address1   string `json:"address1,omitempty"`
	Address2   string `json:"address2,omitempty"`
	City       string `json:"city,omitempty"`
	Region     string `json:"region,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Phone      string `json:"phone,omitempty"`
}
type Card struct {
	Name           string          `json:"name,omitempty"`
	ExpMonth       int             `json:"exp_month,omitempty"`
	ExpYear        int             `json:"exp_year,omitempty"`
	BillingAddress *ContactAddress `json:"billing_address,omitempty"`
	Number         string          `json:"number,omitempty"`
	Brand          string          `json:"brand,omitempty"`
	Last4          string          `json:"last4,omitempty"`
	Cvc            string          `json:"cvc,omitempty"`
}
type PspToken struct {
	Id    string `json:"id,omitempty"`
	Value string `json:"value,omitempty"`
	Psp   string `json:"psp,omitempty"`
}

type FinancialInstrumentData struct {
	ID           string    `json:"id,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	SubAccountID string    `json:"sub_account_id,omitempty"`
	Card         Card      `json:"card,omitempty"`
	PspToken     PspToken  `json:"psp_token,omitempty"`
}

type FinancialInstruments struct {
	Response
	Data []FinancialInstrumentData `json:"data,omitempty"`
}

type FinancialInstrument struct {
	Data FinancialInstrumentData `json:"data,omitempty"`
}

func (c *Client) GetFinancialInstruments() (*FinancialInstruments, error) {
	instruments := &FinancialInstruments{}
	_, err := c.Get("/financial_instruments", instruments)
	if err != nil {
		return nil, err
	}
	return instruments, nil
}

func (c *Client) CreateFinancialInstrument(body interface{}) (*FinancialInstrument, error) {
	resp := &FinancialInstrument{}
	_, err := c.Post("/financial_instruments", body, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type CreatePSPTokenRequest struct {
	PspToken *PspToken `json:"psp_token,omitempty"`
}

func (c *Client) CreatePSPToken(psp, id string) (*FinancialInstrument, error) {
	return c.CreateFinancialInstrument(&CreatePSPTokenRequest{PspToken: &PspToken{Id: id, Psp: psp}})
}

type CreatePaymentCardRequest struct {
	Card *Card `json:"card,omitempty"`
}

func (c *Client) CreatePaymentCard(body *CreatePaymentCardRequest) (*FinancialInstrument, error) {
	return c.CreateFinancialInstrument(body)
}
