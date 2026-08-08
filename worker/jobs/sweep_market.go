package jobs

import (
	"context"
	"time"

	"log/slog"

	"github.com/kudarap/dotagiftx"
)

// SweepMarket represents setting expiration of a market entry job.
type SweepMarket struct {
	marketStg marketStorage
	logger    *slog.Logger
	// job settings
	name     string
	interval time.Duration
}

func NewSweepMarket(ms marketStorage, lg *slog.Logger) *SweepMarket {
	return &SweepMarket{
		marketStg: ms,
		logger:    lg,
		name:      "clean_market",
		interval:  time.Hour * 24,
	}
}

func (cm *SweepMarket) String() string { return cm.name }

func (cm *SweepMarket) Interval() time.Duration { return cm.interval }

func (cm *SweepMarket) Run(ctx context.Context) error {
	const limitPerBatch = 1000
	now := time.Now()

	// Clean up expiring markets.
	t := now.Add(-dayHours * dotagiftx.MarketSweepExpiredDays)
	cm.logger.Info("sweeping old expired market", "cutoff", t)
	if err := cm.marketStg.BulkDeleteByStatus(ctx, dotagiftx.MarketStatusExpired, t, limitPerBatch); err != nil {
		cm.logger.Error("could not clean expired market", "error", err)
		return err
	}
	cm.logger.Info("sweeping old expired market finished!")

	// Clean up removed markets.
	t = now.Add(-dayHours * dotagiftx.MarketSweepRemovedDays)
	cm.logger.Info("sweeping old removed market", "cutoff", t)
	if err := cm.marketStg.BulkDeleteByStatus(ctx, dotagiftx.MarketStatusRemoved, t, limitPerBatch); err != nil {
		cm.logger.Error("could not clean removed market", "error", err)
		return err
	}
	cm.logger.Info("sweeping old removed market finished!")

	return nil
}
