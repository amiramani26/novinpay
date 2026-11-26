# novinpay

A simple Go client for the **NovinPay IPG v2** API.

## Requirements

- Go `1.27.4` or later (uses `encoding/json/v2`)

## Installation

```bash
go get github.com/amiramani26/novinpay
```

## Quick Start

### 1. Create a payment

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/amiramani26/novinpay"
)

func main() {
	client := novinpay.NewNovinPay("test", "https://example.com/callback")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.CreatePayment(ctx, &novinpay.PaymentRequest{
		Amount:      100000, // Amount in Rials
		Description: "Order payment",
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.Data == nil {
		log.Fatalf("gateway error (%s): %s", resp.Status, resp.Message)
	}

	fmt.Println("Redirect user to:", resp.Data.PaymentURL)
}
```

### 2. Verify a payment on callback

After the buyer returns to your `CallbackURL`, read `PaymentStatus` and `Authority` from the request. If `PaymentStatus == "OK"`, call `VerifyPayment` with the same `Amount`:

```go
verifyResp, err := client.VerifyPayment(ctx, &novinpay.VerificationRequest{
	Amount:    100000, // Must match the CreatePayment amount (in Rials)
	Authority: authority,
})
if err != nil {
	log.Fatal(err)
}
if verifyResp.Data != nil {
	fmt.Printf("Status: %s | RefID: %s | TransID: %d\n",
		verifyResp.Status, verifyResp.Data.RefID, verifyResp.Data.TransID)
}
```

## Notes

- Amounts are in **Rials** (not Tomans).
- Use `"test"` as `merchantID` for sandbox testing.
- Both `CreatePayment` and `VerifyPayment` automatically set `MerchantID` (and `CreatePayment` also sets `CallbackURL`) from the `NovinPay` client.
- Check `resp.Status` (`"100"` means success; `"101"` on verification means the payment succeeded and was already verified) and verify `resp.Data != nil` before reading response fields.

## Examples

- [Create payment](examples/create_payment/main.go)
- [Full payment & callback verification flow](examples/payment_flow/main.go)

## Documentation

For more details on endpoints, parameters, status codes, and API behavior, visit the [NovinPay docs](https://novinopay.com/docs).
