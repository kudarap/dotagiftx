package dgx

import (
	"context"
	"net/http"
	"time"

	"github.com/kudarap/dotagiftx/gokit/hash"
)

// Auth error types.
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

// SetDefaults sets auth default values.
func (a *Auth) SetDefaults() {
	a.RefreshToken = a.GenerateRefreshToken()
}

// GenerateRefreshToken generates new refresh token.
func (a *Auth) GenerateRefreshToken() string {
	a.RefreshToken = hash.GenerateSha1()
	return a.RefreshToken
}

// ComposePassword returns composed password.
func (Auth) ComposePassword(steamID, userID string) string {
	return hash.Sha1(steamID + userID)
}

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

// DO NOT CHANGE THIS! unless you know what your doing.
// changing the value of the salt will invalidate all user logins.
func init() {
	hash.Salt = "0stari0n"
}
