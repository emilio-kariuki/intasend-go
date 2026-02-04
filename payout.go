package intasend

import (
	"context"
	"time"
)

// PayoutService handles payout/send money operations.
type PayoutService struct {
	client *Client
}

// Provider represents a payout provider type.
type Provider string

const (
	// ProviderMPesaB2C is for M-Pesa consumer payments.
	ProviderMPesaB2C Provider = "MPESA-B2C"

	// ProviderMPesaB2B is for M-Pesa business payments (PayBill/Till).
	ProviderMPesaB2B Provider = "MPESA-B2B"

	// ProviderPesaLink is for bank transfers.
	ProviderPesaLink Provider = "PESALINK"

	// ProviderIntaSend is for internal wallet transfers.
	ProviderIntaSend Provider = "INTASEND"

	// ProviderAirtime is for airtime top-ups.
	ProviderAirtime Provider = "AIRTIME"
)

// ApprovalStatus represents whether approval is required.
type ApprovalStatus string

const (
	// ApprovalRequired means the payout requires manual approval.
	ApprovalRequired ApprovalStatus = "YES"

	// ApprovalNotRequired means the payout will be processed immediately.
	ApprovalNotRequired ApprovalStatus = "NO"
)

// AccountType represents the type of M-Pesa B2B account.
type AccountType string

const (
	// AccountTypePayBill is for PayBill numbers.
	AccountTypePayBill AccountType = "PayBill"

	// AccountTypeTillNumber is for Till numbers.
	AccountTypeTillNumber AccountType = "TillNumber"
)

// Transaction represents a single payout transaction.
type Transaction struct {
	Name             string `json:"name,omitempty"`
	Account          string `json:"account"`
	Amount           string `json:"amount"`
	Narrative        string `json:"narrative,omitempty"`
	AccountType      string `json:"account_type,omitempty"`
	AccountReference string `json:"account_reference,omitempty"`
	BankCode         string `json:"bank_code,omitempty"`
}

// InitiateRequest represents a request to initiate a payout batch.
type InitiateRequest struct {
	Provider         Provider       `json:"provider"`
	Currency         string         `json:"currency"`
	Transactions     []Transaction  `json:"transactions"`
	CallbackURL      string         `json:"callback_url,omitempty"`
	WalletID         string         `json:"wallet_id,omitempty"`
	RequiresApproval ApprovalStatus `json:"requires_approval,omitempty"`
}

// InitiateResponse represents the response from initiating a payout.
type InitiateResponse struct {
	TrackingID   string              `json:"tracking_id"`
	Status       string              `json:"status"`
	Nonce        string              `json:"nonce"`
	WalletID     string              `json:"wallet_id,omitempty"`
	Transactions []TransactionResult `json:"transactions"`
	CreatedAt    time.Time           `json:"created_at"`
}

