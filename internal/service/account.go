package service

import (
	"context"
	"errors"

	"github.com/andresuarezz26/parkingmanagement/internal/dto"
	"github.com/andresuarezz26/parkingmanagement/internal/model"
	"github.com/andresuarezz26/parkingmanagement/internal/repository"
)

type AccountService struct {
	accountRepo *repository.AccountRepo
	vehicleSvc  *VehicleService
}

func NewAccountService(accountRepo *repository.AccountRepo, vehicleSvc *VehicleService) *AccountService {
	return &AccountService{accountRepo: accountRepo, vehicleSvc: vehicleSvc}
}

func (s *AccountService) GetByUserID(ctx context.Context, userID string) (*model.Account, error) {
	return s.accountRepo.GetByUserID(ctx, userID)
}

func (s *AccountService) Setup(ctx context.Context, userID string, req dto.AccountSetupRequest) (*model.Account, *model.Vehicle, error) {
	// Check if user already has an account
	existing, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if existing != nil {
		return nil, nil, errors.New("user already has an account")
	}

	// Validate
	if req.Name == "" {
		return nil, nil, errors.New("name is required")
	}
	if req.AccountType == "" {
		req.AccountType = "individual"
	}
	if req.AccountType != "individual" && req.AccountType != "company" {
		return nil, nil, errors.New("account_type must be 'individual' or 'company'")
	}

	// Create account
	account, err := s.accountRepo.Create(ctx, req.AccountType, req.Name, req.Email, req.Phone)
	if err != nil {
		return nil, nil, err
	}

	// Link user to account as account_holder
	if err := s.accountRepo.LinkUser(ctx, userID, account.ID, "account_holder"); err != nil {
		return nil, nil, err
	}

	// Optionally register first vehicle
	var vehicle *model.Vehicle
	if req.Vehicle != nil {
		vehicle, err = s.vehicleSvc.Create(ctx, account.ID, *req.Vehicle)
		if err != nil {
			return account, nil, err
		}
	}

	return account, vehicle, nil
}

func (s *AccountService) Update(ctx context.Context, accountID string, req dto.AccountUpdateRequest) (*model.Account, error) {
	fields := make(map[string]interface{})
	if req.Name != nil {
		fields["name"] = *req.Name
	}
	if req.Phone != nil {
		fields["phone"] = *req.Phone
	}
	if req.BillingAddress != nil {
		fields["billing_address"] = *req.BillingAddress
	}
	if req.TaxID != nil {
		fields["tax_id"] = *req.TaxID
	}
	if len(fields) == 0 {
		return nil, errors.New("no fields to update")
	}
	return s.accountRepo.Update(ctx, accountID, fields)
}
