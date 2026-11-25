package main

import (
	"context"
	"fmt"
	"log"
	"time"

	novinpay "github.com/amiramani26/novinPay"
)

func main() {
	// Initialize the NovinPay client with your merchant ID and callback URL.
	client := novinpay.NewNovinPay("test", "https://example.com/callback")

	// Set a 10-second timeout for the payment request.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build the payment request.
	// Note: Setting MerchantID or CallbackURL here has no effect, because
	// CreatePayment overwrites them with the client's values.
	req := &novinpay.PaymentRequest{
		Amount:      100000,
		Description: "Test payment",
	}

	// Send the payment creation request to the NovinPay API.
	resp, err := client.CreatePayment(ctx, req)
	if err != nil {
		log.Fatalf("CreatePayment failed: %v", err)
	}

	// Print the gateway response status and any returned payment or error details.
	fmt.Printf("Status: %s\nMessage: %s\n", resp.Status, resp.Message)
	if resp.Data != nil {
		fmt.Printf("Authority: %s\nTransID: %d\nPaymentURL: %s\n",
			resp.Data.Authority, resp.Data.TransID, resp.Data.PaymentURL)
	}
	if resp.Errors != nil {
		fmt.Printf("Errors: %+v\n", *resp.Errors)
	}
}
