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
	UserErrNotFound Errors = iota + 1100
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
	UserSubscriptionResell    UserSubscription = 1
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
	}

	ManualSubscriptionParam struct {
		UserID string `json:"user_id"`
		Plan   string `json:"plan"`
		Cycles int    `json:"cycles"`
	}

	// UserService provides access to user service.
	UserService interface {
		// Users returns a list of users.
		Users(opts FindOpts) ([]User, error)

		// FlaggedUsers returns a list of flagged/reported users.
		FlaggedUsers(opts FindOpts) ([]User, error)

		// User returns user details by id.
		User(id string) (*User, error)

		// Create saves new user and download profile image to local file.
		Create(*User) error

		// UserFromContext returns user details from context.
		UserFromContext(context.Context) (*User, error)

		// Update saves user changes.
		Update(context.Context, *User) error

		// SteamSync saves updated steam info.
		SteamSync(sp *SteamPlayer) (*User, error)

		// ProcessSubscription validates and process subscription features.
		ProcessSubscription(ctx context.Context, subscriptionID string) (*User, error)

		// UpdateSubscriptionFromWebhook handles user subscription updates form http request.
		UpdateSubscriptionFromWebhook(ctx context.Context, r *http.Request) (*User, error)

		ProcessManualSubscription(ctx context.Context, form ManualSubscriptionParam) (*User, error)
	}

	// UserStorage defines operation for user records.
	UserStorage interface {
		// Find returns a list of users from data store.
		Find(opts FindOpts) ([]User, error)

		// FindFlagged returns a list of flagged users from data store.
		FindFlagged(opts FindOpts) ([]User, error)

		// Get returns user details by id from data store.
		Get(id string) (*User, error)

		// Create persists a new user to data store.
		Create(*User) error

		// Update persists user changes to data store.
		Update(*User) error

		// BaseUpdate persists user changes to data store without updating metadata.
		BaseUpdate(*User) error

		// ExpiringSubscribers return a list of users that has expiring subscription.
		ExpiringSubscribers(ctx context.Context, now time.Time) ([]User, error)

		// PurgeSubscription removes subscription data and boons.
		PurgeSubscription(ctx context.Context, userID string) error
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
func NewUserService(us UserStorage, fm FileManager, sc subscriptionChecker) UserService {
	return &userService{us, fm, sc}
}

type userService struct {
	userStg     UserStorage
	fileMgr     FileManager
	subsChecker subscriptionChecker
}

func (s *userService) Users(opts FindOpts) ([]User, error) {
	return s.userStg.Find(opts)
}

func (s *userService) FlaggedUsers(opts FindOpts) ([]User, error) {
	return s.userStg.FindFlagged(opts)
}

func (s *userService) User(id string) (*User, error) {
	return s.userStg.Get(id)
}

func (s *userService) UserFromContext(ctx context.Context) (*User, error) {
	au := AuthFromContext(ctx)
	if au == nil {
		return nil, UserErrNotFound
	}

	return s.User(au.UserID)
}

func (s *userService) Create(u *User) error {
	url, err := s.downloadProfileImage(u.Avatar)
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

	return s.userStg.Create(u)
}

func (s *userService) Update(ctx context.Context, u *User) error {
	au := AuthFromContext(ctx)
	if au == nil {
		return AuthErrNoAccess
	}
	u.ID = au.UserID

	if err := u.CheckUpdate(); err != nil {
		return err
	}

	return s.userStg.Update(u)
}

func (s *userService) SteamSync(sp *SteamPlayer) (*User, error) {
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

func (s *userService) ProcessSubscription(ctx context.Context, subscriptionID string) (*User, error) {
	au := AuthFromContext(ctx)
	if au == nil {
		return nil, AuthErrNoAccess
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
	userSubs := UserSubscriptionFromString(plan)
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
func (s *userService) UpdateSubscriptionFromWebhook(ctx context.Context, r *http.Request) (*User, error) {
	// get user by steam id and increment their cycles.
	steamID, cancelled, err := s.subsChecker.IsCancelled(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("checking cancelled subscription: %v", err)
	}
	if !cancelled {
		// ignore if not canceled.
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
	if user.Subscription == UserSubscriptionPartner {
		t := time.Now()
		user.SubscriptionEndsAt = &t
	}

	if err = s.userStg.Update(user); err != nil {
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
func (s *userService) ProcessManualSubscription(ctx context.Context, param ManualSubscriptionParam) (*User, error) {
	user, err := s.userStg.Get(param.UserID)
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
	if err = s.userStg.Update(user); err != nil {
		return nil, fmt.Errorf("updating user: %v", err)
	}
	return user, nil
}

// downloadProfileImage saves an image file from url.
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
