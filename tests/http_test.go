package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	intasend "github.com/intasend/intasend-go"
)

func TestHTTP_AuthenticatedGetHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// Verify auth header is present for authenticated requests
		auth := r.Header.Get("Authorization")
		if auth != "Bearer ISSecretKey_test_secret" {
			t.Errorf("expected Bearer token, got %q", auth)
		}
		// Verify public key headers
		if r.Header.Get("X-IntaSend-Public-API-Key") != "ISPubKey_test_abc123" {
			t.Errorf("missing X-IntaSend-Public-API-Key header")
		}
		if r.Header.Get("INTASEND_PUBLIC_API_KEY") != "ISPubKey_test_abc123" {
			t.Errorf("missing INTASEND_PUBLIC_API_KEY header")
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"results": []interface{}{}})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.Wallet().List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTP_AuthenticatedPostHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content type")
		}
		auth := r.Header.Get("Authorization")
		if auth != "Bearer ISSecretKey_test_secret" {
			t.Errorf("expected Bearer token for authenticated POST, got %q", auth)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"invoice": map[string]interface{}{"invoice_id": "INV-1", "state": "PENDING"}})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.Collection().MPesaSTKPush(context.Background(), &intasend.STKPushRequest{
		PhoneNumber: "254712345678",
		Amount:      100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTP_PublicPostNoAuthHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			t.Errorf("public POST should not include Authorization header, got %q", auth)
		}
		// Public key headers should still be present
		if r.Header.Get("X-IntaSend-Public-API-Key") == "" {
			t.Error("expected X-IntaSend-Public-API-Key header for public POST")
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "CHK-123", "url": "", "signature": ""})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.Checkout().Create(context.Background(), &intasend.CreateCheckoutRequest{
		Amount:   100,
		Currency: "KES",
		Customer: intasend.CheckoutCustomer{Email: "test@example.com"},
		Host:     "https://example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTP_UserAgentHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if ua != "intasend-go/"+intasend.Version {
			t.Errorf("expected intasend-go/%s user agent, got %q", intasend.Version, ua)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"results": []interface{}{}})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.Wallet().List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTP_CustomUserAgent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if ua != "my-app/2.0" {
			t.Errorf("expected custom user agent, got %q", ua)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"results": []interface{}{}})
	}))
	defer server.Close()

	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
		intasend.WithSecretKey("ISSecretKey_test_abc"),
		intasend.WithBaseURL(server.URL),
		intasend.WithHTTPClient(server.Client()),
		intasend.WithRetry(0, 0),
		intasend.WithUserAgent("my-app/2.0"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = client.Wallet().List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTP_APIError400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid phone number",
			"errors":  map[string][]string{"phone_number": {"required"}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.Wallet().List(context.Background())
	if err == nil {
		t.Fatal("expected error for 400 response")
	}

	apiErr := intasend.AsAPIError(err)
	if apiErr == nil {
		t.Fatal("expected APIError")
	}
	if apiErr.HTTPStatusCode != 400 {
		t.Errorf("expected status 400, got %d", apiErr.HTTPStatusCode)
	}
	if apiErr.Message != "Invalid phone number" {
		t.Errorf("expected message 'Invalid phone number', got %q", apiErr.Message)
	}
	if !apiErr.IsValidationError() {
		t.Error("expected IsValidationError() to be true")
	}
}

func TestHTTP_APIError401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"detail": "Invalid token"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.Wallet().List(context.Background())

	apiErr := intasend.AsAPIError(err)
	if apiErr == nil {
		t.Fatal("expected APIError")
	}
	if !apiErr.IsAuthenticationError() {
		t.Error("expected IsAuthenticationError() to be true")
	}
}

func TestHTTP_APIError404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"detail": "Not found"})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.Refund().Get(context.Background(), "NONEXISTENT")

	apiErr := intasend.AsAPIError(err)
	if apiErr == nil {
		t.Fatal("expected APIError")
	}
	if !apiErr.IsNotFound() {
		t.Error("expected IsNotFound() to be true")
	}
}

func TestHTTP_NoRetryOnClientError(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"bad"}`))
	}))
	defer server.Close()

	client, _ := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
		intasend.WithSecretKey("ISSecretKey_test_abc"),
		intasend.WithBaseURL(server.URL),
		intasend.WithHTTPClient(server.Client()),
		intasend.WithRetry(3, 1*time.Millisecond),
	)
	_, _ = client.Wallet().List(context.Background())

	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("expected 1 call (no retries on 4xx), got %d", calls)
	}
}

func TestHTTP_RetryOnServerError(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&calls, 1)
		if count <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"server error"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"results": []interface{}{}})
	}))
	defer server.Close()

	client, _ := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
		intasend.WithSecretKey("ISSecretKey_test_abc"),
		intasend.WithBaseURL(server.URL),
		intasend.WithHTTPClient(server.Client()),
		intasend.WithRetry(3, 1*time.Millisecond),
	)

	_, err := client.Wallet().List(context.Background())
	if err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}
	if atomic.LoadInt32(&calls) != 3 {
		t.Errorf("expected 3 calls (2 failures + 1 success), got %d", calls)
	}
}

func TestHTTP_RetryOn429(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&calls, 1)
		if count == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"message":"rate limited"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"results": []interface{}{}})
	}))
	defer server.Close()

	client, _ := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
		intasend.WithSecretKey("ISSecretKey_test_abc"),
		intasend.WithBaseURL(server.URL),
		intasend.WithHTTPClient(server.Client()),
		intasend.WithRetry(2, 1*time.Millisecond),
	)

	_, err := client.Wallet().List(context.Background())
	if err != nil {
		t.Fatalf("expected success after 429 retry, got %v", err)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestHTTP_AllRetriesExhausted(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"always failing"}`))
	}))
	defer server.Close()

	client, _ := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
		intasend.WithSecretKey("ISSecretKey_test_abc"),
		intasend.WithBaseURL(server.URL),
		intasend.WithHTTPClient(server.Client()),
		intasend.WithRetry(2, 1*time.Millisecond),
	)

	_, err := client.Wallet().List(context.Background())
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	apiErr := intasend.AsAPIError(err)
	if apiErr == nil {
		t.Fatal("expected APIError")
	}
	if apiErr.Message != "always failing" {
		t.Errorf("expected 'always failing', got %q", apiErr.Message)
	}
	// 1 initial + 2 retries = 3
	if atomic.LoadInt32(&calls) != 3 {
		t.Errorf("expected 3 calls (1 + 2 retries), got %d", calls)
	}
}

func TestHTTP_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := client.Wallet().List(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestHTTP_NonJSONErrorBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("plain text error"))
	}))
	defer server.Close()

	client := newTestClient(t, server)
	_, err := client.Wallet().List(context.Background())

	apiErr := intasend.AsAPIError(err)
	if apiErr == nil {
		t.Fatal("expected APIError")
	}
	if apiErr.Message != "plain text error" {
		t.Errorf("expected plain text in message, got %q", apiErr.Message)
	}
}
