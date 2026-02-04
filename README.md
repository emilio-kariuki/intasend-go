# IntaSend Go SDK

Official Go SDK for the [IntaSend](https://intasend.com) Payment Gateway API.

IntaSend is an African payment gateway (Kenya-focused) that enables businesses to receive and disburse payments through mobile money (M-Pesa), cards, and bank transfers.

## Installation

```bash
go get github.com/intasend/intasend-go
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/intasend/intasend-go"
)

func main() {
    // Initialize the client
    client, err := intasend.New(
        intasend.WithPublishableKey("ISPubKey_test_xxx"),
        intasend.WithSecretKey("ISSecretKey_test_xxx"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Collect payment via M-Pesa STK Push
    resp, err := client.Collection().MPesaSTKPush(ctx, &intasend.STKPushRequest{
        PhoneNumber: "254712345678",
        Amount:      100,
        APIRef:      "order-123",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Invoice ID: %s, State: %s", resp.Invoice.InvoiceID, resp.Invoice.State)
}
```

## Features

- **Payment Collection**: M-Pesa STK Push, checkout pages, payment status
- **Payouts**: M-Pesa B2C/B2B, bank transfers (PesaLink), airtime top-ups
- **Wallet Management**: Create, list, fund wallets, intra-wallet transfers
- **Refunds**: Create and manage chargebacks
- **Payment Links**: Create shareable payment links

## Configuration Options

```go
client, err := intasend.New(
    // Required: At least one key must be provided
    intasend.WithPublishableKey("ISPubKey_test_xxx"),
    intasend.WithSecretKey("ISSecretKey_test_xxx"),

    // Optional: Override auto-detected environment
    intasend.WithSandbox(),      // Force sandbox
    intasend.WithProduction(),   // Force production

    // Optional: Custom HTTP settings
    intasend.WithTimeout(60 * time.Second),
    intasend.WithHTTPClient(customClient),
    intasend.WithRetry(5, 2*time.Second),

    // Optional: Debug logging
    intasend.WithDebug(true),
)
```

### Environment Detection

The SDK automatically detects the environment from your API key prefixes:
- Keys starting with `ISPubKey_test` or `ISSecretKey_test` → Sandbox
- Keys starting with `ISPubKey_live` or `ISSecretKey_live` → Production

## Services

### Collection Service

Accept payments from customers.

```go
// M-Pesa STK Push - triggers payment prompt on customer's phone
resp, err := client.Collection().MPesaSTKPush(ctx, &intasend.STKPushRequest{
    PhoneNumber: "254712345678",
    Amount:      100,
    APIRef:      "order-123",
    Name:        "John Doe",
    Email:       "john@example.com",
})

// Create checkout page
resp, err := client.Collection().Charge(ctx, &intasend.ChargeRequest{
    Email:    "customer@example.com",
    Host:     "https://yoursite.com",
    Amount:   1000,
    Currency: "KES",
    APIRef:   "order-456",
})

// Check payment status
status, err := client.Collection().Status(ctx, "INV-12345", nil)
```

### Payout Service

Send money to customers, businesses, or buy airtime.

```go
// M-Pesa B2C (send to consumer)
resp, err := client.Payout().MPesa(ctx, &intasend.MPesaRequest{
    Currency: "KES",
    Transactions: []intasend.Transaction{
        {Account: "254712345678", Amount: "1000", Narrative: "Salary payment"},
    },
    RequiresApproval: intasend.ApprovalRequired,
})

// Approve the payout
approved, err := client.Payout().Approve(ctx, &intasend.ApproveRequest{
    TrackingID: resp.TrackingID,
    Nonce:      resp.Nonce,
})

// M-Pesa B2B (send to PayBill/Till)
resp, err := client.Payout().MPesaB2B(ctx, &intasend.MPesaB2BRequest{
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

// Bank transfer via PesaLink
resp, err := client.Payout().Bank(ctx, &intasend.BankRequest{
    Currency: "KES",
    Transactions: []intasend.BankTransaction{
        {
            Name:      "John Doe",
            Account:   "0123456789",
            BankCode:  "2",  // KCB
            Amount:    "10000",
            Narrative: "Supplier payment",
        },
    },
})

// Airtime top-up
resp, err := client.Payout().Airtime(ctx, &intasend.AirtimeRequest{
    Currency: "KES",
    Transactions: []intasend.Transaction{
        {Account: "254712345678", Amount: "100"},
    },
})

// Check payout status
status, err := client.Payout().Status(ctx, "tracking-id-123")
```

### Wallet Service

Manage your IntaSend wallets.

```go
// List all wallets
wallets, err := client.Wallet().List(ctx)

// Create a new wallet
wallet, err := client.Wallet().Create(ctx, &intasend.CreateWalletRequest{
    Currency:    "KES",
    Label:       "Operations Wallet",
    CanDisburse: true,
})

// Get wallet details
wallet, err := client.Wallet().Get(ctx, "WALLET123")

// Get wallet transactions
txns, err := client.Wallet().Transactions(ctx, "WALLET123")

// Transfer between wallets
result, err := client.Wallet().IntraTransfer(ctx, &intasend.IntraTransferRequest{
    SourceID:      "WALLET123",
    DestinationID: "WALLET456",
    Amount:        1000,
    Narrative:     "Commission transfer",
})

// Fund wallet via M-Pesa
result, err := client.Wallet().FundMPesa(ctx, &intasend.FundMPesaRequest{
    WalletID:    "WALLET123",
    PhoneNumber: "254712345678",
    Amount:      5000,
})
```

### Refund Service

Handle refunds and chargebacks.

```go
// List all chargebacks
chargebacks, err := client.Refund().List(ctx)

// Create a refund
chargeback, err := client.Refund().Create(ctx, &intasend.CreateChargebackRequest{
    Invoice:       "INV-123",
    Amount:        500,
    Reason:        intasend.RefundReasonCustomerRequest,
    ReasonDetails: "Customer requested cancellation",
})

// Get chargeback details
chargeback, err := client.Refund().Get(ctx, "CHG-123")
```

### Payment Link Service

Create shareable payment links.

```go
// List payment links
links, err := client.PaymentLink().List(ctx)

// Create a payment link
link, err := client.PaymentLink().Create(ctx, &intasend.CreatePaymentLinkRequest{
    Title:        "Premium Service",
    Currency:     "KES",
    Amount:       5000,
    MobileTariff: intasend.TariffBusinessPays,
    CardTariff:   intasend.TariffBusinessPays,
    IsActive:     true,
})

// Get payment link details
link, err := client.PaymentLink().Get(ctx, "LINK-123")
```

## Error Handling

The SDK provides structured error types for better error handling:

```go
resp, err := client.Collection().MPesaSTKPush(ctx, req)
if err != nil {
    // Check for API errors
    if apiErr := intasend.AsAPIError(err); apiErr != nil {
        fmt.Printf("API Error: %s (HTTP %d)\n", apiErr.Message, apiErr.HTTPStatusCode)

        if apiErr.IsAuthenticationError() {
            // Handle authentication error
        }
        if apiErr.IsValidationError() {
            // Handle validation errors
            for field, errs := range apiErr.Errors {
                fmt.Printf("  %s: %v\n", field, errs)
            }
        }
        if apiErr.IsRateLimited() {
            // Handle rate limiting
        }
        return
    }

    // Check for network errors
    if intasend.IsNetworkError(err) {
        fmt.Println("Network error - please retry")
        return
    }

    // Other errors
    fmt.Printf("Error: %v\n", err)
}
```

## Testing

The SDK automatically uses the sandbox environment when using test API keys. Get your test keys from [IntaSend Sandbox](https://sandbox.intasend.com).

```go
// Test mode is automatically enabled with test keys
client, _ := intasend.New(
    intasend.WithPublishableKey("ISPubKey_test_xxx"),
    intasend.WithSecretKey("ISSecretKey_test_xxx"),
)

fmt.Println(client.IsSandbox()) // true
```

## API Documentation

For detailed API documentation, visit:
- [IntaSend API Docs](https://developers.intasend.com/docs)
- [API Testing & Sandbox](https://developers.intasend.com/docs/api-testing-and-sandbox)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- [IntaSend Documentation](https://developers.intasend.com/docs)
- [IntaSend Support](https://intasend.com/contact)
- [GitHub Issues](https://github.com/intasend/intasend-go/issues)
