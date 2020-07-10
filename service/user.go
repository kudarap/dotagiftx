package service

import (
	"context"
	"net/http"

	"github.com/kudarap/dota2giftables/core"
	"github.com/kudarap/dota2giftables/errors"
)

// NewUser returns a new User service.
func NewUser(us core.UserStorage, fm core.FileManager) core.UserService {
	return &userService{us, fm}
}

type userService struct {
	userStg core.UserStorage
	fileMgr core.FileManager
}

func (s *userService) Users(opts core.FindOpts) ([]core.User, error) {
	return s.userStg.Find(opts)
}

func (s *userService) User(id string) (*core.User, error) {
	return s.userStg.Get(id)
}

func (s *userService) UserFromContext(ctx context.Context) (*core.User, error) {
	au := core.AuthFromContext(ctx)
	if au == nil {
		return nil, core.UserErrNotFound
	}

	return s.User(au.UserID)
}

func (s *userService) Create(u *core.User) error {
	url, err := s.downloadProfileImage(u.Avatar)
	if err != nil {
		return errors.New(core.UserErrProfileImageDL, err)
	}
	u.Avatar = url

	if err := u.CheckCreate(); err != nil {
		return err
	}

	go pingGoogleSitemap()

	return s.userStg.Create(u)
}

func (s *userService) Update(ctx context.Context, u *core.User) error {
	au := core.AuthFromContext(ctx)
	if au == nil {
		return core.AuthErrNoAccess
	}
	u.ID = au.UserID

	if err := u.CheckUpdate(); err != nil {
		return err
	}

	return s.userStg.Update(u)
}

// downloadProfileImage saves image file from a url.
func (s *userService) downloadProfileImage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	f, err := s.fileMgr.Save(resp.Body)
	if err != nil {
		return "", err
	}

	return f, nil
}
