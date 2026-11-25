package novinpay

type PaymentRequest struct {
	MerchantID     string `json:"merchant_id"`
	Amount         int64  `json:"amount"`
	CallbackURL    string `json:"callback_url"`
	CallbackMethod string `json:"callback_method,omitzero"`
	InvoiceID      string `json:"invoice_id,omitzero"`
	Description    string `json:"description,omitzero"`
	Email          string `json:"email,omitzero"`
	Mobile         string `json:"mobile,omitzero"`
	Name           string `json:"name,omitzero"`
	CardPAN        string `json:"card_pan,omitzero"`
}


type PaymentResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    *PaymentData `json:"data"`
	Errors  *[]any        `json:"errors"`
}


type PaymentData struct {
	Wage       int64  `json:"wage"`
	WagePayer  string `json:"wage_payer"`
	Authority  string `json:"authority"`
	TransID    int64  `json:"trans_id"`
	PaymentURL string `json:"payment_url"`
}

type VerificationRequest struct {
	MerchantID string `json:"merchant_id"`
	Amount     int64  `json:"amount"`
	Authority  string `json:"authority"`
}

type VerificationResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    *VerificationData `json:"data"`
	Errors  *[]any            `json:"errors"`
}

type VerificationData struct {
	TransID     int64   `json:"trans_id"`
	RefID       string  `json:"ref_id"`
	Authority   string  `json:"authority"`
	CardPAN     string  `json:"card_pan"`
	Amount      int64   `json:"amount"`
	InvoiceID   *string `json:"invoice_id"`
	BuyerIP     string  `json:"buyer_ip"`
	PaymentTime int64   `json:"payment_time"`
}