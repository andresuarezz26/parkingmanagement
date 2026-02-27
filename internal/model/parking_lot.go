package model

import "time"

type ParkingLot struct {
	ID                  string    `json:"id"`
	OperatorID          string    `json:"operator_id"`
	Name                string    `json:"name"`
	Address             string    `json:"address"`
	City                string    `json:"city"`
	State               string    `json:"state"`
	Latitude            float64   `json:"latitude"`
	Longitude           float64   `json:"longitude"`
	TotalCapacity       int       `json:"total_capacity"`
	HeavyVehicleSpaces  int       `json:"heavy_vehicle_spaces"`
	OperatingHours      string    `json:"operating_hours"` // JSON
	Status              string    `json:"status"` // open | closed | maintenance
	CreatedAt           time.Time `json:"created_at"`
}
