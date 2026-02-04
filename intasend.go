// Package intasend provides a Go client for the IntaSend Payment Gateway API.
//
// IntaSend is an African payment gateway (Kenya-focused) that provides payment
// collection, payouts, wallet management, and more.
//
// Basic usage:
//
//	client, err := intasend.New(
//	    intasend.WithPublishableKey("ISPubKey_test_..."),
//	    intasend.WithSecretKey("ISSecretKey_test_..."),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Use collection service
//	resp, err := client.Collection().MPesaSTKPush(ctx, &intasend.STKPushRequest{
//	    PhoneNumber: "254712345678",
//	    Amount:      100,
//	    APIRef:      "order-123",
//	})
package intasend

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// SandboxBaseURL is the base URL for the sandbox/test environment.
	SandboxBaseURL = "https://sandbox.intasend.com/api/v1"

	// ProductionBaseURL is the base URL for the production environment.
	ProductionBaseURL = "https://payment.intasend.com/api/v1"

	// DefaultTimeout is the default HTTP request timeout.
	DefaultTimeout = 30 * time.Second

	// DefaultMaxRetries is the default number of retry attempts.
	DefaultMaxRetries = 3

	// DefaultRetryWait is the default wait time between retries.
	DefaultRetryWait = 1 * time.Second

	// Version is the SDK version.
	Version = "1.0.0"
)

// Client is the main IntaSend API client.
type Client struct {
	publishableKey string
	secretKey      string
	baseURL        string
	httpClient     *http.Client
	timeout        time.Duration
	maxRetries     int
	retryWait      time.Duration
	userAgent      string
	debug          bool

	// Services (lazily initialized)
	collection  *CollectionService
	payout      *PayoutService
	wallet      *WalletService
	refund      *RefundService
	checkout    *CheckoutService
	paymentLink *PaymentLinkService
}

// New creates a new IntaSend API client with the given options.
//
// At minimum, you must provide either WithPublishableKey or WithSecretKey.
// The environment (sandbox vs production) is automatically detected from the key prefixes.
//
// Example:
//
//	client, err := intasend.New(
//	    intasend.WithPublishableKey("ISPubKey_test_xxx"),
//	    intasend.WithSecretKey("ISSecretKey_test_xxx"),
//	)
func New(opts ...Option) (*Client, error) {
	c := &Client{
		timeout:    DefaultTimeout,
		maxRetries: DefaultMaxRetries,
		retryWait:  DefaultRetryWait,
		userAgent:  fmt.Sprintf("intasend-go/%s", Version),
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// Validate that at least one key is provided
	if c.publishableKey == "" && c.secretKey == "" {
		return nil, ErrNoKeysProvided
	}

	// Auto-detect environment if not explicitly set
	if c.baseURL == "" {
		c.detectEnvironment()
	}

	// Validate environment was detected
	if c.baseURL == "" {
		return nil, ErrInvalidEnvironment
	}

	// Create HTTP client if not provided
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: c.timeout,
		}
	}

	// Initialize services eagerly (they are lightweight, holding only a client pointer).
	c.collection = &CollectionService{client: c}
	c.payout = &PayoutService{client: c}
	c.wallet = &WalletService{client: c}
	c.refund = &RefundService{client: c}
	c.checkout = &CheckoutService{client: c}
	c.paymentLink = &PaymentLinkService{client: c}

	return c, nil
}

// detectEnvironment sets the base URL based on the API key prefixes.
func (c *Client) detectEnvironment() {
	// Check publishable key
	if strings.HasPrefix(c.publishableKey, "ISPubKey_test") {
		c.baseURL = SandboxBaseURL
		return
	}
	if strings.HasPrefix(c.publishableKey, "ISPubKey_live") {
		c.baseURL = ProductionBaseURL
		return
	}

	// Check secret key
	if strings.HasPrefix(c.secretKey, "ISSecretKey_test") {
		c.baseURL = SandboxBaseURL
		return
	}
	if strings.HasPrefix(c.secretKey, "ISSecretKey_live") {
		c.baseURL = ProductionBaseURL
		return
	}
}

// Collection returns the collection service for payment collection operations.
func (c *Client) Collection() *CollectionService { return c.collection }

// Payout returns the payout service for send money operations.
func (c *Client) Payout() *PayoutService { return c.payout }

// Wallet returns the wallet service for wallet management.
func (c *Client) Wallet() *WalletService { return c.wallet }

// Refund returns the refund/chargeback service.
func (c *Client) Refund() *RefundService { return c.refund }

// Checkout returns the checkout service for creating checkout pages.
func (c *Client) Checkout() *CheckoutService { return c.checkout }

// PaymentLink returns the payment link service.
func (c *Client) PaymentLink() *PaymentLinkService { return c.paymentLink }

// PublishableKey returns the client's publishable key.
func (c *Client) PublishableKey() string {
	return c.publishableKey
}

// BaseURL returns the client's base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// IsSandbox returns true if the client is configured for the sandbox environment.
func (c *Client) IsSandbox() bool {
	return c.baseURL == SandboxBaseURL
}

// IsProduction returns true if the client is configured for the production environment.
func (c *Client) IsProduction() bool {
	return c.baseURL == ProductionBaseURL
}
