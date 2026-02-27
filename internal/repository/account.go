package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andresuarezz26/parkingmanagement/internal/model"
)

type AccountRepo struct {
	db *pgxpool.Pool
}

func NewAccountRepo(db *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) GetByID(ctx context.Context, id string) (*model.Account, error) {
	var a model.Account
	err := r.db.QueryRow(ctx,
		`SELECT id, type, name, email, phone, billing_address, tax_id, status, created_at, updated_at
		 FROM accounts WHERE id = $1`, id).
		Scan(&a.ID, &a.Type, &a.Name, &a.Email, &a.Phone, &a.BillingAddress, &a.TaxID, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *AccountRepo) GetByUserID(ctx context.Context, userID string) (*model.Account, error) {
	var a model.Account
	err := r.db.QueryRow(ctx,
		`SELECT a.id, a.type, a.name, a.email, a.phone, a.billing_address, a.tax_id, a.status, a.created_at, a.updated_at
		 FROM accounts a JOIN users u ON u.account_id = a.id WHERE u.id = $1`, userID).
		Scan(&a.ID, &a.Type, &a.Name, &a.Email, &a.Phone, &a.BillingAddress, &a.TaxID, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *AccountRepo) Create(ctx context.Context, accountType, name, email, phone string) (*model.Account, error) {
	var a model.Account
	err := r.db.QueryRow(ctx,
		`INSERT INTO accounts (type, name, email, phone) VALUES ($1, $2, $3, $4)
		 RETURNING id, type, name, email, phone, billing_address, tax_id, status, created_at, updated_at`,
		accountType, name, email, phone).
		Scan(&a.ID, &a.Type, &a.Name, &a.Email, &a.Phone, &a.BillingAddress, &a.TaxID, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	return &a, err
}

func (r *AccountRepo) Update(ctx context.Context, id string, fields map[string]interface{}) (*model.Account, error) {
	setClauses := ""
	args := []interface{}{}
	i := 1
	for k, v := range fields {
		if i > 1 {
			setClauses += ", "
		}
		setClauses += fmt.Sprintf("%s = $%d", k, i)
		args = append(args, v)
		i++
	}
	args = append(args, id)
	setClauses += fmt.Sprintf(", updated_at = now()")

	query := fmt.Sprintf(`UPDATE accounts SET %s WHERE id = $%d
		RETURNING id, type, name, email, phone, billing_address, tax_id, status, created_at, updated_at`, setClauses, i)

	var a model.Account
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&a.ID, &a.Type, &a.Name, &a.Email, &a.Phone, &a.BillingAddress, &a.TaxID, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	return &a, err
}

func (r *AccountRepo) LinkUser(ctx context.Context, userID, accountID, role string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET account_id = $1, role = $2 WHERE id = $3`,
		accountID, role, userID)
	return err
}
