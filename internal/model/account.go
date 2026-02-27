package model

import "time"

type Account struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"` // individual | company
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	BillingAddress string    `json:"billing_address"`
	TaxID          string    `json:"tax_id"`
	Status         string    `json:"status"` // active | suspended | cancelled
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
