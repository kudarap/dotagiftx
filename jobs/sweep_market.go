package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/gokit/log"
)

// SweepMarket represents setting expiration of a market entry job.
type SweepMarket struct {
	marketStg core.MarketStorage
	logger    log.Logger
	// job settings
	name     string
	interval time.Duration
}

func NewSweepMarket(ms core.MarketStorage, lg log.Logger) *SweepMarket {
	return &SweepMarket{
		marketStg: ms,
		logger:    lg,
		name:      "clean_market",
		interval:  defaultJobInterval,
	}
}

func (cm *SweepMarket) String() string { return cm.name }

func (cm *SweepMarket) Interval() time.Duration { return cm.interval }

func (cm *SweepMarket) Run(ctx context.Context) error {
	const limitPerBatch = 1000
	now := time.Now()

	// Clean up expiring markets.
	t := now.Add(-dayHours * core.MarketSweepExpiredDays)
	cm.logger.Println("sweeping old expired market", t)
	if err := cm.marketStg.BulkDeleteByStatus(core.MarketStatusExpired, t, limitPerBatch); err != nil {
		cm.logger.Errorf("could not clean expired market: %s", err)
		return err
	}
	cm.logger.Println("sweeping old expired market finished!")

	// Clean up removed markets.
	t = now.Add(-dayHours * core.MarketSweepRemovedDays)
	cm.logger.Println("sweeping old removed market", t)
	if err := cm.marketStg.BulkDeleteByStatus(core.MarketStatusRemoved, t, limitPerBatch); err != nil {
		cm.logger.Errorf("could not clean removed market: %s", err)
		return err
	}
	cm.logger.Println("sweeping old removed market finished!")

	return nil
}
