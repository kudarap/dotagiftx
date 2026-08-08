package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/kudarap/dotagiftx"
)

const dayHours = time.Hour * 24

// ExpiringMarket represents setting expiration of a market entry job.
type ExpiringMarket struct {
	marketRepo  marketRepository
	catalogRepo catalogRepository
	cache       cacheRemover
	logger      *slog.Logger
	// job settings
	name     string
	interval time.Duration
}

func NewExpiringMarket(ms marketRepository, cs catalogRepository, cc cacheRemover, lg *slog.Logger) *ExpiringMarket {
	return &ExpiringMarket{
		marketRepo:  ms,
		catalogRepo: cs,
		cache:       cc,
		logger:      lg,
		name:        "expiring_market",
		interval:    time.Hour * 24,
	}
}

func (em *ExpiringMarket) String() string { return em.name }

func (em *ExpiringMarket) Interval() time.Duration { return em.interval }

func (em *ExpiringMarket) Run(ctx context.Context) error {
	var itemIDs []string
	now := time.Now()

	// Process expiring bids.
	bidExpr := now.Add(-dayHours * dotagiftx.MarketBidExpirationDays)
	em.logger.Info("updating expiring bids", "cutoff", bidExpr)
	ids, err := em.marketRepo.UpdateExpiring(ctx, dotagiftx.MarketTypeBid, dotagiftx.BoonRefresherShard, bidExpr)
	if err != nil {
		em.logger.Error("could not update expiring bids", "error", err)
		return err
	}
	itemIDs = append(itemIDs, ids...)
	em.logger.Info("updating expiring bids finished!")

	// Process expiring asks.
	askExpr := now.Add(-dayHours * dotagiftx.MarketAskExpirationDays)
	em.logger.Info("updating expiring asks", "cutoff", askExpr)
	ids, err = em.marketRepo.UpdateExpiring(ctx, dotagiftx.MarketTypeAsk, dotagiftx.BoonRefresherOrb, askExpr)
	if err != nil {
		em.logger.Error("could not update expiring asks", "error", err)
		return err
	}
	itemIDs = append(itemIDs, ids...)
	em.logger.Info("updating expiring asks finished!")

	// Process expiring resells.
	em.logger.Info("updating expiring resells", "cutoff", askExpr)
	ids, err = em.marketRepo.UpdateExpiringResell(ctx, dotagiftx.BoonShopKeepersContract)
	if err != nil {
		em.logger.Error("could not update expiring resells", "error", err)
		return err
	}
	itemIDs = append(itemIDs, ids...)
	em.logger.Info("updating expiring resells finished!")

	// Re-index affected items.
	em.logger.Info("indexing affected expire items", "count", len(itemIDs))
	itemIndexed := map[string]struct{}{}
	for _, id := range itemIDs {
		if _, hit := itemIndexed[id]; hit {
			continue
		}
		itemIndexed[id] = struct{}{}

		if _, err = em.catalogRepo.Index(ctx, id); err != nil {
			em.logger.Error("could not index expired item", "error", err)
			continue
		}
	}
	em.logger.Info("affected items indexed", "count", len(itemIndexed))

	// Invalidate market caches.
	em.logger.Info("invalidating market cache...")
	if err = em.cache.BulkDel("catalogs_trend"); err != nil {
		em.logger.Error("could not perform bulk delete on catalog trend cache", "error", err)
		return err
	}
	// svc_market market is the prefixed used for caching market related data.
	if err = em.cache.BulkDel("svc_market"); err != nil {
		em.logger.Error("could not perform bulk delete on market cache", "error", err)
		return err
	}
	em.logger.Info("market cache invalidated!")
	return nil
}
