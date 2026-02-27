package model

import "time"

type DayPass struct {
	ID           string    `json:"id"`
	LotID        string    `json:"lot_id"`
	VisitorName  string    `json:"visitor_name"`
	VisitorPhone string    `json:"visitor_phone"`
	VisitorEmail string    `json:"visitor_email"`
	QRCodeID     *string   `json:"qr_code_id"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidUntil   time.Time `json:"valid_until"`
	Status       string    `json:"status"` // active | used | expired
	AmountPaid   float64   `json:"amount_paid"`
}
