package intasend

import (
	"context"
	"fmt"
	"time"
)

// PaymentLinkService handles payment link operations.
type PaymentLinkService struct {
	client *Client
}

// Tariff represents who pays the transaction fees.
type Tariff string

const (
	// TariffBusinessPays means the business pays the fees.
	TariffBusinessPays Tariff = "BUSINESS-PAYS"

	// TariffCustomerPays means the customer pays the fees.
	TariffCustomerPays Tariff = "CUSTOMER-PAYS"
)

// PaymentLink represents a payment link.
type PaymentLink struct {
	LinkID       string    `json:"link_id"`
	Title        string    `json:"title"`
	Currency     string    `json:"currency"`
	Amount       float64   `json:"amount"`
	URL          string    `json:"url"`
	MobileTariff Tariff    `json:"mobile_tarrif"`
	CardTariff   Tariff    `json:"card_tarrif"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PaymentLinkListResponse represents the response from listing payment links.
type PaymentLinkListResponse struct {
	Results []PaymentLink `json:"results"`
}

// CreatePaymentLinkRequest represents a request to create a payment link.
type CreatePaymentLinkRequest struct {
	Title        string  `json:"title"`
	Currency     string  `json:"currency"`
	Amount       float64 `json:"amount,omitempty"`
	MobileTariff Tariff  `json:"mobile_tarrif,omitempty"`
	CardTariff   Tariff  `json:"card_tarrif,omitempty"`
	IsActive     bool    `json:"is_active"`
}

// List returns all payment links.
//
// Example:
//
//	links, err := client.PaymentLink().List(ctx)
func (s *PaymentLinkService) List(ctx context.Context) (*PaymentLinkListResponse, error) {
	var resp PaymentLinkListResponse
	if err := s.client.get(ctx, "/paymentlinks/", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new payment link.
//
// Example:
//
//	link, err := client.PaymentLink().Create(ctx, &intasend.CreatePaymentLinkRequest{
//	    Title:        "Premium Service",
//	    Currency:     "KES",
//	    Amount:       5000,
//	    MobileTariff: intasend.TariffBusinessPays,
//	    CardTariff:   intasend.TariffBusinessPays,
//	    IsActive:     true,
//	})
func (s *PaymentLinkService) Create(ctx context.Context, req *CreatePaymentLinkRequest) (*PaymentLink, error) {
	var resp PaymentLink
	if err := s.client.post(ctx, "/paymentlinks/", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a specific payment link by ID.
//
// Example:
//
//	link, err := client.PaymentLink().Get(ctx, "LINK-123")
func (s *PaymentLinkService) Get(ctx context.Context, linkID string) (*PaymentLink, error) {
	var resp PaymentLink
	if err := s.client.get(ctx, fmt.Sprintf("/paymentlinks/%s/", linkID), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
