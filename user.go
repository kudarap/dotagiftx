package dotagiftx

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// User error types.
const (
	UserErrNotFound Errors = iota + userErrorIndex
	UserErrRequiredID
	UserErrRequiredFields
	UserErrProfileImageDL
	UserErrSteamSync
	UserErrSuspended
	UserErrBanned
)

// sets error text definition.
func init() {
	appErrorText[UserErrNotFound] = "user not found"
	appErrorText[UserErrRequiredID] = "user id is required"
	appErrorText[UserErrRequiredFields] = "user fields are required"
	appErrorText[UserErrProfileImageDL] = "user profile image could not download"
	appErrorText[UserErrSteamSync] = "user profile steam sync error"
	appErrorText[UserErrSuspended] = "account has been suspended due to scam report"
	appErrorText[UserErrBanned] = "account has been banned due to scam incident"
}

// User statuses.
const (
	UserStatusSuspended UserStatus = 300
	UserStatusBanned    UserStatus = 400
)

const (
	UserSubscriptionSupporter UserSubscription = 100
	UserSubscriptionTrader    UserSubscription = 101
	UserSubscriptionPartner   UserSubscription = 109
)

type (
	UserStatus uint

	UserSubscription uint

	// User represents user information.
	User struct {
		ID        string     `json:"id"         db:"id,omitempty"`
		SteamID   string     `json:"steam_id"   db:"steam_id,indexed,omitempty" valid:"required"`
		Name      string     `json:"name"       db:"name,omitempty"             valid:"required"`
		URL       string     `json:"url"        db:"url,omitempty"              valid:"required"`
		Avatar    string     `json:"avatar"     db:"avatar,omitempty"           valid:"required"`
		Status    UserStatus `json:"status"     db:"status,indexed,omitempty"`
		Notes     string     `json:"notes"      db:"notes,omitempty"`
		Donation  float64    `json:"donation"   db:"donation,omitempty"`
		DonatedAt *time.Time `json:"donated_at" db:"donated_at,omitempty"`
		CreatedAt *time.Time `json:"created_at" db:"created_at,omitempty"`
		UpdatedAt *time.Time `json:"updated_at" db:"updated_at,omitempty"`

		MarketStats MarketStatusCount `json:"market_stats" db:"market_stats,omitempty"`
		RankScore   int               `json:"rank_score"   db:"rank_score,omitempty"`

		Subscription       UserSubscription `json:"subscription"         db:"subscription,indexed,omitempty"`
		SubscribedAt       *time.Time       `json:"subscribed_at"        db:"subscribed_at,omitempty"`
		SubscriptionType   string           `json:"subscription_type"    db:"subscription_type"`
		SubscriptionEndsAt *time.Time       `json:"subscription_ends_at" db:"subscription_ends_at,omitempty"`
		Boons              []string         `json:"boons"                db:"boons,omitempty"`
		Hammer             bool             `json:"hammer"               db:"hammer,omitempty"`

		// paypal metadata
		Paypal struct {
			SubscriptionID        string    `json:"subscription_id"         db:"subscription_id"`
			SubscriptionLastPayed time.Time `json:"subscription_last_payed" db:"subscription_last_payed"`
		} `json:"paypal" db:"paypal"`
	}

	ManualSubscriptionParam struct {
		UserID string `json:"user_id"`
		Plan   string `json:"plan"`
		Cycles int    `json:"cycles"`
	}

	// userRepository defines operation for user records.
	userRepository interface {
		// Find returns a list of users from data store.
		Find(ctx context.Context, opts FindOpts) ([]User, error)

		// FindFlagged returns a list of flagged users from data store.
		FindFlagged(ctx context.Context, opts FindOpts) ([]User, error)

		// Get returns user details by id from data store.
		Get(ctx context.Context, id string) (*User, error)

		// Create persists a new user to data store.
		Create(context.Context, *User) error

		// Update persists user changes to data store.
		Update(context.Context, *User) error

		// BaseUpdate persists user changes to data store without updating metadata.
		BaseUpdate(context.Context, *User) error

		// ExpiringSubscribers return a list of users that has expiring subscription.
		ExpiringSubscribers(ctx context.Context, now time.Time) ([]User, error)

		// PurgeSubscription removes subscription data and boons.
		PurgeSubscription(ctx context.Context, userID string) error

		// ClearSubscriptionEndsAt clears subscription expiration.
		ClearSubscriptionEndsAt(ctx context.Context, id string) error
	}
)

// CheckCreate validates field on creating new user.
func (u User) CheckCreate() error {
	return validator.Struct(u)
}

// CheckUpdate validates field on update user.
func (u User) CheckUpdate() error {
	if u.ID == "" {
		return UserErrRequiredID
	}

	return nil
}

func (u User) TaskPriorityQueue() TaskPriority {
	switch u.Subscription {
	case UserSubscriptionPartner:
		return TaskPriorityHigh
	case UserSubscriptionTrader:
		return TaskPriorityMedium
	}
	return TaskPriorityLow
}

