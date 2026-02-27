package router

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/andresuarezz26/parkingmanagement/internal/config"
	"github.com/andresuarezz26/parkingmanagement/internal/handler"
	mw "github.com/andresuarezz26/parkingmanagement/internal/middleware"
)

func New(cfg *config.Config, logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	// --- Global middleware ---
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(mw.Recovery(logger))
	r.Use(mw.Logging(logger))
	r.Use(mw.CORS(cfg.AllowedOrigins))
	r.Use(mw.RateLimit(cfg.RateLimitRPS))

	// --- Public routes ---
	r.Get("/health", handler.Health)

	// --- API v1 ---
	r.Route("/api/v1", func(r chi.Router) {
		// Public endpoints (no auth)
		// r.Post("/day-passes", ...) — Phase 6
		// r.Get("/lots/{lotId}/plans", ...) — Phase 6
		// r.Post("/webhooks/stripe", ...) — Phase 4

		// Protected endpoints (auth required)
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth(cfg.SupabaseJWKSURL))

			// Driver / Account Holder routes — Phase 3+
			// r.Get("/account", ...)
			// r.Post("/vehicles", ...)
			// r.Get("/sessions", ...)

			// Admin routes (operator / super_admin)
			r.Route("/admin", func(r chi.Router) {
				r.Use(mw.RequireRole("operator", "super_admin"))

				// Phase 7: admin endpoints
				// r.Get("/dashboard", ...)
				// r.Get("/members", ...)
			})
		})

		// Dev-only routes
		if cfg.Env == "development" {
			r.Route("/dev", func(r chi.Router) {
				// Phase 5: simulation endpoints
				// r.Post("/simulate-entry", ...)
				// r.Post("/simulate-exit", ...)
			})
		}
	})

	return r
}
