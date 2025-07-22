package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/logging"
)

type ExpiringSubscription struct {
	userStg dotagiftx.UserStorage
	cache   cache
	logger  logging.Logger
	// job settings
	name     string
	interval time.Duration
}

func NewExpiringSubscription(
	us dotagiftx.UserStorage,
	cache cache,
	lg logging.Logger,
) *ExpiringSubscription {
	return &ExpiringSubscription{
		userStg:  us,
		cache:    cache,
		name:     "expiring_subscription",
		interval: time.Hour * 24,
		logger:   lg,
	}
}

// Run removes subscription status base on its expiration.
func (s *ExpiringSubscription) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		s.logger.Println("EXPIRING SUBSCRIPTION BENCHMARK TIME", time.Since(bs))
	}()

	// get all users that has subscription
	// add leeway of 2 days to process recurring payment.
	// check outstanding days if still valid from last payment and skip.
	withLeeway := time.Now().AddDate(0, 0, -2)
	users, err := s.userStg.ExpiringSubscribers(ctx, withLeeway)
	if err != nil {
		return fmt.Errorf("retrieving subscribers: %w", err)
	}

	// remove boons and subs status
	// clear user cache
	for _, u := range users {
		if err = s.userStg.PurgeSubscription(ctx, u.ID); err != nil {
			s.logger.Errorf("purging subscription: %w", err)
		}

		go func() {
			if err := s.cache.BulkDel(fmt.Sprintf("users/%s*", u.SteamID)); err != nil {
				s.logger.Errorf("invalidate user cache: %w", err)
			}
		}()
	}

	return nil
}

func (s *ExpiringSubscription) String() string {
	return s.name
}

func (s *ExpiringSubscription) Interval() time.Duration {
	return s.interval
}
