package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	intasend "github.com/emilio-kariuki/intasend-go"
)

func TestCollection_Charge(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/checkout/" {
			t.Errorf("expected /checkout/, got %s", r.URL.Path)
		}
		// Should NOT have auth header (postPublic)
		if r.Header.Get("Authorization") != "" {
			t.Error("Charge should not send Authorization header")
		}

		var body chargeRequestBody
		json.NewDecoder(r.Body).Decode(&body)
		if body.PublicKey != "ISPubKey_test_abc123" {
			t.Errorf("expected public key in body, got %q", body.PublicKey)
		}
		if body.Email != "john@example.com" {
			t.Errorf("expected email john@example.com, got %q", body.Email)
		}
		if body.Amount != 100 {
			t.Errorf("expected amount 100, got %v", body.Amount)
		}
		if body.Currency != "KES" {
			t.Errorf("expected currency KES, got %q", body.Currency)
		}
		if body.FirstName != "John" {
			t.Errorf("expected first_name John, got %q", body.FirstName)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.ChargeResponse{
			ID:        "CHK-123",
			URL:       "https://checkout.intasend.com/CHK-123",
			Signature: "sig-abc",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Collection().Charge(context.Background(), &intasend.ChargeRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Host:      "https://example.com",
		Amount:    100,
		Currency:  "KES",
		APIRef:    "order-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "CHK-123" {
		t.Errorf("expected ID CHK-123, got %s", resp.ID)
	}
	if resp.URL != "https://checkout.intasend.com/CHK-123" {
		t.Errorf("unexpected URL: %s", resp.URL)
	}
	if resp.Signature != "sig-abc" {
		t.Errorf("expected signature sig-abc, got %s", resp.Signature)
	}
}

func TestCollection_MPesaSTKPush(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/payment/mpesa-stk-push/" {
			t.Errorf("expected /payment/mpesa-stk-push/, got %s", r.URL.Path)
		}
		// Should have auth header (post)
		if r.Header.Get("Authorization") == "" {
			t.Error("MPesaSTKPush should include Authorization header")
		}

		var body stkPushRequestBody
		json.NewDecoder(r.Body).Decode(&body)
		if body.Method != "M-PESA" {
			t.Errorf("expected method M-PESA, got %q", body.Method)
		}
		if body.Currency != "KES" {
			t.Errorf("expected currency KES, got %q", body.Currency)
		}
		if body.PhoneNumber != "254712345678" {
			t.Errorf("expected phone 254712345678, got %q", body.PhoneNumber)
		}
		if body.Amount != 500 {
			t.Errorf("expected amount 500, got %v", body.Amount)
		}
		if body.PublicKey != "ISPubKey_test_abc123" {
			t.Errorf("expected public key in body, got %q", body.PublicKey)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.STKPushResponse{
			Invoice: &intasend.Invoice{
				InvoiceID: "INV-456",
				State:     "PENDING",
				Provider:  "M-PESA",
				Value:     500,
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Collection().MPesaSTKPush(context.Background(), &intasend.STKPushRequest{
		PhoneNumber: "254712345678",
		Amount:      500,
		APIRef:      "test-ref",
		Name:        "Test User",
		Email:       "test@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Invoice == nil {
		t.Fatal("expected non-nil invoice")
	}
	if resp.Invoice.InvoiceID != "INV-456" {
		t.Errorf("expected invoice ID INV-456, got %s", resp.Invoice.InvoiceID)
	}
	if resp.Invoice.State != "PENDING" {
		t.Errorf("expected state PENDING, got %s", resp.Invoice.State)
	}
	if resp.Invoice.Value != 500 {
		t.Errorf("expected value 500, got %v", resp.Invoice.Value)
	}
}

func TestCollection_Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/payment/status/" {
			t.Errorf("expected /payment/status/, got %s", r.URL.Path)
		}
		// Should NOT have auth header (postPublic)
		if r.Header.Get("Authorization") != "" {
			t.Error("Status should not send Authorization header")
		}

		var body statusRequestBody
		json.NewDecoder(r.Body).Decode(&body)
		if body.InvoiceID != "INV-456" {
			t.Errorf("expected invoice_id INV-456, got %q", body.InvoiceID)
		}
		if body.PublicKey != "ISPubKey_test_abc123" {
			t.Errorf("expected public key, got %q", body.PublicKey)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.StatusResponse{
			Invoice: &intasend.Invoice{
				InvoiceID: "INV-456",
				State:     "COMPLETE",
				Value:     500,
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Collection().Status(context.Background(), "INV-456", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Invoice.State != "COMPLETE" {
		t.Errorf("expected state COMPLETE, got %s", resp.Invoice.State)
	}
}

func TestCollection_StatusWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body statusRequestBody
		json.NewDecoder(r.Body).Decode(&body)
		if body.CheckoutID != "CHK-123" {
			t.Errorf("expected checkout_id CHK-123, got %q", body.CheckoutID)
		}
		if body.Signature != "sig-abc" {
			t.Errorf("expected signature sig-abc, got %q", body.Signature)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.StatusResponse{
			Invoice: &intasend.Invoice{InvoiceID: "INV-456", State: "COMPLETE"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Collection().Status(context.Background(), "INV-456", &intasend.StatusOptions{
		CheckoutID: "CHK-123",
		Signature:  "sig-abc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Invoice.InvoiceID != "INV-456" {
		t.Errorf("expected INV-456, got %s", resp.Invoice.InvoiceID)
	}
}

func TestCollection_ChargeAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid currency"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.Collection().Charge(context.Background(), &intasend.ChargeRequest{
		Email:    "test@example.com",
		Amount:   100,
		Currency: "INVALID",
		Host:     "https://example.com",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr := intasend.AsAPIError(err)
	if apiErr == nil {
		t.Fatal("expected APIError")
	}
	if apiErr.HTTPStatusCode != 400 {
		t.Errorf("expected 400, got %d", apiErr.HTTPStatusCode)
	}
}
