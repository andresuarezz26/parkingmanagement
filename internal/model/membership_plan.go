package model

type MembershipPlan struct {
	ID                    string  `json:"id"`
	LotID                 string  `json:"lot_id"`
	Name                  string  `json:"name"`
	Description           string  `json:"description"`
	BillingCycle          string  `json:"billing_cycle"` // day | month | year | per_use
	Price                 float64 `json:"price"`
	MaxVehicles           int     `json:"max_vehicles"`
	MaxConcurrentEntries  int     `json:"max_concurrent_entries"`
	GracePeriodMinutes    int     `json:"grace_period_minutes"`
	OverstayRatePerHour   float64 `json:"overstay_rate_per_hour"`
	IsActive              bool    `json:"is_active"`
}
