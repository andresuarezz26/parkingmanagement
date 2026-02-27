package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/andresuarezz26/parkingmanagement/internal/dto"
	mw "github.com/andresuarezz26/parkingmanagement/internal/middleware"
	"github.com/andresuarezz26/parkingmanagement/internal/service"
)

type VehicleHandler struct {
	vehicleSvc *service.VehicleService
	accountSvc *service.AccountService
}

func NewVehicleHandler(vehicleSvc *service.VehicleService, accountSvc *service.AccountService) *VehicleHandler {
	return &VehicleHandler{vehicleSvc: vehicleSvc, accountSvc: accountSvc}
}

func (h *VehicleHandler) getAccountID(r *http.Request) (string, error) {
	claims := mw.GetUserClaims(r.Context())
	account, err := h.accountSvc.GetByUserID(r.Context(), claims.UserID)
	if err != nil || account == nil {
		return "", err
	}
	return account.ID, nil
}

// GET /api/v1/vehicles
func (h *VehicleHandler) List(w http.ResponseWriter, r *http.Request) {
	accountID, err := h.getAccountID(r)
	if err != nil || accountID == "" {
		respondErr(w, http.StatusNotFound, "account not found")
		return
	}

	vehicles, err := h.vehicleSvc.ListByAccount(r.Context(), accountID)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, "failed to list vehicles")
		return
	}

	resp := make([]dto.VehicleResponse, 0, len(vehicles))
	for _, v := range vehicles {
		resp = append(resp, dto.VehicleResponse{
			ID:          v.ID,
			AccountID:   v.AccountID,
			Make:        v.Make,
			Model:       v.Model,
			Year:        v.Year,
			PlateNumber: v.PlateNumber,
			VehicleType: v.VehicleType,
			Description: v.Description,
			Status:      v.Status,
			CreatedAt:   v.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	respondJSON(w, http.StatusOK, resp)
}

// POST /api/v1/vehicles
func (h *VehicleHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountID, err := h.getAccountID(r)
	if err != nil || accountID == "" {
		respondErr(w, http.StatusNotFound, "account not found — set up account first")
		return
	}

	var req dto.VehicleCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	vehicle, err := h.vehicleSvc.Create(r.Context(), accountID, req)
	if err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, dto.VehicleResponse{
		ID:          vehicle.ID,
		AccountID:   vehicle.AccountID,
		Make:        vehicle.Make,
		Model:       vehicle.Model,
		Year:        vehicle.Year,
		PlateNumber: vehicle.PlateNumber,
		VehicleType: vehicle.VehicleType,
		Description: vehicle.Description,
		Status:      vehicle.Status,
		CreatedAt:   vehicle.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// GET /api/v1/vehicles/{id}
func (h *VehicleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	vehicle, err := h.vehicleSvc.GetByID(r.Context(), id)
	if err != nil || vehicle == nil {
		respondErr(w, http.StatusNotFound, "vehicle not found")
		return
	}

	// TODO: verify vehicle belongs to user's account

	respondJSON(w, http.StatusOK, dto.VehicleResponse{
		ID:          vehicle.ID,
		AccountID:   vehicle.AccountID,
		Make:        vehicle.Make,
		Model:       vehicle.Model,
		Year:        vehicle.Year,
		PlateNumber: vehicle.PlateNumber,
		VehicleType: vehicle.VehicleType,
		Description: vehicle.Description,
		Status:      vehicle.Status,
		CreatedAt:   vehicle.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// PUT /api/v1/vehicles/{id}
func (h *VehicleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.VehicleUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	vehicle, err := h.vehicleSvc.Update(r.Context(), id, req)
	if err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, dto.VehicleResponse{
		ID:          vehicle.ID,
		AccountID:   vehicle.AccountID,
		Make:        vehicle.Make,
		Model:       vehicle.Model,
		Year:        vehicle.Year,
		PlateNumber: vehicle.PlateNumber,
		VehicleType: vehicle.VehicleType,
		Description: vehicle.Description,
		Status:      vehicle.Status,
		CreatedAt:   vehicle.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// DELETE /api/v1/vehicles/{id}
func (h *VehicleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.vehicleSvc.Delete(r.Context(), id); err != nil {
		respondErr(w, http.StatusInternalServerError, "failed to delete vehicle")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "vehicle removed"})
}
