package dto

type AccountSetupRequest struct {
	AccountType string         `json:"account_type"` // individual | company
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	Phone       string         `json:"phone"`
	Vehicle     *VehicleCreate `json:"vehicle,omitempty"`
}

type AccountUpdateRequest struct {
	Name           *string `json:"name,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	BillingAddress *string `json:"billing_address,omitempty"`
	TaxID          *string `json:"tax_id,omitempty"`
}

type AccountResponse struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	BillingAddress string `json:"billing_address"`
	TaxID          string `json:"tax_id"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
}
