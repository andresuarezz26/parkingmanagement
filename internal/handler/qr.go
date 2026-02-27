package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/andresuarezz26/parkingmanagement/internal/dto"
	mw "github.com/andresuarezz26/parkingmanagement/internal/middleware"
	"github.com/andresuarezz26/parkingmanagement/internal/service"
)

type QRHandler struct {
	qrSvc      *service.QRService
	vehicleSvc *service.VehicleService
	accountSvc *service.AccountService
}

func NewQRHandler(qrSvc *service.QRService, vehicleSvc *service.VehicleService, accountSvc *service.AccountService) *QRHandler {
	return &QRHandler{qrSvc: qrSvc, vehicleSvc: vehicleSvc, accountSvc: accountSvc}
}

// GET /api/v1/vehicles/{id}/qr
func (h *QRHandler) GetByVehicle(w http.ResponseWriter, r *http.Request) {
	vehicleID := chi.URLParam(r, "id")

	qr, err := h.qrSvc.GetByVehicle(r.Context(), vehicleID)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, "failed to get QR code")
		return
	}
	if qr == nil {
		respondErr(w, http.StatusNotFound, "no QR code found for this vehicle")
		return
	}

	resp := dto.QRResponse{
		ID:        qr.ID,
		VehicleID: qr.VehicleID,
		CodeData:  qr.CodeData,
		ImageURL:  qr.ImageURL,
		Status:    qr.Status,
		IssuedAt:  qr.IssuedAt.Format("2006-01-02T15:04:05Z"),
	}
	if qr.ExpiresAt != nil {
		s := qr.ExpiresAt.Format("2006-01-02T15:04:05Z")
		resp.ExpiresAt = &s
	}
	if qr.LastScannedAt != nil {
		s := qr.LastScannedAt.Format("2006-01-02T15:04:05Z")
		resp.LastScannedAt = &s
	}

	respondJSON(w, http.StatusOK, resp)
}

// POST /api/v1/vehicles/{id}/qr/regenerate
func (h *QRHandler) Regenerate(w http.ResponseWriter, r *http.Request) {
	vehicleID := chi.URLParam(r, "id")

	claims := mw.GetUserClaims(r.Context())
	account, err := h.accountSvc.GetByUserID(r.Context(), claims.UserID)
	if err != nil || account == nil {
		respondErr(w, http.StatusNotFound, "account not found")
		return
	}

	// Verify vehicle belongs to this account
	vehicle, err := h.vehicleSvc.GetByID(r.Context(), vehicleID)
	if err != nil || vehicle == nil {
		respondErr(w, http.StatusNotFound, "vehicle not found")
		return
	}
	if vehicle.AccountID != account.ID {
		respondErr(w, http.StatusForbidden, "vehicle does not belong to your account")
		return
	}

	qr, err := h.qrSvc.Regenerate(r.Context(), vehicleID, account.ID)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, "failed to regenerate QR code")
		return
	}

	respondJSON(w, http.StatusCreated, dto.QRResponse{
		ID:        qr.ID,
		VehicleID: qr.VehicleID,
		CodeData:  qr.CodeData,
		ImageURL:  qr.ImageURL,
		Status:    qr.Status,
		IssuedAt:  qr.IssuedAt.Format("2006-01-02T15:04:05Z"),
	})
}