// TransactionResult represents the result of a single transaction.
type TransactionResult struct {
	Status           string    `json:"status"`
	RequestRefID     string    `json:"request_ref_id"`
	Name             string    `json:"name"`
	Account          string    `json:"account"`
	Amount           string    `json:"amount"`
	Narrative        string    `json:"narrative"`
	BankCode         string    `json:"bank_code,omitempty"`
	AccountType      string    `json:"account_type,omitempty"`
	AccountReference string    `json:"account_reference,omitempty"`
	FailedReason     string    `json:"failed_reason,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// MPesaRequest is a simplified request for M-Pesa B2C payouts.
type MPesaRequest struct {
	Currency         string
	Transactions     []Transaction
	CallbackURL      string
	WalletID         string
	RequiresApproval ApprovalStatus
}

// B2BTransaction represents an M-Pesa B2B transaction.
type B2BTransaction struct {
	Name             string
	Account          string
	AccountType      AccountType
	AccountReference string
	Amount           string
	Narrative        string
}

// MPesaB2BRequest is a request for M-Pesa B2B payouts.
type MPesaB2BRequest struct {
	Currency         string
	Transactions     []B2BTransaction
	CallbackURL      string
	WalletID         string
	RequiresApproval ApprovalStatus
}

// BankTransaction represents a bank transfer transaction.
type BankTransaction struct {
	Name      string
	Account   string
	BankCode  string
	Amount    string
	Narrative string
}

// BankRequest is a request for bank payouts.
type BankRequest struct {
	Currency         string
	Transactions     []BankTransaction
	CallbackURL      string
	WalletID         string
	RequiresApproval ApprovalStatus
}

// IntaSendTransferRequest is a request for IntaSend internal transfers.
type IntaSendTransferRequest struct {
	Currency         string
	Transactions     []Transaction
	CallbackURL      string
	WalletID         string
	RequiresApproval ApprovalStatus
}

// AirtimeRequest is a request for airtime top-ups.
type AirtimeRequest struct {
	Currency         string
	Transactions     []Transaction
	CallbackURL      string
	WalletID         string
	RequiresApproval ApprovalStatus
}

// ApproveRequest represents a request to approve a payout batch.
type ApproveRequest struct {
	TrackingID string `json:"tracking_id"`
	Nonce      string `json:"nonce,omitempty"`
	WalletID   string `json:"wallet_id,omitempty"`
}

// ApproveResponse represents the response from approving a payout.
type ApproveResponse struct {
	TrackingID   string              `json:"tracking_id"`
	Status       string              `json:"status"`
	Transactions []TransactionResult `json:"transactions"`
}

// payoutStatusRequest is the internal request for status checks.
type payoutStatusRequest struct {
	TrackingID string `json:"tracking_id"`
}

// PayoutStatusResponse represents a payout status response.
type PayoutStatusResponse struct {
	TrackingID   string              `json:"tracking_id"`
	Status       string              `json:"status"`
	Transactions []TransactionResult `json:"transactions"`
}

// Payout states
const (
	PayoutStatusPending    = "Pending"
	PayoutStatusProcessing = "Processing"
	PayoutStatusCompleted  = "Completed"
	PayoutStatusFailed     = "Failed"
)

// Initiate starts a new payout batch.
// Payouts require approval unless RequiresApproval is set to "NO".
//
// Example:
//
//	resp, err := client.Payout().Initiate(ctx, &intasend.InitiateRequest{
//	    Provider: intasend.ProviderMPesaB2C,
//	    Currency: "KES",
//	    Transactions: []intasend.Transaction{
//	        {Account: "254712345678", Amount: "100", Narrative: "Payment"},
//	    },
//	})
func (s *PayoutService) Initiate(ctx context.Context, req *InitiateRequest) (*InitiateResponse, error) {
	var resp InitiateResponse
	if err := s.client.post(ctx, "/send-money/initiate/", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MPesa initiates an M-Pesa B2C payout (consumer payments).
//
// Example:
//
//	resp, err := client.Payout().MPesa(ctx, &intasend.MPesaRequest{
//	    Currency: "KES",
//	    Transactions: []intasend.Transaction{
//	        {Account: "254712345678", Amount: "100", Narrative: "Salary"},
//	    },
//	})
func (s *PayoutService) MPesa(ctx context.Context, req *MPesaRequest) (*InitiateResponse, error) {
	initReq := &InitiateRequest{
		Provider:         ProviderMPesaB2C,
		Currency:         req.Currency,
		Transactions:     req.Transactions,
		CallbackURL:      req.CallbackURL,
		WalletID:         req.WalletID,
		RequiresApproval: req.RequiresApproval,
	}
	return s.Initiate(ctx, initReq)
}

// MPesaB2B initiates an M-Pesa B2B payout (PayBill or Till Number).
//
// Example:
//
//	resp, err := client.Payout().MPesaB2B(ctx, &intasend.MPesaB2BRequest{
//	    Currency: "KES",
//	    Transactions: []intasend.B2BTransaction{
//	        {
//	            Account:          "247247",
//	            AccountType:      intasend.AccountTypePayBill,
//	            AccountReference: "1001200010",
//	            Amount:           "2000",
//	            Narrative:        "Bill payment",
//	        },
//	    },
//	})
func (s *PayoutService) MPesaB2B(ctx context.Context, req *MPesaB2BRequest) (*InitiateResponse, error) {
	transactions := make([]Transaction, len(req.Transactions))
	for i, t := range req.Transactions {
		transactions[i] = Transaction{
			Name:             t.Name,
			Account:          t.Account,
			AccountType:      string(t.AccountType),
			AccountReference: t.AccountReference,
			Amount:           t.Amount,
			Narrative:        t.Narrative,
		}
	}

	initReq := &InitiateRequest{
		Provider:         ProviderMPesaB2B,
		Currency:         req.Currency,
		Transactions:     transactions,
		CallbackURL:      req.CallbackURL,
		WalletID:         req.WalletID,
		RequiresApproval: req.RequiresApproval,
	}
	return s.Initiate(ctx, initReq)
}

// Bank initiates a bank transfer via PesaLink.
//
// Example:
//
//	resp, err := client.Payout().Bank(ctx, &intasend.BankRequest{
//	    Currency: "KES",
//	    Transactions: []intasend.BankTransaction{
//	        {
//	            Name:      "John Doe",
//	            Account:   "0123456789",
//	            BankCode:  "2",
//	            Amount:    "5000",
//	            Narrative: "Payment",
//	        },
//	    },
//	})
func (s *PayoutService) Bank(ctx context.Context, req *BankRequest) (*InitiateResponse, error) {
	transactions := make([]Transaction, len(req.Transactions))
	for i, t := range req.Transactions {
		transactions[i] = Transaction{
			Name:      t.Name,
			Account:   t.Account,
			BankCode:  t.BankCode,
			Amount:    t.Amount,
			Narrative: t.Narrative,
		}
	}

	initReq := &InitiateRequest{
		Provider:         ProviderPesaLink,
		Currency:         req.Currency,
		Transactions:     transactions,
		CallbackURL:      req.CallbackURL,
		WalletID:         req.WalletID,
		RequiresApproval: req.RequiresApproval,
	}
	return s.Initiate(ctx, initReq)
}

// IntaSend initiates an internal IntaSend wallet transfer.
//
// Example:
//
//	resp, err := client.Payout().IntaSend(ctx, &intasend.IntaSendTransferRequest{
//	    Currency: "KES",
//	    Transactions: []intasend.Transaction{
//	        {Account: "wallet@intasend.com", Amount: "500", Narrative: "Transfer"},
//	    },
//	})
func (s *PayoutService) IntaSend(ctx context.Context, req *IntaSendTransferRequest) (*InitiateResponse, error) {
	initReq := &InitiateRequest{
		Provider:         ProviderIntaSend,
		Currency:         req.Currency,
		Transactions:     req.Transactions,
		CallbackURL:      req.CallbackURL,
		WalletID:         req.WalletID,
		RequiresApproval: req.RequiresApproval,
	}
	return s.Initiate(ctx, initReq)
}

// Airtime initiates an airtime top-up.
//
// Example:
//
//	resp, err := client.Payout().Airtime(ctx, &intasend.AirtimeRequest{
//	    Currency: "KES",
//	    Transactions: []intasend.Transaction{
//	        {Account: "254712345678", Amount: "100", Narrative: "Airtime"},
//	    },
//	})
func (s *PayoutService) Airtime(ctx context.Context, req *AirtimeRequest) (*InitiateResponse, error) {
	initReq := &InitiateRequest{
		Provider:         ProviderAirtime,
		Currency:         req.Currency,
		Transactions:     req.Transactions,
		CallbackURL:      req.CallbackURL,
		WalletID:         req.WalletID,
		RequiresApproval: req.RequiresApproval,
	}
	return s.Initiate(ctx, initReq)
}

// Approve approves a pending payout batch.
// This is required when RequiresApproval is "YES" (default).
//
// Example:
//
//	approved, err := client.Payout().Approve(ctx, &intasend.ApproveRequest{
//	    TrackingID: resp.TrackingID,
//	    Nonce:      resp.Nonce,
//	    WalletID:   resp.WalletID,
//	})
func (s *PayoutService) Approve(ctx context.Context, req *ApproveRequest) (*ApproveResponse, error) {
	var resp ApproveResponse
	if err := s.client.post(ctx, "/send-money/approve/", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Status checks the status of a payout batch.
//
// Example:
//
//	status, err := client.Payout().Status(ctx, "tracking-id-123")
func (s *PayoutService) Status(ctx context.Context, trackingID string) (*PayoutStatusResponse, error) {
	req := &payoutStatusRequest{TrackingID: trackingID}

	var resp PayoutStatusResponse
	if err := s.client.post(ctx, "/send-money/status/", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
