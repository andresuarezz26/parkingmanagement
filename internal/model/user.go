package model

import "time"

type User struct {
	ID           string    `json:"id"`
	AccountID    *string   `json:"account_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"` // account_holder | driver | operator | attendant | super_admin
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time `json:"created_at"`
}
