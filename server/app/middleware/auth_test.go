package middleware

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func b64(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

// signToken mints an RS256 JWT the same way AOAuth does.
func signToken(t *testing.T, key *rsa.PrivateKey, kid string, claims map[string]any) string {
	t.Helper()
	hb, _ := json.Marshal(map[string]string{"alg": "RS256", "typ": "JWT", "kid": kid})
	cb, _ := json.Marshal(claims)
	signingInput := b64(hb) + "." + b64(cb)
	digest := sha256.Sum256([]byte(signingInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	return signingInput + "." + b64(sig)
}

func encodeExp(e int) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(e))
	i := 0
	for i < len(buf)-1 && buf[i] == 0 {
		i++
	}
	return buf[i:]
}

// jwksServer serves an AOAuth-style JWKS for one RSA key.
func jwksServer(kid string, pub *rsa.PublicKey) *httptest.Server {
	doc := map[string]any{"keys": []map[string]string{{
		"kty": "RSA", "use": "sig", "alg": "RS256", "kid": kid,
		"n": b64(pub.N.Bytes()), "e": b64(encodeExp(pub.E)),
	}}}
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(doc)
	})
	return httptest.NewServer(mux)
}

func TestVerify(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	const kid = "test-key-1"
	srv := jwksServer(kid, &key.PublicKey)
	defer srv.Close()

	v := NewJWTVerifier(srv.URL)
	now := time.Now().Unix()

	valid := signToken(t, key, kid, map[string]any{"iss": srv.URL, "sub": "user-123", "exp": now + 60, "nbf": now - 10})
	if sub, err := v.verify(valid); err != nil {
		t.Fatalf("valid token rejected: %v", err)
	} else if sub != "user-123" {
		t.Fatalf("verify subject = %q, want %q", sub, "user-123")
	}

	otherKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	bad := map[string]string{
		"expired":       signToken(t, key, kid, map[string]any{"iss": srv.URL, "exp": now - 10}),
		"wrong issuer":  signToken(t, key, kid, map[string]any{"iss": "http://evil.example", "exp": now + 60}),
		"unknown kid":   signToken(t, key, "rogue-kid", map[string]any{"iss": srv.URL, "exp": now + 60}),
		"wrong key":     signToken(t, otherKey, kid, map[string]any{"iss": srv.URL, "exp": now + 60}),
		"not yet valid": signToken(t, key, kid, map[string]any{"iss": srv.URL, "exp": now + 60, "nbf": now + 60}),
		"malformed":     "not.a.jwt",
		"two parts":     "aaaa.bbbb",
	}
	for name, tok := range bad {
		if _, err := v.verify(tok); err == nil {
			t.Errorf("%s: expected rejection, got nil", name)
		}
	}

	// Tampered signature must fail.
	if _, err := v.verify(valid[:len(valid)-3] + "AAA"); err == nil {
		t.Error("tampered signature accepted")
	}
}

func TestBearerToken(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/page", nil)
	if _, err := bearerToken(r); err == nil {
		t.Error("expected error for missing header")
	}
	r.Header.Set("Authorization", "Bearer abc.def.ghi")
	got, err := bearerToken(r)
	if err != nil || got != "abc.def.ghi" {
		t.Errorf("bearerToken = %q, %v", got, err)
	}
}
