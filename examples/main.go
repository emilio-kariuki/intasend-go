// Package main provides usage examples for the IntaSend Go SDK.
//
// To run these examples, set the following environment variables:
//
//	export INTASEND_PUBLISHABLE_KEY="ISPubKey_test_xxx"
//	export INTASEND_SECRET_KEY="ISSecretKey_test_xxx"
//
// Then run:
//
//	go run examples/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/intasend/intasend-go"
)

func main() {
	// Initialize the client
	client, err := intasend.New(
		intasend.WithPublishableKey(os.Getenv("INTASEND_PUBLISHABLE_KEY")),
		intasend.WithSecretKey(os.Getenv("INTASEND_SECRET_KEY")),
		intasend.WithDebug(true), // Enable debug logging
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("Client initialized (Sandbox: %v)\n\n", client.IsSandbox())

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run examples (comment/uncomment as needed)
	exampleMPesaSTKPush(ctx, client)
	// exampleCreateCheckout(ctx, client)
	// exampleListWallets(ctx, client)
	// exampleSendMoney(ctx, client)
	// examplePaymentLink(ctx, client)
}

// exampleMPesaSTKPush demonstrates M-Pesa STK Push payment collection.
func exampleMPesaSTKPush(ctx context.Context, client *intasend.Client) {
	fmt.Println("=== M-Pesa STK Push Example ===")

	resp, err := client.Collection().MPesaSTKPush(ctx, &intasend.STKPushRequest{
		PhoneNumber: "254712345678", // Replace with test phone number
		Amount:      10,
		APIRef:      "test-order-001",
		Name:        "John Doe",
		Email:       "john@example.com",
	})
	if err != nil {
		// Check if it's an API error
		if apiErr := intasend.AsAPIError(err); apiErr != nil {
			fmt.Printf("API Error: %s (HTTP %d)\n", apiErr.Message, apiErr.HTTPStatusCode)
			return
		}
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Invoice ID: %s\n", resp.Invoice.InvoiceID)
	fmt.Printf("State: %s\n", resp.Invoice.State)
	fmt.Printf("Amount: %.2f\n", resp.Invoice.Value)
	fmt.Println()
}

// exampleCreateCheckout demonstrates creating a checkout page.
func exampleCreateCheckout(ctx context.Context, client *intasend.Client) {
	fmt.Println("=== Create Checkout Example ===")

	resp, err := client.Checkout().Create(ctx, &intasend.CreateCheckoutRequest{
		Amount:   1000,
		Currency: "KES",
		Customer: intasend.CheckoutCustomer{
			Email:     "customer@example.com",
			FirstName: "Jane",
			LastName:  "Doe",
			Country:   "KE",
		},
		Host:        "https://example.com",
		RedirectURL: "https://example.com/callback",
		APIRef:      "order-12345",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Checkout ID: %s\n", resp.ID)
	fmt.Printf("Checkout URL: %s\n", resp.URL)
	fmt.Printf("Signature: %s\n", resp.Signature)
	fmt.Println()
}

// exampleListWallets demonstrates listing wallets.
func exampleListWallets(ctx context.Context, client *intasend.Client) {
	fmt.Println("=== List Wallets Example ===")

	resp, err := client.Wallet().List(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d wallets:\n", len(resp.Results))
	for _, w := range resp.Results {
		fmt.Printf("  - %s (%s): %.2f %s available\n",
			w.Label, w.WalletID, w.AvailableBalance, w.Currency)
	}
	fmt.Println()
}

// exampleSendMoney demonstrates sending money via M-Pesa.
func exampleSendMoney(ctx context.Context, client *intasend.Client) {
	fmt.Println("=== Send Money (M-Pesa B2C) Example ===")

	// Step 1: Initiate the payout
	resp, err := client.Payout().MPesa(ctx, &intasend.MPesaRequest{
		Currency: "KES",
		Transactions: []intasend.Transaction{
			{
				Account:   "254712345678", // Recipient phone number
				Amount:    "100",
				Narrative: "Test payment",
			},
		},
		RequiresApproval: intasend.ApprovalRequired,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Tracking ID: %s\n", resp.TrackingID)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Nonce: %s\n", resp.Nonce)

	// Step 2: Approve the payout (when RequiresApproval is YES)
	approved, err := client.Payout().Approve(ctx, &intasend.ApproveRequest{
		TrackingID: resp.TrackingID,
		Nonce:      resp.Nonce,
	})
	if err != nil {
		fmt.Printf("Approval Error: %v\n", err)
		return
	}

	fmt.Printf("Approved Status: %s\n", approved.Status)
	fmt.Println()
}

// examplePaymentLink demonstrates creating a payment link.
func examplePaymentLink(ctx context.Context, client *intasend.Client) {
	fmt.Println("=== Create Payment Link Example ===")

	link, err := client.PaymentLink().Create(ctx, &intasend.CreatePaymentLinkRequest{
		Title:        "Premium Subscription",
		Currency:     "KES",
		Amount:       2500,
		MobileTariff: intasend.TariffBusinessPays,
		CardTariff:   intasend.TariffBusinessPays,
		IsActive:     true,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Link ID: %s\n", link.LinkID)
	fmt.Printf("Title: %s\n", link.Title)
	fmt.Printf("URL: %s\n", link.URL)
	fmt.Printf("Amount: %.2f %s\n", link.Amount, link.Currency)
	fmt.Println()
}
