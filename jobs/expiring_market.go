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
	marketStg  core.MarketStorage
	catalogStg core.CatalogStorage
	cache      core.Cache
	logger     log.Logger
	// job settings
	name     string
	interval time.Duration
}

func NewExpiringMarket(ms core.MarketStorage, cs core.CatalogStorage, cc core.Cache, lg log.Logger) *ExpiringMarket {
	return &ExpiringMarket{
		marketStg:  ms,
		catalogStg: cs,
		cache:      cc,
		logger:     lg,
		name:       "expiring_market",
		interval:   defaultJobInterval,
	}
}

func (m *ExpiringMarket) String() string { return m.name }

func (m *ExpiringMarket) Interval() time.Duration { return m.interval }

func (m *ExpiringMarket) Run(ctx context.Context) error {
	var itemIDs []string
	now := time.Now()

	bidExpr := now.Add(-dayHours * core.MarketBidExpirationDays)
	m.logger.Println("updating expiring bids", bidExpr)
	ids, err := m.marketStg.UpdateExpiring(core.MarketTypeBid, core.BoonRefresherShard, bidExpr)
	if err != nil {
		m.logger.Errorf("could not update expiring bids: %s", err)
		return err
	}
	itemIDs = append(itemIDs, ids...)
	m.logger.Println("updating expiring bids finished!")

	askExpr := now.Add(-dayHours * core.MarketAskExpirationDays)
	m.logger.Println("updating expiring asks", askExpr)
	ids, err = m.marketStg.UpdateExpiring(core.MarketTypeAsk, core.BoonRefresherOrb, askExpr)
	if err != nil {
		m.logger.Errorf("could not update expiring asks: %s", err)
		return err
	}
	itemIDs = append(itemIDs, ids...)
	m.logger.Println("updating expiring asks finished!")

	m.logger.Println("indexing affected expire items...", len(itemIDs))
	itemIndexed := map[string]struct{}{}
	for _, id := range itemIDs {
		if _, hit := itemIndexed[id]; hit {
			continue
		}
		itemIndexed[id] = struct{}{}

		if _, err = m.catalogStg.Index(id); err != nil {
			m.logger.Errorf("could not index expired item: %s", err)
			continue
		}
	}
	m.logger.Println("affected items indexed!", len(itemIndexed))

	m.logger.Println("invalidating market cache...")
	if err = m.cache.BulkDel("catalogs_trend"); err != nil {
		m.logger.Errorf("could not perform bulk delete on catalog trend cache: %s", err)
		return err
	}
	// svc_market market is the prefixed used for caching market related data.
	if err = m.cache.BulkDel("svc_market"); err != nil {
		m.logger.Errorf("could not perform bulk delete on market cache: %s", err)
		return err
	}
	m.logger.Println("market cache invalidated!")
	return nil
}
