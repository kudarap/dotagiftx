package jobs

import (
	"context"
	"fmt"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/log"
)

type ExpiringSubscription struct {
	userStg dgx.UserStorage
	cache   dgx.Cache
	logger  log.Logger
	// job settings
	name     string
	interval time.Duration
}

func NewExpiringSubscription(
	us dgx.UserStorage,
	cache dgx.Cache,
	lg log.Logger,
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
	// check outstanding days if it's still validate from last payment and skip.
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
