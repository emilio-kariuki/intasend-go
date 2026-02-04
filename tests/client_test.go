package tests

import (
	"testing"

	intasend "github.com/emilio-kariuki/intasend-go"
)

func TestNew_WithTestPublishableKey(t *testing.T) {
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc123"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.BaseURL() != intasend.SandboxBaseURL {
		t.Errorf("expected sandbox URL, got %s", client.BaseURL())
	}
	if !client.IsSandbox() {
		t.Error("expected IsSandbox() to be true")
	}
	if client.IsProduction() {
		t.Error("expected IsProduction() to be false")
	}
}

func TestNew_WithLivePublishableKey(t *testing.T) {
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_live_abc123"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.BaseURL() != intasend.ProductionBaseURL {
		t.Errorf("expected production URL, got %s", client.BaseURL())
	}
	if client.IsSandbox() {
		t.Error("expected IsSandbox() to be false")
	}
	if !client.IsProduction() {
		t.Error("expected IsProduction() to be true")
	}
}

func TestNew_WithTestSecretKey(t *testing.T) {
	client, err := intasend.New(
		intasend.WithSecretKey("ISSecretKey_test_abc123"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.BaseURL() != intasend.SandboxBaseURL {
		t.Errorf("expected sandbox URL, got %s", client.BaseURL())
	}
}

func TestNew_WithLiveSecretKey(t *testing.T) {
	client, err := intasend.New(
		intasend.WithSecretKey("ISSecretKey_live_abc123"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.BaseURL() != intasend.ProductionBaseURL {
		t.Errorf("expected production URL, got %s", client.BaseURL())
	}
}

func TestNew_PublishableKeyTakesPrecedence(t *testing.T) {
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_live_abc"),
		intasend.WithSecretKey("ISSecretKey_test_abc"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Publishable key is checked first, so live should win
	if client.BaseURL() != intasend.ProductionBaseURL {
		t.Errorf("expected production URL (publishable key precedence), got %s", client.BaseURL())
	}
}

func TestNew_NoKeys(t *testing.T) {
	_, err := intasend.New()
	if err != intasend.ErrNoKeysProvided {
		t.Errorf("expected ErrNoKeysProvided, got %v", err)
	}
}

func TestNew_InvalidKeyPrefix(t *testing.T) {
	_, err := intasend.New(
		intasend.WithPublishableKey("INVALID_KEY"),
	)
	if err != intasend.ErrInvalidEnvironment {
		t.Errorf("expected ErrInvalidEnvironment, got %v", err)
	}
}

func TestNew_WithBaseURLOverride(t *testing.T) {
	customURL := "https://custom.example.com/api/v1"
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
		intasend.WithBaseURL(customURL),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.BaseURL() != customURL {
		t.Errorf("expected %s, got %s", customURL, client.BaseURL())
	}
}

func TestNew_WithSandboxOverride(t *testing.T) {
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_live_abc"),
		intasend.WithSandbox(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.BaseURL() != intasend.SandboxBaseURL {
		t.Errorf("expected sandbox URL override, got %s", client.BaseURL())
	}
}

func TestNew_WithProductionOverride(t *testing.T) {
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
		intasend.WithProduction(),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.BaseURL() != intasend.ProductionBaseURL {
		t.Errorf("expected production URL override, got %s", client.BaseURL())
	}
}

func TestNew_Defaults(t *testing.T) {
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.BaseURL() != intasend.SandboxBaseURL {
		t.Errorf("expected sandbox URL, got %s", client.BaseURL())
	}
	if client.PublishableKey() != "ISPubKey_test_abc" {
		t.Errorf("expected ISPubKey_test_abc, got %s", client.PublishableKey())
	}
}

func TestNew_ServicesInitialized(t *testing.T) {
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
		intasend.WithSecretKey("ISSecretKey_test_abc"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Collection() == nil {
		t.Error("expected Collection() to be non-nil")
	}
	if client.Payout() == nil {
		t.Error("expected Payout() to be non-nil")
	}
	if client.Wallet() == nil {
		t.Error("expected Wallet() to be non-nil")
	}
	if client.Refund() == nil {
		t.Error("expected Refund() to be non-nil")
	}
	if client.Checkout() == nil {
		t.Error("expected Checkout() to be non-nil")
	}
	if client.PaymentLink() == nil {
		t.Error("expected PaymentLink() to be non-nil")
	}
}

func TestNew_ServicesSameInstance(t *testing.T) {
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Multiple calls should return the same instance
	c1 := client.Collection()
	c2 := client.Collection()
	if c1 != c2 {
		t.Error("expected same Collection instance on repeated calls")
	}
}

func TestClient_PublishableKey(t *testing.T) {
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.PublishableKey() != "ISPubKey_test_abc" {
		t.Errorf("expected ISPubKey_test_abc, got %s", client.PublishableKey())
	}
}

func TestClient_BaseURL(t *testing.T) {
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.BaseURL() != intasend.SandboxBaseURL {
		t.Errorf("expected %s, got %s", intasend.SandboxBaseURL, client.BaseURL())
	}
}
