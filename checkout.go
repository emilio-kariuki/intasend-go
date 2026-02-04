package intasend

import (
	"context"
)

// CheckoutService handles checkout operations.
type CheckoutService struct {
	client *Client
}

// CheckoutCustomer represents customer information for checkout.
type CheckoutCustomer struct {
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Country     string
	City        string
	Address     string
	State       string
	Zipcode     string
}

// CreateCheckoutRequest represents a request to create a checkout session.
type CreateCheckoutRequest struct {
	Amount       float64
	Currency     string
	Customer     CheckoutCustomer
	Host         string
	RedirectURL  string
	APIRef       string
	Comment      string
	Method       string
	CardTariff   string
	MobileTariff string
	WalletID     string
}

// createCheckoutBody is the internal request body.
type createCheckoutBody struct {
	PublicKey    string  `json:"public_key,omitempty"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Email        string  `json:"email"`
	FirstName    string  `json:"first_name,omitempty"`
	LastName     string  `json:"last_name,omitempty"`
	PhoneNumber  string  `json:"phone_number,omitempty"`
	Country      string  `json:"country,omitempty"`
	Address      string  `json:"address,omitempty"`
	City         string  `json:"city,omitempty"`
	State        string  `json:"state,omitempty"`
	Zipcode      string  `json:"zipcode,omitempty"`
	Host         string  `json:"host"`
	RedirectURL  string  `json:"redirect_url,omitempty"`
	APIRef       string  `json:"api_ref,omitempty"`
	Comment      string  `json:"comment,omitempty"`
	Method       string  `json:"method,omitempty"`
	CardTariff   string  `json:"card_tarrif,omitempty"`
	MobileTariff string  `json:"mobile_tarrif,omitempty"`
	WalletID     string  `json:"wallet_id,omitempty"`
}

// CreateCheckoutResponse represents the response from creating a checkout.
type CreateCheckoutResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Signature string `json:"signature"`
}

// CheckoutStatusRequest represents a request to check checkout status.
type CheckoutStatusRequest struct {
	Signature  string `json:"signature"`
	CheckoutID string `json:"checkout_id"`
	InvoiceID  string `json:"invoice_id"`
}

// CheckoutStatusResponse represents a checkout status response.
type CheckoutStatusResponse struct {
	Invoice  *Invoice      `json:"invoice"`
	Customer *CustomerInfo `json:"customer,omitempty"`
}

// Create creates a new checkout session.
//
// Example:
//
//	session, err := client.Checkout().Create(ctx, &intasend.CreateCheckoutRequest{
//	    Amount:   1000,
//	    Currency: "KES",
//	    Customer: intasend.CheckoutCustomer{
//	        Email:     "john@example.com",
//	        FirstName: "John",
//	        LastName:  "Doe",
//	        Country:   "KE",
//	    },
//	    Host:        "https://yoursite.com",
//	    RedirectURL: "https://yoursite.com/callback",
//	    APIRef:      "order-123",
//	})
func (s *CheckoutService) Create(ctx context.Context, req *CreateCheckoutRequest) (*CreateCheckoutResponse, error) {
	body := &createCheckoutBody{
		PublicKey:    s.client.publishableKey,
		Amount:       req.Amount,
		Currency:     req.Currency,
		Email:        req.Customer.Email,
		FirstName:    req.Customer.FirstName,
		LastName:     req.Customer.LastName,
		PhoneNumber:  req.Customer.PhoneNumber,
		Country:      req.Customer.Country,
		Address:      req.Customer.Address,
		City:         req.Customer.City,
		State:        req.Customer.State,
		Zipcode:      req.Customer.Zipcode,
		Host:         req.Host,
		RedirectURL:  req.RedirectURL,
		APIRef:       req.APIRef,
		Comment:      req.Comment,
		Method:       req.Method,
		CardTariff:   req.CardTariff,
		MobileTariff: req.MobileTariff,
		WalletID:     req.WalletID,
	}

	var resp CreateCheckoutResponse
	if err := s.client.postPublic(ctx, "/checkout/", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CheckStatus checks the status of a checkout session.
//
// Example:
//
//	status, err := client.Checkout().CheckStatus(ctx, &intasend.CheckoutStatusRequest{
//	    Signature:  "xxx",
//	    CheckoutID: "CHK-123",
//	    InvoiceID:  "INV-456",
//	})
func (s *CheckoutService) CheckStatus(ctx context.Context, req *CheckoutStatusRequest) (*CheckoutStatusResponse, error) {
	var resp CheckoutStatusResponse
	if err := s.client.postPublic(ctx, "/payment/status/", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
