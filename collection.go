package intasend

import (
	"context"
	"time"
)

// CollectionService handles payment collection operations.
type CollectionService struct {
	client *Client
}

// ChargeRequest represents a request to create a checkout page.
type ChargeRequest struct {
	// FirstName is the customer's first name.
	FirstName string `json:"first_name,omitempty"`

	// LastName is the customer's last name.
	LastName string `json:"last_name,omitempty"`

	// Email is the customer's email address.
	Email string `json:"email"`

	// PhoneNumber is the customer's phone number.
	PhoneNumber string `json:"phone_number,omitempty"`

	// Host is your website's base URL for CORS.
	Host string `json:"host"`

	// Amount is the payment amount.
	Amount float64 `json:"amount"`

	// Currency is the payment currency (e.g., "KES", "USD").
	Currency string `json:"currency"`

	// APIRef is your unique reference for this transaction.
	APIRef string `json:"api_ref,omitempty"`

	// RedirectURL is the URL to redirect to after payment.
	RedirectURL string `json:"redirect_url,omitempty"`

	// Comment is an optional payment comment/description.
	Comment string `json:"comment,omitempty"`

	// Method limits the payment to a specific method.
	Method string `json:"method,omitempty"`

	// WalletID directs the payment to a specific wallet.
	WalletID string `json:"wallet_id,omitempty"`

	// CardTariff specifies who pays card fees ("BUSINESS-PAYS" or "CUSTOMER-PAYS").
	CardTariff string `json:"card_tarrif,omitempty"`

	// MobileTariff specifies who pays mobile money fees.
	MobileTariff string `json:"mobile_tarrif,omitempty"`

	// Customer address fields
	Country string `json:"country,omitempty"`
	Address string `json:"address,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zipcode string `json:"zipcode,omitempty"`
}

// chargeRequestBody is the internal request body with public_key.
type chargeRequestBody struct {
	PublicKey    string  `json:"public_key,omitempty"`
	FirstName    string  `json:"first_name,omitempty"`
	LastName     string  `json:"last_name,omitempty"`
	Email        string  `json:"email"`
	PhoneNumber  string  `json:"phone_number,omitempty"`
	Host         string  `json:"host"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	APIRef       string  `json:"api_ref,omitempty"`
	RedirectURL  string  `json:"redirect_url,omitempty"`
	Comment      string  `json:"comment,omitempty"`
	Method       string  `json:"method,omitempty"`
	WalletID     string  `json:"wallet_id,omitempty"`
	CardTariff   string  `json:"card_tarrif,omitempty"`
	MobileTariff string  `json:"mobile_tarrif,omitempty"`
	Country      string  `json:"country,omitempty"`
	Address      string  `json:"address,omitempty"`
	City         string  `json:"city,omitempty"`
	State        string  `json:"state,omitempty"`
	Zipcode      string  `json:"zipcode,omitempty"`
}

// ChargeResponse represents the response from creating a checkout.
type ChargeResponse struct {
	// ID is the checkout session ID.
	ID string `json:"id"`

	// URL is the checkout page URL to redirect the customer to.
	URL string `json:"url"`

	// Signature is used for status verification.
	Signature string `json:"signature"`
}

// STKPushRequest represents an M-Pesa STK Push request.
type STKPushRequest struct {
	// PhoneNumber is the M-Pesa phone number (format: 254XXXXXXXXX).
	PhoneNumber string `json:"phone_number"`

	// Amount is the payment amount in KES.
	Amount float64 `json:"amount"`

	// APIRef is your unique reference for this transaction.
	APIRef string `json:"api_ref,omitempty"`

	// Name is the customer's name.
	Name string `json:"name,omitempty"`

	// Email is the customer's email.
	Email string `json:"email,omitempty"`

	// WalletID directs the payment to a specific wallet.
	WalletID string `json:"wallet_id,omitempty"`
}

// stkPushRequestBody is the internal request body.
type stkPushRequestBody struct {
	PublicKey   string  `json:"public_key,omitempty"`
	PhoneNumber string  `json:"phone_number"`
	Amount      float64 `json:"amount"`
	APIRef      string  `json:"api_ref,omitempty"`
	Name        string  `json:"name,omitempty"`
	Email       string  `json:"email,omitempty"`
	WalletID    string  `json:"wallet_id,omitempty"`
	Method      string  `json:"method"`
	Currency    string  `json:"currency"`
}

