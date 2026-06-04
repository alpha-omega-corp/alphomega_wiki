// Package middleware provides HTTP middleware for the wiki backend.
//
// auth.go enforces AOAuth-issued access tokens on protected (mutating) routes.
// It validates RS256 JWTs against the authorization server's published JWKS
// (RFC 7517), so a request can only mutate data if it carries a valid,
// unexpired token signed by the trusted issuer. Read routes stay public.
//
// JWT/JWKS handling is implemented against the standard library to avoid adding
// a dependency; swap for github.com/golang-jwt/jwt + a JWKS cache if richer
// validation is needed.
package middleware

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/uptrace/bunrouter"
)

// jwksMinRefetch rate-limits JWKS refetches triggered by unknown key IDs.
const jwksMinRefetch = 30 * time.Second

// contextKey is unexported so only this package can set the value; consumers
// read it through UserSubFromContext.
type contextKey string

const userSubKey contextKey = "user_sub"

// UserSubFromContext returns the authenticated user's AOAuth subject (the JWT
// `sub` claim) that Middleware stored on the request context, or "" if the
// request was not authenticated.
func UserSubFromContext(ctx context.Context) string {
	sub, _ := ctx.Value(userSubKey).(string)
	return sub
}

// JWTVerifier validates AOAuth access tokens against the issuer's JWKS. Public
// keys are fetched lazily and cached; an unknown `kid` triggers a (rate-limited)
// refetch so key rotation is picked up automatically.
type JWTVerifier struct {
	issuer  string
	jwksURL string
	client  *http.Client

	mu        sync.RWMutex
	keys      map[string]*rsa.PublicKey
	lastFetch time.Time
}

// NewJWTVerifier builds a verifier for tokens issued by the given AOAuth issuer.
func NewJWTVerifier(issuer string) *JWTVerifier {
	issuer = strings.TrimRight(issuer, "/")
	return &JWTVerifier{
		issuer:  issuer,
		jwksURL: issuer + "/.well-known/jwks.json",
		client:  &http.Client{Timeout: 5 * time.Second},
		keys:    map[string]*rsa.PublicKey{},
	}
}

// Middleware authenticates a request; on failure it writes 401 and stops the
// chain, otherwise the wrapped handler runs.
func (v *JWTVerifier) Middleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		token, err := bearerToken(req.Request)
		if err != nil {
			return unauthorized(w, err)
		}
		sub, err := v.verify(token)
		if err != nil {
			return unauthorized(w, err)
		}
		// Expose the verified identity to downstream handlers (e.g. favorites
		// are scoped to the calling user).
		ctx := context.WithValue(req.Context(), userSubKey, sub)
		req.Request = req.Request.WithContext(ctx)
		return next(w, req)
	}
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

type jwtClaims struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	Expiry    int64  `json:"exp"`
	NotBefore int64  `json:"nbf"`
}

// verify checks the token's signature (against the JWKS), algorithm, expiry and
// issuer, returning the `sub` claim (the authenticated user) on success.
func (v *JWTVerifier) verify(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("malformed token")
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", errors.New("invalid token header encoding")
	}
	var header jwtHeader
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return "", errors.New("invalid token header")
	}
	if header.Alg != "RS256" {
		return "", fmt.Errorf("unsupported token algorithm %q", header.Alg)
	}

	key, err := v.keyFor(header.Kid)
	if err != nil {
		return "", err
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return "", errors.New("invalid signature encoding")
	}
	digest := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, digest[:], sig); err != nil {
		return "", errors.New("signature verification failed")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.New("invalid claims encoding")
	}
	var claims jwtClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return "", errors.New("invalid claims")
	}

	now := time.Now().Unix()
	if claims.Expiry != 0 && now >= claims.Expiry {
		return "", errors.New("token expired")
	}
	if claims.NotBefore != 0 && now < claims.NotBefore {
		return "", errors.New("token not yet valid")
	}
	if v.issuer != "" && claims.Issuer != v.issuer {
		return "", errors.New("issuer mismatch")
	}
	return claims.Subject, nil
}

// keyFor returns the cached public key for kid, refetching the JWKS once if the
// key is unknown (handles rotation).
func (v *JWTVerifier) keyFor(kid string) (*rsa.PublicKey, error) {
	v.mu.RLock()
	key := v.keys[kid]
	v.mu.RUnlock()
	if key != nil {
		return key, nil
	}

	if err := v.refresh(); err != nil {
		return nil, err
	}

	v.mu.RLock()
	key = v.keys[kid]
	v.mu.RUnlock()
	if key == nil {
		return nil, fmt.Errorf("unknown signing key %q", kid)
	}
	return key, nil
}

// refresh fetches and caches the JWKS, rate-limited to avoid hammering the
// issuer on a flood of bad tokens.
func (v *JWTVerifier) refresh() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.lastFetch.IsZero() && time.Since(v.lastFetch) < jwksMinRefetch {
		return nil
	}

	resp, err := v.client.Get(v.jwksURL)
	if err != nil {
		return fmt.Errorf("fetch jwks: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jwks endpoint returned %d", resp.StatusCode)
	}

	var doc struct {
		Keys []struct {
			Kty string `json:"kty"`
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return fmt.Errorf("decode jwks: %w", err)
	}

	keys := make(map[string]*rsa.PublicKey, len(doc.Keys))
	for _, k := range doc.Keys {
		if k.Kty != "RSA" {
			continue
		}
		pub, err := rsaPublicKey(k.N, k.E)
		if err != nil {
			continue
		}
		keys[k.Kid] = pub
	}
	v.keys = keys
	v.lastFetch = time.Now()
	return nil
}

// rsaPublicKey reconstructs an RSA public key from the JWK n/e parameters.
func rsaPublicKey(nB64, eB64 string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eB64)
	if err != nil {
		return nil, err
	}
	e := new(big.Int).SetBytes(eBytes)
	if !e.IsInt64() || e.Int64() > 1<<31 {
		return nil, errors.New("invalid exponent")
	}
	return &rsa.PublicKey{N: new(big.Int).SetBytes(nBytes), E: int(e.Int64())}, nil
}

func bearerToken(r *http.Request) (string, error) {
	const prefix = "Bearer "
	h := r.Header.Get("Authorization")
	if len(h) <= len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return "", errors.New("missing bearer token")
	}
	return strings.TrimSpace(h[len(prefix):]), nil
}

func unauthorized(w http.ResponseWriter, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
	w.WriteHeader(http.StatusUnauthorized)
	body, _ := json.Marshal(map[string]string{"error": "unauthorized", "message": err.Error()})
	_, _ = w.Write(body)
	return nil
}
