package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	intasend "github.com/emilio-kariuki/intasend-go"
)

func TestPaymentLink_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/paymentlinks/" {
			t.Errorf("expected /paymentlinks/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.PaymentLinkListResponse{
			Results: []intasend.PaymentLink{
				{LinkID: "LNK-001", Title: "Premium", Currency: "KES", Amount: 5000, IsActive: true},
				{LinkID: "LNK-002", Title: "Basic", Currency: "KES", Amount: 1000, IsActive: false},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.PaymentLink().List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 links, got %d", len(resp.Results))
	}
	if resp.Results[0].LinkID != "LNK-001" {
		t.Errorf("expected LNK-001, got %s", resp.Results[0].LinkID)
	}
	if resp.Results[0].Title != "Premium" {
		t.Errorf("expected Premium, got %s", resp.Results[0].Title)
	}
	if !resp.Results[0].IsActive {
		t.Error("expected first link to be active")
	}
	if resp.Results[1].IsActive {
		t.Error("expected second link to be inactive")
	}
}

func TestPaymentLink_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/paymentlinks/" {
			t.Errorf("expected /paymentlinks/, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}

		var body intasend.CreatePaymentLinkRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Title != "Gold Plan" {
			t.Errorf("expected 'Gold Plan', got %q", body.Title)
		}
		if body.Currency != "KES" {
			t.Errorf("expected KES, got %s", body.Currency)
		}
		if body.Amount != 10000 {
			t.Errorf("expected 10000, got %v", body.Amount)
		}
		if body.MobileTariff != intasend.TariffBusinessPays {
			t.Errorf("expected BUSINESS-PAYS, got %s", body.MobileTariff)
		}
		if !body.IsActive {
			t.Error("expected is_active to be true")
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(intasend.PaymentLink{
			LinkID:   "LNK-NEW",
			Title:    "Gold Plan",
			Currency: "KES",
			Amount:   10000,
			URL:      "https://intasend.com/pay/LNK-NEW",
			IsActive: true,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.PaymentLink().Create(context.Background(), &intasend.CreatePaymentLinkRequest{
		Title:        "Gold Plan",
		Currency:     "KES",
		Amount:       10000,
		MobileTariff: intasend.TariffBusinessPays,
		CardTariff:   intasend.TariffCustomerPays,
		IsActive:     true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.LinkID != "LNK-NEW" {
		t.Errorf("expected LNK-NEW, got %s", resp.LinkID)
	}
	if resp.URL != "https://intasend.com/pay/LNK-NEW" {
		t.Errorf("unexpected URL: %s", resp.URL)
	}
}

func TestPaymentLink_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/paymentlinks/LNK-001/" {
			t.Errorf("expected /paymentlinks/LNK-001/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.PaymentLink{
			LinkID:       "LNK-001",
			Title:        "Premium",
			Currency:     "KES",
			Amount:       5000,
			MobileTariff: intasend.TariffBusinessPays,
			CardTariff:   intasend.TariffBusinessPays,
			IsActive:     true,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.PaymentLink().Get(context.Background(), "LNK-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.LinkID != "LNK-001" {
		t.Errorf("expected LNK-001, got %s", resp.LinkID)
	}
	if resp.MobileTariff != intasend.TariffBusinessPays {
		t.Errorf("expected BUSINESS-PAYS, got %s", resp.MobileTariff)
	}
	if resp.Amount != 5000 {
		t.Errorf("expected 5000, got %v", resp.Amount)
	}
}

func TestPaymentLink_GetNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"detail": "Not found"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.PaymentLink().Get(context.Background(), "NONEXISTENT")
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
