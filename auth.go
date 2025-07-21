package dotagiftx

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	AuthErrNotFound Errors = iota + 1000
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

	// AuthService provides access to service.
	AuthService interface {
		// SteamLogin redirects for authorization and process creation of auth.
		SteamLogin(w http.ResponseWriter, r *http.Request) (*Auth, error)

		// RevokeRefreshToken invalidates refresh token that will prevent on renewing
		// short-lived access token and will result user have to re-login.
		RevokeRefreshToken(refreshToken string) error

		// RenewToken checks refresh token validity that allows to get new short-lived access token.
		RenewToken(refreshToken string) (*Auth, error)

		// Auth returns an auth details by id.
		Auth(id string) (*Auth, error)
	}

	// AuthStorage defines operation for auth records.
	AuthStorage interface {
		// Get returns an auth details by id from data store.
		Get(id string) (*Auth, error)

		// GetByUsername returns an auth details by username from data store.
		GetByUsername(username string) (*Auth, error)

		// GetByUsernameAndPassword returns an auth details by username and password from data store.
		GetByUsernameAndPassword(username, password string) (*Auth, error)

		// GetByRefreshToken returns an auth details by refreshToken from data store.
		GetByRefreshToken(refreshToken string) (*Auth, error)

		// Create persists a new auth to data store.
		Create(*Auth) error

		// Update persists auth changes to data store.
		Update(*Auth) error
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

// NewAuth returns a new Auth service.
func NewAuth(
	salt string,
	sc SteamClient,
	as AuthStorage,
	us UserService,
) AuthService {
	return &authService{sc, as, us, salt}
}

type authService struct {
	steamClient SteamClient
	authStg     AuthStorage
	userSvc     UserService

	salt string
}

func (s *authService) SteamLogin(w http.ResponseWriter, r *http.Request) (*Auth, error) {
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
	au, err := s.authStg.GetByUsername(steamPlayer.ID)
	if err != nil && !errors.Is(err, AuthErrNotFound) {
		return nil, fmt.Errorf("auth not found: %s", err)
	}

	// Account existed and checked login credentials.
	if au != nil {
		if au.Password != s.composePassword(steamPlayer.ID, au.UserID) {
			return nil, AuthErrLogin
		}

		u, _ := s.userSvc.User(au.UserID)
		if err = u.CheckStatus(); err != nil {
			return nil, err
		}

		if _, err = s.userSvc.SteamSync(steamPlayer); err != nil {
			return nil, UserErrSteamSync.X(err)
		}

		return au, nil
	}

	// Process account registration and save details.
	au, err = s.createAccountFromSteam(steamPlayer)
	if err != nil {
		return nil, err
	}

	return au, nil
}

func (s *authService) RenewToken(refreshToken string) (*Auth, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, AuthErrRefreshToken
	}

	au, err := s.authStg.GetByRefreshToken(refreshToken)
	if err != nil {
		return nil, AuthErrRefreshToken
	}

	return au, nil
}

func (s *authService) RevokeRefreshToken(refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return AuthErrRefreshToken
	}

	au, err := s.RenewToken(refreshToken)
	if err != nil {
		return err
	}

	au.RefreshToken = s.generateRefreshToken()
	return s.authStg.Update(au)
}

func (s *authService) Auth(id string) (*Auth, error) {
	u, err := s.authStg.Get(id)
	if err != nil {
		return nil, AuthErrNotFound.X(err)
	}

	return u, nil
}

func (s *authService) createAccountFromSteam(sp *SteamPlayer) (*Auth, error) {
	u := &User{
		SteamID: sp.ID,
		Name:    sp.Name,
		URL:     sp.URL,
		Avatar:  sp.Avatar,
	}
	if err := s.userSvc.Create(u); err != nil {
		return nil, err
	}

	au := &Auth{UserID: u.ID, Username: sp.ID}
	au.RefreshToken = s.generateRefreshToken()
	au.Password = s.composePassword(sp.ID, u.ID)
	if err := s.authStg.Create(au); err != nil {
		return nil, err
	}

	return au, nil
}

func (s *authService) generateRefreshToken() string {
	t := fmt.Sprintf("%d%s", time.Now().UnixNano(), s.salt)
	h := sha1.New()
	h.Write([]byte(t))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *authService) composePassword(steamID, userID string) string {
	h := sha1.New()
	h.Write([]byte(steamID + userID + s.salt))
	return hex.EncodeToString(h.Sum(nil))
}
