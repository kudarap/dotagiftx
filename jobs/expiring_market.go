package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/gokit/log"
)

const dayHours = time.Hour * 24

// ExpiringMarket represents setting expiration of a market entry job.
type ExpiringMarket struct {
	marketStg core.MarketStorage
	cache     core.Cache
	logger    log.Logger
	// job settings
	name     string
	interval time.Duration
}

func NewExpiringMarket(ms core.MarketStorage, cc core.Cache, lg log.Logger) *ExpiringMarket {
	return &ExpiringMarket{
		marketStg: ms,
		cache:     cc,
		logger:    lg,
		name:      "expiring_market",
		interval:  defaultJobInterval,
	}
}

func (m *ExpiringMarket) String() string { return m.name }

func (m *ExpiringMarket) Interval() time.Duration { return m.interval }

func (m *ExpiringMarket) Run(ctx context.Context) error {
	now := time.Now()

	bidExpr := now.Add(-dayHours * core.MarketBidExpirationDays)
	m.logger.Println("updating expiring bids", bidExpr)
	if err := m.marketStg.UpdateExpiring(core.MarketTypeBid, core.BoonRefresherShard, bidExpr); err != nil {
		m.logger.Errorf("could not update expiring bids: %s", err)
		return err
	}
	m.logger.Println("updating expiring bids finished!")

	askExpr := now.Add(-dayHours * core.MarketAskExpirationDays)
	m.logger.Println("updating expiring asks", askExpr)
	if err := m.marketStg.UpdateExpiring(core.MarketTypeAsk, core.BoonRefresherOrb, askExpr); err != nil {
		m.logger.Errorf("could not update expiring asks: %s", err)
		return err
	}
	m.logger.Println("updating expiring asks finished!")

	m.logger.Println("invalidating market cache...")
	// svc_market market is the prefixed used for caching market related data.
	if err := m.cache.BulkDel("svc_market"); err != nil {
		m.logger.Errorf("could not perform bulk delete on cache: %s", err)
		return err
	}
	m.logger.Println("market cache invalidated!")
	return nil
}
