package handler

import (
	"encoding/json"
	"net/http"

	"github.com/andresuarezz26/parkingmanagement/internal/dto"
	mw "github.com/andresuarezz26/parkingmanagement/internal/middleware"
	"github.com/andresuarezz26/parkingmanagement/internal/service"
)

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

// GET /api/v1/account
func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims := mw.GetUserClaims(r.Context())
	if claims == nil {
		respondErr(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	account, err := h.svc.GetByUserID(r.Context(), claims.UserID)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, "failed to get account")
		return
	}
	if account == nil {
		respondErr(w, http.StatusNotFound, "no account found — use POST /account/setup to create one")
		return
	}

	respondJSON(w, http.StatusOK, dto.AccountResponse{
		ID:             account.ID,
		Type:           account.Type,
		Name:           account.Name,
		Email:          account.Email,
		Phone:          account.Phone,
		BillingAddress: account.BillingAddress,
		TaxID:          account.TaxID,
		Status:         account.Status,
		CreatedAt:      account.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// POST /api/v1/account/setup
func (h *AccountHandler) Setup(w http.ResponseWriter, r *http.Request) {
	claims := mw.GetUserClaims(r.Context())
	if claims == nil {
		respondErr(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req dto.AccountSetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondErr(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" {
		req.Email = claims.Email
	}

	account, vehicle, err := h.svc.Setup(r.Context(), claims.UserID, req)
	if err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := map[string]interface{}{
		"account": dto.AccountResponse{
			ID:     account.ID,
			Type:   account.Type,
			Name:   account.Name,
			Email:  account.Email,
			Phone:  account.Phone,
			Status: account.Status,
			CreatedAt: account.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}
	if vehicle != nil {
		resp["vehicle"] = dto.VehicleResponse{
			ID:          vehicle.ID,
			AccountID:   vehicle.AccountID,
			Make:        vehicle.Make,
			Model:       vehicle.Model,
			Year:        vehicle.Year,
			PlateNumber: vehicle.PlateNumber,
			VehicleType: vehicle.VehicleType,
			Status:      vehicle.Status,
			CreatedAt:   vehicle.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	respondJSON(w, http.StatusCreated, resp)
}

// PUT /api/v1/account
func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := mw.GetUserClaims(r.Context())
	if claims == nil {
		respondErr(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	account, err := h.svc.GetByUserID(r.Context(), claims.UserID)
	if err != nil || account == nil {
		respondErr(w, http.StatusNotFound, "account not found")
		return
	}

	var req dto.AccountUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.svc.Update(r.Context(), account.ID, req)
	if err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, dto.AccountResponse{
		ID:             updated.ID,
		Type:           updated.Type,
		Name:           updated.Name,
		Email:          updated.Email,
		Phone:          updated.Phone,
		BillingAddress: updated.BillingAddress,
		TaxID:          updated.TaxID,
		Status:         updated.Status,
		CreatedAt:      updated.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}
