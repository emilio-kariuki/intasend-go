package intasend

import (
	"context"
	"fmt"
	"time"
)

// RefundService handles refund/chargeback operations.
type RefundService struct {
	client *Client
}

// RefundReason represents a refund reason.
type RefundReason string

const (
	// RefundReasonServiceUnavailable indicates the service was unavailable.
	RefundReasonServiceUnavailable RefundReason = "UNAVAILABLE"

	// RefundReasonDuplicatePayment indicates a duplicate payment.
	RefundReasonDuplicatePayment RefundReason = "DUPLICATE"

	// RefundReasonFraudulent indicates a fraudulent transaction.
	RefundReasonFraudulent RefundReason = "FRAUDULENT"

	// RefundReasonCustomerRequest indicates a customer-requested refund.
	RefundReasonCustomerRequest RefundReason = "CUSTOMER_REQUEST"

	// RefundReasonOther indicates another reason.
	RefundReasonOther RefundReason = "OTHER"
)

// Chargeback represents a refund/chargeback record.
type Chargeback struct {
	ChargebackID  string       `json:"chargeback_id"`
	Invoice       string       `json:"invoice"`
	Amount        float64      `json:"amount"`
	Status        string       `json:"status"`
	Reason        RefundReason `json:"reason"`
	ReasonDetails string       `json:"reason_details"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// ChargebackListResponse represents the response from listing chargebacks.
type ChargebackListResponse struct {
	Results []Chargeback `json:"results"`
}

// CreateChargebackRequest represents a request to create a chargeback.
type CreateChargebackRequest struct {
	Invoice       string       `json:"invoice"`
	Amount        float64      `json:"amount"`
	Reason        RefundReason `json:"reason"`
	ReasonDetails string       `json:"reason_details,omitempty"`
}

// Chargeback states
const (
	ChargebackStatusPending  = "PENDING"
	ChargebackStatusApproved = "APPROVED"
	ChargebackStatusRejected = "REJECTED"
	ChargebackStatusComplete = "COMPLETE"
)

// List returns all chargebacks/refunds.
//
// Example:
//
//	refunds, err := client.Refund().List(ctx)
func (s *RefundService) List(ctx context.Context) (*ChargebackListResponse, error) {
	var resp ChargebackListResponse
	if err := s.client.get(ctx, "/chargebacks/", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create initiates a new refund/chargeback request.
//
// Example:
//
//	chargeback, err := client.Refund().Create(ctx, &intasend.CreateChargebackRequest{
//	    Invoice:       "INV-123",
//	    Amount:        500,
//	    Reason:        intasend.RefundReasonCustomerRequest,
//	    ReasonDetails: "Customer requested cancellation",
//	})
func (s *RefundService) Create(ctx context.Context, req *CreateChargebackRequest) (*Chargeback, error) {
	var resp Chargeback
	if err := s.client.post(ctx, "/chargebacks/", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a specific chargeback by ID.
//
// Example:
//
//	chargeback, err := client.Refund().Get(ctx, "CHG-123")
func (s *RefundService) Get(ctx context.Context, chargebackID string) (*Chargeback, error) {
	var resp Chargeback
	if err := s.client.get(ctx, fmt.Sprintf("/chargebacks/%s/", chargebackID), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
