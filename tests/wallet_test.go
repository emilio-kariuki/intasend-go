package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	intasend "github.com/emilio-kariuki/intasend-go"
)

func TestWallet_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/wallets/" {
			t.Errorf("expected /wallets/, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.WalletListResponse{
			Results: []intasend.Wallet{
				{WalletID: "W-001", Label: "Main", Currency: "KES", AvailableBalance: 5000},
				{WalletID: "W-002", Label: "Ops", Currency: "USD", AvailableBalance: 100},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Wallet().List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 wallets, got %d", len(resp.Results))
	}
	if resp.Results[0].WalletID != "W-001" {
		t.Errorf("expected W-001, got %s", resp.Results[0].WalletID)
	}
	if resp.Results[1].AvailableBalance != 100 {
		t.Errorf("expected balance 100, got %v", resp.Results[1].AvailableBalance)
	}
}

func TestWallet_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body intasend.CreateWalletRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Currency != "KES" {
			t.Errorf("expected KES, got %s", body.Currency)
		}
		if body.Label != "Test Wallet" {
			t.Errorf("expected 'Test Wallet', got %q", body.Label)
		}
		if body.WalletType != intasend.WalletTypeWorking {
			t.Errorf("expected WORKING wallet type (default), got %s", body.WalletType)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(intasend.Wallet{
			WalletID:   "W-NEW",
			Label:      "Test Wallet",
			Currency:   "KES",
			WalletType: intasend.WalletTypeWorking,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Wallet().Create(context.Background(), &intasend.CreateWalletRequest{
		Currency: "KES",
		Label:    "Test Wallet",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.WalletID != "W-NEW" {
		t.Errorf("expected W-NEW, got %s", resp.WalletID)
	}
}

func TestWallet_Create_DefaultsWalletType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body intasend.CreateWalletRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.WalletType != intasend.WalletTypeWorking {
			t.Errorf("expected default wallet type WORKING, got %s", body.WalletType)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.Wallet{WalletID: "W-1"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.Wallet().Create(context.Background(), &intasend.CreateWalletRequest{
		Currency: "KES",
		Label:    "Test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWallet_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wallets/W-001/" {
			t.Errorf("expected /wallets/W-001/, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.Wallet{
			WalletID:         "W-001",
			Label:            "Main",
			Currency:         "KES",
			AvailableBalance: 10000,
			CanDisburse:      true,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Wallet().Get(context.Background(), "W-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.WalletID != "W-001" {
		t.Errorf("expected W-001, got %s", resp.WalletID)
	}
	if !resp.CanDisburse {
		t.Error("expected CanDisburse to be true")
	}
}

func TestWallet_Transactions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wallets/W-001/transactions/" {
			t.Errorf("expected /wallets/W-001/transactions/, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.WalletTransactionsResponse{
			Results: []intasend.WalletTransaction{
				{TransactionID: "TXN-1", Amount: 500, TransType: "CREDIT"},
				{TransactionID: "TXN-2", Amount: 200, TransType: "DEBIT"},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Wallet().Transactions(context.Background(), "W-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(resp.Results))
	}
	if resp.Results[0].TransactionID != "TXN-1" {
		t.Errorf("expected TXN-1, got %s", resp.Results[0].TransactionID)
	}
}

func TestWallet_IntraTransfer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wallets/W-001/intra_transfer/" {
			t.Errorf("expected /wallets/W-001/intra_transfer/, got %s", r.URL.Path)
		}

		var body intraTransferRequestBody
		json.NewDecoder(r.Body).Decode(&body)
		if body.WalletID != "W-002" {
			t.Errorf("expected destination W-002, got %s", body.WalletID)
		}
		if body.Amount != 1000 {
			t.Errorf("expected amount 1000, got %v", body.Amount)
		}
		if body.Narrative != "Commission" {
			t.Errorf("expected narrative Commission, got %q", body.Narrative)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.IntraTransferResponse{
			Status:    "success",
			OriginID:  "W-001",
			TargetID:  "W-002",
			Amount:    1000,
			Narrative: "Commission",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Wallet().IntraTransfer(context.Background(), &intasend.IntraTransferRequest{
		SourceID:      "W-001",
		DestinationID: "W-002",
		Amount:        1000,
		Narrative:     "Commission",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got %s", resp.Status)
	}
	if resp.OriginID != "W-001" {
		t.Errorf("expected origin W-001, got %s", resp.OriginID)
	}
}

func TestWallet_FundMPesa(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/payment/mpesa-stk-push/" {
			t.Errorf("expected /payment/mpesa-stk-push/, got %s", r.URL.Path)
		}

		var body fundMPesaRequestBody
		json.NewDecoder(r.Body).Decode(&body)
		if body.Method != "M-PESA" {
			t.Errorf("expected M-PESA, got %s", body.Method)
		}
		if body.Currency != "KES" {
			t.Errorf("expected KES, got %s", body.Currency)
		}
		if body.WalletID != "W-001" {
			t.Errorf("expected W-001, got %s", body.WalletID)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.FundMPesaResponse{
			Invoice: &intasend.Invoice{InvoiceID: "INV-FUND", State: "PENDING"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Wallet().FundMPesa(context.Background(), &intasend.FundMPesaRequest{
		WalletID:    "W-001",
		PhoneNumber: "254712345678",
		Amount:      5000,
		Email:       "test@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Invoice.InvoiceID != "INV-FUND" {
		t.Errorf("expected INV-FUND, got %s", resp.Invoice.InvoiceID)
	}
}

func TestWallet_FundCheckout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/checkout/" {
			t.Errorf("expected /checkout/, got %s", r.URL.Path)
		}
		// Should NOT have auth header (postPublic)
		if r.Header.Get("Authorization") != "" {
			t.Error("FundCheckout should not send Authorization header")
		}

		var body fundCheckoutRequestBody
		json.NewDecoder(r.Body).Decode(&body)
		if body.WalletID != "W-001" {
			t.Errorf("expected W-001, got %s", body.WalletID)
		}
		if body.Email != "john@example.com" {
			t.Errorf("expected john@example.com, got %s", body.Email)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.FundCheckoutResponse{
			ID:  "CHK-FUND",
			URL: "https://checkout.intasend.com/CHK-FUND",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Wallet().FundCheckout(context.Background(), &intasend.FundCheckoutRequest{
		WalletID: "W-001",
		Amount:   5000,
		Currency: "KES",
		Customer: intasend.WalletCustomer{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		},
		Host: "https://example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != "CHK-FUND" {
		t.Errorf("expected CHK-FUND, got %s", resp.ID)
	}
}
