package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andresuarezz26/parkingmanagement/internal/model"
)

type VehicleRepo struct {
	db *pgxpool.Pool
}

func NewVehicleRepo(db *pgxpool.Pool) *VehicleRepo {
	return &VehicleRepo{db: db}
}

func (r *VehicleRepo) GetByID(ctx context.Context, id string) (*model.Vehicle, error) {
	var v model.Vehicle
	err := r.db.QueryRow(ctx,
		`SELECT id, account_id, make, model, year, plate_number, vehicle_type, description, status, created_at
		 FROM vehicles WHERE id = $1`, id).
		Scan(&v.ID, &v.AccountID, &v.Make, &v.Model, &v.Year, &v.PlateNumber, &v.VehicleType, &v.Description, &v.Status, &v.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &v, err
}

func (r *VehicleRepo) ListByAccount(ctx context.Context, accountID string) ([]model.Vehicle, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, account_id, make, model, year, plate_number, vehicle_type, description, status, created_at
		 FROM vehicles WHERE account_id = $1 AND status != 'removed' ORDER BY created_at DESC`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []model.Vehicle
	for rows.Next() {
		var v model.Vehicle
		if err := rows.Scan(&v.ID, &v.AccountID, &v.Make, &v.Model, &v.Year, &v.PlateNumber, &v.VehicleType, &v.Description, &v.Status, &v.CreatedAt); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func (r *VehicleRepo) Create(ctx context.Context, accountID, make_, model_, year, plate, vType, desc string) (*model.Vehicle, error) {
	var v model.Vehicle
	err := r.db.QueryRow(ctx,
		`INSERT INTO vehicles (account_id, make, model, year, plate_number, vehicle_type, description)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, account_id, make, model, year, plate_number, vehicle_type, description, status, created_at`,
		accountID, make_, model_, year, plate, vType, desc).
		Scan(&v.ID, &v.AccountID, &v.Make, &v.Model, &v.Year, &v.PlateNumber, &v.VehicleType, &v.Description, &v.Status, &v.CreatedAt)
	return &v, err
}

func (r *VehicleRepo) Update(ctx context.Context, id string, fields map[string]interface{}) (*model.Vehicle, error) {
	// Build dynamic update (same pattern as account)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	for k, val := range fields {
		_, err := tx.Exec(ctx, `UPDATE vehicles SET `+k+` = $1 WHERE id = $2`, val, id)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *VehicleRepo) SetStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE vehicles SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (r *VehicleRepo) CountByAccount(ctx context.Context, accountID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM vehicles WHERE account_id = $1 AND status = 'active'`, accountID).Scan(&count)
	return count, err
}
