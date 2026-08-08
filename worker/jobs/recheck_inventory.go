package jobs

import (
	"context"
	"time"

	"log/slog"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/verify"
)

// RecheckInventory represents a job that rechecks no-hit items.
// Crawling from SteamInventory.org tends to fail sometimes.
type RecheckInventory struct {
	inventorySvc dotagiftx.InventoryService
	marketStg    dotagiftx.MarketStorage
	source       *verify.Source
	logger       *slog.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dotagiftx.Inventory
}

func NewRecheckInventory(
	is dotagiftx.InventoryService,
	ms dotagiftx.MarketStorage,
	as *verify.Source,
	lg *slog.Logger,
) *RecheckInventory {
	f := dotagiftx.Inventory{Status: dotagiftx.InventoryStatusNoHit}
	return &RecheckInventory{
		is, ms, as, lg,
		"recheck_inventory", time.Hour, f}
}

func (ri *RecheckInventory) String() string { return ri.name }

func (ri *RecheckInventory) Interval() time.Duration { return ri.interval }

func (ri *RecheckInventory) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		ri.logger.Info("RECHECK INVENTORY BENCHMARK TIME", "elapsed", time.Since(bs))
	}()

	opts := dotagiftx.FindOpts{Filter: ri.filter}
	opts.Sort = "updated_at:desc"
	// opts.Limit = 10
	opts.Page = 0
	opts.IndexKey = "status"

	invs, _, err := ri.inventorySvc.Inventories(ctx, opts)
	if err != nil {
		return err
	}

	for _, ii := range invs {
		start := time.Now()

		if ii.RetriesExceeded() {
			continue
		}

		mkt, _ := ri.market(ctx, ii.MarketID)
		if mkt == nil {
			continue
		}

		if mkt.User == nil || mkt.Item == nil {
			ri.logger.Error("skipped process! missing data", "user", mkt.User, "item", mkt.Item)
			continue
		}

		result, err := ri.source.Inventory(ctx, mkt.User.SteamID, mkt.Item.Name)
		if err != nil {
			ri.logger.Error("skipped process! source error", "user", mkt.User, "item", mkt.Item, "error", err)
			continue
		}

		ri.logger.Info("batch", "page", opts.Page, "steam_id", mkt.User.SteamID, "item", mkt.Item.Name, "status", result.Status)
		err = ri.inventorySvc.Set(ctx, &dotagiftx.Inventory{
			MarketID:   mkt.ID,
			Status:     result.Status,
			Assets:     result.Assets,
			VerifiedBy: result.VerifiedBy,
			ElapsedMs:  time.Since(start).Milliseconds(),
		})
		if err != nil {
			ri.logger.Error("inventory update error", "steam_id", mkt.User.SteamID, "item", mkt.Item.Name, "status", result.Status, "error", err)
		}

		time.Sleep(time.Second / 4)
	}

	return nil
}

func (ri *RecheckInventory) market(ctx context.Context, id string) (*dotagiftx.Market, error) {
	f := dotagiftx.FindOpts{Filter: dotagiftx.Market{ID: id}}
	markets, err := ri.marketStg.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(markets) == 0 {
		return nil, nil
	}
	mkt := markets[0]
	return &mkt, nil
}
