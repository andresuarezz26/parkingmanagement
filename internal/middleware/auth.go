package middleware

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserClaimsKey contextKey = "user_claims"

// UserClaims holds the authenticated user's identity extracted from the JWT.
type UserClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// --- JWKS cache ---

type jwksCache struct {
	mu      sync.RWMutex
	keys    map[string]*ecdsa.PublicKey
	fetched time.Time
	url     string
	ttl     time.Duration
}

type jwksResponse struct {
	Keys []jwkKey `json:"keys"`
}

type jwkKey struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
}

func newJWKSCache(url string) *jwksCache {
	return &jwksCache{
		keys: make(map[string]*ecdsa.PublicKey),
		url:  url,
		ttl:  1 * time.Hour,
	}
}

func (c *jwksCache) getKey(kid string) (*ecdsa.PublicKey, error) {
	c.mu.RLock()
	if key, ok := c.keys[kid]; ok && time.Since(c.fetched) < c.ttl {
		c.mu.RUnlock()
		return key, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check
	if key, ok := c.keys[kid]; ok && time.Since(c.fetched) < c.ttl {
		return key, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(c.url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JWKS endpoint returned %d", resp.StatusCode)
	}

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	c.keys = make(map[string]*ecdsa.PublicKey)
	for _, k := range jwks.Keys {
		if k.Kty != "EC" || k.Crv != "P-256" {
			continue
		}
		xBytes, err := base64.RawURLEncoding.DecodeString(k.X)
		if err != nil {
			continue
		}
		yBytes, err := base64.RawURLEncoding.DecodeString(k.Y)
		if err != nil {
			continue
		}
		c.keys[k.Kid] = &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int).SetBytes(xBytes),
			Y:     new(big.Int).SetBytes(yBytes),
		}
	}
	c.fetched = time.Now()

	if key, ok := c.keys[kid]; ok {
		return key, nil
	}
	return nil, fmt.Errorf("key %q not found in JWKS", kid)
}

// --- Middleware ---

// Auth validates Supabase ES256 JWT tokens using JWKS discovery.
func Auth(jwksURL string) func(http.Handler) http.Handler {
	cache := newJWKSCache(jwksURL)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				respondError(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}
			tokenStr := parts[1]

			// Parse unverified to extract kid from header
			parser := jwt.NewParser()
			unverified, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
			if err != nil {
				respondError(w, http.StatusUnauthorized, "malformed token")
				return
			}
			kid, _ := unverified.Header["kid"].(string)
			if kid == "" {
				respondError(w, http.StatusUnauthorized, "token missing kid header")
				return
			}

			// Fetch the signing key
			pubKey, err := cache.getKey(kid)
			if err != nil {
				respondError(w, http.StatusUnauthorized, "unable to verify token signing key")
				return
			}

			// Verify signature and claims
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "ES256" {
					return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
				}
				return pubKey, nil
			}, jwt.WithValidMethods([]string{"ES256"}))

			if err != nil || !token.Valid {
				respondError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				respondError(w, http.StatusUnauthorized, "invalid claims")
				return
			}

			userClaims := UserClaims{
				UserID: claimString(claims, "sub"),
				Email:  claimString(claims, "email"),
				Role:   claimString(claims, "role"),
			}
			if userClaims.Role == "" {
				userClaims.Role = "driver"
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, &userClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserClaims extracts user claims from request context.
func GetUserClaims(ctx context.Context) *UserClaims {
	claims, _ := ctx.Value(UserClaimsKey).(*UserClaims)
	return claims
}

func claimString(claims jwt.MapClaims, key string) string {
	v, ok := claims[key]
	if !ok || v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
