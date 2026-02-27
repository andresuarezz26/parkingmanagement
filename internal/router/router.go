package router

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/andresuarezz26/parkingmanagement/internal/handler"
	"github.com/andresuarezz26/parkingmanagement/internal/config"
)

func New(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Built-in middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)

	// Health check (public)
	r.Get("/health", handler.Health)

	// API v1 routes (to be added in subsequent phases)
	r.Route("/api/v1", func(r chi.Router) {
		// Phase 2+: middleware, auth, routes
	})

	return r
}
