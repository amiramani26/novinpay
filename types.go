package novinpay

// PaymentRequest contains parameters for initiating a payment via CreatePayment.
type PaymentRequest struct {
	MerchantID     string `json:"merchant_id"`              // Overwritten by NovinPay.MerchantID.
	Amount         int64  `json:"amount"`                   // Payment amount in Rials.
	CallbackURL    string `json:"callback_url"`             // Overwritten by NovinPay.CallbackURL.
	CallbackMethod string `json:"callback_method,omitzero"` // Optional: "GET" (default) or "POST".
	InvoiceID      string `json:"invoice_id,omitzero"`      // Optional merchant invoice ID.
	Description    string `json:"description,omitzero"`     // Optional payment description.
	Email          string `json:"email,omitzero"`           // Optional customer email.
	Mobile         string `json:"mobile,omitzero"`          // Optional customer mobile number.
	Name           string `json:"name,omitzero"`            // Optional customer name.
	CardPAN        string `json:"card_pan,omitzero"`        // Optional restricted card number.
}

// PaymentResponse represents the API response returned by CreatePayment.
type PaymentResponse struct {
	Status  string       `json:"status"`  // Gateway status code.
	Message string       `json:"message"` // Gateway status message.
	Data    *PaymentData `json:"data"`    // Payment session details, nil on failure.
	Errors  *[]any       `json:"errors"`  // Optional validation error details.
}

// PaymentData holds the created payment session details returned by CreatePayment.
type PaymentData struct {
	Wage       int64  `json:"wage"`        // Transaction fee in Rials.
	WagePayer  string `json:"wage_payer"`  // Fee payer (e.g., "merchant").
	Authority  string `json:"authority"`   // Unique token passed later to VerifyPayment.
	TransID    int64  `json:"trans_id"`    // NovinPay transaction ID.
	PaymentURL string `json:"payment_url"` // Gateway URL to redirect the buyer to.
}

// VerificationRequest contains parameters for verifying a completed payment.
type VerificationRequest struct {
	MerchantID string `json:"merchant_id"` // Overwritten by NovinPay.MerchantID.
	Amount     int64  `json:"amount"`      // Expected amount in Rials (must match CreatePayment).
	Authority  string `json:"authority"`   // Authority token received from the callback.
}

// VerificationResponse represents the API response returned by VerifyPayment.
type VerificationResponse struct {
	Status  string            `json:"status"`  // Gateway status code.
	Message string            `json:"message"` // Gateway status message.
	Data    *VerificationData `json:"data"`    // Verified transaction details, nil on failure.
	Errors  *[]any            `json:"errors"`  // Optional validation error details.
}

// VerificationData holds the settled transaction details returned by VerifyPayment.
type VerificationData struct {
	TransID     int64   `json:"trans_id"`     // NovinPay transaction ID.
	RefID       string  `json:"ref_id"`       // Bank reference tracking ID.
	Authority   string  `json:"authority"`    // Transaction authority token.
	CardPAN     string  `json:"card_pan"`     // Masked payer card number.
	Amount      int64   `json:"amount"`       // Verified amount in Rials.
	InvoiceID   *string `json:"invoice_id"`   // Merchant invoice ID, if provided.
	BuyerIP     string  `json:"buyer_ip"`     // Payer IP address.
	PaymentTime int64   `json:"payment_time"` // Payment Unix timestamp.
}
