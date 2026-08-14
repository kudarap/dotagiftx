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
	SteamLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (*dotagiftx.Auth, error)

	// RevokeRefreshToken invalidates refresh token that will prevent on renewing
	// short-lived access token and will result user have to re-login.
	RevokeRefreshToken(ctx context.Context, refreshToken string) error

	// RefreshToken checks refresh token validity that allows to get new short-lived access token.
	RefreshToken(ctx context.Context, refreshToken string) (*dotagiftx.Auth, error)
}

type authResp struct {
	UserID       string    `json:"user_id,omitempty"`
	SteamID      string    `json:"steam_id,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	AccessToken  string    `json:"token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func handleAuthSteam(svc authService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle steam auth.
		auth, err := svc.SteamLogin(r.Context(), w, r)
		if err != nil {
			respondError(w, err)
			return
		}
		// Returning nil auth without error means it redirect for
		// authorization
		if auth == nil {
			return
		}

		// Compose new JWT.
		res, err := newAuth(auth)
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, res)
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

		auth, err := svc.RefreshToken(r.Context(), form.RefreshToken)
		if err != nil {
			respond(w, http.StatusUnauthorized, newError(err))
			return
		}

		// Refresh JWT.
		res, err := refreshJWT(auth)
		if err != nil {
			respond(w, http.StatusInternalServerError, newError(err))
			return
		}

		respondOK(w, res)
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

func newAuth(au *dotagiftx.Auth) (*authResp, error) {
	res, err := refreshJWT(au)
	if err != nil {
		return nil, err
	}

	res.UserID = au.UserID
	res.SteamID = au.Username
	res.RefreshToken = au.RefreshToken
	return res, nil
}

const noLevel = ""

func refreshJWT(au *dotagiftx.Auth) (*authResp, error) {
	res := &authResp{}
	res.ExpiresAt = time.Now().Add(defaultTokenExpiration)

	t, err := newAccessToken(au.UserID, noLevel, res.ExpiresAt)
	if err != nil {
		return nil, err
	}
	res.AccessToken = t
	return res, nil
}
