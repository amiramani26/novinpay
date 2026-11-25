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


const Endpoint = "https://api.novinopay.com/payment/ipg/v2/request"


type NovinPay struct {
	MerchantID  string
	CallbackURL string
	HTTPClient  *http.Client
}

// NewNovinPay creates a NovinPay client with the given merchant ID, callback URL,
// and a configured HTTP client (timeouts, connection pooling, and TLS 1.2+).
func NewNovinPay(MerchantID string, CallbackURL string) *NovinPay {
	return &NovinPay{
		MerchantID:  MerchantID,
		CallbackURL: CallbackURL,
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
	tr.ExpectContinueTimeout = 1 * time.Second
	tr.IdleConnTimeout = 90 * time.Second
	tr.MaxIdleConnsPerHost = 50

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

	var body bytes.Buffer
	if err := json.MarshalWrite(&body, paymentReq); err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, Endpoint, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := np.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var paymentResp PaymentResponse
	if err := json.UnmarshalRead(resp.Body, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to decode response (status %d): %w", resp.StatusCode, err)
	}

	return &paymentResp, nil
}
