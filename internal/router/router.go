package router

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/andresuarezz26/parkingmanagement/internal/config"
	"github.com/andresuarezz26/parkingmanagement/internal/handler"
	mw "github.com/andresuarezz26/parkingmanagement/internal/middleware"
	"github.com/andresuarezz26/parkingmanagement/internal/repository"
	"github.com/andresuarezz26/parkingmanagement/internal/service"
)

func New(cfg *config.Config, logger *zap.Logger, db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	// --- Global middleware ---
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(mw.Recovery(logger))
	r.Use(mw.Logging(logger))
	r.Use(mw.CORS(cfg.AllowedOrigins))
	r.Use(mw.RateLimit(cfg.RateLimitRPS))

	// --- Init layers ---
	accountRepo := repository.NewAccountRepo(db)
	vehicleRepo := repository.NewVehicleRepo(db)
	qrRepo := repository.NewQRRepo(db)

	qrSvc := service.NewQRService(qrRepo)
	vehicleSvc := service.NewVehicleService(vehicleRepo, qrSvc)
	accountSvc := service.NewAccountService(accountRepo, vehicleSvc)

	accountH := handler.NewAccountHandler(accountSvc)
	vehicleH := handler.NewVehicleHandler(vehicleSvc, accountSvc)
	qrH := handler.NewQRHandler(qrSvc, vehicleSvc, accountSvc)

	// --- Public routes ---
	r.Get("/health", handler.Health)

	// --- API v1 ---
	r.Route("/api/v1", func(r chi.Router) {
		// Public endpoints
		// r.Post("/day-passes", ...) — Phase 6
		// r.Get("/lots/{lotId}/plans", ...) — Phase 6
		// r.Post("/webhooks/stripe", ...) — Phase 4

		// Protected endpoints
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth(cfg.SupabaseJWKSURL))

			// Account
			r.Get("/account", accountH.Get)
			r.Put("/account", accountH.Update)
			r.Post("/account/setup", accountH.Setup)

			// Vehicles
			r.Get("/vehicles", vehicleH.List)
			r.Post("/vehicles", vehicleH.Create)
			r.Get("/vehicles/{id}", vehicleH.Get)
			r.Put("/vehicles/{id}", vehicleH.Update)
			r.Delete("/vehicles/{id}", vehicleH.Delete)

			// QR Codes
			r.Get("/vehicles/{id}/qr", qrH.GetByVehicle)
			r.Post("/vehicles/{id}/qr/regenerate", qrH.Regenerate)

			// Admin routes
			r.Route("/admin", func(r chi.Router) {
				r.Use(mw.RequireRole("operator", "super_admin"))
				// Phase 7
			})
		})

		// Dev-only routes
		if cfg.Env == "development" {
			r.Route("/dev", func(r chi.Router) {
				// Phase 5
			})
		}
	})

	return r
}
