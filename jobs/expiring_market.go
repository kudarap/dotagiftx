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

func (em *ExpiringMarket) String() string { return em.name }

func (em *ExpiringMarket) Interval() time.Duration { return em.interval }

func (em *ExpiringMarket) Run(ctx context.Context) error {
	var itemIDs []string
	now := time.Now()

	// Process expiring bids.
	bidExpr := now.Add(-dayHours * core.MarketBidExpirationDays)
	em.logger.Println("updating expiring bids", bidExpr)
	ids, err := em.marketStg.UpdateExpiring(core.MarketTypeBid, core.BoonRefresherShard, bidExpr)
	if err != nil {
		em.logger.Errorf("could not update expiring bids: %s", err)
		return err
	}
	itemIDs = append(itemIDs, ids...)
	em.logger.Println("updating expiring bids finished!")

	// Process expiring asks.
	askExpr := now.Add(-dayHours * core.MarketAskExpirationDays)
	em.logger.Println("updating expiring asks", askExpr)
	ids, err = em.marketStg.UpdateExpiring(core.MarketTypeAsk, core.BoonRefresherOrb, askExpr)
	if err != nil {
		em.logger.Errorf("could not update expiring asks: %s", err)
		return err
	}
	itemIDs = append(itemIDs, ids...)
	em.logger.Println("updating expiring asks finished!")

	// Re-index affected items.
	em.logger.Println("indexing affected expire items...", len(itemIDs))
	itemIndexed := map[string]struct{}{}
	for _, id := range itemIDs {
		if _, hit := itemIndexed[id]; hit {
			continue
		}
		itemIndexed[id] = struct{}{}

		if _, err = em.catalogStg.Index(id); err != nil {
			em.logger.Errorf("could not index expired item: %s", err)
			continue
		}
	}
	em.logger.Println("affected items indexed!", len(itemIndexed))

	// Invalidate market caches.
	em.logger.Println("invalidating market cache...")
	if err = em.cache.BulkDel("catalogs_trend"); err != nil {
		em.logger.Errorf("could not perform bulk delete on catalog trend cache: %s", err)
		return err
	}
	// svc_market market is the prefixed used for caching market related data.
	if err = em.cache.BulkDel("svc_market"); err != nil {
		em.logger.Errorf("could not perform bulk delete on market cache: %s", err)
		return err
	}
	em.logger.Println("market cache invalidated!")
	return nil
}
