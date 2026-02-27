package model

import "time"

type Invoice struct {
	ID                 string     `json:"id"`
	AccountID          string     `json:"account_id"`
	SubscriptionID     *string    `json:"subscription_id"`
	InvoiceNumber      string     `json:"invoice_number"`
	BillingPeriodStart string     `json:"billing_period_start"`
	BillingPeriodEnd   string     `json:"billing_period_end"`
	Subtotal           float64    `json:"subtotal"`
	Tax                float64    `json:"tax"`
	Total              float64    `json:"total"`
	Status             string     `json:"status"` // draft | issued | paid | overdue | void
	PDFURL             *string    `json:"pdf_url"`
	IssuedAt           *time.Time `json:"issued_at"`
	DueAt              *time.Time `json:"due_at"`
	PaidAt             *time.Time `json:"paid_at"`
}
