// Package novinpay provides a client for the NovinPay IPG v2 payment gateway.
package novinpay

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json/v2"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	// Endpoint is the NovinPay IPG v2 payment request API URL.
	Endpoint = "https://api.novinopay.com/payment/ipg/v2/request"
	// VerifyEndpoint is the NovinPay IPG v2 payment verification API URL.
	VerifyEndpoint = "https://api.novinopay.com/payment/ipg/v2/verification"
)

// NovinPay is a client for the NovinPay IPG v2 payment gateway.
type NovinPay struct {
	MerchantID  string       // Merchant identifier ("test" for sandbox).
	CallbackURL string       // URL NovinPay redirects back to after payment.
	HTTPClient  *http.Client // HTTP client used for requests; falls back to http.DefaultClient if nil.
}

// NewNovinPay creates a NovinPay client with the given merchant ID, callback URL,
// and a configured HTTP client (timeouts, connection pooling, and TLS 1.2+).
func NewNovinPay(merchantID, callbackURL string) *NovinPay {
	return &NovinPay{
		MerchantID:  merchantID,
		CallbackURL: callbackURL,
		HTTPClient:  newHTTPClient(),
	}
}

func newHTTPClient() *http.Client {
	tr := http.DefaultTransport.(*http.Transport).Clone()

	tr.DialContext = (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext

	tr.TLSHandshakeTimeout = 5 * time.Second
	tr.ResponseHeaderTimeout = 5 * time.Second
	tr.IdleConnTimeout = 90 * time.Second

	tr.MaxIdleConnsPerHost = 50
	tr.MaxConnsPerHost = 100

	tr.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	tr.ForceAttemptHTTP2 = true

	return &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}
}

// CreatePayment sends a payment initiation request to NovinPay.
// It overwrites paymentReq.MerchantID and paymentReq.CallbackURL with the client's values.
func (np *NovinPay) CreatePayment(ctx context.Context, paymentReq *PaymentRequest) (*PaymentResponse, error) {
	if paymentReq == nil {
		return nil, fmt.Errorf("payment request cannot be nil")
	}

	paymentReq.MerchantID = np.MerchantID
	paymentReq.CallbackURL = np.CallbackURL

	var paymentResp PaymentResponse
	if err := np.postJSON(ctx, Endpoint, paymentReq, &paymentResp); err != nil {
		return nil, err
	}

	return &paymentResp, nil
}

// VerifyPayment verifies a completed payment with NovinPay.
// It overwrites vr.MerchantID with the client's MerchantID.
func (np *NovinPay) VerifyPayment(ctx context.Context, vr *VerificationRequest) (*VerificationResponse, error) {
	if vr == nil {
		return nil, fmt.Errorf("verification request cannot be nil")
	}

	vr.MerchantID = np.MerchantID

	var verificationResp VerificationResponse
	if err := np.postJSON(ctx, VerifyEndpoint, vr, &verificationResp); err != nil {
		return nil, err
	}

	return &verificationResp, nil
}

func (np *NovinPay) postJSON(ctx context.Context, url string, reqPayload any, respTarget any) error {
	var body bytes.Buffer
	if err := json.MarshalWrite(&body, reqPayload); err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	httpClient := np.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if err := json.UnmarshalRead(resp.Body, respTarget); err != nil {
		return fmt.Errorf("failed to decode response (status %d): %w", resp.StatusCode, err)
	}

	return nil
}
