package model

import "time"

type Vehicle struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	Make        string    `json:"make"`
	Model       string    `json:"model"`
	Year        string    `json:"year"`
	PlateNumber string    `json:"plate_number"`
	VehicleType string    `json:"vehicle_type"` // semi_truck | tanker | flatbed | construction | mining | other
	Description string    `json:"description"`
	Status      string    `json:"status"` // active | suspended | removed
	CreatedAt   time.Time `json:"created_at"`
}
