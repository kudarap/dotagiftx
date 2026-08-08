package dotagiftx

import (
	"context"
	"crypto/sha1"
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
		ID           string     `json:"id"            db:"id,omitempty"`
		UserID       string     `json:"user_id"       db:"user_id,indexed,omitempty"  valid:"required"`
		Username     string     `json:"username"      db:"username,indexed,omitempty" valid:"required"`
		Password     string     `json:"-"             db:"password,omitempty"         valid:"required"`
		RefreshToken string     `json:"refresh_token" db:"refresh_token,indexed,omitempty"`
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

		// GetByRefreshToken returns an auth details by refreshToken from data store.
		GetByRefreshToken(ctx context.Context, refreshToken string) (*Auth, error)

		// Create persists a new auth to data store.
		Create(context.Context, *Auth) error

		// Update persists auth changes to data store.
		Update(context.Context, *Auth) error
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
	sc SteamClient,
	as authRepository,
	us *UserService,
	logger *slog.Logger,
) *AuthService {
	return &AuthService{salt, sc, as, us, logger}
}

type AuthService struct {
	salt string

	steamClient SteamClient
	authRepo    authRepository
	userSvc     *UserService
	logger      *slog.Logger
}

func (s *AuthService) SteamLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (*Auth, error) {
	// Handle authorization redirect.
	if r.URL.Query().Get("openid.mode") == "" {
		url, err := s.steamClient.AuthorizeURL(r)
		if err != nil {
			return nil, err
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil, nil
	}

	// Validates auth and get player details and use SteamID as auth username.
	steamPlayer, err := s.steamClient.Authenticate(r)
	if err != nil {
		return nil, fmt.Errorf("steam player not found: %s", err)
	}

	// Check account existence.
	authData, err := s.authRepo.GetByUsername(ctx, steamPlayer.ID)
	if err != nil && !errors.Is(err, AuthErrNotFound) {
		return nil, fmt.Errorf("auth not found: %s", err)
	}

	// Account existed and checked login credentials.
	if authData != nil {
		if authData.Password != s.composePassword(steamPlayer.ID, authData.UserID) {
			return nil, AuthErrLogin
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
					return nil, err
				}

				s.logger.DebugContext(ctx, "user re-created", "user", u)
				return authData, nil
			}

			return nil, err
		}
		if err = u.CheckStatus(); err != nil {
			return nil, err
		}
		if _, err = s.userSvc.SteamSync(ctx, steamPlayer); err != nil {
			return nil, UserErrSteamSync.X(err)
		}

		return authData, nil
	}

	// HOTFIX for missing user data
	existingUser, _ := s.userSvc.User(ctx, steamPlayer.ID)

	// Process account registration and save details.
	authData, err = s.createAccountFromSteam(ctx, steamPlayer, existingUser)
	if err != nil {
		return nil, err
	}

	return authData, nil
}

func (s *AuthService) RenewToken(ctx context.Context, refreshToken string) (*Auth, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, AuthErrRefreshToken
	}

	au, err := s.authRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, AuthErrRefreshToken
	}

	return au, nil
}

func (s *AuthService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return AuthErrRefreshToken
	}

	au, err := s.RenewToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	au.RefreshToken = s.generateRefreshToken()
	return s.authRepo.Update(ctx, au)
}

func (s *AuthService) Auth(ctx context.Context, id string) (*Auth, error) {
	u, err := s.authRepo.Get(ctx, id)
	if err != nil {
		return nil, AuthErrNotFound.X(err)
	}

	return u, nil
}

func (s *AuthService) createAccountFromSteam(ctx context.Context, sp *SteamPlayer, user *User) (*Auth, error) {
	if user == nil {
		user = &User{
			SteamID: sp.ID,
			Name:    sp.Name,
			URL:     sp.URL,
			Avatar:  sp.Avatar,
		}
		if err := s.userSvc.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	au := &Auth{UserID: user.ID, Username: sp.ID}
	au.RefreshToken = s.generateRefreshToken()
	au.Password = s.composePassword(sp.ID, user.ID)
	if err := s.authRepo.Create(ctx, au); err != nil {
		return nil, err
	}

	return au, nil
}

func (s *AuthService) generateRefreshToken() string {
	t := fmt.Sprintf("%d%s", time.Now().UnixNano(), s.salt)
	h := sha1.New()
	h.Write([]byte(t))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *AuthService) composePassword(steamID, userID string) string {
	h := sha1.New()
	h.Write([]byte(steamID + userID + s.salt))
	return hex.EncodeToString(h.Sum(nil))
}