// CheckStatus checks for reported and banned status.
func (u User) CheckStatus() error {
	switch u.Status {
	case UserStatusSuspended:
		return UserErrSuspended
	case UserStatusBanned:
		return UserErrBanned
	}

	return nil
}

const (
	userScoreLiveRate        = 1
	userScoreReservedRate    = 2
	userScoreDeliveredRate   = 3
	userScoreBidRate         = 1
	userScoreBidCompleteRate = 4

	userScoreVerifiedInventoryRate      = 2
	userScoreVerifiedDeliveryNameRate   = 1
	userScoreVerifiedDeliverySenderRate = 7

	userScoreResellDeliveryRate = 3
)

// UserBoon represents user perks in an item form.
type UserBoon string

const (
	BoonSupporterBadge = "SUPPORTER_BADGE"
	BoonTraderBadge    = "TRADER_BADGE"
	BoonPartnerBadge   = "PARTNER_BADGE"

	BoonRefresherShard      = "REFRESHER_SHARD"
	BoonRefresherOrb        = "REFRESHER_ORB"
	BoonShopKeepersContract = "SHOPKEEPERS_CONTRACT"
	BoonDedicatedPos5       = "DEDICATED_POS_5"
)

// CalcRankScore return user score base on profile and market activity.
func (u User) CalcRankScore(stats MarketStatusCount) *User {
	u.RankScore = 1
	u.RankScore += (stats.Live - stats.ResellLive) * userScoreLiveRate
	u.RankScore += stats.Reserved * userScoreReservedRate
	u.RankScore += stats.Sold * userScoreDeliveredRate
	u.RankScore += stats.BidCompleted * userScoreBidCompleteRate

	u.RankScore += (stats.InventoryVerified - stats.ResellLive) * userScoreVerifiedInventoryRate
	u.RankScore += stats.DeliveryNameVerified * userScoreVerifiedDeliveryNameRate
	u.RankScore += stats.DeliverySenderVerified * userScoreVerifiedDeliverySenderRate

	u.RankScore += stats.ResellSold * userScoreResellDeliveryRate
	return &u
}

func (u User) HasBoon(ub UserBoon) bool {
	for _, b := range u.Boons {
		if ub == UserBoon(b) {
			return true
		}
	}
	return false
}

var userSubscriptionLabels = map[UserSubscription]string{
	UserSubscriptionSupporter: "SUPPORTER",
	UserSubscriptionTrader:    "TRADER",
	UserSubscriptionPartner:   "PARTNER",
}

var userSubscriptionBoons = map[UserSubscription][]string{
	UserSubscriptionSupporter: {
		BoonSupporterBadge,
		BoonRefresherShard,
	},
	UserSubscriptionTrader: {
		BoonTraderBadge,
		BoonRefresherShard,
		BoonRefresherOrb,
	},
	UserSubscriptionPartner: {
		BoonPartnerBadge,
		BoonRefresherOrb,
		BoonRefresherShard,
		BoonShopKeepersContract,
		BoonDedicatedPos5,
	},
}

func (s UserSubscription) String() string {
	l, ok := userSubscriptionLabels[s]
	if !ok {
		return ""
	}
	return l
}

func (s UserSubscription) Boons() []string {
	bb, ok := userSubscriptionBoons[s]
	if !ok {
		return nil
	}
	return bb
}

func UserSubscriptionFromString(s string) UserSubscription {
	for t, l := range userSubscriptionLabels {
		if strings.EqualFold(s, l) {
			return t
		}
	}
	return 0
}

// NewUserService returns a new User service.
func NewUserService(us userRepository, fm FileManager, sc paymentManager) *UserService {
	return &UserService{us, fm, sc}
}

type UserService struct {
	userRepo userRepository
	fileMgr  FileManager
	payment  paymentManager
}

func (s *UserService) Users(ctx context.Context, opts FindOpts) ([]User, error) {
	return s.userRepo.Find(ctx, opts)
}

func (s *UserService) FlaggedUsers(ctx context.Context, opts FindOpts) ([]User, error) {
	return s.userRepo.FindFlagged(ctx, opts)
}

func (s *UserService) User(ctx context.Context, id string) (*User, error) {
	return s.userRepo.Get(ctx, id)
}

func (s *UserService) UserFromContext(ctx context.Context) (*User, error) {
	au := AuthFromContext(ctx)
	if au == nil {
		return nil, UserErrNotFound
	}

	return s.User(ctx, au.UserID)
}

func (s *UserService) Create(ctx context.Context, u *User) error {
	url, err := s.downloadProfileImage(ctx, u.Avatar)
	if err != nil {
		return NewXError(UserErrProfileImageDL, err)
	}
	u.Avatar = url

	if err = u.CheckCreate(); err != nil {
		return err
	}

	go func() {
		if err = pingGoogleSitemap(); err != nil {
			log.Println("pingGoogleSitemap err:", err)
		}
	}()

	return s.userRepo.Create(ctx, u)
}

func (s *UserService) Update(ctx context.Context, u *User) error {
	au := AuthFromContext(ctx)
	if au == nil {
		return AuthErrNoAccess
	}
	u.ID = au.UserID

	if err := u.CheckUpdate(); err != nil {
		return err
	}

	return s.userRepo.Update(ctx, u)
}

