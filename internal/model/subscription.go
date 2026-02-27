package model

import "time"

type Subscription struct {
	ID              string    `json:"id"`
	AccountID       string    `json:"account_id"`
	PlanID          string    `json:"plan_id"`
	LotID           string    `json:"lot_id"`
	StartDate       string    `json:"start_date"`
	EndDate         *string   `json:"end_date"`
	Status          string    `json:"status"` // pending | active | suspended | expired | cancelled
	PaymentMethodID *string   `json:"payment_method_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
