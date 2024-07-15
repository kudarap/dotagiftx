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
	u, err := s.userStg.Get(sp.ID)
	if err != nil {
		return nil, err
	}

	u.Name = sp.Name
	u.URL = sp.URL
	u.Avatar, err = s.downloadProfileImage(sp.Avatar)
	if err != nil {
		return nil, err
	}
	if err = s.userStg.Update(u); err != nil {
		return nil, err
	}
	return u, nil
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
	user.SubscriptionType = "paypal"
	if err = user.CheckUpdate(); err != nil {
		return nil, err
	}

	return user, s.userStg.Update(user)
}

// UpdateSubscriptionFromWebhook manage updates from webhook payload, most often use in incrementing cycles or
// extending expiration.
func (s *userService) UpdateSubscriptionFromWebhook(ctx context.Context, r *http.Request) (*dgx.User, error) {
	// get user by steam id and increment their cycles.
	steamID, cancelled, err := s.subsChecker.IsCancelled(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("checking cancelled subscription: %v", err)
	}
	if !cancelled {
		// ignore if not cancelled.
		log.Println("ignoring subscription update because its not cancelled:", steamID)
		return nil, nil
	}

	log.Println("cancelling subscription", steamID, "by marking expiration")
	user, err := s.userStg.Get(steamID)
	if err != nil {
		return nil, fmt.Errorf("getting user: %v", err)
	}
	expiresAt := user.SubscribedAt.AddDate(0, 1, 0)
	user.SubscriptionEndsAt = &expiresAt
	if err = s.userStg.Update(user); err != nil {
		return nil, fmt.Errorf("updating user: %v", err)
	}
	return user, nil
}

// ProcessManualSubscription process manual subscription such as one-time payments that process manually, normally
// in bulk and steam items. This function will be used non-recurring payments. ex:
//
//		Manual Partner subscription:
//	    - 3 months (+60% overhead)
//	    - 6 months (+60% overhead)
//	    - 12 months (+60% overhead)
func (s *userService) ProcessManualSubscription(
	ctx context.Context, param dgx.ManualSubscriptionParam,
) (*dgx.User, error) {
	user, err := s.userStg.Get(param.UserID)
	if err != nil {
		return nil, fmt.Errorf("getting user: %v", err)
	}

	subs := dgx.UserSubscriptionFromString(param.Plan)
	user.Subscription = subs
	user.Boons = subs.Boons()
	user.SubscriptionType = "manual"

	now := time.Now()
	end := now.AddDate(0, param.Cycles, 0)
	user.SubscribedAt = &now
	user.SubscriptionEndsAt = &end
	if err = s.userStg.Update(user); err != nil {
		return nil, fmt.Errorf("updating user: %v", err)
	}
	return user, nil
}

// downloadProfileImage saves image file from url.
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
	IsCancelled(ctx context.Context, r *http.Request) (steamID string, cancelled bool, err error)
}
