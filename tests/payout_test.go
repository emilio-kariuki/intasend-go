package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	intasend "github.com/intasend/intasend-go"
)

func TestPayout_Initiate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/send-money/initiate/" {
			t.Errorf("expected /send-money/initiate/, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}

		var body intasend.InitiateRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Provider != intasend.ProviderMPesaB2C {
			t.Errorf("expected MPESA-B2C, got %s", body.Provider)
		}
		if body.Currency != "KES" {
			t.Errorf("expected KES, got %s", body.Currency)
		}
		if len(body.Transactions) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(body.Transactions))
		}
		if body.Transactions[0].Account != "254712345678" {
			t.Errorf("expected account 254712345678, got %s", body.Transactions[0].Account)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.InitiateResponse{
			TrackingID: "TRK-001",
			Status:     "Pending",
			Nonce:      "nonce-abc",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Payout().Initiate(context.Background(), &intasend.InitiateRequest{
		Provider: intasend.ProviderMPesaB2C,
		Currency: "KES",
		Transactions: []intasend.Transaction{
			{Account: "254712345678", Amount: "1000", Narrative: "Salary"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TrackingID != "TRK-001" {
		t.Errorf("expected TRK-001, got %s", resp.TrackingID)
	}
	if resp.Nonce != "nonce-abc" {
		t.Errorf("expected nonce-abc, got %s", resp.Nonce)
	}
}

func TestPayout_MPesa(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body intasend.InitiateRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Provider != intasend.ProviderMPesaB2C {
			t.Errorf("expected MPESA-B2C provider, got %s", body.Provider)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.InitiateResponse{TrackingID: "TRK-MPESA"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Payout().MPesa(context.Background(), &intasend.MPesaRequest{
		Currency: "KES",
		Transactions: []intasend.Transaction{
			{Account: "254712345678", Amount: "100"},
		},
		RequiresApproval: intasend.ApprovalRequired,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TrackingID != "TRK-MPESA" {
		t.Errorf("expected TRK-MPESA, got %s", resp.TrackingID)
	}
}

func TestPayout_MPesaB2B(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body intasend.InitiateRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Provider != intasend.ProviderMPesaB2B {
			t.Errorf("expected MPESA-B2B, got %s", body.Provider)
		}
		if len(body.Transactions) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(body.Transactions))
		}
		if body.Transactions[0].AccountType != string(intasend.AccountTypePayBill) {
			t.Errorf("expected PayBill, got %s", body.Transactions[0].AccountType)
		}
		if body.Transactions[0].AccountReference != "ACC001" {
			t.Errorf("expected ACC001, got %s", body.Transactions[0].AccountReference)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.InitiateResponse{TrackingID: "TRK-B2B"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Payout().MPesaB2B(context.Background(), &intasend.MPesaB2BRequest{
		Currency: "KES",
		Transactions: []intasend.B2BTransaction{
			{
				Account:          "247247",
				AccountType:      intasend.AccountTypePayBill,
				AccountReference: "ACC001",
				Amount:           "5000",
				Narrative:        "Bill payment",
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TrackingID != "TRK-B2B" {
		t.Errorf("expected TRK-B2B, got %s", resp.TrackingID)
	}
}

func TestPayout_Bank(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body intasend.InitiateRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Provider != intasend.ProviderPesaLink {
			t.Errorf("expected PESALINK, got %s", body.Provider)
		}
		if body.Transactions[0].BankCode != "2" {
			t.Errorf("expected bank code 2, got %s", body.Transactions[0].BankCode)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.InitiateResponse{TrackingID: "TRK-BANK"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Payout().Bank(context.Background(), &intasend.BankRequest{
		Currency: "KES",
		Transactions: []intasend.BankTransaction{
			{Name: "John", Account: "123456", BankCode: "2", Amount: "5000"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TrackingID != "TRK-BANK" {
		t.Errorf("expected TRK-BANK, got %s", resp.TrackingID)
	}
}

func TestPayout_IntaSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body intasend.InitiateRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Provider != intasend.ProviderIntaSend {
			t.Errorf("expected INTASEND, got %s", body.Provider)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.InitiateResponse{TrackingID: "TRK-IS"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Payout().IntaSend(context.Background(), &intasend.IntaSendTransferRequest{
		Currency: "KES",
		Transactions: []intasend.Transaction{
			{Account: "user@intasend.com", Amount: "200"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TrackingID != "TRK-IS" {
		t.Errorf("expected TRK-IS, got %s", resp.TrackingID)
	}
}

func TestPayout_Airtime(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body intasend.InitiateRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Provider != intasend.ProviderAirtime {
			t.Errorf("expected AIRTIME, got %s", body.Provider)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.InitiateResponse{TrackingID: "TRK-AIR"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Payout().Airtime(context.Background(), &intasend.AirtimeRequest{
		Currency: "KES",
		Transactions: []intasend.Transaction{
			{Account: "254712345678", Amount: "50"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TrackingID != "TRK-AIR" {
		t.Errorf("expected TRK-AIR, got %s", resp.TrackingID)
	}
}

func TestPayout_Approve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/send-money/approve/" {
			t.Errorf("expected /send-money/approve/, got %s", r.URL.Path)
		}

		var body intasend.ApproveRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.TrackingID != "TRK-001" {
			t.Errorf("expected TRK-001, got %s", body.TrackingID)
		}
		if body.Nonce != "nonce-abc" {
			t.Errorf("expected nonce-abc, got %s", body.Nonce)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.ApproveResponse{
			TrackingID: "TRK-001",
			Status:     "Approved",
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Payout().Approve(context.Background(), &intasend.ApproveRequest{
		TrackingID: "TRK-001",
		Nonce:      "nonce-abc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "Approved" {
		t.Errorf("expected Approved, got %s", resp.Status)
	}
}

func TestPayout_Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/send-money/status/" {
			t.Errorf("expected /send-money/status/, got %s", r.URL.Path)
		}

		var body payoutStatusRequestBody
		json.NewDecoder(r.Body).Decode(&body)
		if body.TrackingID != "TRK-001" {
			t.Errorf("expected TRK-001, got %s", body.TrackingID)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(intasend.PayoutStatusResponse{
			TrackingID: "TRK-001",
			Status:     "Completed",
			Transactions: []intasend.TransactionResult{
				{Status: "Successful", Account: "254712345678", Amount: "1000"},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	resp, err := client.Payout().Status(context.Background(), "TRK-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "Completed" {
		t.Errorf("expected Completed, got %s", resp.Status)
	}
	if len(resp.Transactions) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(resp.Transactions))
	}
	if resp.Transactions[0].Account != "254712345678" {
		t.Errorf("expected account 254712345678, got %s", resp.Transactions[0].Account)
	}
}
