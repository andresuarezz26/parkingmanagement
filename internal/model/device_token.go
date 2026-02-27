package model

import "time"

type DeviceToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	Platform  string    `json:"platform"` // ios | android
	CreatedAt time.Time `json:"created_at"`
}
