package dotagiftx

import (
	"context"
	"crypto/rand"
	"crypto/sha1" //nolint:gosec // legacy password/token hashing for backward compatibility
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	AuthErrNotFound Errors = iota + authErrorIndex
	AuthErrRequiredID
	AuthErrRequiredFields
	AuthErrNoAccess
	AuthErrForbidden
	AuthErrLogin
	AuthErrRefreshToken
)

// defaultAuthSessionTTL is the fallback refresh token session lifetime when no
// session TTL is configured.
const defaultAuthSessionTTL = 30 * 24 * time.Hour

// sets error text definition.
func init() {
	appErrorText[AuthErrNotFound] = "auth not found"
	appErrorText[AuthErrRequiredID] = "auth id is required"
	appErrorText[AuthErrRequiredFields] = "auth fields are required"
	appErrorText[AuthErrNoAccess] = "user has no access"
	appErrorText[AuthErrForbidden] = "user has no access rights"
	appErrorText[AuthErrLogin] = "invalid login credentials"
	appErrorText[AuthErrRefreshToken] = "invalid or revoked refresh token"
}

type (
	// Auth represents access authorization.
	Auth struct {
		ID        string     `json:"id"            db:"id,omitempty"`
		UserID    string     `json:"user_id"       db:"user_id,indexed,omitempty"  valid:"required"`
		Username  string     `json:"username"      db:"username,indexed,omitempty" valid:"required"`
		Password  string     `json:"-"             db:"password,omitempty"         valid:"required"`
		CreatedAt *time.Time `json:"created_at"    db:"created_at,omitempty"`
		UpdatedAt *time.Time `json:"updated_at"    db:"updated_at,omitempty"`
	}

	// AuthSession represents a login session and its refresh token.
	AuthSession struct {
		ID           string     `json:"id"            db:"id,omitempty"`
		AuthID       string     `json:"auth_id"       db:"auth_id,indexed,omitempty"`
		UserID       string     `json:"user_id"       db:"user_id,indexed,omitempty"`
		RefreshToken string     `json:"refresh_token" db:"refresh_token,indexed,omitempty"`
		ExpiresAt    time.Time  `json:"expires_at"    db:"expires_at,omitempty"`
		CreatedAt    *time.Time `json:"created_at"    db:"created_at,omitempty"`
		UpdatedAt    *time.Time `json:"updated_at"    db:"updated_at,omitempty"`
	}

	// authRepository defines operation for auth records.	// authRepository defines operation for auth records.
	authRepository interface {
		// Get returns an auth details by id from data store.
		Get(ctx context.Context, id string) (*Auth, error)

		// GetByUsername returns an auth details by username from data store.
		GetByUsername(ctx context.Context, username string) (*Auth, error)

		// GetByUsernameAndPassword returns an auth details by username and password from data store.
		GetByUsernameAndPassword(ctx context.Context, username, password string) (*Auth, error)

		// Create persists a new auth to data store.
		Create(context.Context, *Auth) error

		// Update persists auth changes to data store.
		Update(context.Context, *Auth) error
	}

	// authSessionRepository defines operation for login sessions.
	authSessionRepository interface {
		// GetByRefreshToken returns a session by its refreshToken from data store.
		GetByRefreshToken(ctx context.Context, refreshToken string) (*AuthSession, error)

		// Create persists a new session to data store.
		Create(context.Context, *AuthSession) error

		// Update persists session changes to data store.
		Update(context.Context, *AuthSession) error

		// Delete removes a session from data store.
		Delete(ctx context.Context, id string) error
	}
)

type ctxKey int

const authKey ctxKey = iota

// AuthToContext sets auth details to context.
func AuthToContext(parent context.Context, au *Auth) context.Context {
	return context.WithValue(parent, authKey, au)
}

// AuthFromContext returns an auth details from the given context if one is present.
// Return nil if auth detail cannot be found.
func AuthFromContext(ctx context.Context) *Auth {
	if ctx == nil {
		return nil
	}
	if au, ok := ctx.Value(authKey).(*Auth); ok {
		return au
	}
	return nil
}

