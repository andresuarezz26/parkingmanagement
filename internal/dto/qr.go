package dto

type QRResponse struct {
	ID            string  `json:"id"`
	VehicleID     string  `json:"vehicle_id"`
	CodeData      string  `json:"code_data"`
	ImageURL      string  `json:"image_url"`
	Status        string  `json:"status"`
	IssuedAt      string  `json:"issued_at"`
	ExpiresAt     *string `json:"expires_at"`
	LastScannedAt *string `json:"last_scanned_at"`
}
