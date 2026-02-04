package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	intasend "github.com/emilio-kariuki/intasend-go"
)

func TestCheckout_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/checkout/" {
			t.Errorf("expected /checkout/, got %s", r.URL.Path)
		}
		// Should NOT have auth header (postPublic)
		if r.Header.Get("Authorization") != "" {
			t.Error("Checkout.Create should not send Authorization header")
		}

		var body createCheckoutRequestBody
		json.NewDecoder(r.Body).Decode(&body)
		if body.PublicKey != "ISPubKey_test_abc123" {
			t.Errorf("expected public key, got %q", body.PublicKey)
		}
		if body.Amount != 1000 {
			t.Errorf("expected amount 1000, got %v", body.Amount)
		}
		if body.Currency != "KES" {
			t.Errorf("expected KES, got %s", body.Currency)
		}
		if body.Email != "jane@example.com" {
			t.Errorf("expected jane@example.com, got %s", body.Email)
		}
		if body.FirstName != "Jane" {
			t.Errorf("expected Jane, got %s", body.FirstName)
		}
		if body.Host != "https://mysite.com" {
			t.Errorf("expected https://mysite.com, got %s", body.Host)
		}
		if body.RedirectURL != "https://mysite.com/callback" {
			t.Errorf("expected redirect URL, got %s", body.RedirectURL)
		}
		if body.APIRef != "order-99" {
			t.Errorf("expected order-99, got %s", body.APIRef)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.CreateCheckoutResponse{
			ID:        "CHK-999",
			URL:       "https://checkout.intasend.com/CHK-999",
			Signature: "sig-xyz",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Checkout().Create(context.Background(), &intasend.CreateCheckoutRequest{
		Amount:   1000,
		Currency: "KES",
		Customer: intasend.CheckoutCustomer{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "jane@example.com",
			Country:   "KE",
		},
		Host:        "https://mysite.com",
		RedirectURL: "https://mysite.com/callback",
		APIRef:      "order-99",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "CHK-999" {
		t.Errorf("expected CHK-999, got %s", resp.ID)
	}
	if resp.URL != "https://checkout.intasend.com/CHK-999" {
		t.Errorf("unexpected URL: %s", resp.URL)
	}
	if resp.Signature != "sig-xyz" {
		t.Errorf("expected sig-xyz, got %s", resp.Signature)
	}
}

func TestCheckout_CheckStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/payment/status/" {
			t.Errorf("expected /payment/status/, got %s", r.URL.Path)
		}
		// Should NOT have auth header (postPublic)
		if r.Header.Get("Authorization") != "" {
			t.Error("CheckStatus should not send Authorization header")
		}

		var body intasend.CheckoutStatusRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Signature != "sig-xyz" {
			t.Errorf("expected sig-xyz, got %s", body.Signature)
		}
		if body.CheckoutID != "CHK-999" {
			t.Errorf("expected CHK-999, got %s", body.CheckoutID)
		}
		if body.InvoiceID != "INV-999" {
			t.Errorf("expected INV-999, got %s", body.InvoiceID)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.CheckoutStatusResponse{
			Invoice: &intasend.Invoice{
				InvoiceID: "INV-999",
				State:     intasend.StateComplete,
				Value:     1000,
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Checkout().CheckStatus(context.Background(), &intasend.CheckoutStatusRequest{
		Signature:  "sig-xyz",
		CheckoutID: "CHK-999",
		InvoiceID:  "INV-999",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Invoice == nil {
		t.Fatal("expected non-nil invoice")
	}
	if resp.Invoice.State != intasend.StateComplete {
		t.Errorf("expected COMPLETE, got %s", resp.Invoice.State)
	}
	if resp.Invoice.Value != 1000 {
		t.Errorf("expected 1000, got %v", resp.Invoice.Value)
	}
}

func TestCheckout_CreateWithAllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body createCheckoutRequestBody
		json.NewDecoder(r.Body).Decode(&body)
		if body.Country != "KE" {
			t.Errorf("expected KE, got %s", body.Country)
		}
		if body.Address != "123 Main St" {
			t.Errorf("expected address, got %s", body.Address)
		}
		if body.WalletID != "W-001" {
			t.Errorf("expected W-001, got %s", body.WalletID)
		}
		if body.CardTariff != "BUSINESS-PAYS" {
			t.Errorf("expected BUSINESS-PAYS, got %s", body.CardTariff)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.CreateCheckoutResponse{ID: "CHK-FULL"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Checkout().Create(context.Background(), &intasend.CreateCheckoutRequest{
		Amount:   2000,
		Currency: "KES",
		Customer: intasend.CheckoutCustomer{
			FirstName: "John",
			Email:     "john@example.com",
			Country:   "KE",
			Address:   "123 Main St",
			City:      "Nairobi",
			State:     "Nairobi",
			Zipcode:   "00100",
		},
		Host:         "https://example.com",
		WalletID:     "W-001",
		CardTariff:   "BUSINESS-PAYS",
		MobileTariff: "CUSTOMER-PAYS",
		Comment:      "Order #123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "CHK-FULL" {
		t.Errorf("expected CHK-FULL, got %s", resp.ID)
	}
}
