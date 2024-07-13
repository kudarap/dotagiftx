package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
)

// NewUser returns a new User service.
func NewUser(us dgx.UserStorage, fm dgx.FileManager, sc subscriptionChecker) dgx.UserService {
	return &userService{us, fm, sc}
}

type userService struct {
	userStg     dgx.UserStorage
	fileMgr     dgx.FileManager
	subsChecker subscriptionChecker
}

func (s *userService) Users(opts dgx.FindOpts) ([]dgx.User, error) {
	return s.userStg.Find(opts)
}

func (s *userService) FlaggedUsers(opts dgx.FindOpts) ([]dgx.User, error) {
	return s.userStg.FindFlagged(opts)
}

func (s *userService) User(id string) (*dgx.User, error) {
	return s.userStg.Get(id)
}

func (s *userService) UserFromContext(ctx context.Context) (*dgx.User, error) {
	au := dgx.AuthFromContext(ctx)
	if au == nil {
		return nil, dgx.UserErrNotFound
	}

	return s.User(au.UserID)
}

func (s *userService) Create(u *dgx.User) error {
	url, err := s.downloadProfileImage(u.Avatar)
	if err != nil {
		return errors.New(dgx.UserErrProfileImageDL, err)
	}
	u.Avatar = url

	if err = u.CheckCreate(); err != nil {
		return err
	}

	go func() {
		err := pingGoogleSitemap()
		if err != nil {
			log.Println("pingGoogleSitemap err:", err)
		}
	}()

	return s.userStg.Create(u)
}

func (s *userService) Update(ctx context.Context, u *dgx.User) error {
	au := dgx.AuthFromContext(ctx)
	if au == nil {
		return dgx.AuthErrNoAccess
	}
	u.ID = au.UserID

	if err := u.CheckUpdate(); err != nil {
		return err
	}

	return s.userStg.Update(u)
}

func (s *userService) SteamSync(sp *dgx.SteamPlayer) (*dgx.User, error) {
	opts := dgx.FindOpts{Filter: dgx.User{SteamID: sp.ID}, IndexKey: "steam_id"}
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

func (s *userService) ProcessSubscription(ctx context.Context, subscriptionID string) (*dgx.User, error) {
	au := dgx.AuthFromContext(ctx)
	if au == nil {
		return nil, dgx.AuthErrNoAccess
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
	userSubs := dgx.UserSubscriptionFromString(plan)
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

func (s *userService) ProcessManualSubscription(
	ctx context.Context, p dgx.ManualSubscriptionParam,
) (*dgx.User, error) {
	panic("implement me")
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
