package jobs

import (
	"context"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/log"
)

// SweepMarket represents setting expiration of a market entry job.
type SweepMarket struct {
	marketStg dgx.MarketStorage
	logger    log.Logger
	// job settings
	name     string
	interval time.Duration
}

func NewSweepMarket(ms dgx.MarketStorage, lg log.Logger) *SweepMarket {
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
	t := now.Add(-dayHours * dgx.MarketSweepExpiredDays)
	cm.logger.Println("sweeping old expired market", t)
	if err := cm.marketStg.BulkDeleteByStatus(dgx.MarketStatusExpired, t, limitPerBatch); err != nil {
		cm.logger.Errorf("could not clean expired market: %s", err)
		return err
	}
	cm.logger.Println("sweeping old expired market finished!")

	// Clean up removed markets.
	t = now.Add(-dayHours * dgx.MarketSweepRemovedDays)
	cm.logger.Println("sweeping old removed market", t)
	if err := cm.marketStg.BulkDeleteByStatus(dgx.MarketStatusRemoved, t, limitPerBatch); err != nil {
		cm.logger.Errorf("could not clean removed market: %s", err)
		return err
	}
	cm.logger.Println("sweeping old removed market finished!")

	return nil
}