// NewAuthService returns a new Auth service.
func NewAuthService(
	salt string,
	sessionTTL time.Duration,
	sc SteamClient,
	as authRepository,
	ss authSessionRepository,
	us *UserService,
	logger *slog.Logger,
) *AuthService {
	if sessionTTL <= 0 {
		sessionTTL = defaultAuthSessionTTL
	}
	return &AuthService{salt, sessionTTL, sc, as, ss, us, logger}
}

type AuthService struct {
	salt       string
	sessionTTL time.Duration

	steamClient SteamClient
	authRepo    authRepository
	sessionRepo authSessionRepository
	userSvc     *UserService
	logger      *slog.Logger
}

func (s *AuthService) SteamLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (*Auth, string, error) {
	// Handle authorization redirect.
	if r.URL.Query().Get("openid.mode") == "" {
		url, err := s.steamClient.AuthorizeURL(r)
		if err != nil {
			return nil, "", err
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect) //nolint:gosec // steam openid login redirect flow
		return nil, "", nil
	}

	// Validates auth and get player details and use SteamID as auth username.
	steamPlayer, err := s.steamClient.Authenticate(r)
	if err != nil {
		return nil, "", fmt.Errorf("steam player not found: %s", err)
	}

	// Check account existence.
	authData, err := s.authRepo.GetByUsername(ctx, steamPlayer.ID)
	if err != nil && !errors.Is(err, AuthErrNotFound) {
		return nil, "", fmt.Errorf("auth not found: %s", err)
	}

	// Account existed and checked login credentials.
	if authData != nil {
		// try logging with both password v1 and v2
		passwordV1 := s.composePassword(steamPlayer.ID, authData.UserID)
		passwordV2 := s.composePasswordV2(steamPlayer.ID, authData.UserID)
		if authData.Password != passwordV1 && authData.Password != passwordV2 {
			return nil, "", AuthErrLogin
		}

		// upgrade password to v2
		if authData.Password == passwordV1 {
			if err := s.authRepo.Update(ctx, &Auth{
				ID:       authData.ID,
				Password: passwordV2,
			}); err != nil {
				s.logger.ErrorContext(ctx, "upgrade password v2 failed", "error", err)
			}
		}

		u, err := s.userSvc.User(ctx, authData.UserID)
		if err != nil {
			if errors.Is(err, UserErrNotFound) {
				s.logger.WarnContext(ctx, "user not found, but auth exists. re-creating user.",
					"auth_id", authData.ID,
					"user_id", authData.UserID,
					"username", authData.Username,
				)
				// create user data with auth user id and creation date to preserve previous user data.
				u = &User{
					ID:        authData.UserID,
					CreatedAt: authData.CreatedAt,
					SteamID:   steamPlayer.ID,
					Name:      steamPlayer.Name,
					URL:       steamPlayer.URL,
					Avatar:    steamPlayer.Avatar,
				}
				if err = s.userSvc.Create(ctx, u); err != nil {
					return nil, "", err
				}

				s.logger.DebugContext(ctx, "user re-created", "user", u)

				refreshToken, err := s.newSession(ctx, authData)
				if err != nil {
					return nil, "", err
				}
				return authData, refreshToken, nil
			}

			return nil, "", err
		}
		if err = u.CheckStatus(); err != nil {
			return nil, "", err
		}
		if _, err = s.userSvc.SteamSync(ctx, steamPlayer); err != nil {
			return nil, "", UserErrSteamSync.X(err)
		}

		refreshToken, err := s.newSession(ctx, authData)
		if err != nil {
			return nil, "", err
		}
		return authData, refreshToken, nil
	}

	// HOTFIX for missing user data
	existingUser, _ := s.userSvc.User(ctx, steamPlayer.ID)

	// Process account registration and save details.
	authData, refreshToken, err := s.createAccountFromSteam(ctx, steamPlayer, existingUser)
	if err != nil {
		return nil, "", err
	}

	return authData, refreshToken, nil
}

