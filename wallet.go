package intasend

import (
	"context"
	"fmt"
	"time"
)

// WalletService handles wallet operations.
type WalletService struct {
	client *Client
}

// WalletType represents the type of wallet.
type WalletType string

const (
	// WalletTypeWorking is a standard working wallet.
	WalletTypeWorking WalletType = "WORKING"
)

// Wallet represents an IntaSend wallet.
type Wallet struct {
	WalletID         string     `json:"wallet_id"`
	Label            string     `json:"label"`
	Currency         string     `json:"currency"`
	WalletType       WalletType `json:"wallet_type"`
	CurrentBalance   float64    `json:"current_balance"`
	AvailableBalance float64    `json:"available_balance"`
	CanDisburse      bool       `json:"can_disburse"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// WalletListResponse represents the response from listing wallets.
type WalletListResponse struct {
	Results []Wallet `json:"results"`
}

// CreateWalletRequest represents a request to create a wallet.
type CreateWalletRequest struct {
	Currency    string     `json:"currency"`
	Label       string     `json:"label"`
	WalletType  WalletType `json:"wallet_type,omitempty"`
	CanDisburse bool       `json:"can_disburse,omitempty"`
}

// WalletTransaction represents a wallet transaction.
type WalletTransaction struct {
	TransactionID  string    `json:"transaction_id"`
	WalletID       string    `json:"wallet_id"`
	TransType      string    `json:"trans_type"`
	Amount         float64   `json:"amount"`
	Narrative      string    `json:"narrative"`
	RunningBalance float64   `json:"running_balance"`
	CreatedAt      time.Time `json:"created_at"`
}

// WalletTransactionsResponse represents the response from listing wallet transactions.
type WalletTransactionsResponse struct {
	Results []WalletTransaction `json:"results"`
}

// IntraTransferRequest represents a request to transfer between wallets.
type IntraTransferRequest struct {
	SourceID      string
	DestinationID string
	Amount        float64
	Narrative     string
}

// intraTransferBody is the internal request body.
type intraTransferBody struct {
	WalletID  string  `json:"wallet_id"`
	Amount    float64 `json:"amount"`
	Narrative string  `json:"narrative"`
}

// IntraTransferResponse represents the response from an intra-wallet transfer.
type IntraTransferResponse struct {
	Status    string  `json:"status"`
	OriginID  string  `json:"origin_wallet_id"`
	TargetID  string  `json:"target_wallet_id"`
	Amount    float64 `json:"amount"`
	Narrative string  `json:"narrative"`
}

// WalletCustomer represents customer information for wallet funding.
type WalletCustomer struct {
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

// FundMPesaRequest represents a request to fund a wallet via M-Pesa.
type FundMPesaRequest struct {
	WalletID    string
	PhoneNumber string
	Amount      float64
	Email       string
	APIRef      string
}

// fundMPesaBody is the internal request body.
type fundMPesaBody struct {
	PublicKey   string  `json:"public_key,omitempty"`
	WalletID    string  `json:"wallet_id"`
	PhoneNumber string  `json:"phone_number"`
	Amount      float64 `json:"amount"`
	Email       string  `json:"email,omitempty"`
	APIRef      string  `json:"api_ref,omitempty"`
	Method      string  `json:"method"`
	Currency    string  `json:"currency"`
}

// FundMPesaResponse represents the response from funding via M-Pesa.
type FundMPesaResponse struct {
	Invoice  *Invoice      `json:"invoice"`
	Customer *CustomerInfo `json:"customer,omitempty"`
}

// FundCheckoutRequest represents a request to fund a wallet via checkout.
type FundCheckoutRequest struct {
	WalletID     string
	Amount       float64
	Currency     string
	Customer     WalletCustomer
	Host         string
	RedirectURL  string
	APIRef       string
	CardTariff   string
	MobileTariff string
}

// fundCheckoutBody is the internal request body.
type fundCheckoutBody struct {
	PublicKey    string  `json:"public_key,omitempty"`
	WalletID     string  `json:"wallet_id"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Email        string  `json:"email"`
	FirstName    string  `json:"first_name,omitempty"`
	LastName     string  `json:"last_name,omitempty"`
	PhoneNumber  string  `json:"phone_number,omitempty"`
	Country      string  `json:"country,omitempty"`
	Host         string  `json:"host"`
	RedirectURL  string  `json:"redirect_url,omitempty"`
	APIRef       string  `json:"api_ref,omitempty"`
	CardTariff   string  `json:"card_tarrif,omitempty"`
	MobileTariff string  `json:"mobile_tarrif,omitempty"`
}

// FundCheckoutResponse represents the response from creating a checkout.
type FundCheckoutResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Signature string `json:"signature"`
}

