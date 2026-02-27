package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port               int    `envconfig:"PORT" default:"8080"`
	Env                string `envconfig:"ENV" default:"development"`
	DatabaseURL        string `envconfig:"DATABASE_URL" required:"true"`
	SupabaseURL        string `envconfig:"SUPABASE_URL"`
	SupabaseJWKSURL    string `envconfig:"SUPABASE_JWKS_URL" default:"https://vssfudgcsqkvaqvqymnd.supabase.co/auth/v1/.well-known/jwks.json"`
	PaymentMockEnabled bool   `envconfig:"PAYMENT_MOCK_ENABLED" default:"true"`
	AllowedOrigins     string `envconfig:"ALLOWED_ORIGINS" default:"http://localhost:3000,http://localhost:5173"`
	RateLimitRPS       int    `envconfig:"RATE_LIMIT_RPS" default:"100"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}
