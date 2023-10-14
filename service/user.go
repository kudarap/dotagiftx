package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
)

// NewUser returns a new User service.
func NewUser(us core.UserStorage, fm core.FileManager, sc subscriptionChecker) core.UserService {
	return &userService{us, fm, sc}
}

type userService struct {
	userStg     core.UserStorage
	fileMgr     core.FileManager
	subsChecker subscriptionChecker
}

func (s *userService) Users(opts core.FindOpts) ([]core.User, error) {
	return s.userStg.Find(opts)
}

func (s *userService) FlaggedUsers(opts core.FindOpts) ([]core.User, error) {
	return s.userStg.FindFlagged(opts)
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

	if err = u.CheckCreate(); err != nil {
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

func (s *userService) SteamSync(sp *core.SteamPlayer) (*core.User, error) {
	opts := core.FindOpts{Filter: core.User{SteamID: sp.ID}, IndexKey: "steam_id"}
	res, err := s.userStg.Find(opts)
	if err != nil {
		return nil, err
	}
	u := res[0]
	u.Name = sp.Name
	u.URL = sp.URL
	u.Avatar, err = s.downloadProfileImage(sp.Avatar)
	if err != nil {
		return nil, err
	}

	if err = s.userStg.Update(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *userService) ProcSubscription(ctx context.Context, subscriptionID string) (*core.User, error) {
	au := core.AuthFromContext(ctx)
	if au == nil {
		return nil, core.AuthErrNoAccess
	}
	user, err := s.userStg.Get(au.UserID)
	if err != nil {
		return nil, err
	}

	plan, steamID, err := s.subsChecker.Subscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	if user.SteamID != steamID {
		return nil, fmt.Errorf("could not validate subscription steam id")
	}
	userSubs := core.UserSubscriptionFromString(plan)
	if userSubs == 0 {
		return nil, fmt.Errorf("could not validate subscription plan")
	}

	if user.SubscribedAt != nil && user.Subscription == userSubs {
		return user, nil
	}

	t := time.Now()
	user.Subscription = userSubs
	user.SubscribedAt = &t
	user.Boons = userSubs.Boons()
	if err = user.CheckUpdate(); err != nil {
		return nil, err
	}

	return user, s.userStg.Update(user)
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

type subscriptionChecker interface {
	Subscription(id string) (plan, steamID string, err error)
}
