package model

import "time"

type Incident struct {
	ID          string     `json:"id"`
	LotID       string     `json:"lot_id"`
	GateID      *string    `json:"gate_id"`
	VehicleID   *string    `json:"vehicle_id"`
	ReportedBy  *string    `json:"reported_by"`
	Type        string     `json:"type"` // unauthorized_entry | tailgating | gate_fault | scan_failure | overstay
	Description string     `json:"description"`
	ImageURL    *string    `json:"image_url"`
	Status      string     `json:"status"` // open | investigating | resolved
	OccurredAt  time.Time  `json:"occurred_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
}
