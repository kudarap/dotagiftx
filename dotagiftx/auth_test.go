package dotagiftx

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
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

// fakeAuthRepo implements authRepository for service tests.
type fakeAuthRepo struct {
	auth *Auth
}

func (f *fakeAuthRepo) Get(ctx context.Context, id string) (*Auth, error) {
	if f.auth == nil || f.auth.ID != id {
		return nil, AuthErrNotFound
	}
	return f.auth, nil
}

func (f *fakeAuthRepo) GetByUsername(ctx context.Context, username string) (*Auth, error) {
	if f.auth == nil || f.auth.Username != username {
		return nil, AuthErrNotFound
	}
	return f.auth, nil
}

func (f *fakeAuthRepo) GetByUsernameAndPassword(ctx context.Context, username, password string) (*Auth, error) {
	if f.auth == nil || f.auth.Username != username || f.auth.Password != password {
		return nil, AuthErrNotFound
	}
	return f.auth, nil
}

func (f *fakeAuthRepo) Create(ctx context.Context, in *Auth) error {
	f.auth = in
	return nil
}

func (f *fakeAuthRepo) Update(ctx context.Context, in *Auth) error {
	if f.auth == nil {
		return AuthErrNotFound
	}
	f.auth = in
	return nil
}

// fakeSessionRepo implements authSessionRepository for service tests.
type fakeSessionRepo struct {
	sessions map[string]*AuthSession
	nextID   int
}