func (s *UserService) SteamSync(ctx context.Context, sp *SteamPlayer) (*User, error) {
	u, err := s.userRepo.Get(ctx, sp.ID)
	if err != nil {
		return nil, err
	}

	u.Name = sp.Name
	u.URL = sp.URL
	u.Avatar, err = s.downloadProfileImage(ctx, sp.Avatar)
	if err != nil {
		return nil, err
	}
	if err = s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) CreateSubscription(ctx context.Context, planID string) (subscriptionID string, err error) {
	au := AuthFromContext(ctx)
	if au == nil {
		return "", AuthErrNoAccess
	}
	user, err := s.userRepo.Get(ctx, au.UserID)
	if err != nil {
		return "", err
	}

	subID, err := s.payment.CreateSubscription(ctx, planID, user.SteamID)
	if err != nil {
		return "", err
	}
	return subID, nil
}

func (s *UserService) ProcessSubscription(ctx context.Context, subscriptionID string) (*User, error) {
	au := AuthFromContext(ctx)
	if au == nil {
		return nil, AuthErrNoAccess
	}
	user, err := s.userRepo.Get(ctx, au.UserID)
	if err != nil {
		return nil, err
	}

	plan, steamID, subscriptionID, lastPayed, err := s.payment.Subscription(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	if user.SteamID != steamID {
		return nil, fmt.Errorf("could not validate subscription steam id: %s", steamID)
	}
	userSubs := UserSubscriptionFromString(plan)
	if userSubs == 0 {
		return nil, fmt.Errorf("could not validate subscription plan: %s", plan)
	}

	if user.SubscribedAt != nil && user.Subscription == userSubs {
		return user, nil
	}

	t := time.Now()
	user.Subscription = userSubs
	user.SubscribedAt = &t
	user.Boons = userSubs.Boons()
	user.SubscriptionType = "paypal"
	user.Paypal.SubscriptionID = subscriptionID
	user.Paypal.SubscriptionLastPayed = lastPayed
	if err = user.CheckUpdate(); err != nil {
		return nil, err
	}
	if err = s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	if err = s.userRepo.ClearSubscriptionEndsAt(ctx, user.ID); err != nil {
		return nil, err
	}

	user.SubscriptionEndsAt = nil
	return user, nil
}

// UpdateSubscriptionFromWebhook manage updates from webhook payload, most often use in incrementing cycles or
// extending expiration.
func (s *UserService) UpdateSubscriptionFromWebhook(ctx context.Context, r *http.Request) (*User, error) {
	// get user by steam id and increment their cycles.
	steamID, cancelled, lastPayment, err := s.payment.IsCancelled(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("checking cancelled subscription: %v", err)
	}
	if !cancelled {
		// ignore if not canceled.
		log.Println("ignoring subscription update because its not cancelled:", steamID)
		return nil, nil
	}

	log.Println("cancelling subscription", steamID, "by marking expiration")
	user, err := s.userRepo.Get(ctx, steamID)
	if err != nil {
		return nil, fmt.Errorf("getting user %s: %w", steamID, err)
	}

	ex := lastPayment.AddDate(0, 1, 0)
	user.SubscriptionEndsAt = &ex
	if err = s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("updating user: %v", err)
	}
	return user, nil
}

// ProcessManualSubscription process manual subscription such as one-time payments that process manually, normally
// in bulk and steam items. This function will be used for non-recurring payments. ex:
//
//		Manual Partner subscription:
//	    - 3 months (+60% overhead)
//	    - 6 months (+60% overhead)
//	    - 12 months (+60% overhead)
func (s *UserService) ProcessManualSubscription(ctx context.Context, param ManualSubscriptionParam) (*User, error) {
	user, err := s.userRepo.Get(ctx, param.UserID)
	if err != nil {
		return nil, fmt.Errorf("getting user: %v", err)
	}

	subs := UserSubscriptionFromString(param.Plan)
	user.Subscription = subs
	user.Boons = subs.Boons()
	user.SubscriptionType = "manual"

	now := time.Now()
	end := now.AddDate(0, param.Cycles, 0)
	user.SubscribedAt = &now
	user.SubscriptionEndsAt = &end
	if err = s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("updating user: %v", err)
	}
	return user, nil
}

// downloadProfileImage saves an image file from url.
func (s *UserService) downloadProfileImage(ctx context.Context, url string) (filename string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	uu := strings.Split(url, "/")
	name := uu[len(uu)-1]
	name = strings.TrimSuffix(name, ".jpg")
	name = strings.TrimSuffix(name, "_full")
	filename, err = s.fileMgr.SaveWithName(resp.Body, name)
	if err != nil {
		return "", err
	}
	return filename, nil
}

type paymentManager interface {
	Subscription(ctx context.Context, id string) (plan, steamID, subscriptionID string, lastPayment time.Time, err error)
	IsCancelled(ctx context.Context, r *http.Request) (steamID string, cancelled bool, lastPayment time.Time, err error)
	CreateSubscription(ctx context.Context, planID, customID string) (subscriptionID string, err error)
}