// List returns all wallets in the account.
//
// Example:
//
//	wallets, err := client.Wallet().List(ctx)
func (s *WalletService) List(ctx context.Context) (*WalletListResponse, error) {
	var resp WalletListResponse
	if err := s.client.get(ctx, "/wallets/", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new wallet.
//
// Example:
//
//	wallet, err := client.Wallet().Create(ctx, &intasend.CreateWalletRequest{
//	    Currency:    "KES",
//	    Label:       "Operations Wallet",
//	    CanDisburse: true,
//	})
func (s *WalletService) Create(ctx context.Context, req *CreateWalletRequest) (*Wallet, error) {
	if req.WalletType == "" {
		req.WalletType = WalletTypeWorking
	}

	var resp Wallet
	if err := s.client.post(ctx, "/wallets/", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a specific wallet by ID.
//
// Example:
//
//	wallet, err := client.Wallet().Get(ctx, "WALLET123")
func (s *WalletService) Get(ctx context.Context, walletID string) (*Wallet, error) {
	var resp Wallet
	if err := s.client.get(ctx, fmt.Sprintf("/wallets/%s/", walletID), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Transactions retrieves transactions for a specific wallet.
//
// Example:
//
//	txns, err := client.Wallet().Transactions(ctx, "WALLET123")
func (s *WalletService) Transactions(ctx context.Context, walletID string) (*WalletTransactionsResponse, error) {
	var resp WalletTransactionsResponse
	if err := s.client.get(ctx, fmt.Sprintf("/wallets/%s/transactions/", walletID), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IntraTransfer transfers funds between two wallets in the same account.
//
// Example:
//
//	result, err := client.Wallet().IntraTransfer(ctx, &intasend.IntraTransferRequest{
//	    SourceID:      "WALLET123",
//	    DestinationID: "WALLET456",
//	    Amount:        1000,
//	    Narrative:     "Commission transfer",
//	})
func (s *WalletService) IntraTransfer(ctx context.Context, req *IntraTransferRequest) (*IntraTransferResponse, error) {
	body := &intraTransferBody{
		WalletID:  req.DestinationID,
		Amount:    req.Amount,
		Narrative: req.Narrative,
	}

	var resp IntraTransferResponse
	path := fmt.Sprintf("/wallets/%s/intra_transfer/", req.SourceID)
	if err := s.client.post(ctx, path, body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FundMPesa initiates an M-Pesa STK Push to fund a wallet.
//
// Example:
//
//	result, err := client.Wallet().FundMPesa(ctx, &intasend.FundMPesaRequest{
//	    WalletID:    "WALLET123",
//	    PhoneNumber: "254712345678",
//	    Amount:      1000,
//	    Email:       "customer@example.com",
//	    APIRef:      "fund-wallet-001",
//	})
func (s *WalletService) FundMPesa(ctx context.Context, req *FundMPesaRequest) (*FundMPesaResponse, error) {
	body := &fundMPesaBody{
		PublicKey:   s.client.publishableKey,
		WalletID:    req.WalletID,
		PhoneNumber: req.PhoneNumber,
		Amount:      req.Amount,
		Email:       req.Email,
		APIRef:      req.APIRef,
		Method:      "M-PESA",
		Currency:    "KES",
	}

	var resp FundMPesaResponse
	if err := s.client.post(ctx, "/payment/mpesa-stk-push/", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FundCheckout creates a checkout session to fund a wallet.
//
// Example:
//
//	result, err := client.Wallet().FundCheckout(ctx, &intasend.FundCheckoutRequest{
//	    WalletID: "WALLET123",
//	    Amount:   1000,
//	    Currency: "KES",
//	    Customer: intasend.WalletCustomer{
//	        FirstName: "John",
//	        LastName:  "Doe",
//	        Email:     "john@example.com",
//	    },
//	    Host:        "https://yoursite.com",
//	    RedirectURL: "https://yoursite.com/callback",
//	})
func (s *WalletService) FundCheckout(ctx context.Context, req *FundCheckoutRequest) (*FundCheckoutResponse, error) {
	body := &fundCheckoutBody{
		PublicKey:    s.client.publishableKey,
		WalletID:     req.WalletID,
		Amount:       req.Amount,
		Currency:     req.Currency,
		Email:        req.Customer.Email,
		FirstName:    req.Customer.FirstName,
		LastName:     req.Customer.LastName,
		PhoneNumber:  req.Customer.PhoneNumber,
		Country:      req.Customer.Country,
		Host:         req.Host,
		RedirectURL:  req.RedirectURL,
		APIRef:       req.APIRef,
		CardTariff:   req.CardTariff,
		MobileTariff: req.MobileTariff,
	}

	var resp FundCheckoutResponse
	if err := s.client.postPublic(ctx, "/checkout/", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
