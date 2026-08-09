package dotagiftx

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"
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

type mockAuthRepo struct {
	auths map[string]*Auth
}

func newMockAuthRepo() *mockAuthRepo {
	return &mockAuthRepo{auths: map[string]*Auth{}}
}

func (m *mockAuthRepo) Get(_ context.Context, id string) (*Auth, error) {
	if a, ok := m.auths[id]; ok {
		return a, nil
	}
	return nil, AuthErrNotFound
}

func (m *mockAuthRepo) GetByUsername(_ context.Context, username string) (*Auth, error) {
	for _, a := range m.auths {
		if a.Username == username {
			return a, nil
		}
	}
	return nil, AuthErrNotFound
}

func (m *mockAuthRepo) GetByUsernameAndPassword(context.Context, string, string) (*Auth, error) {
	return nil, AuthErrNotFound
}

func (m *mockAuthRepo) Create(_ context.Context, in *Auth) error {
	in.ID = fmt.Sprintf("auth-%d", len(m.auths))
	m.auths[in.ID] = in
	return nil
}

func (m *mockAuthRepo) Update(_ context.Context, in *Auth) error {
	m.auths[in.ID] = in
	return nil
}

type mockTokenRepo struct {
	tokens map[string]*RefreshToken
	byHash map[string]*RefreshToken
}

func newMockTokenRepo() *mockTokenRepo {
	return &mockTokenRepo{
		tokens: map[string]*RefreshToken{},
		byHash: map[string]*RefreshToken{},
	}
}

func (m *mockTokenRepo) GetByTokenHash(_ context.Context, hash string) (*RefreshToken, error) {
	if rt, ok := m.byHash[hash]; ok {
		return rt, nil
	}
	return nil, AuthErrRefreshToken
}

func (m *mockTokenRepo) Create(_ context.Context, in *RefreshToken) error {
	in.ID = fmt.Sprintf("tok-%d", len(m.tokens))
	m.tokens[in.ID] = in
	m.byHash[in.TokenHash] = in
	return nil
}

func (m *mockTokenRepo) Update(_ context.Context, in *RefreshToken) error {
	m.tokens[in.ID] = in
	m.byHash[in.TokenHash] = in
	return nil
}

func (m *mockTokenRepo) RevokeFamily(_ context.Context, familyID string) error {
	for _, rt := range m.tokens {
		if rt.FamilyID == familyID {
			rt.Revoked = true
		}
	}
	return nil
}

func newTestAuthService(authRepo authRepository, tokenRepo refreshTokenRepository) *AuthService {
	return &AuthService{
		salt:      "test-salt",
		authRepo:  authRepo,
		tokenRepo: tokenRepo,
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	ctx := context.Background()
	authRepo := newMockAuthRepo()
	tokenRepo := newMockTokenRepo()
	svc := newTestAuthService(authRepo, tokenRepo)

	authRepo.auths["a1"] = &Auth{ID: "a1", UserID: "u1", Username: "s1"}
	raw, err := svc.issueRefreshToken(ctx, "a1", "u1")
	if err != nil {
		t.Fatalf("issueRefreshToken() returned error: %v", err)
	}

	t.Run("rotates token on successful refresh", func(t *testing.T) {
		got, err := svc.RefreshToken(ctx, raw)
		if err != nil {
			t.Fatalf("RefreshToken() returned error: %v", err)
		}
		if got.RefreshToken == raw || got.RefreshToken == "" {
			t.Errorf("RefreshToken() did not rotate token, got same token")
		}
		if got.UserID != "u1" {
			t.Errorf("RefreshToken() UserID = %q, want u1", got.UserID)
		}

		old, err := tokenRepo.GetByTokenHash(ctx, svc.hash(raw))
		if err != nil {
			t.Fatalf("GetByTokenHash(old) returned error: %v", err)
		}
		if !old.Revoked {
			t.Error("presented token should be revoked after rotation")
		}

		if _, err := svc.RefreshToken(ctx, got.RefreshToken); err != nil {
			t.Errorf("new token should be valid for next refresh: %v", err)
		}
	})

	t.Run("rejects reused rotated token and revokes family", func(t *testing.T) {
		// rotation above left family with several tokens; reuse of the first
		// token must invalidate the whole family.
		if _, err := svc.RefreshToken(ctx, raw); err != AuthErrRefreshToken {
			t.Errorf("RefreshToken() reused token error = %v, want AuthErrRefreshToken", err)
		}
		for _, rt := range tokenRepo.tokens {
			if !rt.Revoked {
				t.Errorf("token %s in family %s should be revoked after reuse detection", rt.ID, rt.FamilyID)
			}
		}
	})

	t.Run("rejects expired token", func(t *testing.T) {
		exp := time.Now().Add(-time.Minute)
		expiredRaw, err := svc.issueRefreshToken(ctx, "a1", "u1")
		if err != nil {
			t.Fatalf("issueRefreshToken() returned error: %v", err)
		}
		expired := tokenRepo.byHash[svc.hash(expiredRaw)]
		expired.ExpiresAt = &exp

		if _, err := svc.RefreshToken(ctx, expiredRaw); err != AuthErrRefreshToken {
			t.Errorf("RefreshToken() expired token error = %v, want AuthErrRefreshToken", err)
		}
	})

	t.Run("rejects empty token", func(t *testing.T) {
		if _, err := svc.RefreshToken(ctx, ""); err != AuthErrRefreshToken {
			t.Errorf("RefreshToken() empty token error = %v, want AuthErrRefreshToken", err)
		}
	})
}

func TestAuthService_RevokeRefreshToken(t *testing.T) {
	ctx := context.Background()
	authRepo := newMockAuthRepo()
	tokenRepo := newMockTokenRepo()
	svc := newTestAuthService(authRepo, tokenRepo)

	raw, err := svc.issueRefreshToken(ctx, "a1", "u1")
	if err != nil {
		t.Fatalf("issueRefreshToken() returned error: %v", err)
	}

	if err := svc.RevokeRefreshToken(ctx, raw); err != nil {
		t.Fatalf("RevokeRefreshToken() returned error: %v", err)
	}
	for _, rt := range tokenRepo.tokens {
		if !rt.Revoked {
			t.Errorf("token %s should be revoked after RevokeRefreshToken", rt.ID)
		}
	}
	if _, err := svc.RefreshToken(ctx, raw); err != AuthErrRefreshToken {
		t.Errorf("RefreshToken() revoked token error = %v, want AuthErrRefreshToken", err)
	}
}
