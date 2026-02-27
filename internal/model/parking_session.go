package model

import "time"

type ParkingSession struct {
	ID                string     `json:"id"`
	VehicleID         string     `json:"vehicle_id"`
	LotID             string     `json:"lot_id"`
	ZoneID            *string    `json:"zone_id"`
	EntryGateID       string     `json:"entry_gate_id"`
	ExitGateID        *string    `json:"exit_gate_id"`
	SubscriptionID    *string    `json:"subscription_id"`
	EntryTime         time.Time  `json:"entry_time"`
	ExitTime          *time.Time `json:"exit_time"`
	DurationMinutes   *int       `json:"duration_minutes"`
	BaseAmount        float64    `json:"base_amount"`
	OverstayAmount    float64    `json:"overstay_amount"`
	TotalAmount       float64    `json:"total_amount"`
	Status            string     `json:"status"` // active | overstay | completed | disputed
	EntryScanImageURL *string    `json:"entry_scan_image_url"`
	ExitScanImageURL  *string    `json:"exit_scan_image_url"`
}
