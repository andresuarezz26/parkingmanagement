package dto

type VehicleCreate struct {
	Make        string `json:"make"`
	Model       string `json:"model"`
	Year        string `json:"year"`
	PlateNumber string `json:"plate_number"`
	VehicleType string `json:"vehicle_type"` // semi_truck | tanker | flatbed | construction | mining | other
	Description string `json:"description"`
}

type VehicleUpdate struct {
	Make        *string `json:"make,omitempty"`
	Model       *string `json:"model,omitempty"`
	Year        *string `json:"year,omitempty"`
	PlateNumber *string `json:"plate_number,omitempty"`
	VehicleType *string `json:"vehicle_type,omitempty"`
	Description *string `json:"description,omitempty"`
}

type VehicleResponse struct {
	ID          string      `json:"id"`
	AccountID   string      `json:"account_id"`
	Make        string      `json:"make"`
	Model       string      `json:"model"`
	Year        string      `json:"year"`
	PlateNumber string      `json:"plate_number"`
	VehicleType string      `json:"vehicle_type"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	CreatedAt   string      `json:"created_at"`
	QRCode      *QRResponse `json:"qr_code,omitempty"`
}
