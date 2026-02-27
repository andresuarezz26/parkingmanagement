package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andresuarezz26/parkingmanagement/internal/model"
)

type QRRepo struct {
	db *pgxpool.Pool
}

func NewQRRepo(db *pgxpool.Pool) *QRRepo {
	return &QRRepo{db: db}
}

func (r *QRRepo) GetByVehicleID(ctx context.Context, vehicleID string) (*model.QRCode, error) {
	var q model.QRCode
	err := r.db.QueryRow(ctx,
		`SELECT id, vehicle_id, code_data, image_url, status, issued_at, expires_at, last_scanned_at
		 FROM qr_codes WHERE vehicle_id = $1 AND status IN ('generated', 'active')
		 ORDER BY issued_at DESC LIMIT 1`, vehicleID).
		Scan(&q.ID, &q.VehicleID, &q.CodeData, &q.ImageURL, &q.Status, &q.IssuedAt, &q.ExpiresAt, &q.LastScannedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &q, err
}

func (r *QRRepo) GetByCodeData(ctx context.Context, codeData string) (*model.QRCode, error) {
	var q model.QRCode
	err := r.db.QueryRow(ctx,
		`SELECT id, vehicle_id, code_data, image_url, status, issued_at, expires_at, last_scanned_at
		 FROM qr_codes WHERE code_data = $1`, codeData).
		Scan(&q.ID, &q.VehicleID, &q.CodeData, &q.ImageURL, &q.Status, &q.IssuedAt, &q.ExpiresAt, &q.LastScannedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &q, err
}

func (r *QRRepo) Create(ctx context.Context, vehicleID, codeData, imageURL string) (*model.QRCode, error) {
	var q model.QRCode
	err := r.db.QueryRow(ctx,
		`INSERT INTO qr_codes (vehicle_id, code_data, image_url, status)
		 VALUES ($1, $2, $3, 'active')
		 RETURNING id, vehicle_id, code_data, image_url, status, issued_at, expires_at, last_scanned_at`,
		vehicleID, codeData, imageURL).
		Scan(&q.ID, &q.VehicleID, &q.CodeData, &q.ImageURL, &q.Status, &q.IssuedAt, &q.ExpiresAt, &q.LastScannedAt)
	return &q, err
}

func (r *QRRepo) RevokeByVehicle(ctx context.Context, vehicleID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE qr_codes SET status = 'revoked' WHERE vehicle_id = $1 AND status IN ('generated', 'active')`,
		vehicleID)
	return err
}

func (r *QRRepo) SetStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE qr_codes SET status = $1 WHERE id = $2`, status, id)
	return err
}
