package tests

import (
	"net/http/httptest"
	"testing"

	intasend "github.com/emilio-kariuki/intasend-go"
)

// newTestClient creates a Client pointed at the given httptest.Server.
func newTestClient(t *testing.T, server *httptest.Server) *intasend.Client {
	t.Helper()
	client, err := intasend.New(
		intasend.WithPublishableKey("ISPubKey_test_abc123"),
		intasend.WithSecretKey("ISSecretKey_test_secret"),
		intasend.WithBaseURL(server.URL),
		intasend.WithHTTPClient(server.Client()),
		intasend.WithRetry(0, 0), // no retries by default in tests
	)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}
	return client
}

// Local struct definitions that mirror internal JSON request bodies.
// Used in test server handlers to decode and validate request payloads.

type chargeRequestBody struct {
	PublicKey    string  `json:"public_key"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	PhoneNumber  string  `json:"phone_number"`
	Host         string  `json:"host"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	APIRef       string  `json:"api_ref"`
	RedirectURL  string  `json:"redirect_url"`
	Comment      string  `json:"comment"`
	Method       string  `json:"method"`
	WalletID     string  `json:"wallet_id"`
	CardTariff   string  `json:"card_tarrif"`
	MobileTariff string  `json:"mobile_tarrif"`
	Country      string  `json:"country"`
	Address      string  `json:"address"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	Zipcode      string  `json:"zipcode"`
}

type stkPushRequestBody struct {
	PublicKey   string  `json:"public_key"`
	PhoneNumber string  `json:"phone_number"`
	Amount      float64 `json:"amount"`
	APIRef      string  `json:"api_ref"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	WalletID    string  `json:"wallet_id"`
	Method      string  `json:"method"`
	Currency    string  `json:"currency"`
}

type statusRequestBody struct {
	InvoiceID  string `json:"invoice_id"`
	PublicKey  string `json:"public_key"`
	CheckoutID string `json:"checkout_id"`
	Signature  string `json:"signature"`
}

type payoutStatusRequestBody struct {
	TrackingID string `json:"tracking_id"`
}

type intraTransferRequestBody struct {
	WalletID  string  `json:"wallet_id"`
	Amount    float64 `json:"amount"`
	Narrative string  `json:"narrative"`
}

type fundMPesaRequestBody struct {
	PublicKey   string  `json:"public_key"`
	WalletID    string  `json:"wallet_id"`
	PhoneNumber string  `json:"phone_number"`
	Amount      float64 `json:"amount"`
	Email       string  `json:"email"`
	APIRef      string  `json:"api_ref"`
	Method      string  `json:"method"`
	Currency    string  `json:"currency"`
}

type fundCheckoutRequestBody struct {
	PublicKey    string  `json:"public_key"`
	WalletID    string  `json:"wallet_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Email       string  `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber string  `json:"phone_number"`
	Country     string  `json:"country"`
	Host        string  `json:"host"`
	RedirectURL string  `json:"redirect_url"`
	APIRef      string  `json:"api_ref"`
	CardTariff  string  `json:"card_tarrif"`
	MobileTariff string `json:"mobile_tarrif"`
}

type createCheckoutRequestBody struct {
	PublicKey    string  `json:"public_key"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Email       string  `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber string  `json:"phone_number"`
	Country     string  `json:"country"`
	Address     string  `json:"address"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Zipcode     string  `json:"zipcode"`
	Host        string  `json:"host"`
	RedirectURL string  `json:"redirect_url"`
	APIRef      string  `json:"api_ref"`
	Comment     string  `json:"comment"`
	Method      string  `json:"method"`
	CardTariff  string  `json:"card_tarrif"`
	MobileTariff string `json:"mobile_tarrif"`
	WalletID    string  `json:"wallet_id"`
}
