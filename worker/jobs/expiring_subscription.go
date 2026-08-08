package jobs

import (
	"context"
	"fmt"
	"time"

	"log/slog"
)

type ExpiringSubscription struct {
	userRepo userRepository
	cache    cacheRemover
	logger   *slog.Logger
	// job settings
	name     string
	interval time.Duration
}

func NewExpiringSubscription(
	us userRepository,
	cache cacheRemover,
	lg *slog.Logger,
) *ExpiringSubscription {
	return &ExpiringSubscription{
		userRepo: us,
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
		s.logger.Info("EXPIRING SUBSCRIPTION BENCHMARK TIME", "elapsed", time.Since(bs))
	}()

	// get all users that has subscription
	// add leeway of 2 days to process recurring payment.
	// check outstanding days if still valid from last payment and skip.
	withLeeway := time.Now().AddDate(0, 0, -2)
	users, err := s.userRepo.ExpiringSubscribers(ctx, withLeeway)
	if err != nil {
		return fmt.Errorf("retrieving subscribers: %w", err)
	}

	// remove boons and subs status
	// clear user cache
	for _, u := range users {
		if err = s.userRepo.PurgeSubscription(ctx, u.ID); err != nil {
			s.logger.Error("purging subscription", "user_id", u.ID, "error", err)
		}

		go func() {
			if err := s.cache.BulkDel(fmt.Sprintf("users/%s*", u.SteamID)); err != nil {
				s.logger.Error("invalidate user cache", "steam_id", u.SteamID, "error", err)
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