// RefreshToken validates a session by its refresh token, rotates it to a new
// one and returns the auth details with the new raw refresh token.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*Auth, string, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, "", AuthErrRefreshToken
	}

	sess, err := s.sessionRepo.GetByRefreshToken(ctx, s.hash(refreshToken))
	if err != nil {
		return nil, "", AuthErrRefreshToken
	}
	if sess == nil || !sess.ExpiresAt.After(time.Now()) {
		// purge expired session on access attempt.
		if sess != nil {
			if err := s.sessionRepo.Delete(ctx, sess.ID); err != nil {
				s.logger.ErrorContext(ctx, "delete expired session failed", "error", err)
			}
		}
		return nil, "", AuthErrRefreshToken
	}

	newToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, "", err
	}
	sess.RefreshToken = s.hash(newToken)
	sess.ExpiresAt = time.Now().Add(s.sessionTTL)
	if err := s.sessionRepo.Update(ctx, sess); err != nil {
		return nil, "", err
	}

	au, err := s.authRepo.Get(ctx, sess.AuthID)
	if err != nil {
		return nil, "", AuthErrRefreshToken
	}

	return au, newToken, nil
}

// RevokeRefreshToken invalidates the session of the given refresh token so it
// can no longer be renewed. Other sessions of the same user stay valid.
func (s *AuthService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return AuthErrRefreshToken
	}

	sess, err := s.sessionRepo.GetByRefreshToken(ctx, s.hash(refreshToken))
	if err != nil {
		return AuthErrRefreshToken
	}
	if sess == nil {
		return AuthErrRefreshToken
	}

	return s.sessionRepo.Delete(ctx, sess.ID)
}

func (s *AuthService) Auth(ctx context.Context, id string) (*Auth, error) {
	u, err := s.authRepo.Get(ctx, id)
	if err != nil {
		return nil, AuthErrNotFound.X(err)
	}

	return u, nil
}

// newSession creates a new login session for the given auth and returns the
// raw refresh token bound to it.
func (s *AuthService) newSession(ctx context.Context, au *Auth) (string, error) {
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return "", err
	}
	sess := &AuthSession{
		AuthID:       au.ID,
		UserID:       au.UserID,
		RefreshToken: s.hash(refreshToken),
		ExpiresAt:    time.Now().Add(s.sessionTTL),
	}
	if err := s.sessionRepo.Create(ctx, sess); err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (s *AuthService) createAccountFromSteam(ctx context.Context, sp *SteamPlayer, user *User) (*Auth, string, error) {
	if user == nil {
		user = &User{
			SteamID: sp.ID,
			Name:    sp.Name,
			URL:     sp.URL,
			Avatar:  sp.Avatar,
		}
		if err := s.userSvc.Create(ctx, user); err != nil {
			return nil, "", err
		}
	}

	au := &Auth{UserID: user.ID, Username: sp.ID}
	au.Password = s.composePasswordV2(sp.ID, user.ID)
	if err := s.authRepo.Create(ctx, au); err != nil {
		return nil, "", err
	}

	refreshToken, err := s.newSession(ctx, au)
	if err != nil {
		return nil, "", err
	}
	return au, refreshToken, nil
}

func (s *AuthService) generateRefreshToken() (string, error) {
	buf := make([]byte, 48)
	if _, err := rand.Read(buf); err != nil {
		// try again
		if _, err := rand.Read(buf); err != nil {
			return "", err
		}
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func (s *AuthService) composePassword(steamID, userID string) string {
	h := sha1.New() //nolint:gosec // legacy password hashing for backward compatibility
	h.Write([]byte(steamID + userID + s.salt))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *AuthService) composePasswordV2(steamID, userID string) string {
	return s.hash(steamID, userID)
}

func (s *AuthService) hash(a ...string) string {
	h := sha256.New()
	a = append(a, s.salt)
	h.Write([]byte(strings.Join(a, "")))
	return hex.EncodeToString(h.Sum(nil))
}
