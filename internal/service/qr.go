package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"

	"github.com/andresuarezz26/parkingmanagement/internal/model"
	"github.com/andresuarezz26/parkingmanagement/internal/repository"
)

type QRService struct {
	qrRepo *repository.QRRepo
}

func NewQRService(qrRepo *repository.QRRepo) *QRService {
	return &QRService{qrRepo: qrRepo}
}

func (s *QRService) GetByVehicle(ctx context.Context, vehicleID string) (*model.QRCode, error) {
	return s.qrRepo.GetByVehicleID(ctx, vehicleID)
}

func (s *QRService) Generate(ctx context.Context, vehicleID, accountID string) (*model.QRCode, error) {
	// Revoke any existing active QR for this vehicle
	if err := s.qrRepo.RevokeByVehicle(ctx, vehicleID); err != nil {
		return nil, fmt.Errorf("failed to revoke existing QR: %w", err)
	}

	// Generate a signed token as QR payload
	qrID := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"qr_id":      qrID,
		"vehicle_id": vehicleID,
		"account_id": accountID,
	})

	// Sign with a deterministic key (in production, use a proper secret)
	// For now, sign with the qr_id itself to make it unique but verifiable
	codeData, err := token.SignedString([]byte("heavypark-qr-signing-key"))
	if err != nil {
		return nil, fmt.Errorf("failed to sign QR token: %w", err)
	}

	// Generate QR image as base64 data URI
	png, err := qrcode.Encode(codeData, qrcode.Medium, 512)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR image: %w", err)
	}
	imageURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)

	// Store in DB
	qr, err := s.qrRepo.Create(ctx, vehicleID, codeData, imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to store QR: %w", err)
	}

	return qr, nil
}

func (s *QRService) Regenerate(ctx context.Context, vehicleID, accountID string) (*model.QRCode, error) {
	return s.Generate(ctx, vehicleID, accountID)
}

func (s *QRService) RevokeByVehicle(ctx context.Context, vehicleID string) error {
	return s.qrRepo.RevokeByVehicle(ctx, vehicleID)
}
