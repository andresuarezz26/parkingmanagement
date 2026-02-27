package model

import "time"

type QRCode struct {
	ID            string     `json:"id"`
	VehicleID     string     `json:"vehicle_id"`
	CodeData      string     `json:"code_data"`
	ImageURL      string     `json:"image_url"`
	Status        string     `json:"status"` // generated | active | suspended | revoked
	IssuedAt      time.Time  `json:"issued_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	LastScannedAt *time.Time `json:"last_scanned_at"`
}
