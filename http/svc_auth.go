package http

import (
	"context"
	"net/http"
	"time"

	"github.com/kudarap/dotagiftx/dotagiftx"
)

const defaultTokenExpiration = time.Minute * 5

// authService provides access to auth service methods used by http handlers.
type authService interface {
	// SteamLogin redirects for authorization and process creation of auth.
	// Returns the auth details and its raw refresh token.
	SteamLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (*dotagiftx.Auth, error)

	// RefreshToken checks refresh token validity that allows to get new short-lived
	// access token and rotates the refresh token to a new one.
	RefreshToken(ctx context.Context, refreshToken string) (*dotagiftx.Auth, error)

	// RevokeRefreshToken invalidates refresh token that will prevent on renewing
	// short-lived access token and will result user have to re-login.
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
}

type authResp struct {
	UserID       string    `json:"user_id,omitempty"`
	SteamID      string    `json:"steam_id,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Token        string    `json:"token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func handleAuthSteam(svc authService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle steam auth.
		au, err := svc.SteamLogin(r.Context(), w, r)
		if err != nil {
			respondError(w, err)
			return
		}
		// Returning nil auth without error means it redirect for
		// authorization
		if au == nil {
			return
		}

		// Compose new JWT.
		a, err := newAuth(au, au.RefreshToken)
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, a)
	}
}

func handleAuthRenew(svc authService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := new(struct {
			RefreshToken string `json:"refresh_token"`
		})
		if err := parseForm(r, form); err != nil {
			respondError(w, err)
			return
		}

		au, err := svc.RefreshToken(r.Context(), form.RefreshToken)
		if err != nil {
			respond(w, http.StatusUnauthorized, newError(err))
			return
		}

		// Refresh JWT and rotate refresh token.
		a, err := newAuth(au, au.RefreshToken)
		if err != nil {
			respond(w, http.StatusInternalServerError, newError(err))
			return
		}

		respondOK(w, a)
	}
}

func handleAuthRevoke(svc authService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := new(struct {
			RefreshToken string `json:"refresh_token"`
		})
		if err := parseForm(r, form); err != nil {
			respondError(w, err)
			return
		}

		if err := svc.RevokeRefreshToken(r.Context(), form.RefreshToken); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, struct {
			Msg string `json:"msg"`
		}{
			"refresh token successfully revoked",
		})
	}
}

func newAuth(au *dotagiftx.Auth, refreshToken string) (*authResp, error) {
	a, err := refreshJWT(au)
	if err != nil {
		return nil, err
	}

	a.UserID = au.UserID
	a.SteamID = au.Username
	a.RefreshToken = refreshToken
	return a, nil
}

const noLevel = ""

func refreshJWT(au *dotagiftx.Auth) (*authResp, error) {
	a := &authResp{}
	a.ExpiresAt = time.Now().Add(defaultTokenExpiration)

	t, err := newAccessToken(au.UserID, noLevel, a.ExpiresAt)
	if err != nil {
		return nil, err
	}
	a.Token = t

	return a, nil
}
