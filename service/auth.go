package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/xerrors"
)

// NewAuth returns a new Auth service.
func NewAuth(salt string, sc dotagiftx.SteamClient, as dotagiftx.AuthStorage, us dotagiftx.UserService) dotagiftx.AuthService {
	return &authService{salt, sc, as, us}
}

type authService struct {
	salt        string
	steamClient dotagiftx.SteamClient
	authStg     dotagiftx.AuthStorage
	userSvc     dotagiftx.UserService
}

func (s *authService) SteamLogin(w http.ResponseWriter, r *http.Request) (*dotagiftx.Auth, error) {
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
	if err != nil && !errors.Is(err, dotagiftx.AuthErrNotFound) {
		return nil, fmt.Errorf("auth not found: %s", err)
	}

	// Account existed and checked login credentials.
	if au != nil {
		if au.Password != s.composePassword(steamPlayer.ID, au.UserID) {
			return nil, dotagiftx.AuthErrLogin
		}

		u, _ := s.userSvc.User(au.UserID)
		if err = u.CheckStatus(); err != nil {
			return nil, err
		}

		if _, err = s.userSvc.SteamSync(steamPlayer); err != nil {
			return nil, xerrors.New(dotagiftx.UserErrSteamSync, err)
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

func (s *authService) RenewToken(refreshToken string) (*dotagiftx.Auth, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, dotagiftx.AuthErrRefreshToken
	}

	au, err := s.authStg.GetByRefreshToken(refreshToken)
	if err != nil {
		return nil, xerrors.New(dotagiftx.AuthErrRefreshToken, err)
	}

	return au, nil
}

func (s *authService) RevokeRefreshToken(refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return dotagiftx.AuthErrRefreshToken
	}

	au, err := s.RenewToken(refreshToken)
	if err != nil {
		return err
	}

	au.RefreshToken = s.generateRefreshToken()
	return s.authStg.Update(au)
}

func (s *authService) Auth(id string) (*dotagiftx.Auth, error) {
	u, err := s.authStg.Get(id)
	if err != nil {
		return nil, xerrors.New(dotagiftx.AuthErrNotFound, err)
	}

	return u, nil
}

func (s *authService) createAccountFromSteam(sp *dotagiftx.SteamPlayer) (*dotagiftx.Auth, error) {
	u := &dotagiftx.User{
		SteamID: sp.ID,
		Name:    sp.Name,
		URL:     sp.URL,
		Avatar:  sp.Avatar,
	}
	if err := s.userSvc.Create(u); err != nil {
		return nil, err
	}

	au := &dotagiftx.Auth{UserID: u.ID, Username: sp.ID}
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
