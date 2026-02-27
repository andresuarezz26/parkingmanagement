package model

import "time"

type Payment struct {
	ID                   string     `json:"id"`
	AccountID            string     `json:"account_id"`
	SessionID            *string    `json:"session_id"`
	InvoiceID            *string    `json:"invoice_id"`
	Amount               float64    `json:"amount"`
	Currency             string     `json:"currency"`
	PaymentType          string     `json:"payment_type"` // subscription | session | overstay | day_pass
	Status               string     `json:"status"` // pending | succeeded | failed | refunded
	GatewayTransactionID *string    `json:"gateway_transaction_id"`
	ProcessedAt          *time.Time `json:"processed_at"`
}