// STKPushResponse represents the response from an STK Push request.
type STKPushResponse struct {
	// Invoice contains the invoice details.
	Invoice *Invoice `json:"invoice"`

	// Customer contains customer details.
	Customer *CustomerInfo `json:"customer,omitempty"`
}

// Invoice represents an IntaSend invoice.
type Invoice struct {
	InvoiceID    string    `json:"invoice_id"`
	State        string    `json:"state"`
	Provider     string    `json:"provider"`
	Value        float64   `json:"value"`
	Account      string    `json:"account"`
	APIRef       string    `json:"api_ref"`
	FailedReason string    `json:"failed_reason,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CustomerInfo represents a customer record.
type CustomerInfo struct {
	CustomerID  string `json:"customer_id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

// StatusOptions contains optional parameters for status checks.
type StatusOptions struct {
	CheckoutID string
	Signature  string
}

// statusRequest is the internal request structure for status checks.
type statusRequest struct {
	InvoiceID  string `json:"invoice_id"`
	PublicKey  string `json:"public_key,omitempty"`
	CheckoutID string `json:"checkout_id,omitempty"`
	Signature  string `json:"signature,omitempty"`
}

// StatusResponse represents a payment status response.
type StatusResponse struct {
	Invoice  *Invoice      `json:"invoice"`
	Customer *CustomerInfo `json:"customer,omitempty"`
}

// Payment states
const (
	StateNew        = "NEW"
	StatePending    = "PENDING"
	StateProcessing = "PROCESSING"
	StateComplete   = "COMPLETE"
	StateFailed     = "FAILED"
)

// Charge creates a checkout page for payment collection.
// This method does not require the secret key.
//
// Example:
//
//	resp, err := client.Collection().Charge(ctx, &intasend.ChargeRequest{
//	    FirstName: "John",
//	    LastName:  "Doe",
//	    Email:     "john@example.com",
//	    Host:      "https://yoursite.com",
//	    Amount:    100,
//	    Currency:  "KES",
//	    APIRef:    "order-123",
//	})
func (s *CollectionService) Charge(ctx context.Context, req *ChargeRequest) (*ChargeResponse, error) {
	body := &chargeRequestBody{
		PublicKey:    s.client.publishableKey,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Host:         req.Host,
		Amount:       req.Amount,
		Currency:     req.Currency,
		APIRef:       req.APIRef,
		RedirectURL:  req.RedirectURL,
		Comment:      req.Comment,
		Method:       req.Method,
		WalletID:     req.WalletID,
		CardTariff:   req.CardTariff,
		MobileTariff: req.MobileTariff,
		Country:      req.Country,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		Zipcode:      req.Zipcode,
	}

	var resp ChargeResponse
	if err := s.client.postPublic(ctx, "/checkout/", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MPesaSTKPush initiates an M-Pesa STK Push request.
// This triggers a payment prompt on the customer's phone.
//
// Example:
//
//	resp, err := client.Collection().MPesaSTKPush(ctx, &intasend.STKPushRequest{
//	    PhoneNumber: "254712345678",
//	    Amount:      100,
//	    APIRef:      "order-123",
//	    Name:        "John Doe",
//	    Email:       "john@example.com",
//	})
func (s *CollectionService) MPesaSTKPush(ctx context.Context, req *STKPushRequest) (*STKPushResponse, error) {
	body := &stkPushRequestBody{
		PublicKey:   s.client.publishableKey,
		PhoneNumber: req.PhoneNumber,
		Amount:      req.Amount,
		APIRef:      req.APIRef,
		Name:        req.Name,
		Email:       req.Email,
		WalletID:    req.WalletID,
		Method:      "M-PESA",
		Currency:    "KES",
	}

	var resp STKPushResponse
	if err := s.client.post(ctx, "/payment/mpesa-stk-push/", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Status checks the payment status for an invoice.
// This method does not require the secret key.
//
// Example:
//
//	status, err := client.Collection().Status(ctx, "INV-12345", nil)
func (s *CollectionService) Status(ctx context.Context, invoiceID string, opts *StatusOptions) (*StatusResponse, error) {
	req := &statusRequest{
		InvoiceID: invoiceID,
		PublicKey: s.client.publishableKey,
	}

	if opts != nil {
		req.CheckoutID = opts.CheckoutID
		req.Signature = opts.Signature
	}

	var resp StatusResponse
	if err := s.client.postPublic(ctx, "/payment/status/", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
