package model

type ParkingZone struct {
	ID               string `json:"id"`
	LotID            string `json:"lot_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Capacity         int    `json:"capacity"`
	ZoneType         string `json:"zone_type"` // standard | oversized | hazmat | refrigerated
	CurrentOccupancy int    `json:"current_occupancy"`
}
