package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
)

// NewUser returns a new User service.
func NewUser(us dotagiftx.UserStorage, fm dotagiftx.FileManager, sc subscriptionChecker) dotagiftx.UserService {
	return &userService{us, fm, sc}
}

type userService struct {
	userStg     dotagiftx.UserStorage
	fileMgr     dotagiftx.FileManager
	subsChecker subscriptionChecker
}

func (s *userService) Users(opts dotagiftx.FindOpts) ([]dotagiftx.User, error) {
	return s.userStg.Find(opts)
}

func (s *userService) FlaggedUsers(opts dotagiftx.FindOpts) ([]dotagiftx.User, error) {
	return s.userStg.FindFlagged(opts)
}

func (s *userService) User(id string) (*dotagiftx.User, error) {
	return s.userStg.Get(id)
}

func (s *userService) UserFromContext(ctx context.Context) (*dotagiftx.User, error) {
	au := dotagiftx.AuthFromContext(ctx)
	if au == nil {
		return nil, dotagiftx.UserErrNotFound
	}

	return s.User(au.UserID)
}

func (s *userService) Create(u *dotagiftx.User) error {
	url, err := s.downloadProfileImage(u.Avatar)
	if err != nil {
		return errors.New(dotagiftx.UserErrProfileImageDL, err)
	}
	u.Avatar = url

	if err = u.CheckCreate(); err != nil {
		return err
	}

	go pingGoogleSitemap()

	return s.userStg.Create(u)
}

func (s *userService) Update(ctx context.Context, u *dotagiftx.User) error {
	au := dotagiftx.AuthFromContext(ctx)
	if au == nil {
		return dotagiftx.AuthErrNoAccess
	}
	u.ID = au.UserID

	if err := u.CheckUpdate(); err != nil {
		return err
	}

	return s.userStg.Update(u)
}

func (s *userService) SteamSync(sp *dotagiftx.SteamPlayer) (*dotagiftx.User, error) {
	opts := dotagiftx.FindOpts{Filter: dotagiftx.User{SteamID: sp.ID}, IndexKey: "steam_id"}
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

func (s *userService) ProcSubscription(ctx context.Context, subscriptionID string) (*dotagiftx.User, error) {
	au := dotagiftx.AuthFromContext(ctx)
	if au == nil {
		return nil, dotagiftx.AuthErrNoAccess
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
	userSubs := dotagiftx.UserSubscriptionFromString(plan)
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