func (f *fakeSessionRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (*AuthSession, error) {
	if f.sessions == nil {
		return nil, AuthErrNotFound
	}
	sess, ok := f.sessions[refreshToken]
	if !ok {
		return nil, AuthErrNotFound
	}
	return sess, nil
}

func (f *fakeSessionRepo) Create(ctx context.Context, in *AuthSession) error {
	if f.sessions == nil {
		f.sessions = make(map[string]*AuthSession)
	}
	f.nextID++
	in.ID = fmt.Sprintf("sess-%d", f.nextID)
	f.sessions[in.RefreshToken] = in
	return nil
}

func (f *fakeSessionRepo) Update(ctx context.Context, in *AuthSession) error {
	if f.sessions == nil {
		return AuthErrNotFound
	}
	if _, ok := f.sessions[in.RefreshToken]; !ok {
		// refresh token key may have changed on rotation; locate by ID.
		for k, s := range f.sessions {
			if s.ID == in.ID {
				delete(f.sessions, k)
				f.sessions[in.RefreshToken] = in
				return nil
			}
		}
		return AuthErrNotFound
	}
	f.sessions[in.RefreshToken] = in
	return nil
}

func (f *fakeSessionRepo) Delete(ctx context.Context, id string) error {
	if f.sessions == nil {
		return nil
	}
	for k, s := range f.sessions {
		if s.ID == id {
			delete(f.sessions, k)
			return nil
		}
	}
	return nil
}

func newTestAuthService(sess *fakeSessionRepo) *AuthService {
	if sess.sessions == nil {
		sess.sessions = make(map[string]*AuthSession)
	}
	return &AuthService{
		salt:        "test-salt",
		sessionTTL:  time.Hour,
		authRepo:    &fakeAuthRepo{auth: &Auth{ID: "a1", UserID: "u1", Username: "s1"}},
		sessionRepo: sess,
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	sess := &fakeSessionRepo{}
	svc := newTestAuthService(sess)

	orig := "original-raw-token"
	sess.sessions[svc.hash(orig)] = &AuthSession{
		ID:           "sess-1",
		AuthID:       "a1",
		UserID:       "u1",
		RefreshToken: svc.hash(orig),
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	au, rotated, err := svc.RefreshToken(context.Background(), orig)
	if err != nil {
		t.Fatalf("RefreshToken() returned error: %v", err)
	}

	if rotated == orig {
		t.Error("RefreshToken() expected rotated refresh token")
	}
	if au.UserID != "u1" {
		t.Errorf("RefreshToken() user_id = %q, want u1", au.UserID)
	}
	if au.Username != "s1" {
		t.Errorf("RefreshToken() username = %q, want s1", au.Username)
	}
	if _, ok := sess.sessions[svc.hash(orig)]; ok {
		t.Error("RefreshToken() old token should be removed")
	}
	stored, ok := sess.sessions[svc.hash(rotated)]
	if !ok {
		t.Error("RefreshToken() rotated token not persisted")
	}
	if !stored.ExpiresAt.After(time.Now()) {
		t.Error("RefreshToken() rotated session should have fresh expiration")
	}

	// already rotated token should fail.
	if _, _, err := svc.RefreshToken(context.Background(), orig); !errors.Is(err, AuthErrRefreshToken) {
		t.Errorf("RefreshToken() after rotation error = %v, want AuthErrRefreshToken", err)
	}
}

func TestAuthService_RefreshTokenExpired(t *testing.T) {
	sess := &fakeSessionRepo{}
	svc := newTestAuthService(sess)

	expTok := "expired-token"
	sess.sessions[svc.hash(expTok)] = &AuthSession{
		ID:           "sess-exp",
		AuthID:       "a1",
		UserID:       "u1",
		RefreshToken: svc.hash(expTok),
		ExpiresAt:    time.Now().Add(-time.Hour),
	}

	if _, _, err := svc.RefreshToken(context.Background(), expTok); !errors.Is(err, AuthErrRefreshToken) {
		t.Errorf("RefreshToken() expired error = %v, want AuthErrRefreshToken", err)
	}
	if _, ok := sess.sessions[svc.hash(expTok)]; ok {
		t.Error("RefreshToken() expired session should be deleted")
	}
}

func TestAuthService_RefreshTokenNotFound(t *testing.T) {
	svc := newTestAuthService(&fakeSessionRepo{})

	if _, _, err := svc.RefreshToken(context.Background(), "no-such-token"); !errors.Is(err, AuthErrRefreshToken) {
		t.Errorf("RefreshToken() missing error = %v, want AuthErrRefreshToken", err)
	}
	if _, _, err := svc.RefreshToken(context.Background(), ""); !errors.Is(err, AuthErrRefreshToken) {
		t.Errorf("RefreshToken() empty error = %v, want AuthErrRefreshToken", err)
	}
}

func TestAuthService_MultipleSessions(t *testing.T) {
	sess := &fakeSessionRepo{}
	svc := newTestAuthService(sess)

	tokA := "token-device-a"
	tokB := "token-device-b"
	sess.sessions[svc.hash(tokA)] = &AuthSession{
		ID:           "sess-a",
		AuthID:       "a1",
		UserID:       "u1",
		RefreshToken: svc.hash(tokA),
		ExpiresAt:    time.Now().Add(time.Hour),
	}
	sess.sessions[svc.hash(tokB)] = &AuthSession{
		ID:           "sess-b",
		AuthID:       "a1",
		UserID:       "u1",
		RefreshToken: svc.hash(tokB),
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	if _, _, err := svc.RefreshToken(context.Background(), tokA); err != nil {
		t.Fatalf("RefreshToken() session A error: %v", err)
	}
	if _, _, err := svc.RefreshToken(context.Background(), tokB); err != nil {
		t.Fatalf("RefreshToken() session B error: %v", err)
	}
}

func TestAuthService_RevokeRefreshToken(t *testing.T) {
	sess := &fakeSessionRepo{}
	svc := newTestAuthService(sess)

	tok := "revoke-token"
	sess.sessions[svc.hash(tok)] = &AuthSession{
		ID:           "sess-revoke",
		AuthID:       "a1",
		UserID:       "u1",
		RefreshToken: svc.hash(tok),
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	if err := svc.RevokeRefreshToken(context.Background(), tok); err != nil {
		t.Fatalf("RevokeRefreshToken() returned error: %v", err)
	}
	if _, ok := sess.sessions[svc.hash(tok)]; ok {
		t.Error("RevokeRefreshToken() should delete the session")
	}

	if err := svc.RevokeRefreshToken(context.Background(), tok); !errors.Is(err, AuthErrRefreshToken) {
		t.Errorf("RevokeRefreshToken() already revoked error = %v, want AuthErrRefreshToken", err)
	}
}
