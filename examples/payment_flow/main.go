package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/amiramani26/novinpay"
)

const amount = int64(100000)

func main() {
	// Initialize the client with the test merchant ID and local callback URL.
	client := novinpay.NewNovinPay("test", "http://localhost:8080/callback")

	// Set a 10-second timeout for creating the payment session.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Step 1: Create the payment and get the gateway redirect URL.
	resp, err := client.CreatePayment(ctx, &novinpay.PaymentRequest{
		Amount:      amount,
		Description: "Test payment",
	})
	if err != nil {
		log.Fatalf("CreatePayment failed: %v", err)
	}
	if resp.Data == nil {
		log.Fatalf("CreatePayment error (status %s): %s", resp.Status, resp.Message)
	}

	// Step 2: Open this URL in a browser to complete (or cancel) the test payment.
	fmt.Println("Open this URL in your browser to pay:")
	fmt.Println(resp.Data.PaymentURL)

	// Step 3: Handle the gateway callback after the user is redirected back.
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// r.FormValue reads parameters from both GET query strings and POST forms.
		status := r.FormValue("PaymentStatus")
		authority := r.FormValue("Authority")

		// Only verify if the gateway reports PaymentStatus == "OK".
		if status != "OK" {
			fmt.Fprintf(w, "Payment not completed (PaymentStatus=%s)\n", status)
			return
		}

		verifyCtx, cancelVerify := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancelVerify()

		// Step 4: Verify the transaction server-to-server; amount must match CreatePayment.
		verifyResp, err := client.VerifyPayment(verifyCtx, &novinpay.VerificationRequest{
			Amount:    amount,
			Authority: authority,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		// Print the verification result
		fmt.Printf("Verified -> Status: %s | Message: %s\n", verifyResp.Status, verifyResp.Message)
		if verifyResp.Data != nil {
			fmt.Printf("RefID: %s | TransID: %d\n", verifyResp.Data.RefID, verifyResp.Data.TransID)
		}
		fmt.Fprintf(w, "Status: %s\nMessage: %s\n", verifyResp.Status, verifyResp.Message)
	})

	// Start the local HTTP server to wait for the callback redirect.
	fmt.Println("Waiting for callback on http://localhost:8080/callback ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
