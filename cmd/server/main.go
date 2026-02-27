package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/andresuarezz26/parkingmanagement/internal/config"
	"github.com/andresuarezz26/parkingmanagement/internal/router"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		sugar.Fatalw("failed to load config", "error", err)
	}

	sugar.Infow("config loaded",
		"port", cfg.Port,
		"env", cfg.Env,
		"payment_mock", cfg.PaymentMockEnabled,
	)

	// Initialize router
	r := router.New(cfg)

	// Create HTTP server
	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		sugar.Infow("🚀 HeavyPark API starting", "addr", addr, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalw("server failed", "error", err)
		}
	}()

	<-done
	sugar.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalw("forced shutdown", "error", err)
	}

	sugar.Info("server stopped")
}
