package service

import (
	"context"
	"errors"

	"github.com/andresuarezz26/parkingmanagement/internal/dto"
	"github.com/andresuarezz26/parkingmanagement/internal/model"
	"github.com/andresuarezz26/parkingmanagement/internal/repository"
)

type VehicleService struct {
	vehicleRepo *repository.VehicleRepo
	qrSvc       *QRService
}

func NewVehicleService(vehicleRepo *repository.VehicleRepo, qrSvc *QRService) *VehicleService {
	return &VehicleService{vehicleRepo: vehicleRepo, qrSvc: qrSvc}
}

func (s *VehicleService) ListByAccount(ctx context.Context, accountID string) ([]model.Vehicle, error) {
	return s.vehicleRepo.ListByAccount(ctx, accountID)
}

func (s *VehicleService) GetByID(ctx context.Context, id string) (*model.Vehicle, error) {
	return s.vehicleRepo.GetByID(ctx, id)
}

func (s *VehicleService) Create(ctx context.Context, accountID string, req dto.VehicleCreate) (*model.Vehicle, error) {
	if req.Make == "" || req.Model == "" {
		return nil, errors.New("make and model are required")
	}
	if req.VehicleType == "" {
		req.VehicleType = "other"
	}

	vehicle, err := s.vehicleRepo.Create(ctx, accountID, req.Make, req.Model, req.Year, req.PlateNumber, req.VehicleType, req.Description)
	if err != nil {
		return nil, err
	}

	// Auto-generate QR code
	_, err = s.qrSvc.Generate(ctx, vehicle.ID, accountID)
	if err != nil {
		// Vehicle created but QR failed — log but don't fail
		return vehicle, nil
	}

	return vehicle, nil
}

func (s *VehicleService) Update(ctx context.Context, id string, req dto.VehicleUpdate) (*model.Vehicle, error) {
	fields := make(map[string]interface{})
	if req.Make != nil {
		fields["make"] = *req.Make
	}
	if req.Model != nil {
		fields["model"] = *req.Model
	}
	if req.Year != nil {
		fields["year"] = *req.Year
	}
	if req.PlateNumber != nil {
		fields["plate_number"] = *req.PlateNumber
	}
	if req.VehicleType != nil {
		fields["vehicle_type"] = *req.VehicleType
	}
	if req.Description != nil {
		fields["description"] = *req.Description
	}
	if len(fields) == 0 {
		return nil, errors.New("no fields to update")
	}
	return s.vehicleRepo.Update(ctx, id, fields)
}

func (s *VehicleService) Delete(ctx context.Context, id string) error {
	// Revoke QR first
	if err := s.qrSvc.RevokeByVehicle(ctx, id); err != nil {
		return err
	}
	return s.vehicleRepo.SetStatus(ctx, id, "removed")
}
