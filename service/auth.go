package service

import (
	"net/http"
	"strings"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
)

// NewAuth returns a new Auth service.
func NewAuth(sc core.SteamClient, as core.AuthStorage, us core.UserService) core.AuthService {
	return &authService{sc, as, us}
}

type authService struct {
	steamClient core.SteamClient
	authStg     core.AuthStorage
	userSvc     core.UserService
}

func (s *authService) SteamLogin(w http.ResponseWriter, r *http.Request) (*core.Auth, error) {
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
		return nil, err
	}

	// Check account existence.
	au, err := s.authStg.GetByUsername(steamPlayer.ID)
	if err != nil && err != core.AuthErrNotFound {
		return nil, err
	}
	// Account existed and checks login credentials.
	if au != nil {
		if au.Password != au.ComposePassword(steamPlayer.ID, au.UserID) {
			return nil, core.AuthErrLogin
		}

		u, _ := s.userSvc.User(au.UserID)
		if err = u.CheckStatus(); err != nil {
			return nil, err
		}

		if _, err = s.userSvc.SteamSync(steamPlayer); err != nil {
			return nil, errors.New(core.UserErrSteamSync, err)
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

func (s *authService) RenewToken(refreshToken string) (*core.Auth, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, core.AuthErrRefreshToken
	}

	au, err := s.authStg.GetByRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New(core.AuthErrRefreshToken, err)
	}

	return au, nil
}

func (s *authService) RevokeRefreshToken(refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return core.AuthErrRefreshToken
	}

	au, err := s.RenewToken(refreshToken)
	if err != nil {
		return err
	}

	au.GenerateRefreshToken()
	return s.authStg.Update(au)
}

func (s *authService) Auth(id string) (*core.Auth, error) {
	u, err := s.authStg.Get(id)
	if err != nil {
		return nil, errors.New(core.AuthErrNotFound, err)
	}

	return u, nil
}

func (s *authService) createAccountFromSteam(sp *core.SteamPlayer) (*core.Auth, error) {
	u := &core.User{
		SteamID: sp.ID,
		Name:    sp.Name,
		URL:     sp.URL,
		Avatar:  sp.Avatar,
	}
	if err := s.userSvc.Create(u); err != nil {
		return nil, err
	}

	au := &core.Auth{UserID: u.ID, Username: sp.ID}
	au.SetDefaults()
	au.Password = au.ComposePassword(sp.ID, u.ID)
	if err := s.authStg.Create(au); err != nil {
		return nil, err
	}

	return au, nil
}
