package model

import "time"

type Gate struct {
	ID            string    `json:"id"`
	LotID         string    `json:"lot_id"`
	Name          string    `json:"name"`
	HardwareID    string    `json:"hardware_id"`
	GateType      string    `json:"gate_type"` // entry | exit | both
	Status        string    `json:"status"` // online | offline | fault | maintenance
	IPAddress     string    `json:"ip_address"`
	LastHeartbeat *time.Time `json:"last_heartbeat"`
}
