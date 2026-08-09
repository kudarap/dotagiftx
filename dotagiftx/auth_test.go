package dotagiftx

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestAuthToContextAndFromContext(t *testing.T) {
	au := &Auth{ID: "a1", UserID: "u1", Username: "s1"}

	ctx := AuthToContext(context.Background(), au)
	if got := AuthFromContext(ctx); got != au {
		t.Errorf("AuthFromContext() = %v, want %v", got, au)
	}

	if got := AuthFromContext(context.Background()); got != nil {
		t.Errorf("AuthFromContext() without auth = %v, want nil", got)
	}
}

func TestAuthService_composePassword(t *testing.T) {
	svc := &AuthService{salt: "test-salt"}

	a := svc.composePassword("s1", "u1")
	if a != svc.composePassword("s1", "u1") {
		t.Error("composePassword() should be deterministic")
	}
	if a == svc.composePassword("s2", "u1") {
		t.Error("composePassword() should differ for different inputs")
	}
	if len(a) != 40 {
		t.Errorf("composePassword() len = %d, want 40 (sha1 hex)", len(a))
	}
}

func TestAuthService_composePasswordV2(t *testing.T) {
	svc := &AuthService{salt: "test-salt"}

	got := svc.composePasswordV2("s1", "u1")

	sum := sha256.Sum256([]byte("s1u1test-salt"))
	want := hex.EncodeToString(sum[:])
	if got != want {
		t.Errorf("composePasswordV2() = %q, want %q (sha256 of concatenated steamID+userID+salt)", got, want)
	}

	if svc.composePasswordV2("s1", "u1") != got {
		t.Error("composePasswordV2() should be deterministic")
	}
}

func TestAuthService_hash(t *testing.T) {
	svc := &AuthService{salt: "test-salt"}

	// hash(a, b) must equal hash(a+b): it is a digest of the concatenated args.
	if joined := svc.hash("s1", "u1"); joined != svc.hash("s1u1") {
		t.Errorf("hash(a, b) = %q, want equal to hash(a+b) = %q", joined, svc.hash("s1u1"))
	}

	tok := "some-token"
	if svc.hash(tok) == svc.hash(tok+"x") {
		t.Error("hash() should differ for different inputs")
	}
}

func TestAuthService_generateRefreshToken(t *testing.T) {
	svc := &AuthService{}

	const samples = 10
	seen := make(map[string]bool, samples)
	for range samples {
		tok, err := svc.generateRefreshToken()
		if err != nil {
			t.Fatalf("generateRefreshToken() returned error: %v", err)
		}
		if len(tok) != 64 {
			t.Errorf("generateRefreshToken() len = %d, want 64", len(tok))
		}
		b, err := base64.RawURLEncoding.DecodeString(tok)
		if err != nil {
			t.Fatalf("generateRefreshToken() not base64url: %v", err)
		}
		if len(b) != 48 {
			t.Errorf("generateRefreshToken() decoded len = %d, want 48", len(b))
		}
		if seen[tok] {
			t.Errorf("generateRefreshToken() produced duplicate token %q", tok)
		}
		seen[tok] = true
	}
}
