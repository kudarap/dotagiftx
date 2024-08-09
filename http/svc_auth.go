package http

import (
	"net/http"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/http/jwt"
)

const defaultTokenExpiration = time.Minute * 5

type authResp struct {
	UserID       string    `json:"user_id,omitempty"`
	SteamID      string    `json:"steam_id,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Token        string    `json:"token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
}

func handleAuthSteam(svc dgx.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle steam auth.
		au, err := svc.SteamLogin(w, r)
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
		a, err := newAuth(au)
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, a)
	}
}

func handleAuthRenew(svc dgx.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := new(struct {
			RefreshToken string `json:"refresh_token"`
		})
		if err := parseForm(r, form); err != nil {
			respondError(w, err)
			return
		}

		au, err := svc.RenewToken(form.RefreshToken)
		if err != nil {
			respond(w, http.StatusUnauthorized, newError(err))
			return
		}

		// Refresh JWT.
		a, err := refreshJWT(au)
		if err != nil {
			respond(w, http.StatusInternalServerError, newError(err))
			return
		}

		respondOK(w, a)
	}
}

func handleAuthRevoke(svc dgx.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := new(struct {
			RefreshToken string `json:"refresh_token"`
		})
		if err := parseForm(r, form); err != nil {
			respondError(w, err)
			return
		}

		if err := svc.RevokeRefreshToken(form.RefreshToken); err != nil {
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

func newAuth(au *dgx.Auth) (*authResp, error) {
	a, err := refreshJWT(au)
	if err != nil {
		return nil, err
	}

	a.UserID = au.UserID
	a.SteamID = au.Username
	a.RefreshToken = au.RefreshToken
	return a, nil
}

const noLevel = ""

func refreshJWT(au *dgx.Auth) (*authResp, error) {
	a := &authResp{}
	a.ExpiresAt = time.Now().Add(defaultTokenExpiration)

	t, err := jwt.New(au.UserID, noLevel, a.ExpiresAt)
	if err != nil {
		return nil, err
	}
	a.Token = t

	return a, nil
}
