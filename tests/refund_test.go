package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	intasend "github.com/intasend/intasend-go"
)

func TestRefund_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/chargebacks/" {
			t.Errorf("expected /chargebacks/, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.ChargebackListResponse{
			Results: []intasend.Chargeback{
				{ChargebackID: "CHG-001", Invoice: "INV-100", Amount: 500, Status: intasend.ChargebackStatusPending},
				{ChargebackID: "CHG-002", Invoice: "INV-200", Amount: 300, Status: intasend.ChargebackStatusApproved},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Refund().List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 chargebacks, got %d", len(resp.Results))
	}
	if resp.Results[0].ChargebackID != "CHG-001" {
		t.Errorf("expected CHG-001, got %s", resp.Results[0].ChargebackID)
	}
	if resp.Results[1].Status != intasend.ChargebackStatusApproved {
		t.Errorf("expected APPROVED, got %s", resp.Results[1].Status)
	}
}

func TestRefund_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/chargebacks/" {
			t.Errorf("expected /chargebacks/, got %s", r.URL.Path)
		}

		var body intasend.CreateChargebackRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Invoice != "INV-100" {
			t.Errorf("expected INV-100, got %s", body.Invoice)
		}
		if body.Amount != 500 {
			t.Errorf("expected 500, got %v", body.Amount)
		}
		if body.Reason != intasend.RefundReasonCustomerRequest {
			t.Errorf("expected CUSTOMER_REQUEST, got %s", body.Reason)
		}
		if body.ReasonDetails != "Customer cancelled" {
			t.Errorf("expected 'Customer cancelled', got %q", body.ReasonDetails)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(intasend.Chargeback{
			ChargebackID:  "CHG-NEW",
			Invoice:       "INV-100",
			Amount:        500,
			Status:        intasend.ChargebackStatusPending,
			Reason:        intasend.RefundReasonCustomerRequest,
			ReasonDetails: "Customer cancelled",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Refund().Create(context.Background(), &intasend.CreateChargebackRequest{
		Invoice:       "INV-100",
		Amount:        500,
		Reason:        intasend.RefundReasonCustomerRequest,
		ReasonDetails: "Customer cancelled",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ChargebackID != "CHG-NEW" {
		t.Errorf("expected CHG-NEW, got %s", resp.ChargebackID)
	}
	if resp.Reason != intasend.RefundReasonCustomerRequest {
		t.Errorf("expected CUSTOMER_REQUEST, got %s", resp.Reason)
	}
}

func TestRefund_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/chargebacks/CHG-001/" {
			t.Errorf("expected /chargebacks/CHG-001/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.Chargeback{
			ChargebackID: "CHG-001",
			Invoice:      "INV-100",
			Amount:       500,
			Status:       intasend.ChargebackStatusComplete,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Refund().Get(context.Background(), "CHG-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ChargebackID != "CHG-001" {
		t.Errorf("expected CHG-001, got %s", resp.ChargebackID)
	}
	if resp.Status != intasend.ChargebackStatusComplete {
		t.Errorf("expected COMPLETE, got %s", resp.Status)
	}
}

func TestRefund_GetNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"detail": "Not found"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.Refund().Get(context.Background(), "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr := intasend.AsAPIError(err)
	if apiErr == nil {
		t.Fatal("expected APIError")
	}
	if !apiErr.IsNotFound() {
		t.Error("expected IsNotFound() to be true")
	}
}
